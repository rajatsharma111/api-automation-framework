package report

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type GrpcService struct {
	serviceName string
	err         error
	meta        []string
	request     interface{}
	response    interface{}
}

// SetGrpcReqRes sets the request and response can be used after httpexpect function calls
func SetGrpcReqRes(serviceName string, err error, request interface{}, response interface{}, meta ...string) {
	grpcReport := &GrpcService{
		serviceName: serviceName,
		err:         err,
		meta:        meta,
		request:     request,
		response:    response,
	}
	HTTPRequest = grpcReport
}

// GetModuleName gets the module name
func (grpc *GrpcService) GetModuleName() string {
	return "GRPC"
}

//ReqestDetail Set request body, header, and method in a map
func (grpc *GrpcService) ReqestDetail() map[string]string {
	request := map[string]string{}
	request["ServiceName"] = grpc.serviceName
	out, _ := json.Marshal(grpc.request)
	request["Request"] = string(out)
	request["Meta"] = grpc.createHeader()
	return request
}

func (grpc *GrpcService) createHeader() string {
	b := new(bytes.Buffer)
	if grpc.meta == nil {
		fmt.Fprintf(b, "\"%s\"\n", "No metadata provided")
	} else {
		for i := 0; i < len(grpc.meta); {
			fmt.Fprintf(b, "%s=\"%s\"\n", grpc.meta[i], grpc.meta[i+1])
			i += 2
		}
	}
	return b.String()
}

//ResponseDetail Set response body and status code in a map
func (grpc *GrpcService) ResponseDetail() map[string]string {
	response := map[string]string{}
	out, _ := json.Marshal(grpc.response)
	response["Response Body"] = string(out)
	return response
}
