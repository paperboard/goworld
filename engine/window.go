package engine

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/johanhenriksson/goworld/engine/keys"
	"github.com/johanhenriksson/goworld/engine/mouse"
	"github.com/johanhenriksson/goworld/render"
)

func init() {
	// glfw event handling must run on the main OS thread
	runtime.LockOSThread()

	// init glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
}

// UpdateCallback defines the window update callback function
type UpdateCallback func(float32)

// RenderCallback defines the window render callback function
type RenderCallback func(*Window, float32)

// ResizeCallback defines the window resize callback function
type ResizeCallback func(*Window, int, int)

// Window represents the main engine window
type Window struct {
	Wnd           *glfw.Window
	Width         int
	Height        int
	HighDPI       bool
	FPS           float32
	Scale         float32
	focused       bool
	updateCb      UpdateCallback
	renderCb      RenderCallback
	resizeCb      ResizeCallback
	maxFrameTime  float64
	lastFrameTime float64
}

// CreateWindow creates the main engine window, and the OpenGL context
func CreateWindow(title string, width int, height int, highDPI bool) *Window {
	/* GLFW Window settings */
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 1)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(10)

	/* Initialize OpenGL */
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// get actual window size
	width, height = window.GetSize()

	w := &Window{
		Wnd:           window,
		Width:         width,
		Height:        height,
		HighDPI:       highDPI,
		maxFrameTime:  0.0,
		lastFrameTime: glfw.GetTime(),
	}
	w.SetMaxFps(60)

	// set the dimensions of the engine output buffer
	buffw, buffh := w.GetBufferSize()
	render.ScreenBuffer.Width = buffw
	render.ScreenBuffer.Height = buffh
	w.Scale = float32(buffw) / float32(width)

	log.Printf("Created window of size %dx%d, framebuffer %dx%d, scale: %.2f", width, height, buffw, buffh, w.Scale)

	window.SetKeyCallback(keys.KeyCallback)
	window.SetMouseButtonCallback(mouse.ButtonCallback)
	window.SetCursorPosCallback(func(wnd *glfw.Window, x, y float64) {
		mouse.MoveCallback(wnd, x, y, w.Scale)
	})
	window.SetFocusCallback(func(wnd *glfw.Window, focused bool) {
		w.focused = focused
	})

	window.SetSizeCallback(func(wnd *glfw.Window, width, height int) {
		log.Printf("window resized to %dx%d\n", width, height)
		w.Width = width
		w.Height = height
		w.Scale = float32(render.ScreenBuffer.Width) / float32(width)
	})

	window.SetFramebufferSizeCallback(func(wnd *glfw.Window, width, height int) {
		render.ScreenBuffer.Width = width
		render.ScreenBuffer.Height = height
		w.Scale = float32(width) / float32(w.Width)
		log.Printf("framebuffer resized to %dx%d\n", width, height)

		if w.resizeCb != nil {
			w.resizeCb(w, width, height)
		}
	})

	w.LockCursor()
	return w
}

// Closed returns true if the window will close.
func (wnd *Window) Closed() bool {
	return wnd.Wnd.ShouldClose()
}

// SetMaxFps adjusts the maximum frame rate.
func (wnd *Window) SetMaxFps(fps int) {
	wnd.maxFrameTime = 1.0 / float64(fps)
}

// LockCursor locks cursor movement.
func (wnd *Window) LockCursor() {
	wnd.Wnd.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

// ReleaseCursor unlocks cursor movement.
func (wnd *Window) ReleaseCursor() {
	wnd.Wnd.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

// SetRenderCallback sets the engine render callback.
func (wnd *Window) SetRenderCallback(cb RenderCallback) {
	wnd.renderCb = cb
}

// SetUpdateCallback sets the engine update callback.
func (wnd *Window) SetUpdateCallback(cb UpdateCallback) {
	wnd.updateCb = cb
}

// SetResizeCallback sets the window resize callback.
func (wnd *Window) SetResizeCallback(cb ResizeCallback) {
	wnd.resizeCb = cb
}

// Loop runs the main engine loop.
func (wnd *Window) Loop() {
	fps := FpsCounter{}

	for !wnd.Closed() {
		// calculate frame delta time
		t := glfw.GetTime()
		dt := float32(t - wnd.lastFrameTime)
		wnd.lastFrameTime = t

		// todo: this is part of the FPS camera.
		// move it somewhere more appropriate
		if mouse.Down(mouse.Button1) {
			wnd.LockCursor()
		} else {
			wnd.ReleaseCursor()
		}

		// update scene
		if wnd.updateCb != nil {
			wnd.updateCb(dt)
		}

		// render scene
		if wnd.focused {
			if wnd.renderCb != nil {
				wnd.renderCb(wnd, dt)
			}

			// end scene
			wnd.Wnd.SwapBuffers()
		}

		// get events
		glfw.PollEvents()

		mouse.Update(dt)
		keys.Update(dt)

		elapsed := glfw.GetTime() - t
		wnd.FPS = float32(fps.Append(elapsed))

		// wait a bit if we're running faster than the maximum fps
		if wnd.maxFrameTime > 0 {
			dur := wnd.maxFrameTime - elapsed

			if dur > 0 {
				time.Sleep(time.Duration(dur) * time.Second)
			}
		}
	}
}

// Terminate kills the GLFW window.
func (wnd *Window) Terminate() {
	glfw.Terminate()
}

// GetBufferSize returns the window framebuffer size.
func (wnd *Window) GetBufferSize() (int, int) {
	return wnd.Wnd.GetFramebufferSize()
}

// FpsCounter keeps a ring buffer of frame times to compute frames per second.
type FpsCounter struct {
	idx     int
	samples []float64
}

// Append a frame time to the buffer. Returns the current FPS.
func (fps *FpsCounter) Append(sample float64) float64 {
	length := 10
	if len(fps.samples) < length {
		fps.samples = append(fps.samples, sample)
	} else {
		fps.samples[fps.idx%length] = sample
	}
	fps.idx++

	sum := 0.0
	for _, sample := range fps.samples {
		sum += sample
	}
	avgTime := sum / float64(len(fps.samples))
	return 1.0 / avgTime
}
