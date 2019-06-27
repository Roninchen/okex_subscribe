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
	"okex/ichat"
	"okex/model"
	"okex/scheduler"
	"okex/utils"
	"strconv"
	"sync"
	"time"
)
const WXURL = "https://u.ifeige.cn/api/message/send"


var maps = make(map[string]string)
var client *okex.Client
var TIME = make(map[string]CoinTimeMap)

func init() {
	maps["BTC"] = ""
	maps["ETH"] = ""
	maps["BCH"] = ""
	maps["LTC"] = ""
	maps["EOS"] = ""

	TIME["BTC"] = CoinTimeMap{Data:map[string]time.Time{"BTC":time.Now()},Mutex:&sync.Mutex{}}
	TIME["ETH"] = CoinTimeMap{Data:map[string]time.Time{"ETH":time.Now()},Mutex:&sync.Mutex{}}
	TIME["BCH"] = CoinTimeMap{Data:map[string]time.Time{"BCH":time.Now()},Mutex:&sync.Mutex{}}
	TIME["LTC"] = CoinTimeMap{Data:map[string]time.Time{"LTC":time.Now()},Mutex:&sync.Mutex{}}
	TIME["EOS"] = CoinTimeMap{Data:map[string]time.Time{"EOS":time.Now()},Mutex:&sync.Mutex{}}
}

type CoinTimeMap struct {
	*sync.Mutex
	Data map[string]time.Time
}
func main() {
	//init config
	conf.Init("config")

	if viper.GetInt("wechat.robot") == 1 {
		go ichat.Weixin()
	}
	// Load log.
	scheduler.SetLogger("logConfig.xml")
	defer seelog.Flush()

	// init okex client
	client = NewOKExClient()


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

	//dingTicker := time.NewTicker(time.Second * 5)
	//go func() {
	//	for _ = range dingTicker.C {
	//		testWeiXin()
	//	}
	//}()

	verifyTickerFuture := time.NewTicker(time.Second * 10)
	if viper.GetInt("future.position.enable") ==1 {
		go FutureContractPositionWorker(verifyTickerFuture,viper.GetString("coin.eos"))
		go FutureContractPositionWorker(verifyTickerFuture,viper.GetString("coin.eth"))
		go FutureContractPositionWorker(verifyTickerFuture,viper.GetString("coin.ltc"))
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
	verifyTicker := time.NewTicker(time.Millisecond * time.Duration(rate) )
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
	for _=range t.C {
		if viper.GetInt("future.position.okex_huobi") == 1 {
			futuresPosition, err := client.GetFuturesInstrumentPosition(coin)
			position, err := client.GetSwapPositionByInstrument("EOS-USD-SWAP")
			if err!=nil {
				seelog.Error("okex GetFuturesInstrumentPosition err:",err)
				continue
			}
			seelog.Info("futures:",&futuresPosition.FixedPosition)
			seelog.Info("swap:",position)

		} else {
			var result struct {
				Status string                       `json:"status"`
				Data   []model.ContractPositionInfo `json:"data"`
				TS     int64                        `json:"ts"`
			}
			var Contract []model.ContractPositionInfo
			jsonStr, response, err := services.FutureContractPositionInfo(coin)
			if err != nil {
				seelog.Info("FutureContractPositionInfo err:", err)
				continue
			}
			seelog.Info("future:", jsonStr)
			err = json.NewDecoder(response.Body).Decode(&result)
			if err != nil {
				seelog.Info("json2Future err:", err)
				continue
			}
			seelog.Info("res:", result)
			Contract = result.Data
			for k, v := range Contract {
				seelog.Info("==========第", k, "个订单==========")
				seelog.Info(v)
				seelog.Info("币种：", v.ContractCode, "收益率：", v.ProfitRate)
				if v.ProfitRate < 0 {
					seelog.Info("亏损警报	>>>", "币种：", v.ContractCode, "收益率：", v.ProfitRate)
				}
			}
		}
	}
}

func NewOKExClient() *okex.Client {
	var config okex.Config
	config.Endpoint = "https://www.okex.me/"
	//config.ApiKey = viper.GetString("okex.api_key")
	//config.SecretKey = viper.GetString("okex.secret_key")
	//config.Passphrase = "okex1qaz"
	config.TimeoutSecond = 45
	config.IsPrint = false
	config.I18n = okex.ENGLISH

	req := new(model.Req)
	req.Init()

	client := okex.NewClient(config)
	return client
}

func MarketRun(ch chan<- *okex.FuturesInstrumentLiquidationResult,CoinId string,coin string,n int) {
	// To avoid deadlock, channel must be closed.
	//defer close(ch)
	list, err := client.GetFuturesInstrumentLiquidation(CoinId, 1, 1, 0, 5)
	if err != nil {
		seelog.Error("获取订单：", err)
		return
	}
	//seelog.Info(CoinId)
	//seelog.Info("create", list.LiquidationList[0].CreatedAt)
	//seelog.Info(coin,":",list.LiquidationList)
	if len(list.LiquidationList) < 1 {
		seelog.Error("长度为空")
		return
	}
	if maps[coin] != list.LiquidationList[0].CreatedAt {
		maps[coin] = list.LiquidationList[0].CreatedAt
	} else {
		return
	}
	if n <= 2 {
		//seelog.Info("coin",coin,timeMAP[coin])
		//seelog.Info("new:",utils.StrToTime(list.LiquidationList[0].CreatedAt))
		TIME[coin].Data[coin] = utils.StrToTime(list.LiquidationList[0].CreatedAt)
		return
	}
	if TIME[coin].Data[coin].Before(utils.StrToTime(list.LiquidationList[len(list.LiquidationList)-1].CreatedAt)) {
		seelog.Info("time before 5")
		TIME[coin].Lock()
		TIME[coin].Data[coin] = utils.StrToTime(list.LiquidationList[0].CreatedAt)
		TIME[coin].Unlock()
		seelog.Info("list:",list.LiquidationList)
		for _, v := range list.LiquidationList {
			ch <- &v
		}
	} else if TIME[coin].Data[coin].Before(utils.StrToTime(list.LiquidationList[0].CreatedAt)){
		seelog.Info("time before 1")
		seelog.Info(list.LiquidationList)
		for _, v := range list.LiquidationList {
			now := utils.StrToTime(v.CreatedAt)
			before := TIME[coin].Data[coin].Before(now)
			if !before {
				continue
			}
			seelog.Info("timeMap:",TIME[coin].Data[coin],"now:",now)
			ch <- &v
		}
		TIME[coin].Lock()
		TIME[coin].Data[coin] = utils.StrToTime(list.LiquidationList[0].CreatedAt)
		TIME[coin].Unlock()
	}
	return
}

func sendWork(ch <-chan *okex.FuturesInstrumentLiquidationResult,max int){
	var total int
	var sizeTotal int64
	ticker := time.NewTicker(time.Second * time.Duration(viper.GetInt("message.reset_total_time")))
	go func() {
		for range ticker.C{
			if total<=viper.GetInt("message.reset_less_then") {
				total = 0
				sizeTotal = 0
				seelog.Info("清零:",total,sizeTotal)
			}
		}
	}()
	for {
		select {
		case  v:=<-ch :
			total++
			seelog.Info("有:",v.InstrumentId,"爆单信息:爆单数量为",v.Size)
			seelog.Info("有爆单信息,累计total为:",total)
			size, _ := strconv.ParseInt(v.Size, 10, 64)
			sizeTotal+=size
			if total>= viper.GetInt("message.send_size"){
				send(ch,v,max,sizeTotal)
				time.Sleep(time.Duration(viper.GetInt64("message.sleep"))*time.Second)
				total = 0
				sizeTotal = 0
			}
		}

	}
}

func send(ch <-chan *okex.FuturesInstrumentLiquidationResult,result *okex.FuturesInstrumentLiquidationResult,max int,sizeTotal int64)  {
	req := new(model.Req)
	req.Init()
	req.Make(ch,*result,max,sizeTotal)
	data, err := json.Marshal(req)
	logs.Info("json:/n",string(data))
	go req.DingDing()
	go req.WeiXin(&ichat.LoginMap)
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

func testDingDing()  {
	req := new(model.Req)
	req.Init()
	req.TestDingDing()
	go req.DingDing()
}
func testWeiXin()  {
	req := new(model.Req)
	req.Init()
	req.TestDingDing()
	time.Sleep(3)
	req.WeiXin(&ichat.LoginMap)
}