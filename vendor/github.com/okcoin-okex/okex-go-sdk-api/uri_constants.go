package okex

/*
 OKEX api config info
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

const (
	OKEX_TIME_URI                         = "/api/general/v3/time"
	FUTURES_RATE                          = "/api/futures/v3/rate"
	FUTURES_INSTRUMENTS                   = "/api/futures/v3/instruments"
	FUTURES_CURRENCIES                    = "/api/futures/v3/instruments/currencies"
	FUTURES_INSTRUMENT_BOOK               = "/api/futures/v3/instruments/{instrument_id}/book"
	FUTURES_TICKERS                       = "/api/futures/v3/instruments/ticker"
	FUTURES_INSTRUMENT_TICKER             = "/api/futures/v3/instruments/{instrument_id}/ticker"
	FUTURES_INSTRUMENT_TRADES             = "/api/futures/v3/instruments/{instrument_id}/trades"
	FUTURES_INSTRUMENT_CANDLES            = "/api/futures/v3/instruments/{instrument_id}/candles"
	FUTURES_INSTRUMENT_INDEX              = "/api/futures/v3/instruments/{instrument_id}/index"
	FUTURES_INSTRUMENT_ESTIMATED_PRICE    = "/api/futures/v3/instruments/{instrument_id}/estimated_price"
	FUTURES_INSTRUMENT_OPEN_INTEREST      = "/api/futures/v3/instruments/{instrument_id}/open_interest"
	FUTURES_INSTRUMENT_PRICE_LIMIT        = "/api/futures/v3/instruments/{instrument_id}/price_limit"
	FUTURES_INSTRUMENT_LIQUIDATION        = "/api/futures/v3/instruments/{instrument_id}/liquidation"
	FUTURES_POSITION                      = "/api/futures/v3/position"
	FUTURES_INSTRUMENT_POSITION           = "/api/futures/v3/{instrument_id}/position"
	FUTURES_ACCOUNTS                      = "/api/futures/v3/accounts"
	FUTURES_ACCOUNT_CURRENCY_INFO         = "/api/futures/v3/accounts/{currency}"
	FUTURES_ACCOUNT_CURRENCY_LEDGER       = "/api/futures/v3/accounts/{currency}/ledger"
	FUTURES_ACCOUNT_INSTRUMENT_HOLDS      = "/api/futures/v3/accounts/{instrument_id}/holds"
	FUTURES_ORDER                         = "/api/futures/v3/order"
	FUTURES_ORDERS                        = "/api/futures/v3/orders"
	FUTURES_INSTRUMENT_ORDER_LIST         = "/api/futures/v3/orders/{instrument_id}"
	FUTURES_INSTRUMENT_ORDER_INFO         = "/api/futures/v3/orders/{instrument_id}/{order_id}"
	FUTURES_INSTRUMENT_ORDER_CANCEL       = "/api/futures/v3/cancel_order/{instrument_id}/{order_id}"
	FUTURES_INSTRUMENT_ORDER_BATCH_CANCEL = "/api/futures/v3/cancel_batch_orders/{instrument_id}"
	FUTURES_FILLS                         = "/api/futures/v3/fills"

	SWAP_INSTRUMENT_ACCOUNT                 = "/api/swap/v3/{instrument_id}/accounts"
	SWAP_INSTRUMENT_POSITION                = "/api/swap/v3/{instrument_id}/position"
	SWAP_ACCOUNTS                           = "/api/swap/v3/accounts"
	SWAP_ACCOUNTS_HOLDS                     = "/api/swap/v3/accounts/{instrument_id}/holds"
	SWAP_ACCOUNTS_LEDGER                    = "/api/swap/v3/accounts/{instrument_id}/ledger"
	SWAP_ACCOUNTS_LEVERAGE                  = "/api/swap/v3/accounts/{instrument_id}/leverage"
	SWAP_ACCOUNTS_SETTINGS                  = "/api/swap/v3/accounts/{instrument_id}/settings"
	SWAP_FILLS                              = "/api/swap/v3/fills"
	SWAP_INSTRUMENTS                        = "/api/swap/v3/instruments"
	SWAP_INSTRUMENTS_TICKER                 = "/api/swap/v3/instruments/ticker"
	SWAP_INSTRUMENT_CANDLES                 = "/api/swap/v3/instruments/{instrument_id}/candles"
	SWAP_INSTRUMENT_DEPTH                   = "/api/swap/v3/instruments/{instrument_id}/depth"
	SWAP_INSTRUMENT_FUNDING_TIME            = "/api/swap/v3/instruments/{instrument_id}/funding_time"
	SWAP_INSTRUMENT_HISTORICAL_FUNDING_RATE = "/api/swap/v3/instruments/{instrument_id}/historical_funding_rate"
	SWAP_INSTRUMENT_INDEX                   = "/api/swap/v3/instruments/{instrument_id}/index"
	SWAP_INSTRUMENT_LIQUIDATION             = "/api/swap/v3/instruments/{instrument_id}/liquidation"
	SWAP_INSTRUMENT_MARK_PRICE              = "/api/swap/v3/instruments/{instrument_id}/mark_price"
	SWAP_INSTRUMENT_OPEN_INTEREST           = "/api/swap/v3/instruments/{instrument_id}/open_interest"
	SWAP_INSTRUMENT_PRICE_LIMIT             = "/api/swap/v3/instruments/{instrument_id}/price_limit"
	SWAP_INSTRUMENT_TICKER                  = "/api/swap/v3/instruments/{instrument_id}/ticker"
	SWAP_INSTRUMENT_TRADES                  = "/api/swap/v3/instruments/{instrument_id}/trades"
	SWAP_INSTRUMENT_ORDER_LIST              = "/api/swap/v3/orders/{instrument_id}"
	SWAP_INSTRUMENT_ORDER_INFO              = "/api/swap/v3/orders/{instrument_id}/{order_id}"
	SWAP_RATE                               = "/api/swap/v3/rate"
	SWAP_ORDER                              = "/api/swap/v3/order"
	SWAP_ORDERS                             = "/api/swap/v3/orders"

	SWAP_CANCEL_BATCH_ORDERS = "/api/swap/v3/cancel_batch_orders/{instrument_id}"
	SWAP_CANCEL_ORDER        = "/api/swap/v3/cancel_order/{instrument_id}/{order_id}"
)
