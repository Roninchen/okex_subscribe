package okex

/*
 OKEX websocket api constants
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import "errors"

const (
	WS_API_HOST = "okexcomreal.bafang.com:10442"
	WS_API_URL  = "wss://real.okex.com:10442/ws/v3"

	CHNL_SWAP_TICKER        = "swap/ticker"        //行情数据频道
	CHNL_SWAP_CANDLE60S     = "swap/candle60s"     //1分钟k线数据频道
	CHNL_SWAP_CANDLE180S    = "swap/candle180s"    //3分钟k线数据频道
	CHNL_SWAP_CANDLE300S    = "swap/candle300s"    //5分钟k线数据频道
	CHNL_SWAP_CANDLE900S    = "swap/candle900s"    //15分钟k线数据频道
	CHNL_SWAP_CANDLE1800S   = "swap/candle1800s"   //30分钟k线数据频道
	CHNL_SWAP_CANDLE3600S   = "swap/candle3600s"   //1小时k线数据频道
	CHNL_SWAP_CANDLE7200S   = "swap/candle7200s"   //2小时k线数据频道
	CHNL_SWAP_CANDLE14400S  = "swap/candle14400s"  //4小时k线数据频道
	CHNL_SWAP_CANDLE21600   = "swap/candle21600"   //6小时k线数据频道
	CHNL_SWAP_CANDLE43200S  = "swap/candle43200s"  //12小时k线数据频道
	CHNL_SWAP_CANDLE86400S  = "swap/candle86400s"  //1day
	CHNL_SWAP_CANDLE604800S = "swap/candle604800s" //1week
	CHNL_SWAP_TRADE         = "swap/trade"         //交易信息频道
	CHNL_SWAP_FUNDING_RATE  = "swap/funding_rate"  //资金费率频道
	CHNL_SWAP_PRICE_RANGE   = "swap/price_range"   //限价范围频道
	CHNL_SWAP_DEPTH         = "swap/depth"         //深度数据频道，首次200档，后续增量
	CHNL_SWAP_DEPTH5        = "swap/depth5"        //深度数据频道，每次返回前5档
	CHNL_SWAP_MARK_PRICE    = "swap/mark_price"    //标记价格频道

	CHNL_SWAP_ACCOUNT  = "swap/account"  //用户账户信息频道
	CHNL_SWAP_POSITION = "swap/position" //用户持仓信息频道
	CHNL_SWAP_ORDER    = "swap/order"    //用户交易数据频道

	CHNL_EVENT_SUBSCRIBE   = "subscribe"
	CHNL_EVENT_UNSUBSCRIBE = "unsubscribe"
)

var (
	ERR_WS_SUBSCRIOTION_PARAMS = errors.New(`ws subscription parameter error`)
)

var (
	DefaultDataCallBack = defaultPrintData
)
