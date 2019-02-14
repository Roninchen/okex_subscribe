package okex

import (
	"net/http"
	"strings"
)

/*
 OKEX futures contract api
 @author Tony Tian
 @date 2018-03-17
 @version 1.0.0
*/

/*
 =============================== Futures market api ===============================
*/
/*
 The exchange rate of legal tender pairs
*/
func (client *Client) GetFuturesExchangeRate() (ExchangeRate, error) {
	var exchangeRate ExchangeRate
	_, err := client.Request(GET, FUTURES_RATE, nil, &exchangeRate)
	return exchangeRate, err
}

/*
  Get all of futures contract list
*/
func (client *Client) GetFuturesInstruments() ([]FuturesInstrumentsResult, error) {
	var Instruments []FuturesInstrumentsResult
	_, err := client.Request(GET, FUTURES_INSTRUMENTS, nil, &Instruments)
	return Instruments, err
}

/*
 Get the futures contract currencies
*/
func (client *Client) GetFuturesInstrumentCurrencies() ([]FuturesInstrumentCurrenciesResult, error) {
	var currencies []FuturesInstrumentCurrenciesResult
	_, err := client.Request(GET, FUTURES_CURRENCIES, nil, &currencies)
	return currencies, err
}

/*
 Get the futures contract Instrument book
 depth value：1-200
 merge value：1(merge depth)
*/
func (client *Client) GetFuturesInstrumentBook(InstrumentId string, size int) (FuturesInstrumentBookResult, error) {
	var book FuturesInstrumentBookResult
	params := NewParams()
	params["size"] = Int2String(size)
	requestPath := BuildParams(GetInstrumentIdUri(FUTURES_INSTRUMENT_BOOK, InstrumentId), params)
	_, err := client.Request(GET, requestPath, nil, &book)
	return book, err
}

/*
 Get the futures contract Instrument all ticker
*/
func (client *Client) GetFuturesInstrumentAllTicker() ([]FuturesInstrumentTickerResult, error) {
	var tickers []FuturesInstrumentTickerResult
	_, err := client.Request(GET, FUTURES_TICKERS, nil, &tickers)
	return tickers, err
}

/*
 Get the futures contract Instrument ticker
*/
func (client *Client) GetFuturesInstrumentTicker(InstrumentId string) (FuturesInstrumentTickerResult, error) {
	var ticker FuturesInstrumentTickerResult
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_TICKER, InstrumentId), nil, &ticker)
	return ticker, err
}

/*
 Get the futures contract Instrument trades
*/
func (client *Client) GetFuturesInstrumentTrades(InstrumentId string) ([]FuturesInstrumentTradesResult, error) {
	var trades []FuturesInstrumentTradesResult
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_TRADES, InstrumentId), nil, &trades)
	return trades, err
}

/*
 Get the futures contract Instrument candles
 granularity: @see  file: futures_constants.go
*/
func (client *Client) GetFuturesInstrumentCandles(InstrumentId, start, end string, granularity, count int) ([][]float64, error) {
	var candles [][]float64
	params := NewParams()
	params["start"] = start
	params["end"] = end
	params["granularity"] = Int2String(granularity)
	params["count"] = Int2String(count)
	requestPath := BuildParams(GetInstrumentIdUri(FUTURES_INSTRUMENT_CANDLES, InstrumentId), params)
	_, err := client.Request(GET, requestPath, nil, &candles)
	return candles, err
}

/*
 Get the futures contract Instrument index
*/
func (client *Client) GetFuturesInstrumentIndex(InstrumentId string) (FuturesInstrumentIndexResult, error) {
	var index FuturesInstrumentIndexResult
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_INDEX, InstrumentId), nil, &index)
	return index, err
}

/*
 Get the futures contract Instrument estimated price
*/
func (client *Client) GetFuturesInstrumentEstimatedPrice(InstrumentId string) (FuturesInstrumentEstimatedPriceResult, error) {
	var estimatedPrice FuturesInstrumentEstimatedPriceResult
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_ESTIMATED_PRICE, InstrumentId), nil, &estimatedPrice)
	return estimatedPrice, err
}

/*
 Get the futures contract Instrument holds
*/
func (client *Client) GetFuturesInstrumentOpenInterest(InstrumentId string) (FuturesInstrumentOpenInterestResult, error) {
	var openInterest FuturesInstrumentOpenInterestResult
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_OPEN_INTEREST, InstrumentId), nil, &openInterest)
	return openInterest, err
}

/*
 Get the futures contract Instrument limit price
*/
func (client *Client) GetFuturesInstrumentPriceLimit(InstrumentId string) (FuturesInstrumentPriceLimitResult, error) {
	var priceLimit FuturesInstrumentPriceLimitResult
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_PRICE_LIMIT, InstrumentId), nil, &priceLimit)
	return priceLimit, err
}

/*
 Get the futures contract liquidation
*/
func (client *Client) GetFuturesInstrumentLiquidation(InstrumentId string, status, from, to, limit int) (FuturesInstrumentLiquidationListResult, error) {
	var liquidation []FuturesInstrumentLiquidationResult
	params := NewParams()
	params["status"] = Int2String(status)
	params["from"] = Int2String(from)
	params["to"] = Int2String(to)
	params["limit"] = Int2String(limit)
	requestPath := BuildParams(GetInstrumentIdUri(FUTURES_INSTRUMENT_LIQUIDATION, InstrumentId), params)
	response, err := client.Request(GET, requestPath, nil, &liquidation)
	var list FuturesInstrumentLiquidationListResult
	if err != nil {
		return list,err
	}
	page := parsePage(response)
	list.Page = page
	list.LiquidationList = liquidation
	return list, err
}

/*
 =============================== Futures trade api ===============================
*/

/*
 Get all of futures contract position list.
 return struct: FuturesPositions
*/
func (client *Client) GetFuturesPositions() (FuturesPosition, error) {
	response, err := client.Request(GET, FUTURES_POSITION, nil, nil)
	return parsePositions(response, err)
}

/*
 Get all of futures contract position list.
 return struct: FuturesPositions
*/
func (client *Client) GetFuturesInstrumentPosition(InstrumentId string) (FuturesPosition, error) {
	response, err := client.Request(GET, GetInstrumentIdUri(FUTURES_INSTRUMENT_POSITION, InstrumentId), nil, nil)
	return parsePositions(response, err)
}

/*
 Get all of futures contract account list
 return struct: FuturesAccounts
*/
func (client *Client) GetFuturesAccounts() (FuturesAccount, error) {
	response, err := client.Request(GET, FUTURES_ACCOUNTS, nil, nil)
	return parseAccounts(response, err)
}

/*
 Get the futures contract currency account @see file : futures_constants.go
 return struct: FuturesCurrencyAccounts
*/
func (client *Client) GetFuturesAccountsByCurrency(currency string) (FuturesCurrencyAccount, error) {
	response, err := client.Request(GET, GetCurrencyUri(FUTURES_ACCOUNT_CURRENCY_INFO, currency), nil, nil)
	return parseCurrencyAccounts(response, err)
}

/*
 Get the futures contract currency ledger
*/
func (client *Client) GetFuturesAccountsLedgerByCurrency(currency string, from, to, limit int) ([]FuturesCurrencyLedger, error) {
	var ledger [] FuturesCurrencyLedger
	params := NewParams()
	params["from"] = Int2String(from)
	params["to"] = Int2String(to)
	params["limit"] = Int2String(limit)
	requestPath := BuildParams(GetCurrencyUri(FUTURES_ACCOUNT_CURRENCY_LEDGER, currency), params)
	_, err := client.Request(GET, requestPath, nil, &ledger)
	return ledger, err
}

/*
 Get the futures contract Instrument holds
*/
func (client *Client) GetFuturesAccountsHoldsByInstrumentId(InstrumentId string) (FuturesAccountsHolds, error) {
	var holds FuturesAccountsHolds
	_, err := client.Request(GET, GetInstrumentIdUri(FUTURES_ACCOUNT_INSTRUMENT_HOLDS, InstrumentId), nil, &holds)
	return holds, err
}

/*
 Create a new order
*/
func (client *Client) FuturesOrder(newOrderParams FuturesNewOrderParams) (FuturesNewOrderResult, error) {
	var newOrderResult FuturesNewOrderResult
	_, err := client.Request(POST, FUTURES_ORDER, newOrderParams, &newOrderResult)
	return newOrderResult, err
}

/*
 Batch create new order.(Max of 5 orders are allowed per request)
*/
func (client *Client) FuturesOrders(batchNewOrder FuturesBatchNewOrderParams) (FuturesBatchNewOrderResult, error) {
	var batchNewOrderResult FuturesBatchNewOrderResult
	_, err := client.Request(POST, FUTURES_ORDERS, batchNewOrder, &batchNewOrderResult)
	return batchNewOrderResult, err
}

/*
 Get all of futures contract order list
*/
func (client *Client) GetFuturesOrders(InstrumentId string, status, from, to, limit int) (FuturesGetOrdersResult, error) {
	var ordersResult FuturesGetOrdersResult
	params := NewParams()
	params["status"] = Int2String(status)
	params["from"] = Int2String(from)
	params["to"] = Int2String(to)
	params["limit"] = Int2String(limit)
	requestPath := BuildParams(GetInstrumentIdUri(FUTURES_INSTRUMENT_ORDER_LIST, InstrumentId), params)
	_, err := client.Request(GET, requestPath, nil, &ordersResult)
	return ordersResult, err
}

/*
 Get all of futures contract a order by order id
*/
func (client *Client) GetFuturesOrder(InstrumentId string, orderId int64) (FuturesGetOrderResult, error) {
	var getOrderResult FuturesGetOrderResult
	_, err := client.Request(GET, GetInstrumentIdOrdersUri(FUTURES_INSTRUMENT_ORDER_INFO, InstrumentId, orderId), nil, &getOrderResult)
	return getOrderResult, err
}

/*
 Batch Cancel the orders
*/
func (client *Client) BatchCancelFuturesInstrumentOrders(InstrumentId, orderIds string) (FuturesBatchCancelInstrumentOrdersResult, error) {
	var cancelInstrumentOrdersResult FuturesBatchCancelInstrumentOrdersResult
	params := NewParams()
	params["order_ids"] = orderIds
	_, err := client.Request(POST, GetInstrumentIdUri(FUTURES_INSTRUMENT_ORDER_BATCH_CANCEL, InstrumentId), params, &cancelInstrumentOrdersResult)
	return cancelInstrumentOrdersResult, err
}

/*
 Cancel the order
*/
func (client *Client) CancelFuturesInstrumentOrder(InstrumentId string, orderId int64) (FuturesCancelInstrumentOrderResult, error) {
	var cancelInstrumentOrderResult FuturesCancelInstrumentOrderResult
	_, err := client.Request(POST, GetInstrumentIdOrdersUri(FUTURES_INSTRUMENT_ORDER_CANCEL, InstrumentId, orderId), nil,
		&cancelInstrumentOrderResult)
	return cancelInstrumentOrderResult, err
}

/*
 Get all of futures contract transactions.
*/
func (client *Client) GetFuturesFills(InstrumentId string, orderId int64, from, to, limit int) ([]FuturesFillResult, error) {
	var fillsResult []FuturesFillResult
	params := NewParams()
	params["order_id"] = Int64ToString(orderId)
	params["instrument_id"] = InstrumentId
	params["from"] = Int2String(from)
	params["to"] = Int2String(to)
	params["limit"] = Int2String(limit)
	requestPath := BuildParams(FUTURES_FILLS, params)
	_, err := client.Request(GET, requestPath, nil, &fillsResult)
	return fillsResult, err
}

func parsePositions(response *http.Response, err error) (FuturesPosition, error) {
	var position FuturesPosition
	if err != nil {
		return position, err
	}
	var result Result
	result.Result = false
	jsonString := GetResponseDataJsonString(response)
	if strings.Contains(jsonString, "\"margin_mode\":\"fixed\"") {
		var fixedPosition FuturesFixedPosition
		err = JsonString2Struct(jsonString, &fixedPosition)
		if err != nil {
			return position, err
		} else {
			position.Result = fixedPosition.Result
			position.MarginMode = fixedPosition.MarginMode
			position.FixedPosition = fixedPosition.FixedPosition
		}
	} else if strings.Contains(jsonString, "\"margin_mode\":\"crossed\"") {
		var crossPosition FuturesCrossPosition
		err = JsonString2Struct(jsonString, &crossPosition)
		if err != nil {
			return position, err
		} else {
			position.Result = crossPosition.Result
			position.MarginMode = crossPosition.MarginMode
			position.CrossPosition = crossPosition.CrossPosition
		}
	} else if strings.Contains(jsonString, "\"code\":") {
		JsonString2Struct(jsonString, &position)
		position.Result = result
	} else {
		position.Result = result
	}

	return position, nil
}

func parseAccounts(response *http.Response, err error) (FuturesAccount, error) {
	var account FuturesAccount
	if err != nil {
		return account, err
	}
	var result Result
	result.Result = false
	jsonString := GetResponseDataJsonString(response)
	if strings.Contains(jsonString, "\"contracts\"") {
		var fixedAccount FuturesFixedAccountInfo
		err = JsonString2Struct(jsonString, &fixedAccount)
		if err != nil {
			return account, err
		} else {
			account.Result = fixedAccount.Result
			account.FixedAccount = fixedAccount.Info
			account.MarginMode = "fixed"
		}
	} else if strings.Contains(jsonString, "\"realized_pnl\"") {
		var crossAccount FuturesCrossAccountInfo
		err = JsonString2Struct(jsonString, &crossAccount)
		if err != nil {
			return account, err
		} else {
			account.Result = crossAccount.Result
			account.MarginMode = "crossed"
			account.CrossAccount = crossAccount.Info
		}
	} else if strings.Contains(jsonString, "\"code\":") {
		JsonString2Struct(jsonString, &account)
		account.Result = result
	} else {
		account.Result = result
	}
	return account, nil
}

func parseCurrencyAccounts(response *http.Response, err error) (FuturesCurrencyAccount, error) {
	var currencyAccount FuturesCurrencyAccount
	if err != nil {
		return currencyAccount, err
	}
	jsonString := GetResponseDataJsonString(response)
	var result Result
	result.Result = true
	if strings.Contains(jsonString, "\"margin_mode\":\"fixed\"") {
		var fixedAccount FuturesFixedAccount
		err = JsonString2Struct(jsonString, &fixedAccount)
		if err != nil {
			return currencyAccount, err
		} else {
			currencyAccount.Result = result
			currencyAccount.MarginMode = fixedAccount.MarginMode
			currencyAccount.FixedAccount = fixedAccount
		}
	} else if strings.Contains(jsonString, "\"margin_mode\":\"crossed\"") {
		var crossAccount FuturesCrossAccount
		err = JsonString2Struct(jsonString, &crossAccount)
		if err != nil {
			return currencyAccount, err
		} else {
			currencyAccount.Result = result
			currencyAccount.MarginMode = crossAccount.MarginMode
			currencyAccount.CrossAccount = crossAccount
		}
	} else if strings.Contains(jsonString, "\"code\":") {
		result.Result = true
		JsonString2Struct(jsonString, &currencyAccount)
		currencyAccount.Result = result
	} else {
		result.Result = true
		currencyAccount.Result = result
	}
	return currencyAccount, nil
}

func parsePage(response *http.Response) PageResult {
	var page PageResult
	jsonString := GetResponsePageJsonString(response)
	JsonString2Struct(jsonString, &page)
	return page
}
