package camera

import (
	"context"
	"fmt"
	"opengl/metric"
	"time"
)

type debug struct {
	lockAt *string
	fov    *string
	pos    *string
}

func (d *debug) run(ctx context.Context, cam *Camera) {
	ticker := time.NewTicker(metric.TickerResolution)
	for {
		select {
		case <-ticker.C:
			*d.lockAt = fmt.Sprintf("camera pos:%v", cam.GetPos())
			*d.fov = fmt.Sprintf("look at:%v", cam.GetLookAt())
			*d.pos = fmt.Sprintf("fov:%v", cam.GetFov())
		case <-ctx.Done():
			return
		}
	}
}

func (d *debug) GetLookAtString() *string {
	return d.lockAt
}
func (d *debug) GetPosString() *string {
	return d.pos
}
func (d *debug) GetFovString() *string {
	return d.fov
}
