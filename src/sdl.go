package main

import "github.com/veandco/go-sdl2/sdl"

type SDLDisplay struct {
	surface       *sdl.Surface
	window        *sdl.Window
	scalingFactor int32
	bg            uint32
	fg            uint32
}

// SDLInit : Initialise SDL window with scaling factor
func (display *SDLDisplay) init(scalingFactor int32, bg uint32, fg uint32) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	check(err)
	display.window, err = sdl.CreateWindow("chip8go", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		scalingFactor*64, scalingFactor*32, sdl.WINDOW_SHOWN)
	check(err)

	display.surface, err = display.window.GetSurface()
	check(err)
	err = display.surface.FillRect(nil, bg)
	check(err)

	display.bg = bg
	display.fg = fg
	display.scalingFactor = scalingFactor
}

func (display *SDLDisplay) drawPixel(x int32, y int32) {
	rect := sdl.Rect{
		X: x * display.scalingFactor,
		Y: y * display.scalingFactor,
		W: display.scalingFactor,
		H: display.scalingFactor}
	err := display.surface.FillRect(&rect, display.fg)
	check(err)
}
func (display *SDLDisplay) clearDisplay() {
	err := display.surface.FillRect(nil, display.bg)
	check(err)

}

func (display *SDLDisplay) updateDisplay() {
	err := display.window.UpdateSurface()
	check(err)
}

func (display *SDLDisplay) Destroy() {
	err := display.window.Destroy()
	check(err)
}
