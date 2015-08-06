package main

import (
	"game_engine/logs"
)

func main() {

	log := logs.NewLogger(10000)

	log.SetLogger("file", `{"filename":"test.log"}`)
	log.EnableFuncCallDepth(true)

	for i := 0; i < 1000; i++ {
		log.Trace("trace %s %s", "param1", "param2")
		log.Debug("debug")
		log.Info("info")
		log.Warn("warning")
		log.Error("error")
		log.Critical("critical")
	}

}
