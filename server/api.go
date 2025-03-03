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
	"fmt"
	"net/http"
	"strings"

	"github.com/auula/wiredb/clog"
	"github.com/gin-gonic/gin"
)

const version = "wiredb/1.0.0"

var (
	root         *gin.Engine
	authPassword string
	allowIpList  []string
)

// http://192.168.101.225:2668/{types}/{key}
// POST 创建 http://192.168.101.225:2668/zset/user-01-score
// PUT  更新 http://192.168.101.225:2668/zset/user-01-score
// GET  获取 http://192.168.101.225:2668/table/user-01-shop-cart

func init() {
	gin.SetMode(gin.ReleaseMode)
	root = gin.New()

	root.Use(authMiddleware())
	root.NoRoute(Error404Handler)
	root.GET("/", GetHealthController)

	set := root.Group("/set")
	{
		set.GET("/:key", GetSetController)
		set.PUT("/:key", PutSetController)
		set.DELETE("/:key", DeleteSetController)
	}

	zset := root.Group("/zset")
	{
		zset.GET("/:key", GetZsetController)
		zset.PUT("/:key", PutZsetController)
		zset.DELETE("/:key", DeleteZsetController)
	}

	list := root.Group("/list")
	{
		list.GET("/:key", GetListController)
		list.PUT("/:key", PutListController)
		list.DELETE("/:key", DeleteListController)
	}

	text := root.Group("/text")
	{
		text.GET("/:key", GetTextController)
		text.PUT("/:key", PutTextController)
		text.DELETE("/:key", DeleteTextController)
	}

	table := root.Group("/table")
	{
		table.GET("/:key", GetTableController)
		table.PUT("/:key", PutTableController)
		table.DELETE("/:key", DeleteTableController)
	}

	number := root.Group("/number")
	{
		number.GET("/:key", GetNumberController)
		number.PUT("/:key", PutNumberController)
		number.DELETE("/:key", DeleteNumberController)
	}
}

type SystemInfo struct {
	KeyCount    int    `json:"key_count"`
	Version     string `json:"version"`
	GCState     int8   `json:"gc_state"`
	DiskFree    string `json:"disk_free"`
	DiskUsed    string `json:"disk_used"`
	DiskTotal   string `json:"disk_total"`
	MemoryFree  string `json:"mem_free"`
	MemoryTotal string `json:"mem_total"`
	DiskPercent string `json:"disk_percent"`
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", version)

		// 从请求头中获取 "Auth-Token" 字段的值
		auth := c.GetHeader("Auth-Token")
		clog.Debugf("HTTP request header authorization: %v", c.Request)

		// 获取客户端 IP 地址
		ip := c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip = c.ClientIP()
		}

		// 检查 IP 白名单
		if len(allowIpList) > 0 {
			ok := false
			for _, allowedIP := range allowIpList {
				// 只要找到匹配的 IP，就终止循环
				if allowedIP == strings.Split(ip, ":")[0] {
					ok = true
					break
				}
			}
			if !ok {
				clog.Warnf("Unauthorized IP address: %s", ip)
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": fmt.Sprintf("client IP %s is not allowed!", ip),
				})
				c.Abort()
				return
			}
		}

		if auth != authPassword {
			clog.Warnf("Unauthorized access attempt from client %s", ip)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "access not authorised!",
			})
			c.Abort()
			return
		}

		clog.Infof("Client %s connection successfully", ip)

		// 如果验证通过，继续执行后续的处理程序
		c.Next()
	}
}
