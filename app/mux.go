package app

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strings"
)

func NewGrpcGatewayMux(gatewayHandler *runtime.ServeMux) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.ToLower(strings.Split(request.Header.Get("Content-Type"), ";")[0]) == "application/x-www-form-urlencoded" {
			if err := request.ParseForm(); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			jsonMap := make(map[string]interface{}, len(request.Form))
			for k, v := range request.Form {
				if len(v) > 0 {
					jsonMap[k] = v[0]
				}
			}
			jsonBody, err := jsoniter.Marshal(jsonMap)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}
			request.Body = io.NopCloser(bytes.NewReader(jsonBody))
			request.ContentLength = int64(len(jsonBody))
			request.Header.Set("Content-Type", "application/json")
		}
		gatewayHandler.ServeHTTP(writer, request)
	})
}

func NewHttpAndGrpcMux(gin *gin.Engine, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
			grpcHandler.ServeHTTP(writer, request)
			return
		}
		gin.ServeHTTP(writer, request)
	})
}
