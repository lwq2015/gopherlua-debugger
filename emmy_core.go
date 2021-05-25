package lua_debugger

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

const (
	KeyDebuggerFcd = "__Debugger_Fcd"
)

func TcpConnect(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.CheckNumber(2)

	fcd := newFacade()
	fcdUd := L.NewUserData()
	fcdUd.Value = fcd
	L.SetField(L.Get(lua.RegistryIndex), KeyDebuggerFcd, fcdUd)

	if err := fcd.TcpConnect(L, host, int(port)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LTrue)
	return 1
}

func TcpClose(L *lua.LState) int {
	if fcdUd, ok := L.GetField(L.Get(lua.RegistryIndex), KeyDebuggerFcd).(*lua.LUserData); ok {
		if fcd, ok := fcdUd.Value.(*Facade); ok {
			if err := fcd.TcpClose(L); err != nil {
				L.Push(lua.LString(err.Error()))
				return 1
			}
		}
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
