package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/cihub/seelog"
	"github.com/okcoin-okex/okex-go-sdk-api"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"okex/conf"
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
	FSB := viper.GetInt("message.first_send")
	seelog.Info("first send button ",FSB)
	// Verify every 5min.
	verifyTicker1 := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range verifyTicker1.C {
			seelog.Info("heartbeat")
		}
	}()
	btcChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	ethChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	bchChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	eosChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	ltcChan :=make(chan *okex.FuturesInstrumentLiquidationResult,20)
	max :=viper.GetInt("message.max")
	seelog.Info("max: ",max)
	go sendWork(ethChan,max)
	go sendWork(bchChan,max)
	go sendWork(ltcChan,max)
	go sendWork(eosChan,max)
	go sendWork(btcChan,max)

	verifyTicker := time.NewTicker(time.Second * 1 )
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

func NewOKExClient() *okex.Client {
	var config okex.Config
	config.Endpoint = "https://www.okex.me/"
	config.ApiKey = viper.GetString("okex.api_key")
	config.SecretKey = viper.GetString("okex.secret_key")
	config.Passphrase = ""
	config.TimeoutSecond = 45
	config.IsPrint = false
	config.I18n = okex.ENGLISH

	req := new(Req)
	req.Init()

	client := okex.NewClient(config)
	return client
}

type Req struct {
	Secret string `json:"secret"`
	AppKey string `json:"app_key"`
	TemplateId string `json:"template_id"`
	Url      string `json:"url"`
	Data    data `json:"data"`
}

type data struct {
	First first `json:"first"`
	Keyword1 keyword1 `json:"keyword1"`
	Keyword2 keyword2 `json:"keyword2"`
	Keyword3 keyword3 `json:"keyword3"`
	Remark remark `json:"remark"`
}
type first struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
type keyword1 struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
type keyword2 struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
type keyword3 struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
type remark struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

func (req *Req)Init() *Req {
	req.Secret = viper.GetString("ifeige2.secret")
	req.AppKey = viper.GetString("ifeige2.app_key")
	req.TemplateId = viper.GetString("ifeige2.template_id")
	req.Data.First.Color = "#173177"
	req.Data.Keyword1.Color = "#173177"
	req.Data.Keyword2.Color = "#173177"
	req.Data.Keyword3.Color = "#173177"
	req.Data.Remark.Color = "#173177"
	return req
}

func (req *Req)Make(ch <-chan *okex.FuturesInstrumentLiquidationResult,result okex.FuturesInstrumentLiquidationResult,max int) *Req{
	req.Data.First.Value = result.InstrumentId
	if result.Type == 3 {
		req.Data.Keyword1.Value = "卖出平多"
	}else {
		req.Data.Keyword1.Value = "买入平空"
	}
	req.Data.Keyword2.Value = "chauncy"
	req.Data.Keyword3.Value = fmt.Sprintf("%s",time.Now().Format("2006/1/2 15:04:05"))
	req.Data.Remark.Value = "行情爆仓推送 "+fmt.Sprintf("价格:%v 数量:%v \n",result.Price,result.Size)
	i := 0
	for  {
		if i > max {
			break
		}
		if len(ch) == 0 {
			break
		}
		req.Data.Remark.Value =req.Data.Remark.Value + LiquidationResult2String(<-ch)
		i++
	}

	return req
}

func LiquidationResult2String(result *okex.FuturesInstrumentLiquidationResult) string {
	s := fmt.Sprintf("%s","=======================\n")
	s = s+fmt.Sprintf("币对:%v \n",result.InstrumentId)
	if result.Type == 3 {
		s = s+fmt.Sprintf("爆仓类型:%v \n","卖出平多")
	}else {
		s = s+fmt.Sprintf("爆仓类型:%v \n","买入平空")
	}
	s = s+fmt.Sprintf("时间:%v \n",time.Now().Format("2006/1/2 15:04:05"))
	s = s+fmt.Sprintf("价格:%v 数量:%v \n",result.Price,result.Size)
	return s
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
	if maps[coin] != list.LiquidationList[0].CreatedAt {
		maps[coin] = list.LiquidationList[0].CreatedAt
	}else {
		return
	}
	if n == 1 {
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
			time.Sleep(2*time.Second)
		}
	}
}

func send(ch <-chan *okex.FuturesInstrumentLiquidationResult,result *okex.FuturesInstrumentLiquidationResult,max int)  {
	req := new(Req)
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