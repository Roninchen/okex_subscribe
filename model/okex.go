package model

import (
	"fmt"
	"github.com/okcoin-okex/okex-go-sdk-api"
	"github.com/spf13/viper"
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
	req.Data.Keyword2.Value = viper.GetString("message.version")
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