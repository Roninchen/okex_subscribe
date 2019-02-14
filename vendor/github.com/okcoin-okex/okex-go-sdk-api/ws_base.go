package okex

/*
 OKEX websocket api wrapper
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import (
	"fmt"
)

type BaseOp struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

func subscribeOp(sts []*SubscriptionTopic) (op *BaseOp, err error) {

	strArgs := []string{}

	for i := 0; i < len(sts); i++ {
		channel, err := sts[i].ToString()
		if err != nil {
			return nil, err
		}
		strArgs = append(strArgs, channel)
	}

	b := BaseOp{
		Op:   "subscribe",
		Args: strArgs,
	}
	return &b, nil
}

func unsubscribeOp(sts []*SubscriptionTopic) (op *BaseOp, err error) {

	strArgs := []string{}

	for i := 0; i < len(sts); i++ {
		channel, err := sts[i].ToString()
		if err != nil {
			return nil, err
		}
		strArgs = append(strArgs, channel)
	}

	b := BaseOp{
		Op:   "unsubscribe",
		Args: strArgs,
	}
	return &b, nil
}

func loginOp(apiKey string, passphrase string, timestamp string, sign string) (op *BaseOp, err error) {
	b := BaseOp{
		Op:   "login",
		Args: []string{apiKey, passphrase, timestamp, sign},
	}
	return &b, nil
}

type SubscriptionTopic struct {
	channel string
	filter  string `default:""`
}

func (st *SubscriptionTopic) ToString() (topic string, err error) {
	if len(st.channel) == 0 {
		return "", ERR_WS_SUBSCRIOTION_PARAMS
	}

	if len(st.filter) > 0 {
		return st.channel + ":" + st.filter, nil
	} else {
		return st.channel, nil
	}
}

type WSEventResponse struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

func (r *WSEventResponse) Valid() bool {
	return len(r.Event) > 0 && len(r.Channel) > 0
}

type WSTableResponse struct {
	Table  string        `json:"table"`
	Action string        `json:"action",default:""`
	Data   []interface{} `json:"data"`
}

func (r *WSTableResponse) Valid() bool {
	return (len(r.Table) > 0 || len(r.Action) > 0) && len(r.Data) > 0
}

type WSErrorResponse struct {
	Event     string `json:"event"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

func (r *WSErrorResponse) Valid() bool {
	return len(r.Event) > 0 && len(r.Message) > 0 && r.ErrorCode >= 30000
}

func loadResponse(rspMsg []byte) (interface{}, error) {

	//log.Printf("%s", rspMsg)

	evtR := WSEventResponse{}
	err := JsonBytes2Struct(rspMsg, &evtR)
	if err == nil && evtR.Valid() {
		return &evtR, nil
	}

	tr := WSTableResponse{}
	err = JsonBytes2Struct(rspMsg, &tr)
	if err == nil && tr.Valid() {
		return &tr, nil
	}

	er := WSErrorResponse{}
	err = JsonBytes2Struct(rspMsg, &er)
	if err == nil && er.Valid() {
		return &er, nil
	}

	if string(rspMsg) == "pong" {
		return string(rspMsg), nil
	}

	return nil, err

}

type ReceivedDataCallback func(interface{}) error

func defaultPrintData(obj interface{}) error {
	switch obj.(type) {
	case string:
		fmt.Println(obj)
	default:
		msg, err := Struct2JsonString(obj)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println(msg)

	}
	return nil
}
