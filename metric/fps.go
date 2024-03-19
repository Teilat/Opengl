package metric

import (
	"context"
	"fmt"
	"time"
)

const TickerResolution = time.Millisecond

type Fps struct {
	frames uint64
	fps    *string
}

func NewFPSMeter() *Fps {
	s := fmt.Sprintf("fps:%d", 0)
	f := Fps{
		frames: 0,
		fps:    &s,
	}
	return &f
}

func (f *Fps) run(ctx context.Context, str *string, start time.Time) {
	ticker := time.NewTicker(TickerResolution)
	for {
		select {
		case <-ticker.C:
			*str = fmt.Sprintf("fps:%f", float64(f.frames)/time.Since(start).Seconds())
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (f *Fps) Start(ctx context.Context) {
	go f.run(ctx, f.fps, time.Now())
}

func (f *Fps) GetString() *string {
	return f.fps
}

func (f *Fps) Tick() {
	f.frames++
}
