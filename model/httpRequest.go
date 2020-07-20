package model

type HttpUrl struct {
	BaseUrl    string
	EndPoint   string
	QueryParam map[string]interface{}
}

type HttpRequest struct {
	HttpUrl HttpUrl
	Method  string
	Header  map[string]string
	Body    map[string]interface{}
	Module  string
}
