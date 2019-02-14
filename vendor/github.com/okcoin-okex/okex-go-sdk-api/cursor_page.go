package okex

/*
 OKEX uses cursor pagination for all REST requests which return arrays
*/
type CursorPage struct {
	// Request page before (newer) this pagination id.
	Before int
	// Request page after (older) this pagination id.
	After int
	// Number of results per request. Maximum 100. (default 100)
	Limit int
}
