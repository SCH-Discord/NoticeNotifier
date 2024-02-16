package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type NotOk struct {
	Code int
}

func (o *NotOk) Error() string {
	return fmt.Sprintf("NotOk: StatusCode is %d", o.Code)
}

func SendMessage(url string, message *Message) error {
	payload := new(bytes.Buffer)

	if err := json.NewEncoder(payload).Encode(message); err != nil {
		return err
	}

	client := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(payload.Bytes())

	resp := fasthttp.AcquireResponse()

	if err := client.Do(req, resp); err != nil {
		return err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return &NotOk{resp.StatusCode()}
	}

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return nil
}
