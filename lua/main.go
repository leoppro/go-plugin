package main

import (
	"fmt"
	"os"
	"time"

	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"github.com/pingcap/log"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func main() {
	pluginPath := "./plugin.lua"
	if len(os.Args) >= 2 {
		pluginPath = os.Args[1]
	}
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(pluginPath); err != nil {
		panic(err)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("max"),
		NRet:    1,
		Protect: true,
	}, lua.LNumber(10), lua.LNumber(20)); err != nil {
		panic(err)
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)         // remove received value

	fmt.Printf("%s", ret.String())
	table := L.NewTable()
	table.RawSet(lua.LString("table"), lua.LString("t1"))
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("emit_row"),
		NRet:    1,
		Protect: true,
	}, table); err != nil {
		panic(err)
	}
	ret = L.Get(-1) // returned value
	L.Pop(1)        // remove received value
	fmt.Printf("%s", ret.String())

	log.Info("run benckmark =================")
	benckmark(L)
	log.Info("finished benckmark ============")
}

func tranRow2Lua(L *lua.LState, row *sink.RowChangedEvent) *lua.LTable {
	table := L.NewTable()
	table.RawSet(lua.LString("startTs"), lua.LNumber(row.StartTs))
	table.RawSet(lua.LString("commitTs"), lua.LNumber(row.CommitTs))
	table.RawSet(lua.LString("table"), lua.LString(row.Table.Table))
	columns := L.NewTable()
	for columnName, columnValue := range row.Columns {

		columns.RawSet(lua.LString(columnName), lua.LString(fmt.Sprintf("%s", columnValue.Value)))
	}
	table.RawSet(lua.LString("column"), columns)

	return table
}

func benckmark(L *lua.LState) {
	emitRowFunc := lua.P{
		Fn:      L.GetGlobal("emit_row"),
		NRet:    1,
		Protect: true,
	}

	startTime := time.Now()
	for i := int64(0); i < 5_000_000; i++ {
		if err := L.CallByParam(emitRowFunc, tranRow2Lua(L, newRow(i))); err != nil {
			panic(err)
		}
		L.Pop(1)
	}
	log.Info("-", zap.Duration("cost", time.Since(startTime)), zap.Duration("op", time.Since(startTime)/5_000_000))
}

func newRow(mark int64) *sink.RowChangedEvent {
	return &sink.RowChangedEvent{
		RowID:   mark,
		StartTs: uint64(mark),
		Table: &sink.TableName{
			Table: "123",
		},
	}
}
