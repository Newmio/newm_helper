package newm_helper

type Param struct {
	Url       string
	Body      interface{}
	Method    string
	Headers   map[string]interface{}
	BodyType  string
	CreateLog bool
}