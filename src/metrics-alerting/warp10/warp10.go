package warp10

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Warp10Client struct {
	ExecEndpoint string
	ReadToken    string
}

func (w *Warp10Client) ReadBool(script string) (b bool, err error) {
	resp, err := w.sendRequest(script)
	if err != nil {
		return
	}

	var respBody []bool
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return
	}

	b = respBody[0]
	return
}

func (w *Warp10Client) ReadNumber(script string) (f float64, err error) {
	resp, err := w.sendRequest(script)
	if err != nil {
		return
	}

	var respBody []float64
	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return
	}

	f = respBody[0]
	return
}

func (w *Warp10Client) ReadSeriesOfNumbers(script string) (f []FloatTimeSerie, err error) {
	resp, err := w.sendRequest(script)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&f)
	return
}

func (w *Warp10Client) appendToken(script string) string {
	return fmt.Sprintf("'%s' 'token' STORE\n%s", w.ReadToken, script)
}

func (w *Warp10Client) sendRequest(script string) (*http.Response, error) {
	script = w.appendToken(script)

	client := http.Client{}

	req, err := http.NewRequest("POST", w.ExecEndpoint, strings.NewReader(script))
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}
