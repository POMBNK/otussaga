package order

import (
	orderService "github.com/POMBNK/orderservice/internal/service/order"
	mw "github.com/oapi-codegen/nethttp-middleware"
	"net/http"
)

type Server struct {
	orderService *orderService.Service
}

func NewServer(orderService *orderService.Service) *Server {
	return &Server{
		orderService: orderService,
	}
}

func (s *Server) Register(mux *http.ServeMux, baseURL string) http.Handler {
	swagger, _ := GetSwagger()
	mws := append([]MiddlewareFunc{}, mw.OapiRequestValidator(swagger))

	return HandlerWithOptions(s, StdHTTPServerOptions{
		BaseURL:     baseURL,
		BaseRouter:  mux,
		Middlewares: mws,
	})
}
