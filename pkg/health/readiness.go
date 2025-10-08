package health

import (
	"context"
	"sync/atomic"
	"time"
)

// StartupGate
type StartupGate struct {
	ready atomic.Bool
	//TTL cache Ready (so as not to pull too often)
	lastOKAt atomic.Int64 // unix nano
	cacheTTL time.Duration
}

func NewStartupGate(cacheTTL time.Duration) *StartupGate {
	g := &StartupGate{cacheTTL: cacheTTL}
	g.ready.Store(false)
	return g
}

func (g *StartupGate) MarkReady() { g.ready.Store(true) }

func (g *StartupGate) Ready(ctx context.Context) error {
	if !g.ready.Load() {
		return ErrNotReady
	}
	// little cache
	if ttl := g.cacheTTL; ttl > 0 {
		now := time.Now().UnixNano()
		if last := g.lastOKAt.Load(); now-last < ttl.Nanoseconds() {
			return nil
		}
		g.lastOKAt.Store(now)
	}
	return nil
}

var ErrNotReady = errNotReady("startup not completed")

type errNotReady string

func (e errNotReady) Error() string { return string(e) }
