/**
* @Author: Tanxf
* @Date: 2021/5/17 13:47
 */
package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	lua "github.com/yuin/gopher-lua"
)

type Service struct {
	pool *sync.Pool
}

func NewService() *Service {
	return &Service{
		pool: &sync.Pool{New: func() interface{} {
			L := lua.NewState()
			L.SetGlobal("request", L.NewFunction(Request))
			return L
		}},
	}
}

const PINGURL = "http://localhost:10000/ping"

func Request(l *lua.LState) int {
	_, err := req.Get(PINGURL)
	if err != nil {
		return 1
	}
	return 0
}

const luascript = `
request()
`

func (s *Service) callLua() {
	// lua.OpenLibraries(L)
	L := s.pool.Get().(*lua.LState)
	defer func() {
		s.pool.Put(L)
	}()

	if err := L.DoString(luascript); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	svc := NewService()
	go func() {
		web := gin.Default()
		web.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		web.Run(":10000")
	}()
	e := gin.Default()
	e.GET("/test", func(c *gin.Context) {
		svc.callLua()
		c.String(http.StatusOK, "ok")
	})
	pprof.Register(e)
	e.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404")
	})
	e.Run(":3000")
}
