package webhook

import "fmt"

type NotOk struct {
	Code int
}

func (o *NotOk) Error() string {
	return fmt.Sprintf("NotOk: StatusCode is %d", o.Code)
}
