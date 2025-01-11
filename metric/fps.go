package metric

import (
	"context"
	"fmt"
	"time"
)

const TickerResolution = time.Millisecond

type Fps struct {
	frames         float64
	totalFrameTime time.Duration
	str            *string
}

func NewFPSMeter() *Fps {
	s := fmt.Sprintf("str:%d", 0)
	f := Fps{
		frames: 0,
		str:    &s,
	}
	return &f
}

func (f *Fps) run(ctx context.Context, ticker *time.Ticker, str *string, start time.Time) {
	for {
		select {
		case <-ticker.C:
			*str = fmt.Sprintf("fps:%f frametime:%f", f.frames/time.Since(start).Seconds(), float64(f.totalFrameTime.Milliseconds())/f.frames)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (f *Fps) Start(ctx context.Context, ticker *time.Ticker) {
	go f.run(ctx, ticker, f.str, time.Now())
}

func (f *Fps) GetString() *string {
	return f.str
}

func (f *Fps) Tick(t time.Time) {
	f.frames++
	f.totalFrameTime += time.Since(t)
}
