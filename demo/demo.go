package main

import (
	"ffmpegwriter"
	"image/color"
)

func main() {
	m, _ := ffmpegwriter.OpenVideo("out.mp4", 30, 300, 300, ffmpegwriter.DefaultCRF)
	for i := 0; i < 300; i++ {
		frame := m.MakeFrame()
		for x := 0; x < 300; x++ {
			var c color.RGBA
			if x < i {
				c = color.RGBA{0, 0, 255, 255}
			} else {
				c = color.RGBA{255, 0, 0, 255}
			}

			for y := 0; y < 300; y++ {
				frame.SetRGBA(x, y, c)
			}
		}

		m.SaveFrame(frame)
	}
	m.Done()
}
