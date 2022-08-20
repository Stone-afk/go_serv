package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go-serv/constant"
	"log"
	"net"
	"strings"
	"sync"
)

type TcpServer struct {
	Name         string
	Addr         string
	shutdownFlag chan struct{}
	wg           *sync.WaitGroup
	failCount    int32
	maxFailCount int32
}

func (s *TcpServer) Stop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			s.shutdownFlag <- struct{}{}
			return ctx.Err()
		default:
			continue
		}
	}

}

func (s *TcpServer) process(conn net.Conn) error {
	for {
		select {
		case <-s.shutdownFlag:
			conn.Close()
			return nil
		default:
			defer conn.Close()
			reader := bufio.NewReader(conn)
			var buf [constant.TcpBufferSize]byte
			n, err := reader.Read(buf[:])
			if err != nil {
				return errors.New(constant.ReadFailed)
			}
			recv := string(buf[:n])
			log.Println("收到来自后台管理系统的数据: ", recv)
			conn.Write([]byte(recv))
		}

	}

}

func (s *TcpServer) Serve(ctx context.Context) error {
	listen, err := net.Listen("tcp", s.Addr)
	fmt.Println("listen", s.Addr)
	if err != nil {
		return errors.New(constant.ListenFailed)
	}

	for {
		select {
		case <-s.shutdownFlag:
			break
		default:
			conn, err := listen.Accept()
			if err != nil {
				if s.failCount >= s.maxFailCount {
					log.Println(ctx.Err(), constant.ConnectFailed)
					return ctx.Err()
				}
				log.Println(ctx.Err(), constant.ConnectFailed)
				s.failCount += 1
				continue
			}
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				err := s.process(conn)
				if err != nil {
					log.Println(err)
				}
			}()

		}

	}
	s.wg.Wait()
	return nil
}

func BuildTcpinStance() *TcpServer {
	tcpServAddr := []string{constant.Host, constant.TcpServerPort}
	return &TcpServer{
		Name:         constant.TcpServName,
		Addr:         strings.Join(tcpServAddr, ":"),
		shutdownFlag: make(chan struct{}),
		wg:           &sync.WaitGroup{},
		maxFailCount: constant.MaxConnFailCount,
	}
}
