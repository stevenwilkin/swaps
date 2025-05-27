package bybit

type wsCommand struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

type tickerMessage struct {
	Topic string `json:"topic"`
	Data  []struct {
		Price string `json:"P"`
	} `json:"data"`
}
