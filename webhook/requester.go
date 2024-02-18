package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

const webhookSuccess = 204

// 디스코드 웹후크를 전송
func SendMessage(url string, message *Message) error {
	newM := checkMessage(message)

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

	if newM != nil {
		go SendMessage(url, newM)
	}

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return nil
}

func checkMessage(message *Message) *Message {
	if message.Embeds == nil || len(*message.Embeds) <= 10 {
		return nil
	}

	a := (*message.Embeds)[10:]
	b := (*message.Embeds)[:10]

	newMessage := *message
	newMessage.Embeds = &a

	message.Embeds = &b

	return &newMessage
}
