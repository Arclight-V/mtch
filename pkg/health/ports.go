package health

import "context"

type LivenessChecker interface {
	Alive(ctx context.Context) error
}

type ReadinessChecker interface {
	Ready(ctx context.Context) error
}
