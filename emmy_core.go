package lua_debugger

import (
	lua "github.com/yuin/gopher-lua"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TcpConnect(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.CheckNumber(2)
	if err := Fcd.TcpConnect(L, host, int(port)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LTrue)
	return 1
}

func TcpClose(L *lua.LState) int {
	if err := Fcd.TcpClose(L); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

var coreApi = map[string]lua.LGFunction{
	"tcpConnect": TcpConnect,
	"tcpClose":   TcpClose,
}

func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, coreApi)
	L.Push(t)
	return 1
}

func Preload(L *lua.LState) {
	L.PreloadModule("emmy_core", Loader)
}