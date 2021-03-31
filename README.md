# ffmpeg-writer

This is a simple go package for encoding video by creating individual frames. It is not particularly productionized, and is mostly useful for playing around with generating animations in go.

The ffmpeg process is managed in a goroutine and streamed data as you create it. On my 2019 MBP, simple animations at 1080p can be streamed out at ~60fps from the main thread.

This is extremely opinionated; you can specify the crf parameter, dimensions, and framerate, but otherwise, you get main-profile libx264.

Requires ffmpeg. If ffmpeg is not on your \$PATH, you can give an absolute path by changing `ffmpegwriter.Executable`.

# Usage

First, create a `Manager` by calling `OpenVideo`:

```go
m, _ := ffmpegwriter.OpenVideo(
    "movie.mp4", // output path
    30, // fps
    1024, // width
    768, // height
    ffmpegwriter.DefaultCRF, // constant rate factor; lower is higher quality
)
```

Next, request a frame from the `Manager`:

```go
frame := m.MakeFrame()
```

This will give you an `*image.RGBA` of the correct size to draw into. Once you've modified the frame the way you want it, save it to the video:

```go
m.SaveFrame(frame)
```

Of course, there's not much point in a video with a single frame. You probably want to use a loop. Note: when you call `SaveFrame`, you're relinquishing ownership of the frame. It is not safe to reuse the frame; ask the manager for a new one instead.

When you're done writing frames, end the output stream:

```
err := m.Done()
```

That's it! `Done` will block until ffmpeg finishes encoding.

See the demo for an example of a program that generates a 10 second 30fps video wiping the screen horizontally from red to blue.
