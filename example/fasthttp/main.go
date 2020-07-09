package main

import (
	"github.com/valyala/fasthttp"
	"github.com/wosai/ultron"
)

const (
	api = "http://10.0.0.30/benchmark"
)

func main() {
	attacker := ultron.NewFastHTTPAttacker("benchmark", func(r *fasthttp.Request) error {
		r.SetRequestURI(api)
		return nil
	})
	task := ultron.NewTask()
	task.Add(attacker, 1)

	//ultron.LocalEventHook.Concurrency = 0

	ultron.LocalRunner.Config.Concurrence = 100
	ultron.LocalRunner.Config.HatchRate = 10
	ultron.LocalRunner.Config.MinWait = ultron.ZeroDuration
	ultron.LocalRunner.Config.MaxWait = ultron.ZeroDuration

	ultron.LocalRunner.WithTask(task)
	ultron.LocalRunner.Start()
}
