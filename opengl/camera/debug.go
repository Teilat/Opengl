package camera

import (
	"context"
	"fmt"
	"opengl/metric"
	"time"
)

type Debug struct {
	lockAt *string
	fov    *string
	pos    *string
}

func (d *Debug) run(ctx context.Context, cam *Camera) {
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

func (d *Debug) GetLookAtString() *string {
	return d.lockAt
}
func (d *Debug) GetPosString() *string {
	return d.pos
}
func (d *Debug) GetFovString() *string {
	return d.fov
}
