package main

import (
	"time"

	"github.com/qastub/ultron"
	"github.com/valyala/fasthttp"
)

const (
	api = "http://10.0.0.30/benchmark"
)

func main() {
	attacker := ultron.NewFastHTTPAttacker(
		"benchmark",
		func(r *fasthttp.Request) error {
			r.SetRequestURI(api)
			return nil
		})
	task := ultron.NewTask()
	task.Add(attacker, 1)

	ultron.LocalRunner.Config.MinWait = 2 * time.Second
	ultron.LocalRunner.Config.MaxWait = 3 * time.Second
	ultron.LocalRunner.Config.AppendStages(
		&ultron.Stage{Concurrence: 100, Duration: 3 * time.Minute, HatchRate: 10},
		&ultron.Stage{Concurrence: 500, Requests: 1000 * 1000, HatchRate: 30},
		&ultron.Stage{Concurrence: 300, HatchRate: 20},
	)

	ultron.LocalRunner.WithTask(task)
	ultron.LocalRunner.Start()
}
