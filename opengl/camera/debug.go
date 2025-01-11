package camera

import (
	"context"
	"fmt"
	"time"
)

type Debug interface {
	Start(ctx context.Context)
	Stop()

	GetLookAtString() *string
	GetFovString() *string
	GetPosString() *string
}

type debug struct {
	ctx    context.Context
	cancel context.CancelFunc

	ticker *time.Ticker
	cam    *Camera

	lockAt *string
	fov    *string
	pos    *string
}

func NewCameraDebug(cam *Camera, ticker *time.Ticker) Debug {
	dCamPos := ""
	dLookAt := ""
	dFov := ""
	return &debug{
		lockAt: &dLookAt,
		fov:    &dFov,
		pos:    &dCamPos,

		ticker: ticker,
		cam:    cam,
	}
}

func (d *debug) Start(ctx context.Context) {
	d.ctx, d.cancel = context.WithCancel(ctx)
	go d.run()
}

func (d *debug) Stop() {
	d.cancel()
}

func (d *debug) GetLookAtString() *string {
	if d == nil {
		return nil
	}
	return d.lockAt
}

func (d *debug) GetPosString() *string {
	if d == nil {
		return nil
	}
	return d.pos
}

func (d *debug) GetFovString() *string {
	if d == nil {
		return nil
	}
	return d.fov
}

func (d *debug) run() {
	for {
		select {
		case <-d.ticker.C:
			*d.lockAt = fmt.Sprintf("camera pos:%v", d.cam.GetPos())
			*d.fov = fmt.Sprintf("look at:%v", d.cam.GetLookAt())
			*d.pos = fmt.Sprintf("fov:%v", d.cam.GetFov())
		case <-d.ctx.Done():
			return
		}
	}
}
