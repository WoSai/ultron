package ultron

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLimitedSizeMap(t *testing.T) {
	m := newLimitedSizeMap(3)
	assert.EqualValues(t, 3, m.size)
}

func TestLimitedSizeMap_accumulateOK(t *testing.T) {
	m := newLimitedSizeMap(3)
	m.accumulate(100, 1)
	assert.EqualValues(t, 1, m.content[100])

	m.accumulate(100, 3)
	assert.EqualValues(t, 4, m.content[100])
}

func TestLimitedSizeMap_accumulateToRemoveEle(t *testing.T) {
	m := newLimitedSizeMap(3)
	m.accumulate(100, 1)
	assert.EqualValues(t, 1, m.content[100])
	assert.EqualValues(t, 1, len(m.content))

	m.accumulate(103, 1)
	assert.EqualValues(t, 2, len(m.content))

	m.accumulate(104, 1)
	assert.EqualValues(t, 2, len(m.content))
	_, ok := m.content[100]
	assert.False(t, ok)
}

func TestNewAttackerStatistics(t *testing.T) {
	s := newAttackerStatistics("foobar")
	assert.Equal(t, s.name, "foobar")
	assert.Equal(t, s.interval, 12*time.Second)
	assert.EqualValues(t, s.trendFailures.size, 20)
	assert.EqualValues(t, s.trendSuccess.size, 20)
	assert.True(t, s.lastRequestTime.IsZero())
	assert.True(t, s.startTime.IsZero())
}

func TestAttackerStatistics_logSuccess(t *testing.T) {
	name := "foobar"
	stats := newAttackerStatistics(name)

	ret := &Result{Name: name, Duration: int64(10 * time.Millisecond)}
	stats.logSuccess(ret)

	assert.False(t, stats.startTime.IsZero())
	assert.False(t, stats.lastRequestTime.IsZero())
	assert.EqualValues(t, stats.minResponseTime, 10*time.Millisecond)
	assert.EqualValues(t, stats.minResponseTime, stats.maxResponseTime)
	assert.EqualValues(t, stats.totalResponseTime, 10*time.Millisecond)
	assert.EqualValues(t, stats.numRequests, 1)

	ret = &Result{Name: name, Duration: int64(20 * time.Millisecond)}
	time.Sleep(time.Second)
	stats.log(ret)
	assert.True(t, stats.lastRequestTime.After(stats.startTime))
	assert.EqualValues(t, stats.maxResponseTime, 20*time.Millisecond)
	assert.EqualValues(t, stats.minResponseTime, 10*time.Millisecond)
	assert.EqualValues(t, stats.totalResponseTime, 30*time.Millisecond)
	assert.EqualValues(t, stats.numRequests, 2)
	assert.EqualValues(t, 2, len(stats.trendSuccess.content))

	ret = &Result{Name: name, Duration: int64(5 * time.Millisecond)}
	stats.log(ret)
	assert.True(t, stats.lastRequestTime.After(stats.startTime))
	assert.EqualValues(t, stats.minResponseTime, 5*time.Millisecond)
	assert.EqualValues(t, stats.totalResponseTime, 35*time.Millisecond)
	assert.EqualValues(t, stats.numRequests, 3)
}

func TestAttackerStatistics_logFailure(t *testing.T) {
	name := "foobar"
	stats := newAttackerStatistics(name)

	ret := &Result{Name: name, Duration: int64(10 * time.Millisecond), Error: newAttackerError(name, errors.New("error"))}
	stats.log(ret)
	assert.False(t, stats.startTime.IsZero())
	assert.Equal(t, stats.startTime, stats.lastRequestTime)
	assert.EqualValues(t, 1, stats.numFailures)
	assert.EqualValues(t, 0, stats.numRequests)
	assert.EqualValues(t, 1, stats.failuresTimes["error"])
	assert.EqualValues(t, 1, len(stats.trendFailures.content))
}

func TestAttackerStatistics_totalQPS(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.logSuccess(newResult("foobar", 100*time.Millisecond, nil))
	stats.logSuccess(newResult("foobar", 100*time.Millisecond, nil))
	assert.True(t, stats.totalQPS() > 0)
}

func TestAttackerStatistics_totalQPSNoSucceed(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.logFailure(newResult("foobar", 100*time.Millisecond, newAttackerError("foobar", errors.New("bad"))))
	stats.logFailure(newResult("foobar", 100*time.Millisecond, newAttackerError("foobar", errors.New("bad"))))
	assert.EqualValues(t, 0, stats.totalQPS())
}

func TestAttackerStatistics_currentQPSNoRequests(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	assert.EqualValues(t, 0, stats.currentQPS())
}

func TestAttackerStatistics_currentQPSToLastSecond(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.log(newResult("foobar", 100*time.Millisecond, nil))
	time.Sleep(1 * time.Second)
	stats.log(newResult("foobar", 100*time.Millisecond, nil))
	assert.EqualValues(t, 1, stats.currentQPS())

}

func TestAttackerStatistics_min(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.log(newResult("foobar", 100*time.Second, nil))
	stats.log(newResult("foobar", 10*time.Millisecond, nil))
	stats.log(newResult("foobar", time.Millisecond, errors.New("bad")))
	assert.Equal(t, stats.min(), 10*time.Millisecond)
}

func TestAttackerStatistics_max(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.log(newResult("foobar", 100*time.Second, nil))
	stats.log(newResult("foobar", 10*time.Millisecond, nil))
	stats.log(newResult("foobar", 101*time.Second, errors.New("bad")))
	assert.Equal(t, stats.max(), 100*time.Second)
}

func TestAttackerStatistics_avg(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	assert.EqualValues(t, 0, stats.average())

	stats.log(newResult("foobar", 1*time.Millisecond, nil))
	stats.log(newResult("foobar", 3*time.Millisecond, nil))
	assert.Equal(t, stats.average(), 2*time.Millisecond)
}

func TestAttackerStatistics_median(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.log(newResult("foobar", 1*time.Millisecond, nil))
	stats.log(newResult("foobar", 2*time.Millisecond, nil))
	stats.log(newResult("foobar", 5*time.Millisecond, nil))
	assert.Equal(t, stats.median(), 2*time.Millisecond)
}

func TestAttackerStatistics_failRation(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.log(newResult("foobar", 1*time.Millisecond, nil))
	assert.EqualValues(t, stats.failRatio(), 0)

	stats.log(newResult("foobar", 1*time.Millisecond, errors.New("bad")))
	assert.EqualValues(t, stats.failRatio(), 0.5)
}

func TestAttackerStatistics_report(t *testing.T) {
	stats := newAttackerStatistics("foobar")
	stats.log(newResult("foobar", 1*time.Millisecond, nil))
	rep := stats.report(false)
	assert.NotNil(t, rep)
}

func TestSummaryStatistics_record(t *testing.T) {
	ss := newSummaryStats()
	ss.record(newResult("foobar", 10*time.Millisecond, nil))
	var exists bool
	ss.nodes.Range(func(key, value interface{}) bool {
		if key.(string) == "foobar" {
			exists = true
			return false
		}
		return true
	})

	assert.True(t, exists)
}

func TestSummaryStatistics_report(t *testing.T) {
	ss := newSummaryStats()
	ss.record(newResult("foobar", 10*time.Millisecond, nil))
	rep := ss.report(true)
	_, ok := rep["foobar"]
	assert.True(t, ok)
}

func TestSummaryStatistics_rest(t *testing.T) {
	ss := newSummaryStats()
	ss.record(newResult("foobar", 10*time.Millisecond, nil))
	var counts int
	ss.nodes.Range(func(key, value interface{}) bool {
		counts++
		return true
	})
	assert.EqualValues(t, counts, 1)

	counts = 0
	ss.reset()
	ss.nodes.Range(func(key, value interface{}) bool {
		counts++
		return true
	})
	assert.EqualValues(t, counts, 0)
}

func BenchmarkAttackerStatistics_logSuccess(b *testing.B) {
	stats := newAttackerStatistics("foobar")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stats.log(newResult("foobar", 10*time.Millisecond, nil))
		}
	})
}
