package server

import (
	"context"
	"go-serv/constant"
)

type Service interface {
	Serve(c context.Context) error
	Stop(c context.Context) error
}

//type serveMux struct {
//	reject bool
//	*http.ServeMux
//}
//
//func (m *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if m.reject {
//		log.Println("服务不可用")
//		w.WriteHeader(http.StatusServiceUnavailable)
//		_, _ = w.Write([]byte(http.StatusText(http.StatusServiceUnavailable)))
//		return
//	}
//	m.ServeMux.ServeHTTP(w, r)
//}
//
//func (s *Server) newRouter() *serveMux {
//	return &serveMux{ServeMux: http.NewServeMux()}
//}
//
//type Server struct {
//	srv    *http.Server
//	wg     *sync.WaitGroup
//	Name   string
//	Addr   string
//	reject bool
//	mux    *serveMux
//}
//
//func (s *Server) Stop(ctx context.Context) error {
//	s.reject = true
//	err := s.srv.Shutdown(ctx)
//	return err
//}
//
//func (s *Server) Serve(ctx context.Context) error {
//
//	router := s.newRouter()
//
//	s.srv = &http.Server{Addr: s.Addr, Handler: router}
//
//	//s.BuildinStance("httpserver","0.0.0.0:8000")
//	fmt.Printf("serve %s, listen %s", s.Name, s.Addr)
//	err := s.srv.ListenAndServe()
//
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func BuildinStance(host, port, servName string) *Server {
//	ServAddr := []string{host, port}
//	return &Server{
//		Name: servName,
//		Addr: strings.Join(ServAddr, ":"),
//	}
//
//}

func BuidServers() []Service {
	servers := make([]Service, 0, constant.ServiceSliceCap)
	servers = append(servers, BuildTcpinStance())
	servers = append(servers, BuildHttpinStance())
	return servers
}

//func BuidServers() []*Server {
//	servers := make([]*Server, 0, constant.ServiceSliceCap)
//
//	servers = append(servers, BuildinStance(constant.Host, constant.AdminServerPort, constant.AdminServName))
//	servers = append(servers, BuildinStance(constant.Host, constant.AppServerPort, constant.AppServName))
//	return servers
//}
