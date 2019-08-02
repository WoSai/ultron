# Ultron
a http load testing tool in go

[![Go Report Card](https://goreportcard.com/badge/github.com/qastub/ultron)](https://goreportcard.com/report/github.com/qastub/ultron) [![Build Status](https://travis-ci.org/qastub/ultron.svg?branch=master)](https://travis-ci.org/qastub/ultron) [![codecov](https://codecov.io/gh/qastub/ultron/branch/master/graph/badge.svg)](https://codecov.io/gh/qastub/ultron)  [![GoDoc](https://godoc.org/github.com/qastub/ultron?status.svg)](https://godoc.org/github.com/qastub/ultron)

## Requirements

Go 1.12+

## Example

### **Script**

file path: `example/fasthttp/main.go`

```go
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
```

### Report

```json
{
  "benchmark": {
    "name": "benchmark",
    "requests": 1917994,
    "failures": 0,
    "min": 0,
    "max": 23,
    "median": 2,
    "average": 2,
    "qps": 50211,
    "distributions": {
      "0.50": 2,
      "0.60": 2,
      "0.70": 2,
      "0.80": 2,
      "0.90": 2,
      "0.95": 2,
      "0.97": 2,
      "0.98": 3,
      "0.99": 4,
      "1.00": 23
    },
    "failure_details": {},
    "full_history": false
  }
}
```
