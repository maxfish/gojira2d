package utils

type FPSCounter struct {
	frames      uint32
	accumulator float64
	fps         uint32
}

// deltaTime: seconds since the previous frame
// updateRate: rate, in seconds, at which the FPS are computed
func (f *FPSCounter) Update(deltaTime float64, updateRate uint32) {
	f.frames++
	f.accumulator += deltaTime // seconds
	if f.accumulator > float64(updateRate) {
		f.fps = f.frames / updateRate
		f.accumulator -= float64(updateRate)
		f.frames = 0
	}
}

func (f *FPSCounter) FPS() uint32 {
	return f.fps
}
