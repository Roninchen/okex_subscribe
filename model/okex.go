package model

import (
	"bytes"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin/json"
	"github.com/okcoin-okex/okex-go-sdk-api"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	m "okex/ichat/model"
	s "okex/ichat/service"
	"strconv"
	"time"
)

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

type MarkDownMsg struct {
	MsgType string `json:"msgtype"`
	MarkDownT MarkDown `json:"markdown"`
	At At `json:"at"`
}
type MarkDown struct {
	Title string `json:"title"`
	Text string `json:"text"`
}
type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll bool `json:"isAtAll"`
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

func (req *Req)Make(ch <-chan *okex.FuturesInstrumentLiquidationResult,result okex.FuturesInstrumentLiquidationResult,max int,sizeTotal int64) *Req{
	req.Data.First.Value = result.InstrumentId
	if result.Type == "3" {
		req.Data.Keyword1.Value = "底部可开多多多多多多多"
	}else {
		req.Data.Keyword1.Value = "顶部可开空空空空空空空"
	}
	req.Data.Keyword2.Value = viper.GetString("message.version")
	req.Data.Keyword3.Value = fmt.Sprintf("%s",time.Now().Format("2006/1/2 15:04:05"))
	size, err := strconv.ParseInt(result.Size,10,64)
	if err != nil {
		seelog.Info(err)
	}
	req.Data.Remark.Value = "行情推送 "+fmt.Sprintf("价格:%v 反弹指数:+%v \n",result.Price,size+sizeTotal)
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
func (req *Req)TestDingDing() {
	req.Data.First.Value = "Coin"
	req.Data.Keyword1.Value = "BTC测试1"
	req.Data.Keyword2.Value = "LTC测试2"
	req.Data.Keyword3.Value = "ETH测试2"
	req.Data.Remark.Value = "7日，权游8惊现星巴克受到网友热议，《权力的游戏》第八季目前正在热播中，在昨晚播出的第四集中，竟然出现了星巴克咖啡的镜头，不得不说，这植入广告真的是很隐蔽了。"
}
func (req *Req) DingDing() {
	var md MarkDownMsg
	md.MsgType = "markdown"
	md.MarkDownT.Title = req.Data.Keyword1.Value
	md.MarkDownT.Text = fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		"#### "+req.Data.First.Value,
		"##### "+req.Data.Keyword1.Value,
		"##### "+req.Data.Keyword2.Value,
		"##### "+req.Data.Keyword3.Value,
		"##### "+req.Data.Remark.Value)
	mdByte, err := json.Marshal(&md)
	if err != nil {
		seelog.Error("dingding marshal err",err)
		return
	}
	dingURL := viper.GetString("ding.url")
	dingRequest, err := http.NewRequest("POST",
		dingURL,
		bytes.NewReader(mdByte))
	//dingRequest, err := http.NewRequest("POST", "https://oapi.dingtalk.com/robot/send?access_token=ab6893afba86b066067cb898d0f5df44ccc56395e02887ac12b159acbb6a74c5", bytes.NewReader(mdByte))
	if err != nil {
		seelog.Info(dingRequest)
	}
	dingRequest.Header.Set("Content-Type", "application/json")
	dingHTTPClient := http.Client{}
	dingResp, err := dingHTTPClient.Do(dingRequest)
	readAll, err := ioutil.ReadAll(dingResp.Body)
	seelog.Info(string(readAll))
	return
}

func(req *Req) WeiXin(WeixinLogin *m.LoginMap)  {
	if ! WeixinLogin.IsLogin {
		seelog.Info("微信未登陆... 退出")
		return
	}
	seelog.Info("微信已登陆... ")
	wxSendMsg := m.WxSendMsg{}
	wxSendMsg.Type = 1
	wxSendMsg.Content = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		""+req.Data.First.Value,
		""+req.Data.Keyword1.Value,
		""+req.Data.Keyword2.Value,
		""+req.Data.Keyword3.Value,
		""+req.Data.Remark.Value)
	wxSendMsg.FromUserName = WeixinLogin.FormMe
	wxSendMsg.ToUserName = WeixinLogin.SendTo
	fmt.Println("打印 To UserName:",wxSendMsg.ToUserName)
	//zbc = wxSendMsg.ToUserName
	wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
	wxSendMsg.ClientMsgId = wxSendMsg.LocalID
	msg := s.SendMsg(WeixinLogin, wxSendMsg)
	seelog.Info("weixin err",msg)
}


func LiquidationResult2String(result *okex.FuturesInstrumentLiquidationResult) string {
	s := fmt.Sprintf("%s","====================\n")
	s = s+fmt.Sprintf("币对:%v \n",result.InstrumentId)
	if result.Type == "3" {
		s = s+fmt.Sprintf("行情推送类型:%v \n","底部可开多")
	}else {
		s = s+fmt.Sprintf("行情推送类型:%v \n","顶部可开空")
	}
	s = s+fmt.Sprintf("时间:%v \n",time.Now().Format("2006/1/2 15:04:05"))
	s = s+fmt.Sprintf("价格:%v 反弹指数:+%v \n",result.Price,result.Size)
	return s
}