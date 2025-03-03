// Copyright 2022 Leon Ding <ding@ibyte.me> https://wiredb.github.io

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/auula/wiredb/clog"
	"github.com/auula/wiredb/vfs"
)

var (
	// ipv4 return local IPv4 address
	ipv4 string = "127.0.0.1"
)

const (
	minPort = 1024
	maxPort = 1 << 16
	timeout = time.Second * 3
)

func init() {
	// Initialized local server ip address
	addrs, err := net.Interfaces()
	if err != nil {
		clog.Errorf("get server IPv4 address failed: %s", err)
	}

	for _, face := range addrs {
		adders, err := face.Addrs()
		if err != nil {
			clog.Errorf("get server IPv4 address failed: %s", err)
		}

		for _, addr := range adders {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ipv4 = ipNet.IP.String()
					return
				}
			}
		}
	}
}

type HttpServer struct {
	serv *http.Server
	port int
}

type Options struct {
	Port int
	Auth string
	// CertMagic *tls.Config
}

// New 创建一个新的 HTTP 服务器
func New(opt *Options) (*HttpServer, error) {
	if opt.Port < minPort || opt.Port > maxPort {
		return nil, errors.New("HTTP server port illegal")
	}

	if opt.Auth != "" {
		authPassword = opt.Auth
	}

	hs := HttpServer{
		serv: &http.Server{
			Handler:      root,
			Addr:         net.JoinHostPort("0.0.0.0", strconv.Itoa(opt.Port)),
			WriteTimeout: timeout,
			ReadTimeout:  timeout,
		},
		port: opt.Port,
	}

	// 开启 HTTP Keep-Alive 长连接
	hs.serv.SetKeepAlivesEnabled(true)

	return &hs, nil
}

func (hs *HttpServer) SetupFS(fss *vfs.LogStructuredFS) {
	storage = fss
}

func (hs *HttpServer) SetAllowIP(allowd []string) {
	allowIpList = allowd
}

func (hs *HttpServer) Port() int {
	return hs.port
}

// IPv4 return local IPv4 address
func (hs *HttpServer) IPv4() string {
	return ipv4
}

// Startup blocking goroutine
func (hs *HttpServer) Startup() error {
	if storage == nil {
		return errors.New("file storage system is not initialized")
	}

	// 这个函数是一个阻塞函数
	err := hs.serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http api server :%w", err)
	}

	return nil
}

func (hs *HttpServer) Shutdown() error {
	// 先关闭 http 服务器停止接受数据请求
	err := hs.serv.Shutdown(context.Background())
	if err != nil && err != http.ErrServerClosed {
		// 这里发生了错误，外层处理这个错误时也要关闭文件存储系统
		inner := closeStorage()
		if inner != nil {
			return fmt.Errorf("failed to shutdown the server: %w", errors.Join(err, inner))
		}
		return err
	}
	return closeStorage()
}

func closeStorage() error {
	if storage != nil {
		err := storage.CloseFS()
		if err != nil {
			return err
		}
	}
	return nil
}
