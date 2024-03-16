package newm_helper

import (
	"fmt"
	"time"
)

type Param struct {
	Url       string
	Body      interface{}
	Method    string
	Headers   map[string]interface{}
	BodyType  string
	CreateLog bool
	RequestId string
}

func NewRequestId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
