package main

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/cihub/seelog"
	"github.com/okcoin-okex/okex-go-sdk-api"
	"github.com/spf13/viper"
	"hbfuture/config"
	"hbfuture/services"
	"io/ioutil"
	"net/http"
	"okex/conf"
	"okex/model"
	"okex/scheduler"
	"time"
)
const WXURL = "https://u.ifeige.cn/api/message/send"


var maps = make(map[string]string)


func init() {
	maps["BTC"] = ""
	maps["ETH"] = ""
	maps["BCH"] = ""
	maps["LTC"] = ""
	maps["EOS"] = ""
}
func main() {
	//init config
	conf.Init("config")
	// Load log.
	scheduler.SetLogger("logConfig.xml")
	defer seelog.Flush()

	//first send button
	FSB := viper.GetInt("message.first_send")
	seelog.Info("first send button ",FSB)

	//API
	config.ACCESS_KEY = viper.GetString("huobi.api.access_key")
	config.SECRET_KEY = viper.GetString("huobi.api.secret_key")

	// Verify every 5min.
	verifyTicker1 := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range verifyTicker1.C {
			seelog.Info("heartbeat")
		}
	}()


	verifyTickerFuture := time.NewTicker(time.Second * 10)
	if viper.GetInt("huobi.position.enable") ==1 {
		go FutureContractPositionWorker(verifyTickerFuture,"EOS")
	}

	btcChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	ethChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	bchChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	eosChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	ltcChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	max :=viper.GetInt("message.max")
	seelog.Info("max: ",max)


	if viper.GetInt("okex.enable") ==1 {
		go sendWork(ethChan, max)
		go sendWork(bchChan, max)
		go sendWork(ltcChan, max)
		go sendWork(eosChan, max)
		go sendWork(btcChan, max)
	}

	rate := viper.GetInt64("message.rate")
	verifyTicker := time.NewTicker(time.Second * time.Duration(rate) )
	seelog.Info("监控开始")

	for _ = range verifyTicker.C {
		go MarketRun(ethChan,viper.GetString("coin.eth"), "ETH", FSB)
		go MarketRun(bchChan,viper.GetString("coin.bch"), "BCH", FSB)
		go MarketRun(ltcChan,viper.GetString("coin.ltc"), "LTC", FSB)
		go MarketRun(eosChan,viper.GetString("coin.eos"), "EOS", FSB)
		go MarketRun(btcChan,viper.GetString("coin.btc"), "BTC", FSB)
		if FSB > 1000000 {
			FSB--
		} else {
			FSB++
		}
	}

}

func FutureContractPositionWorker(t *time.Ticker,coin string) {
	for _=range t.C{
		var result struct{
			Status string `json:"status"`
			Data   []model.ContractPositionInfo `json:"data"`
			TS     int64 `json:"ts"`
		}
		var Contract []model.ContractPositionInfo
		jsonStr, response, err := services.FutureContractPositionInfo(coin)
		if err!=nil {
			seelog.Info("FutureContractPositionInfo err:",err)
			continue
		}
		seelog.Info("future:",jsonStr)
		err=json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			seelog.Info("json2Future err:",err)
			continue
		}
		seelog.Info("res:",result)
		Contract = result.Data
		for k,v:=range Contract {
			seelog.Info("==========第",k,"个订单==========")
			seelog.Info(v)
		}
	}
}

func NewOKExClient() *okex.Client {
	var config okex.Config
	config.Endpoint = "https://www.okex.me/"
	config.ApiKey = viper.GetString("okex.api_key")
	config.SecretKey = viper.GetString("okex.secret_key")
	config.Passphrase = ""
	config.TimeoutSecond = 45
	config.IsPrint = false
	config.I18n = okex.ENGLISH

	req := new(model.Req)
	req.Init()

	client := okex.NewClient(config)
	return client
}

func MarketRun(ch chan<- *okex.FuturesInstrumentLiquidationResult,CoinId string,coin string,n int)  {
	// To avoid deadlock, channel must be closed.
	//defer close(ch)

	client := NewOKExClient()
	list, err := client.GetFuturesInstrumentLiquidation(CoinId, 1,1,0,1)
	if err!=nil {
		seelog.Error("爆仓订单：",err)
		return
	}
	if len(list.LiquidationList)<1 {
		seelog.Error("长度为空")
		return
	}
	//seelog.Info("create",list.LiquidationList[0].CreatedAt)
	if maps[coin] != list.LiquidationList[0].CreatedAt {
		maps[coin] = list.LiquidationList[0].CreatedAt
	}else {
		return
	}
	if n <= 2 {
		return
	}
	ch <- &list.LiquidationList[0]
	return
}

func sendWork(ch <-chan *okex.FuturesInstrumentLiquidationResult,max int){
	for {
		select {
		case  v:=<-ch :
			send(ch,v,max)
			time.Sleep(time.Duration(viper.GetInt64("message.sleep"))*time.Second)
		}
	}
}

func send(ch <-chan *okex.FuturesInstrumentLiquidationResult,result *okex.FuturesInstrumentLiquidationResult,max int)  {
	req := new(model.Req)
	req.Init()
	req.Make(ch,*result,max)
	data, err := json.Marshal(req)
	logs.Info("json:/n",string(data))
	bytes.NewReader(data)
	request, err := http.NewRequest("POST", WXURL, bytes.NewReader(data))
	if err != nil {
		seelog.Error(err)
	}
	request.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		seelog.Error(err)
		return
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		seelog.Error(err)
	}
	logs.Info("all:/n",string(all))
	if err != nil {
		seelog.Error(err)
	}

	seelog.Info("list:/n",*result)
}