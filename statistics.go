package ultron

import (
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	timeDistributions = [10]float64{0.5, 0.6, 0.7, 0.8, 0.9, 0.95, 0.97, 0.98, 0.99, 1.0}
)

const (
	// ZeroDuration 0，用于一些特殊判断
	ZeroDuration = time.Duration(0)
	// StatsReportInterval 统计报表输出间隔
	StatsReportInterval = 5 * time.Second
)

type (
	attackerStatistics struct {
		name              string                       // 统计对象名称
		numRequests       int64                        // 成功请求次数
		numFailures       int64                        // 失败次数
		totalResponseTime time.Duration                // 成功请求的response time总和，基于原始的响应时间
		minResponseTime   time.Duration                // 最快的响应时间
		maxResponseTime   time.Duration                // 最慢的响应时间
		trendSuccess      *limitedSizeMap              // 按时间轴（秒级）记录成功请求次数
		trendFailures     *limitedSizeMap              // 按时间轴（秒级）记录错误的请求次数
		responseTimes     map[roundedMillisecond]int64 // 按优化后的响应时间记录成功请求次数
		failuresTimes     map[string]int64             // 记录不同错误的次数
		startTime         time.Time                    // 第一次收到请求的时间(包括成功以及失败)
		lastRequestTime   time.Time                    // 最后一次收到请求的时间(包括成功以及失败)
		interval          time.Duration                // 统计时间间隔，影响 CurrentQPS()
		lock              sync.Mutex
	}

	summaryStatistics struct {
		nodes sync.Map
	}

	limitedSizeMap struct {
		content map[int64]int64
		size    int64
	}

	// AttackerReport Attacker级别的报告
	AttackerReport struct {
		Name           string           `json:"name"`
		Requests       int64            `json:"requests"`
		Failures       int64            `json:"failures"`
		Min            int64            `json:"min"`
		Max            int64            `json:"max"`
		Median         int64            `json:"median"`
		Average        int64            `json:"average"`
		QPS            int64            `json:"qps"`
		Distributions  map[string]int64 `json:"distributions"`
		FailRatio      string           `json:"fail_ratio"`
		FailureDetails map[string]int64 `json:"failure_details"`
		FullHistory    bool             `json:"full_history"`
	}

	// Report Task级别的报告
	Report map[string]*AttackerReport

	roundedMillisecond int64
)

func timeDurationToRoundedMillisecond(t time.Duration) roundedMillisecond {
	ms := int64(t.Seconds()*1000 + 0.5)
	var rm roundedMillisecond
	if ms < 100 {
		rm = roundedMillisecond(ms)
	} else if ms < 1000 {
		rm = roundedMillisecond(ms + 5 - (ms+5)%10)
	} else {
		rm = roundedMillisecond(ms + 50 - (ms+50)%100)
	}
	return rm
}

func roundedMillisecondToDuration(r roundedMillisecond) time.Duration {
	return time.Duration(r * 1000 * 1000)
}

func timeDurationToMillisecond(t time.Duration) int64 {
	return int64(t) / int64(time.Millisecond)
}

func newLimitedSizeMap(s int64) *limitedSizeMap {
	return &limitedSizeMap{
		content: map[int64]int64{},
		size:    s,
	}
}

func (ls *limitedSizeMap) accumulate(k, v int64) {
	// 调用方自身来保证线程安全
	if _, ok := ls.content[k]; ok {
		ls.content[k] += v
	} else {
		for key := range ls.content {
			if (k - key) > ls.size {
				delete(ls.content, key)
			}
		}
		ls.content[k] = v
	}
}

func newAttackerStatistics(n string) *attackerStatistics {
	return &attackerStatistics{
		name:          n,
		trendSuccess:  newLimitedSizeMap(20),
		trendFailures: newLimitedSizeMap(20),
		responseTimes: map[roundedMillisecond]int64{},
		failuresTimes: map[string]int64{},
		interval:      12 * time.Second,
	}
}

func (as *attackerStatistics) logSuccess(ret *Result) {
	as.lock.Lock()

	now := time.Now()
	t := time.Duration(ret.Duration)
	as.numRequests++

	if as.startTime.IsZero() {
		as.startTime = now
		as.minResponseTime = t
	}
	as.lastRequestTime = now

	if t < as.minResponseTime {
		as.minResponseTime = t
	}
	if t > as.maxResponseTime {
		as.maxResponseTime = t
	}

	as.totalResponseTime += t
	//as.trendSuccess[now.Unix()]++
	as.trendSuccess.accumulate(now.Unix(), 1)
	as.responseTimes[timeDurationToRoundedMillisecond(t)]++

	as.lock.Unlock()
}

func (as *attackerStatistics) log(ret *Result) {
	if ret.Error == nil {
		as.logSuccess(ret) // 请求成功
		return
	}
	as.logFailure(ret) // 请求失败
}

func (as *attackerStatistics) logFailure(ret *Result) {
	as.lock.Lock()

	now := time.Now()
	if as.startTime.IsZero() {
		as.startTime = now
		as.minResponseTime = time.Duration(ret.Duration)
	}
	as.lastRequestTime = now

	//atomic.AddInt64(&as.numFailures, 1)
	as.numFailures++
	as.failuresTimes[ret.Error.Error()]++
	as.trendFailures.accumulate(now.Unix(), 1)

	as.lock.Unlock()
}

// totalQPS 获取总的QPS
func (as *attackerStatistics) totalQPS() float64 {
	return float64(as.numRequests) / as.lastRequestTime.Sub(as.startTime).Seconds()
}

// currentQPS 最近12秒的QPS
func (as *attackerStatistics) currentQPS() float64 {
	if as.startTime.IsZero() || as.lastRequestTime.IsZero() {
		return 0
	}

	now := time.Now().Add(-time.Second)            // 当前一秒可能未完成，统计时，往前推一秒
	start := now.Add(-(as.interval - time.Second)) // 比如当前15秒，往回推5秒，起点是11秒而不是10秒

	if start.Before(as.startTime) {
		start = as.startTime
	}

	//if now.Unix() == start.Unix() {
	//	return 0 // 相减会是0，不处理这种情况
	//}

	var total int64

	for i := start.Unix(); i <= now.Unix(); i++ {
		if v, ok := as.trendSuccess.content[i]; ok {
			total += v
		}
	}
	return float64(total) / float64(now.Unix()-start.Unix()+1)
}

// percentile 获取x%的响应时间
func (as *attackerStatistics) percentile(f float64) time.Duration {
	if f <= 0.0 {
		return as.minResponseTime
	}

	if f >= 1.0 {
		return as.maxResponseTime
	}

	hit := int64(float64(as.numRequests)*f + .5)
	if hit == as.numRequests {
		return as.maxResponseTime
	}

	var times []int
	for k := range as.responseTimes {
		times = append(times, int(k))
	}

	// sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	sort.Ints(times)

	for _, val := range times {
		counts := as.responseTimes[roundedMillisecond(val)]
		hit -= counts
		if hit <= 0 {
			return roundedMillisecondToDuration(roundedMillisecond(val))
		}
	}
	Logger.Error("unreachable code")
	return ZeroDuration // occur error
}

// min 最快响应时间
func (as *attackerStatistics) min() time.Duration {
	return as.minResponseTime
}

// max 最慢响应时间
func (as *attackerStatistics) max() time.Duration {
	return as.maxResponseTime
}

// average 平均响应时间
func (as *attackerStatistics) average() time.Duration {
	if as.numRequests == 0 {
		return ZeroDuration
	}
	return time.Duration(int64(as.totalResponseTime) / int64(as.numRequests))
}

// median 响应时间中位数
func (as *attackerStatistics) median() time.Duration {
	return as.percentile(.5)
}

// failRatio 错误率
func (as *attackerStatistics) failRatio() float64 {
	total := as.numFailures + as.numRequests
	if total == 0 {
		return 0.0
	}
	return float64(as.numFailures) / float64(total)
}

// report 打印统计结果
func (as *attackerStatistics) report(full bool) *AttackerReport {
	as.lock.Lock()

	r := &AttackerReport{
		Name:           as.name,
		Requests:       as.numRequests,
		Failures:       as.numFailures,
		Min:            timeDurationToMillisecond(as.min()),
		Max:            timeDurationToMillisecond(as.max()),
		Median:         timeDurationToMillisecond(as.median()),
		Average:        timeDurationToMillisecond(as.average()),
		Distributions:  map[string]int64{},
		FailRatio:      strconv.FormatFloat(as.failRatio()*100, 'f', 2, 64) + " %%",
		FailureDetails: as.failuresTimes,
		FullHistory:    full,
	}

	if full {
		r.QPS = int64(as.totalQPS())
	} else {
		r.QPS = int64(as.currentQPS())
	}

	for _, percent := range timeDistributions {
		r.Distributions[strconv.FormatFloat(percent, 'f', 2, 64)] = timeDurationToMillisecond(as.percentile(percent))
	}
	as.lock.Unlock()
	return r
}

func newSummaryStats() *summaryStatistics {
	return &summaryStatistics{}
}

func (ss *summaryStatistics) record(ret *Result) {
	val, _ := ss.nodes.LoadOrStore(ret.Name, newAttackerStatistics(ret.Name))
	val.(*attackerStatistics).log(ret)
}

func (ss *summaryStatistics) report(full bool) Report {
	rep := map[string]*AttackerReport{}

	ss.nodes.Range(func(key, value interface{}) bool {
		rep[key.(string)] = value.(*attackerStatistics).report(full)
		return true
	})

	return rep
}

// summaryStatistics 重置所有统计
func (ss *summaryStatistics) reset() {
	ss.nodes.Range(func(key, value interface{}) bool {
		ss.nodes.Delete(key)
		return true
	})
}
