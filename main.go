package main

import (
	lua "github.com/Shopify/go-lua"
	"github.com/imroc/req"
)

const PINGURL = "http://localhost:8080/ping"

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

func main() {
	L := lua.NewState()
	lua.OpenLibraries(L)
	lua.NewMetaTable(L, "req")
	L.PushValue(-1)
	L.SetField(-2, "__index")

	lua.SetFunctions(L, []lua.RegistryFunction{
		{"request", Request},
	}, 0)
	L.SetGlobal("req")

	if err := lua.DoString(L, luascript); err != nil {
		panic(err.Error())
	}

}
