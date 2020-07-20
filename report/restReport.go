package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cldcvr/api-automation-framework/model"
	"github.com/gavv/httpexpect"
	"strconv"
)

type RestService struct {
	httpRequest *model.HttpRequest
	response    *httpexpect.Response
}

// SetRestReqRes sets the request and response can be used after httpexpect function calls
func SetRestReqRes(httpRequest *model.HttpRequest, response *httpexpect.Response) {
	restreport := &RestService{
		httpRequest: httpRequest,
		response:    response,
	}
	HTTPRequest = restreport
}

// GetModuleName gets the module name
func (rest *RestService) GetModuleName() string {
	return rest.httpRequest.Module
}

//ReqestDetail Set request body, header, and method in a map
func (rest *RestService) ReqestDetail() map[string]string {
	request := map[string]string{}
	request[rest.httpRequest.Method] = rest.createURL()
	request["Header"] = rest.createHeader()
	request["Body"] = rest.createBody()
	return request
}

//ResponseDetail Set response body and status code in a map
func (rest *RestService) ResponseDetail() map[string]string {
	response := map[string]string{}
	response["Response Body"] = rest.setResponseBody()
	response["Code"] = rest.setResponseCode()
	return response
}

// Set the request url and return as a string
func (rest *RestService) createURL() string {
	var url string
	if HTTPRequest != nil {
		if rest.httpRequest.HttpUrl.BaseUrl != "" {
			url += rest.httpRequest.HttpUrl.BaseUrl
		}
		if rest.httpRequest.HttpUrl.EndPoint != "" {
			url += rest.httpRequest.HttpUrl.EndPoint
		}
		if rest.httpRequest.HttpUrl.QueryParam != nil {
			url += "/"
			for key, val := range rest.httpRequest.HttpUrl.QueryParam {
				url += key + "=" + val.(string)
			}
		}
	}
	return url
}

// Set the request body and return as a string
func (rest *RestService) createBody() string {
	var body string
	if HTTPRequest != nil {
		b, err := json.Marshal(rest.httpRequest.Body)
		if err != nil {
			return "No Body"
		}
		body = string(b)
	}
	return body
}

// Set the header and return as a string
func (rest *RestService) createHeader() string {
	b := new(bytes.Buffer)
	if rest.httpRequest.Header == nil {
		return "No Header"
	} else {
		for key, value := range rest.httpRequest.Header {
			fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
		}
	}
	return b.String()
}

// Create response body and return as a string or 'No Response'
func (rest *RestService) setResponseBody() string {
	if rest.response != nil {
		re1 := rest.response.JSON()
		b, _ := json.Marshal(re1.Raw())
		return string(b)
	}
	return "No Response"
}

// Create a string which consists status code otherwise return empty string
func (rest *RestService) setResponseCode() string {
	if rest.response != nil {
		return strconv.Itoa(rest.response.Raw().StatusCode)
	}
	return ""
}
