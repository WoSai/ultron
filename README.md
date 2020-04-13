# Ultron
a http load testing tool in go

[![Go Report Card](https://goreportcard.com/badge/github.com/qastub/ultron)](https://goreportcard.com/report/github.com/qastub/ultron) [![Build Status](https://travis-ci.org/qastub/ultron.svg?branch=master)](https://travis-ci.org/qastub/ultron) [![codecov](https://codecov.io/gh/qastub/ultron/branch/master/graph/badge.svg)](https://codecov.io/gh/qastub/ultron)  [![GoDoc](https://godoc.org/github.com/qastub/ultron?status.svg)](https://godoc.org/github.com/qastub/ultron)

## Requirements

Go 1.12+

## Example

### **Script**

file path: `example/http/main.go`

```go
attacker := ultron.NewHTTPAttacker("benchmark", func() (*http.Request, error) { return http.NewRequest(http.MethodGet, "http://127.0.0.1/", nil) })
task := ultron.NewTask()
task.Add(attacker, 1)

ultron.LocalRunner.Config.Concurrence = 1000
ultron.LocalRunner.Config.HatchRate = 10
ultron.LocalRunner.Config.MinWait = ultron.ZeroDuration
ultron.LocalRunner.Config.MaxWait = ultron.ZeroDuration

ultron.LocalRunner.WithTask(task)
ultron.LocalRunner.Start()
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

## Architecture

### Local Runner
![LocalRunner](https://my-storage.oss-cn-shanghai.aliyuncs.com/picgo/20200413155141.png)

### Master Runner

![MasterRunner](https://my-storage.oss-cn-shanghai.aliyuncs.com/picgo/20200413155207.png)

### Slave Runer

![SlaveRunner](https://my-storage.oss-cn-shanghai.aliyuncs.com/picgo/20200413155607.png)