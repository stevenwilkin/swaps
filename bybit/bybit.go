package bybit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Bybit struct {
	ApiKey    string
	ApiSecret string
	Testnet   bool
}

func (b *Bybit) hostname() string {
	if b.Testnet {
		return "api-testnet.bybit.com"
	} else {
		return "api.bybit.com"
	}
}

func (b *Bybit) websocketHostname() string {
	if b.Testnet {
		return "stream-testnet.bybit.com"
	} else {
		return "stream.bybit.com"
	}
}

func (b *Bybit) subscribe(channels []string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   b.websocketHostname(),
		Path:   "/realtime"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &websocket.Conn{}, err
	}

	command := wsCommand{Op: "subscribe", Args: channels}
	if err = c.WriteJSON(command); err != nil {
		return &websocket.Conn{}, err
	}

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		heartbeat := wsCommand{Op: "ping"}

		for {
			if err = c.WriteJSON(heartbeat); err != nil {
				return
			}
			<-ticker.C
		}
	}()

	return c, nil
}

func (b *Bybit) _Price() chan float64 {
	ch := make(chan float64)
	tradeTopic := "trade.BTCUSD"

	c, err := b.subscribe([]string{tradeTopic})
	if err != nil {
		close(ch)
		return ch
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				c.Close()
				close(ch)
				return
			}

			var ticker tickerMessage
			json.Unmarshal(message, &ticker)

			if ticker.Topic != tradeTopic {
				continue
			}

			if len(ticker.Data) == 0 {
				continue
			}

			ch <- ticker.Data[0].Price
		}
	}()

	return ch
}

func (b *Bybit) getNoAuth(path string, params url.Values, result interface{}) error {
	u := url.URL{
		Scheme:   "https",
		Host:     b.hostname(),
		Path:     path,
		RawQuery: params.Encode()}

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, result)

	return nil
}

func (b *Bybit) GetPrice() float64 {
	v := url.Values{
		"category": {"inverse"},
		"symbol":   {"BTCUSD"}}

	var response tickersResponse
	err := b.getNoAuth("/v5/market/tickers", v, &response)

	if err != nil {
		return 0
	}

	if len(response.Result.List) != 1 {
		return 0
	}

	price, _ := strconv.ParseFloat(response.Result.List[0].LastPrice, 64)
	return price
}

func (b *Bybit) Price() chan float64 {
	ch := make(chan float64)

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			ch <- b.GetPrice()
			<-t.C
		}
	}()

	return ch
}
