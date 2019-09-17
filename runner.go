package ultron

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const (
	// StatusIdle 空闲状态
	StatusIdle Status = iota
	// StatusBusy 执行中状态
	StatusBusy
	// StatusStopped 已经停止状态
	StatusStopped
)

var (
	// LocalRunner 单机执行入口
	LocalRunner *localRunner

	// SlaveRunner 分布式执行，节点执行入口
	//SlaveRunner *slaveRunner
)

type (
	// Runner 定义执行器接口
	Runner interface {
		WithConfig(*RunnerConfig)
		WithTask(*Task)
		GetConfig() *RunnerConfig
		GetStatus() Status
		Start()
		Done()
	}

	// Status runer状态
	Status = uint32

	baseRunner struct {
		Config      *RunnerConfig
		task        *Task
		status      Status
		workerCount uint32
		cancelCh    chan context.CancelFunc // worker的cancelFunc队列，用于通知结束工作
		wg          sync.WaitGroup
	}

	localRunner struct {
		stats *summaryStatistics
		once  sync.Once
		*baseRunner
	}
)

func newBaseRunner() *baseRunner {
	return &baseRunner{status: StatusIdle, Config: DefaultRunnerConfig}
}

func (br *baseRunner) WithConfig(rc *RunnerConfig) {
	br.Config = rc
}

func (br *baseRunner) WithTask(t *Task) {
	br.task = t
}

func (br *baseRunner) GetConfig() *RunnerConfig {
	return br.Config
}

func (br *baseRunner) Done() {
	atomic.StoreUint32(&br.status, StatusStopped)
}

func (br *baseRunner) GetStatus() Status {
	return atomic.LoadUint32(&br.status)
}

func newLocalRunner(ss *summaryStatistics) *localRunner {
	return &localRunner{stats: ss, baseRunner: newBaseRunner()}
}

func checkRunner(br *baseRunner) error {
	if br.task == nil {
		return errors.New("no Task provided")
	}

	if br.Config == nil {
		return errors.New("no RunnerConfig provided")
	}

	if err := br.Config.check(); err != nil {
		return err
	}

	return nil
}

func (lr *localRunner) record(r *Result) {
	lr.stats.record(r)
}

func (lr *localRunner) Start() {
	if err := checkRunner(lr.baseRunner); err != nil {
		panic(err)
	}

	// 初始化取消函数队列，size大小为最大并发数
	lr.cancelCh = make(chan context.CancelFunc, lr.Config.findMaxConcurrence())

	Logger.Info("start to attack")
	atomic.StoreUint32(&lr.status, StatusBusy)

	lr.once.Do(func() {
		localReportPipeline = newReportPipeline(LocalReportPipelineBufferSize)
		localResultPipeline = newResultPipeline(LocalResultPipelineBufferSize)
		LocalEventHook.AddResultHandleFunc(lr.stats.record)
		LocalEventHook.listen(localResultPipeline, localReportPipeline)

		// ctrl+c退出,输出信号
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)
		go func() {
			<-signalCh
			Logger.Error("captured interrupt signal")
			printReportToConsole(lr.stats.report(true))
			os.Exit(1)
		}()
	})
	lr.stats.reset()

	// 定时输出压测数据
	feedTicker := time.NewTicker(StatsReportInterval)
	go func() {
		for range feedTicker.C {
			localReportPipeline <- lr.stats.report(false) // 管道传输统计结果
		}
	}()

	nextStage := make(chan struct{}, 1)
	go func() {
		t := time.NewTicker(200 * time.Millisecond)
		for range t.C {
			sf, tf := lr.isFinishedCurrentStage()
			if sf {
				nextStage <- struct{}{} // 发送信号，开启下一阶段
			}

			if tf {
				t.Stop()
				break
			}
		}
	}()

	for _, stage := range lr.Config.Stages {
		Logger.Info("current stage info", zap.Any("stage", stage))
		lr.hatchWorkersOnStage(stage, localResultPipeline)
		<-nextStage // 阻塞，直到当前stage结束
	}

	// 等待所有worker的goroutine退出
	lr.wg.Wait()

	feedTicker.Stop()
	localReportPipeline <- lr.stats.report(true)

	Logger.Info("task done")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

// hatchWorkersOnStage 每阶段增压、减压逻辑
func (br *baseRunner) hatchWorkersOnStage(s *Stage, ch resultPipeline) {
	var batch int
	var wg sync.WaitGroup

	ticker := time.NewTicker(time.Second) // 不适用time.Sleep的原因，是因为cancelFunc存储于channel中，会存在阻塞的可能
	batches := s.hatchWorkerCounts()
	Logger.Info("will hatch many workers", zap.Ints("batches", batches))

	if batch == 0 {
		go func(b int) {
			wg.Add(1)
			defer wg.Done()
			br.hatchOrKillWorker(batches[b], ch)
			atomic.AddUint32(&br.workerCount, uint32(batches[b]))
			Logger.Info(fmt.Sprintf("hatched %d workers", atomic.LoadUint32(&br.workerCount)))
		}(batch)
		batch++
	}

	for batch < len(batches) {
		select {
		case <-ticker.C:
			if batch >= len(batches) {
				break
			}
			go func(b int) {
				wg.Add(1)
				defer wg.Done()
				br.hatchOrKillWorker(batches[b], ch)
				atomic.AddUint32(&br.workerCount, uint32(batches[b]))
				Logger.Info(fmt.Sprintf("hatched %d workers", atomic.LoadUint32(&br.workerCount)))
			}(batch)
			batch++

		default:
		}
	}

	ticker.Stop()
	wg.Wait()

	if s.Duration > ZeroDuration {
		go func(s *Stage) {
			Logger.Info(fmt.Sprintf("current stage will expried at %s", time.Now().Add(s.Duration).String()))
			t := time.NewTimer(s.Duration)
			<-t.C
			atomic.StoreUint32(&s.expired, StageExpired)
			t.Stop()
		}(s)
	}
}

// hatchOrKillWorker 增压、减压的具体实现
func (br *baseRunner) hatchOrKillWorker(n int, ch resultPipeline) {
	if n > 0 { // 该阶段加压
		for i := 0; i < n; i++ {
			go br.doCancelableWork(ch)
		}
	} else if n < 0 { // 降压阶段
		for i := 0; i > n; i-- {
			cancel := <-br.cancelCh
			cancel()
		}
	}
}

// doCancelableWork 生成一个goroutine，执行Attacker.Fire
func (br *baseRunner) doCancelableWork(ch resultPipeline) {
	br.wg.Add(1)
	defer func() {
		if rec := recover(); rec != nil {
			debug.PrintStack()
			Logger.Error("recovered", zap.Any("recover", rec))
		}
		br.wg.Done()
	}()

	if ch == nil {
		panic("invalid result pipeline")
	}

	ctx, cancel := context.WithCancel(context.Background())
	br.cancelCh <- cancel

	for {
		q := br.task.pickUp()
		start := time.Now()
		if br.GetStatus() == StatusStopped {
			return
		}

		select {
		case <-ctx.Done():
			//Logger.Info("this work was canceled")
			return
		default:
		}

		err := q.Fire()
		duration := time.Since(start)
		_, stage := br.Config.CurrentStage()
		atomic.AddUint64(&stage.counts, 1)
		ret := newResult(q.Name(), duration, err)
		ch <- ret

		if err != nil {
			Logger.Warn("occur error: " + err.Error())
		}

		if br.GetStatus() == StatusStopped {
			return
		}
		select {
		case <-ctx.Done():
			//Logger.Info("this work was canceled")
			return
		default:
		}

		br.Config.block()
	}
}

// isFinishedCurrentStage 判断当前stage是否满足退出条件，第一个返回值表示当前stage是否结束，第二个返回标识全局任务是否结束
func (br *baseRunner) isFinishedCurrentStage() (bool, bool) {
	if br.GetStatus() == StatusStopped {
		return true, true
	}

	index, stage := br.Config.CurrentStage()
	if (stage.Requests > 0 && atomic.LoadUint64(&stage.counts) >= stage.Requests) ||
		(atomic.LoadUint32(&stage.expired) == StageExpired) {
		_, _, f := br.Config.finishCurrentStage(index)
		Logger.Info("current stage is finished")
		if f {
			br.Done()
			return true, true
		}
		return true, false
	}
	return false, false
}
