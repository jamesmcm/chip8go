package main

import (
	"flag"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/vharitonsky/iniflags"
)

func main() {
	wrapX := flag.String("wrapX", "on",
		"Wrap screen horizontally: on, off, error")
	wrapY := flag.String("wrapY", "on",
		"Wrap screen vertically: on, off, error")
	clockSpeed := flag.Int("clock-speed", 1300,
		"Approximate cycle speed in Hz (default: 1300)")
	timerSpeed := flag.Int("timer-speed", 60,
		"Approximate timer speed in Hz (default: 60)")
	screenBuffer := flag.Int("screen-buffer", 1,
		"Number of frames to merge for output to prevent flickering (default: 1)")
	scalingFactor := flag.Int("scaling-factor", 8,
		"Scaling factor for pixels (sets screen size) (default: 8)")
	fgColour := flag.String("fg", "0xFFFFFFFF",
		"Colour for foreground (active pixels) as hexadecimal string (default: 0xFFFFFFFF)")
	bgColour := flag.String("bg", "0x00000000",
		"Colour for background (active pixels) as hexadecimal string (default: 0x00000000)")
	debug := flag.Bool("debug", false, "Produce output for debugging")
	iniflags.Parse()

	filename := flag.Arg(0)
	rombytes := readROM(filename)

	fg, err := strconv.ParseUint(*fgColour, 0, 32)
	check(err)
	bg, err := strconv.ParseUint(*bgColour, 0, 32)
	check(err)

	vm := VM{}
	if *debug {
		PrintROM(rombytes)
	}
	vm.init(rombytes, *wrapX, *wrapY, *clockSpeed, *timerSpeed, *screenBuffer)
	// SDL init
	display := SDLDisplay{}
	display.init(int32(*scalingFactor), uint32(bg), uint32(fg))

	keyboard := SDLKeyboard{}
	keyboard.generateKeymaps()

	vm.loop(&display, &keyboard)

	display.Destroy()
	sdl.Quit()

	if *debug {
		vm.printState()
	}
}
