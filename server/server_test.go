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
	"io/fs"
	"testing"
	"time"

	"github.com/auula/wiredb/vfs"
	"github.com/stretchr/testify/assert"
)

// 测试 New 方法
func TestNewHttpServer(t *testing.T) {
	opt := &Options{Port: 8080, Auth: "secret"}
	server, err := New(opt)
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, 8080, server.Port())

	// 测试端口非法情况
	opt = &Options{Port: 80} // 端口小于 1024
	server, err = New(opt)
	assert.Error(t, err)
	assert.Nil(t, server)
}

// 测试 HttpServer 的 IPv4 方法
func TestHttpServer_IPv4(t *testing.T) {
	server, _ := New(&Options{Port: 8080})
	assert.NotEmpty(t, server.IPv4())
}

// 测试 HttpServer 的 Port 方法
func TestHttpServer_Port(t *testing.T) {
	server, _ := New(&Options{Port: 8080})
	assert.Equal(t, 8080, server.Port())
}

// 测试 Startup 方法（非阻塞）
func TestHttpServer_Startup(t *testing.T) {
	server, err := New(&Options{Port: 8081})

	assert.NoError(t, err)
	// 启动服务器（在 goroutine 中运行）
	go func() {
		fss, err := vfs.OpenFS(&vfs.Options{
			FSPerm:    fs.FileMode(0755),
			Path:      "/tmp/wiredb",
			Threshold: 3,
		})
		assert.NoError(t, err)

		server.SetupFS(fss)
		err = server.Startup()
		assert.NoError(t, err)
	}()

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)
	assert.NoError(t, err)

	// 关闭服务器
	err = server.Shutdown()
	assert.NoError(t, err)
}

// 测试 SetupFS 方法
func TestHttpServer_SetupFS(t *testing.T) {
	hts, err := New(&Options{
		Port: 6379,
		Auth: "secret",
	})
	if err != nil {
		assert.NoError(t, err)
	}

	assert.NotNil(t, hts)

	fss, err := vfs.OpenFS(&vfs.Options{
		FSPerm:    fs.FileMode(0755),
		Path:      "/tmp/wiredb",
		Threshold: 3,
	})

	assert.NotNil(t, fss)

	if err != nil {
		assert.NoError(t, err)
	}

	hts.SetupFS(fss)
}

// 测试 Shutdown 方法
func TestHttpServer_Shutdown(t *testing.T) {
	hts, err := New(&Options{
		Port: 6379,
		Auth: "secret",
	})
	if err != nil {
		assert.NoError(t, err)
	}

	assert.NotNil(t, hts)

	go func() {
		err := hts.Startup()
		assert.NoError(t, err)
	}()

	err = hts.Shutdown()
	assert.NoError(t, err)
}
