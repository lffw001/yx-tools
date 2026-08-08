package app

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Event struct {
	Type     string   `json:"type"`
	Stage    string   `json:"stage,omitempty"`
	Message  string   `json:"message,omitempty"`
	Current  int      `json:"current,omitempty"`
	Total    int      `json:"total,omitempty"`
	Results  []Result `json:"results,omitempty"`
	Finished bool     `json:"finished"`
	At       int64    `json:"at"`
}

type Runner struct {
	mu       sync.RWMutex
	running  bool
	cancel   context.CancelFunc
	history  []Event
	results  []Result
	subs     map[chan Event]struct{}
	lastOpts Options
}

func NewRunner() *Runner {
	return &Runner{subs: make(map[chan Event]struct{})}
}

func (r *Runner) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 64)
	r.mu.Lock()
	for _, e := range r.history {
		select {
		case ch <- e:
		default:
		}
	}
	r.subs[ch] = struct{}{}
	r.mu.Unlock()
	return ch, func() {
		r.mu.Lock()
		if _, ok := r.subs[ch]; ok {
			delete(r.subs, ch)
			close(ch)
		}
		r.mu.Unlock()
	}
}

func (r *Runner) broadcast(e Event) {
	e.At = time.Now().UnixMilli()
	r.mu.Lock()
	r.history = append(r.history, e)
	if len(r.history) > 200 {
		r.history = r.history[len(r.history)-200:]
	}
	for ch := range r.subs {
		select {
		case ch <- e:
		default:
		}
	}
	r.mu.Unlock()
}

func (r *Runner) Running() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

func (r *Runner) Results() []Result {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Result, len(r.results))
	copy(out, r.results)
	return out
}

func (r *Runner) LastStage() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := len(r.history) - 1; i >= 0; i-- {
		if r.history[i].Stage != "" {
			return r.history[i].Stage
		}
	}
	return ""
}

func (r *Runner) LastOptions() Options {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastOpts
}

func (r *Runner) Start(o Options) bool {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return false
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.running = true
	r.cancel = cancel
	r.history = nil
	r.lastOpts = o
	r.mu.Unlock()

	go func() {
		defer func() {
			r.mu.Lock()
			r.running = false
			r.cancel = nil
			r.mu.Unlock()
			cancel()
		}()
		rs, err := Run(ctx, o, func(p Progress) {
			r.broadcast(Event{Type: "progress", Stage: p.Stage, Message: p.Message, Current: p.Current, Total: p.Total})
		})
		if err != nil {
			if ctx.Err() != nil {
				r.broadcast(Event{Type: "error", Message: "已取消", Finished: true})
				return
			}
			r.broadcast(Event{Type: "error", Message: err.Error(), Finished: true})
			return
		}
		r.mu.Lock()
		r.results = rs
		r.mu.Unlock()
		if err := WriteCSV(ResultFile, rs); err != nil {
			r.broadcast(Event{Type: "log", Message: "结果写入失败: " + err.Error()})
		}
		r.broadcast(Event{Type: "done", Message: "测速完成", Results: rs, Total: len(rs), Current: len(rs), Finished: true})
	}()
	return true
}

func (r *Runner) Cancel() {
	r.mu.RLock()
	c := r.cancel
	r.mu.RUnlock()
	if c != nil {
		c()
	}
}

func SystemInfo() (map[string]any, error) {
	selfPath, _ := os.Executable()
	cronSupported := false
	if _, err := exec.LookPath("crontab"); err == nil {
		cronSupported = true
	}
	return map[string]any{
		"cron_supported": cronSupported,
		"self_path":      selfPath,
		"result_file":    ResultFile,
		"proxy_file":     ProxyListFile,
		"default_url":    DefaultTestURL,
		"local_ipv4":     LocalIPv4(),
		"local_ipv6":     LocalIPv6(),
		"public_ipv4":    publicIPv4,
		"public_ipv6":    publicIPv6,
	}, nil
}
