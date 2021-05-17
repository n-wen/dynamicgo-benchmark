package main

import (
	"fmt"
	"net/http"

	lua "github.com/Shopify/go-lua"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

const PINGURL = "http://localhost:10000/ping"

func Request(l *lua.State) int {
	_, err := req.Get(PINGURL)
	if err != nil {
		return 1
	}
	return 0
}

const luascript = `
req.request()
`

func callLua() {
	L := lua.NewState()
	// lua.OpenLibraries(L)
	lua.NewMetaTable(L, "req")
	L.PushValue(-1)
	L.SetField(-2, "__index")

	lua.SetFunctions(L, []lua.RegistryFunction{
		{"request", Request},
	}, 0)
	L.SetGlobal("req")

	if err := lua.DoString(L, luascript); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	go func() {
		web := gin.Default()
		web.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		web.Run(":10000")
	}()
	e := gin.Default()
	e.GET("/test", func(c *gin.Context) {
		callLua()
		c.String(http.StatusOK, "ok")
	})
	pprof.Register(e)
	e.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404")
	})
	e.Run(":3000")
}
