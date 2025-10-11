package prober

import "sync"

type combined struct {
	mu     sync.Mutex
	probes []Probe
}

func Combine(probes ...Probe) Probe {
	return &combined{probes: probes}
}

func (c *combined) Healthy() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, probe := range c.probes {
		probe.Healthy()
	}
}

func (c *combined) NotHealthy(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, probe := range c.probes {
		probe.NotHealthy(err)
	}
}

func (c *combined) Ready() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, probe := range c.probes {
		probe.Ready()
	}
}

func (c *combined) NotReady(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, probe := range c.probes {
		probe.NotReady(err)
	}
}
