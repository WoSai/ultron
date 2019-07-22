package ultron

import (
	"time"
)

func newAttackerError(err error) *AttackerError {
	return &AttackerError{CausedBy: err.Error()}
}

func (m *AttackerError) Error() string {
	return m.CausedBy
}

func newResult(n string, d time.Duration, err error) *Result {
	if err == nil {
		return &Result{Name: n, Duration: int64(d)}
	}
	return &Result{Name: n, Duration: int64(d), Error: newAttackerError(err)}
}
