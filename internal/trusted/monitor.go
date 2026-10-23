package trusted

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ZeroDos1988/godeb-tunnel-arch/internal/config"
	"github.com/ZeroDos1988/godeb-tunnel-arch/internal/tunnel"
)

// Monitor disconnects tunnels whose GoDeb-Trusted-Endpoint is reachable.
type Monitor struct {
	Interval time.Duration
}

func (m *Monitor) Run(ctx context.Context) {
	if m.Interval <= 0 {
		m.Interval = 15 * time.Second
	}
	log.Print("godeb-wg: trusted monitor running")
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()
	last := map[string]bool{}
	var mu sync.Mutex

	eval := func() {
		names, err := tunnel.List()
		if err != nil {
			return
		}
		for _, name := range names {
			text, err := tunnel.Read(name)
			if err != nil {
				continue
			}
			meta := config.ParseMeta(text)
			if meta.TrustedEndpoint == "" {
				continue
			}
			cctx, cancel := context.WithTimeout(ctx, 4*time.Second)
			ok := IsTrusted(cctx, meta.TrustedEndpoint)
			cancel()

			mu.Lock()
			prev, seen := last[name]
			last[name] = ok
			mu.Unlock()
			if seen && prev == ok {
				continue
			}
			log.Printf("godeb-wg: %s trusted=%v", name, ok)
			if ok {
				if tunnel.IsUp(name) {
					if err := tunnel.Down(name); err != nil {
						log.Printf("godeb-wg: down %s: %v", name, err)
					}
				}
			} else if seen && prev {
				if !tunnel.IsUp(name) {
					if err := tunnel.Up(name); err != nil {
						log.Printf("godeb-wg: up %s: %v", name, err)
					}
				}
			}
		}
	}

	eval()
	for {
		select {
		case <-ctx.Done():
			log.Print("godeb-wg: trusted monitor stopped")
			return
		case <-ticker.C:
			eval()
		}
	}
}
