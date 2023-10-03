package main

import (
	"flag"
	"runtime/debug"

	"github.com/polevpn/elog"
)

var glog *elog.EasyLogger
var handler *RequestHandler

func handlePanic() {
	if err := recover(); err != nil {
		elog.Error("Panic Exception:", err)
		elog.Error(string(debug.Stack()))
	}
}

func main() {
	flag.Parse()
	defer elog.Flush()
	glog = elog.GetLogger()

	handler = NewRequestHandler()
	server := NewHttpServer(handler)

	err := server.Listen("127.0.0.1:9011")

	if err != nil {
		elog.Error("start server fail:", err)
		return
	}
}
