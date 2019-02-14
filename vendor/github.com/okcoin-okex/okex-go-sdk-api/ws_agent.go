package okex

/*
 OKEX websocket API agent
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import (
	"bytes"
	"compress/flate"
	"github.com/gorilla/websocket"
	"io/ioutil"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type OKWSAgent struct {
	baseUrl string
	config  *Config
	conn    *websocket.Conn

	wsEvtCh  chan interface{}
	wsErrCh  chan interface{}
	wsTbCh   chan interface{}
	stopCh   chan interface{}
	errCh    chan error
	signalCh chan os.Signal

	subMap         map[string][]ReceivedDataCallback
	activeChannels map[string]bool

	processMut sync.Mutex
}

func (a *OKWSAgent) Start(config *Config) error {
	a.baseUrl = config.WSEndpoint + "ws/v3?compress=true"
	log.Printf("Connecting to %s", a.baseUrl)
	c, _, err := websocket.DefaultDialer.Dial(a.baseUrl, nil)

	if err != nil {
		log.Fatalf("dial:%+v", err)
		return err
	} else {
		log.Printf("Connected to %s", a.baseUrl)
		a.conn = c
		a.config = config

		a.wsEvtCh = make(chan interface{})
		a.wsErrCh = make(chan interface{})
		a.wsTbCh = make(chan interface{})
		a.errCh = make(chan error)
		a.stopCh = make(chan interface{}, 16)
		a.signalCh = make(chan os.Signal)
		a.activeChannels = make(map[string]bool)
		a.subMap = make(map[string][]ReceivedDataCallback)

		signal.Notify(a.signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		go a.work()
		go a.receive()
		go a.finalize()
	}

	return nil
}

func (a *OKWSAgent) Subscribe(channel, filter string, cb ReceivedDataCallback) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	st := SubscriptionTopic{channel, filter}
	bo, err := subscribeOp([]*SubscriptionTopic{&st})
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	log.Printf("Send Msg: %s", msg)
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return err
	}

	cbs := a.subMap[st.channel]
	if cbs == nil {
		cbs = []ReceivedDataCallback{}
		a.activeChannels[st.channel] = false
	}
	cbs = append(cbs, cb)
	a.subMap[st.channel] = cbs
	fullTopic, err := st.ToString()
	a.subMap[fullTopic] = cbs

	return nil
}

func (a *OKWSAgent) UnSubscribe(channel, filter string) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	st := SubscriptionTopic{channel, filter}
	bo, err := unsubscribeOp([]*SubscriptionTopic{&st})
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	log.Printf("Send Msg: %s", msg)
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return err
	}

	a.subMap[channel] = nil
	a.activeChannels[channel] = false

	return nil
}

func (a *OKWSAgent) Login(apiKey, passphrase string) error {

	timestamp := EpochTime()

	preHash := PreHashString(timestamp, GET, "/users/self/verify", "")
	if sign, err := HmacSha256Base64Signer(preHash, a.config.SecretKey); err != nil {
		return err
	} else {
		op, err := loginOp(apiKey, passphrase, timestamp, sign)
		data, err := Struct2JsonString(op)
		log.Printf("Send Msg: %s", data)
		err = a.conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}

func (a *OKWSAgent) keepalive() {
	a.ping()
}

func (a *OKWSAgent) Stop() error {
	defer func() {
		a := recover()
		log.Println(a)
	}()

	a.stopCh <- nil
	return nil
}

func (a *OKWSAgent) finalize() error {
	defer func() {
		log.Printf("Connection to WebSocket is closed.")
	}()

	select {
	case <-a.stopCh:
		if a.conn != nil {
			close(a.stopCh)
			close(a.errCh)
			close(a.wsTbCh)
			close(a.wsEvtCh)
			close(a.wsErrCh)
			return a.conn.Close()
		}
	}

	return nil
}

func (a *OKWSAgent) ping() {
	msg := "ping"
	//log.Printf("Send Msg: %s", msg)
	a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (a *OKWSAgent) GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func (a *OKWSAgent) handleErrResponse(r interface{}) error {
	log.Printf("handleErrResponse %+v \n", r)
	return nil
}

func (a *OKWSAgent) handleEventResponse(r interface{}) error {
	er := r.(*WSEventResponse)
	a.activeChannels[er.Channel] = (er.Event == CHNL_EVENT_SUBSCRIBE)
	return nil
}

func (a *OKWSAgent) handleTableResponse(r interface{}) error {
	tr := r.(*WSTableResponse)
	cbs := a.subMap[tr.Table]
	if cbs != nil {
		for i := 0; i < len(cbs); i++ {
			cb := cbs[i]
			if err := cb(tr); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *OKWSAgent) work() {
	defer func() {
		a := recover()
		log.Println(a)
	}()

	defer a.Stop()

	ticker := time.NewTicker(29 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.keepalive()
		case errR := <-a.wsErrCh:
			a.handleErrResponse(errR)
		case evtR := <-a.wsEvtCh:
			a.handleEventResponse(evtR)
		case tb := <-a.wsTbCh:
			a.handleTableResponse(tb)
		case <-a.signalCh:
			break
		case err := <-a.errCh:
			DefaultDataCallBack(err)
			break
		}
	}
}

func (a *OKWSAgent) receive() {
	defer func() {
		a := recover()
		log.Println(a)
	}()

	for {
		messageType, message, err := a.conn.ReadMessage()
		if err != nil {
			a.errCh <- err
			break
		}

		txtMsg := message
		switch messageType {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			txtMsg, err = a.GzipDecode(message)
		}

		rsp, err := loadResponse(txtMsg)
		if rsp != nil {
			log.Printf("LoadedRep: %+v, err: %+v", rsp, err)
		} else {
			log.Printf("TextMessg: %s", txtMsg)
		}

		if err != nil {
			break
		}

		switch rsp.(type) {
		case *WSErrorResponse:
			a.wsErrCh <- rsp
		case *WSEventResponse:
			er := rsp.(*WSEventResponse)
			a.wsEvtCh <- er
		case *WSTableResponse:
			tb := rsp.(*WSTableResponse)
			a.wsTbCh <- tb
		default:
			//log.Println(rsp)
		}
	}
}
