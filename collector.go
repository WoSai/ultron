package ultron

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type statsCollector struct {
	entries map[string]*statsEntry
	lock    sync.RWMutex
}

var defaultStatsCollector *statsCollector

func newStatsCollector() *statsCollector {
	return &statsCollector{
		entries: map[string]*statsEntry{},
	}
}

func (c *statsCollector) logSuccess(name string, t time.Duration) {
	if _, ok := c.entries[name]; !ok {
		c.lock.Lock()
		c.entries[name] = newStatsEntry(name)
		c.lock.Unlock()
	}
	c.entries[name].logSuccess(t)
}

func (c *statsCollector) logFailure(name string, err error) {
	if _, ok := c.entries[name]; !ok {
		c.lock.Lock()
		c.entries[name] = newStatsEntry(name)
		c.lock.Unlock()
	}
	Logger.Warn("occure error", zap.String("error", err.Error()))
	c.entries[name].logFailure(err)
}

// Receiving 主函数，监听channel进行统计
func (c *statsCollector) log(ret *QueryResult) {
	if ret.Error == nil {
		c.logSuccess(ret.Name, ret.Duration)
	} else {
		c.logFailure(ret.Name, ret.Error)
	}
}

func (c *statsCollector) report(full bool) map[string]*StatsReport {
	c.lock.RLock()
	defer c.lock.RUnlock()

	r := map[string]*StatsReport{}
	for k, v := range c.entries {
		r[k] = v.Report(full)
	}
	return r
}

func init() {
	defaultStatsCollector = newStatsCollector()
}