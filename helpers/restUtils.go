package helpers

import (
	"encoding/json"
	"github.com/cldcvr/api-automation-framework/model"
	"github.com/cldcvr/api-automation-framework/report"

	"github.com/gavv/httpexpect"
	"os"
	"testing"
)

// RestUtility consist httpexpect.Expect and model.HttpRequest
type RestUtility struct {
	e           *httpexpect.Expect
	httpRequest *model.HttpRequest
}

// Request create a request with baseUrl, assertReporter and NewDebugPrinter
func Request(t *testing.T, httpRequest *model.HttpRequest) *RestUtility {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  httpRequest.HttpUrl.BaseUrl,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	return &RestUtility{e, httpRequest}
}

// TriggerPostAPIWithMultipart POST request with endpoint, header, query and body and send request detail to  the report
func (rest *RestUtility) TriggerPostAPIWithMultipart(name, imagepath string) *httpexpect.Response {
	inrec, _ := json.Marshal(rest.httpRequest.Body)

	reader, _ := os.Open(imagepath)
	request := rest.e.POST(rest.httpRequest.HttpUrl.EndPoint)

	request.WithMultipart().WithFile("image", name, reader).WithFormField("body", string(inrec))
	rest.setHeader(request)
	res := request.Expect()
	report.SetRestReqRes(rest.httpRequest, res)
	return res
}

// TriggerPutAPIWithMultipart PUT request with endpoint, header, query and body and send request detail to  the report
func (rest *RestUtility) TriggerPutAPIWithMultipart(name, imagepath string) *httpexpect.Response {
	inrec, _ := json.Marshal(rest.httpRequest.Body)

	reader, _ := os.Open(imagepath)
	request := rest.e.PUT(rest.httpRequest.HttpUrl.EndPoint)

	request.WithMultipart().WithFile("image", name, reader).WithFormField("body", string(inrec))
	rest.setHeader(request)
	res := request.Expect()
	report.SetRestReqRes(rest.httpRequest, res)
	return res
}

// TriggerPostAPI POST request with endpoint, header, query and body and send request detail to  the report
func (rest *RestUtility) TriggerPostAPI() *httpexpect.Response {
	request := rest.e.POST(rest.httpRequest.HttpUrl.EndPoint)
	rest.setHeader(request).setQueryParam(request).setBody(request)
	res := request.Expect()
	report.SetRestReqRes(rest.httpRequest, res)
	return res
}

// TriggerGetAPI GET request with endpoint, header and query and send request detail to  the report
func (rest *RestUtility) TriggerGetAPI() *httpexpect.Response {
	request := rest.e.GET(rest.httpRequest.HttpUrl.EndPoint)
	rest.setHeader(request).setQueryParam(request)
	res := request.Expect()
	report.SetRestReqRes(rest.httpRequest, res)
	return res
}

// TriggerPutAPI PUT request with endpoint, header,  query and body and send request detail to  the report
func (rest *RestUtility) TriggerPutAPI() *httpexpect.Response {
	request := rest.e.PUT(rest.httpRequest.HttpUrl.EndPoint)
	rest.setHeader(request).setQueryParam(request).setBody(request)
	res := request.Expect()
	report.SetRestReqRes(rest.httpRequest, res)
	return res
}

// TriggerDeleteAPI delete request with endpoint, header,  query or body and send request detail to  the report
func (rest *RestUtility) TriggerDeleteAPI() *httpexpect.Response {
	request := rest.e.DELETE(rest.httpRequest.HttpUrl.EndPoint)
	rest.setHeader(request).setQueryParam(request).setBody(request)
	res := request.Expect()
	report.SetRestReqRes(rest.httpRequest, res)
	return res
}

// Set header to the request
func (rest *RestUtility) setHeader(req *httpexpect.Request) *RestUtility {
	for key, val := range rest.httpRequest.Header {
		req.WithHeader(key, val)
	}
	return rest
}

// Set query params to the request
func (rest *RestUtility) setQueryParam(req *httpexpect.Request) *RestUtility {
	for key, val := range rest.httpRequest.HttpUrl.QueryParam {
		req.WithQuery(key, val)
	}
	return rest
}

// Set body to the request
func (rest *RestUtility) setBody(req *httpexpect.Request) *RestUtility {
	req.WithJSON(rest.httpRequest.Body)
	return rest
}
