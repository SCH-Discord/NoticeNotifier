package webhook

import "fmt"

// 웹후크 전송을 성공한 코드가 아닐 경우 발생

type NotOk struct {
	Code    int
	Message string
}

func (o *NotOk) Error() string {
	return fmt.Sprintf("NotOk: StatusCode is %d Message %s", o.Code, o.Message)
}
