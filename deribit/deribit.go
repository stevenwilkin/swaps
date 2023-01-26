package deribit

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Deribit struct {
	ApiId     string
	ApiSecret string
	Test      bool
}

func (d *Deribit) hostname() string {
	if d.Test {
		return "test.deribit.com"
	} else {
		return "www.deribit.com"
	}
}

func (d *Deribit) subscribe(channels []string) (*websocket.Conn, error) {
	socketUrl := url.URL{Scheme: "wss", Host: d.hostname(), Path: "/ws/api/v2"}

	c, _, err := websocket.DefaultDialer.Dial(socketUrl.String(), nil)
	if err != nil {
		return &websocket.Conn{}, err
	}

	request := requestMessage{
		Method: "/public/subscribe",
		Params: map[string]interface{}{
			"channels": channels}}

	if err = c.WriteJSON(request); err != nil {
		return &websocket.Conn{}, err
	}

	ticker := time.NewTicker(10 * time.Second)
	testMessage := requestMessage{Method: "/public/test"}

	go func() {
		for {
			if err = c.WriteJSON(testMessage); err != nil {
				return
			}
			<-ticker.C
		}
	}()

	return c, nil
}

func (d *Deribit) Price() chan float64 {
	ch := make(chan float64)

	c, err := d.subscribe([]string{"trades.BTC-PERPETUAL.100ms"})
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

			var trade tradeMessage
			json.Unmarshal(message, &trade)

			if len(trade.Params.Data) == 0 {
				continue
			}

			ch <- trade.Params.Data[0].Price
		}
	}()

	return ch
}
