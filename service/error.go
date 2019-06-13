package service

import "fmt"

const (
	ErrInvalidRequest     = 1001
	ErrNotFoundProductId  = 1002
	ErrServiceBusy        = 1004
	ErrActivityNotStart   = 1005
	ErrActivityAlreadyEnd = 1006
	ErrProductSaleOut     = 1007
	ErrProcessTimeout     = 1008
	ErrClientClosed       = 1009
	ErrRetry              = 1010
	ErrAlreadyBuy         = 1011
	ErrSoldout            = 1012
)

type BusinessError struct {
	Code    int
	Message string
}

func NewBusinessError(code int) (err error) {
	e := &BusinessError{Code: code, Message: ""}
	msg := ""
	switch code {
	case ErrInvalidRequest:
		msg = "invalid request"
	case ErrNotFoundProductId:
		msg = "product not found"
	case ErrServiceBusy:
		msg = "service busy"
	case ErrActivityNotStart:
		msg = "activity not start"
	case ErrActivityAlreadyEnd:
		msg = "activity already end"
	case ErrProductSaleOut:
		msg = "product sale out"
	case ErrProcessTimeout:
		msg = "process timeout"
	case ErrClientClosed:
		msg = "client closed"
	case ErrRetry:
		msg = "please try again later"
	case ErrAlreadyBuy:
		msg = "already buy"
	case ErrSoldout:
		msg = "sold out"
	}

	e.Message = msg
	err = e
	return
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("BusinessError,code[%v],message[%v]", e.Code, e.Message)
}
