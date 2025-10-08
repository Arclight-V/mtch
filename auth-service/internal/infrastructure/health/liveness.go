package health

import (
	"context"
	healt "github.com/Arclight-V/mtch/pkg/health"
)

type Liveness struct {
	l *healt.Liveness
}

func NewLiveness() *Liveness {
	return &Liveness{l: &healt.Liveness{}}
}

func (l *Liveness) Alive(ctx context.Context) error {
	return l.Alive(ctx)
}
