package model

type ContractPositionInfo struct {
	Symbol 			string `json:"symbol" desc:"品种代码 BTC ETH"`
	ContractCode 	string `json:"contract_code" desc:"合约代码"`
	ContractType 	string `json:"contract_type" desc:"合约类型 当周 次周 季度"`
	Volume  		float64 `json:"volume" desc:"持仓量"`
	Available 		float64 `json:"available" desc:"可平仓数量"`
	Frozen 			float64 `json:"frozen" desc:"冻结数量	"`
	CostOpen 		float64 `json:"cost_open" desc:"开仓均价"`
	CostHold 		float64 `json:"cost_hold" desc:"持仓均价"`
	ProfitUnreal 	float64 `json:"profit_unreal" desc:"未实现盈亏"`
	ProfitRate 		float64 `json:"profit_rate" desc:"收益率	"`
	Profit 			float64 `json:"profit" desc:"收益"`
	PositionMargin 	float64 `json:"position_margin" desc:"持仓保证金"`
	LeverRate 		int `json:"lever_rate" desc:"杠杠倍数	"`
	Direction 		string `json:"direction" desc:"buy:买入开多 sell:卖出开空"`
}

type HBResponse struct {
	Status string `json:"status"`
	Date   []ContractPositionInfo `json:"date"`
	TS     int64  `json:"ts"`
}
/**
"symbol": "BTC",
"contract_code": "BTC180914",
"contract_type": "this_week",
"volume": 1,
"available": 0,
"frozen": 0.3,
"cost_open": 422.78,
"cost_hold": 422.78,
"profit_unreal": 0.00007096,
"profit_rate": 0.07,
"profit": 0.97,
"position_margin": 3.4,
"lever_rate": 10,
"direction":"buy"
**/

