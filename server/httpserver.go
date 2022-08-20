package server

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-serv/constant"
	"net/http"
	"strings"
)

type HttpServer struct {
	srv    *http.Server
	reject bool
	Name   string
	Addr   string
}

func (s *HttpServer) Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务不可用"))
		return
	}
	fmt.Fprintf(w, "Welcome index! \n")
}

func (s *HttpServer) Hello(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务不可用"))
		return
	}
	fmt.Fprintf(w, "Hello %s ! \n", ps.ByName("name"))
}

func (s *HttpServer) Stop(ctx context.Context) error {
	s.reject = true
	err := s.srv.Shutdown(ctx)
	return err
}

func (s *HttpServer) Serve(ctx context.Context) error {

	router := s.newRouter()
	router.GET("/", s.Index)
	router.GET("/hello/:name", s.Hello)

	s.srv = &http.Server{Addr: s.Addr, Handler: router}

	//s.BuildinStance("httpserver","0.0.0.0:8000")
	fmt.Println("listen", s.Addr)
	err := s.srv.ListenAndServe()

	if err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) newRouter() *httprouter.Router {
	return httprouter.New()
}

func BuildHttpinStance() *HttpServer {
	httpServAddr := []string{constant.Host, constant.HttpServerPort}
	return &HttpServer{
		Name: constant.HttpServName,
		Addr: strings.Join(httpServAddr, ":"),
	}

}
