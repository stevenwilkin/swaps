package bybit

import (
	"encoding/json"
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
		Path:   "/v5/public/inverse"}

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

func (b *Bybit) Price() chan float64 {
	ch := make(chan float64)
	tradeTopic := "publicTrade.BTCUSD"

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

			price, _ := strconv.ParseFloat(ticker.Data[0].Price, 64)
			ch <- price
		}
	}()

	return ch
}
