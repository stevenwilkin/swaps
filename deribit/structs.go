package deribit

type requestMessage struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type tradeMessage struct {
	Method string `json:"method"`
	Params struct {
		Data []struct {
			Price float64 `json:"price"`
		} `json:"data"`
	} `json:"params"`
}
