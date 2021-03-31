package ffmpegwriter

import (
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"sync"
)

var Executable = "ffmpeg"

const DefaultCRF = 20

type Manager struct {
	frameChannel chan *image.RGBA
	width        int
	height       int
	pipe         io.WriteCloser
	cmd          *exec.Cmd
	waitgroup    *sync.WaitGroup
}

func OpenVideo(filename string, rate, width, height, crf int) (*Manager, error) {
	w := &Manager{
		frameChannel: make(chan *image.RGBA, 30), // buffer up to 30 frames
		width:        width,
		height:       height,
		waitgroup:    &sync.WaitGroup{},
	}

	cmd := exec.Command(
		Executable,
		"-y",             // overwrite
		"-f", "rawvideo", // stream bitmap data over STDIN
		"-pix_fmt", "rgba", // this lets us write directly from the image.RGBA pixels
		"-s", fmt.Sprintf("%dx%d", width, height), // size
		"-r", fmt.Sprintf("%d", rate), // framerate
		"-i", "-", // stdin
		"-pix_fmt", "yuv420p", "-profile:v", "main", "-c:v", "libx264", "-crf", fmt.Sprintf("%d", crf),
		filename,
	)

	w.waitgroup.Add(1)

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	w.cmd = cmd
	w.pipe = pipe
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()

	if err != nil {
		return nil, err
	}

	go w.startWriting()

	return w, nil
}

func (w *Manager) SaveFrame(frame *image.RGBA) {
	w.frameChannel <- frame
}

func (w *Manager) MakeFrame() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w.width, w.height))
}

func (w *Manager) Done() error {
	close(w.frameChannel)
	w.waitgroup.Wait()
	return w.cmd.Wait()
}

func (w *Manager) startWriting() {
	for img := range w.frameChannel {
		w.writeFrame(img)
	}
	w.pipe.Close()
	w.waitgroup.Done()
}

func (w *Manager) writeFrame(img *image.RGBA) {
	w.pipe.Write(img.Pix)
}
