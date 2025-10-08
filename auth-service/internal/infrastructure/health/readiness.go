package health

import (
	"context"
	healt "github.com/Arclight-V/mtch/pkg/health"
	"time"
)

type StartupGate struct {
	sg *healt.StartupGate
}

func NewStartupGate(cacheTTL time.Duration) *StartupGate {
	return &StartupGate{sg: healt.NewStartupGate(cacheTTL)}
}

func (g *StartupGate) MarkReady() {
	g.sg.MarkReady()
}

func (g *StartupGate) Ready(ctx context.Context) error {
	if err := g.sg.Ready(ctx); err != nil {
		return err
	}

	return nil
}
