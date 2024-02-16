package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

const webhookSuccess = 204

func SendMessage(url string, message *Message) error {
	payload := new(bytes.Buffer)

	if err := json.NewEncoder(payload).Encode(message); err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(payload.Bytes())

	resp := fasthttp.AcquireResponse()

	if err := fasthttp.Do(req, resp); err != nil {
		return err
	}

	if resp.StatusCode() != webhookSuccess {
		return &NotOk{resp.StatusCode()}
	}

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return nil
}
