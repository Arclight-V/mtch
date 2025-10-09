package health

import "context"

type Liveness struct{}

func (l *Liveness) Alive(ctx context.Context) error {
	return nil
}
