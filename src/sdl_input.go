package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"gopkg.in/ini.v1"
)

type SDLKeyboard struct {
	keycodeMap       map[uint16]uint8
	scancodeMap      map[uint16]uint8
	scancodeReversed map[uint8]uint16
	specialMap       map[string]uint16
}

func (keyboard *SDLKeyboard) generateKeymaps() {
	keycfg, err := ini.Load("keys.ini")
	if err != nil {
		keycfg, err = ini.Load([]byte(`1 = 1
2 = 2
3 = 3
C = 4
4 = Q
5 = W
6 = E
D = R
7 = A
8 = S
9 = D
E = F
A = Z
0 = X
B = C
F = V
PAUSE = Space
QUIT = Escape
`))
	}
	check(err)

	// Create keycodeMap
	keycodeMap := make(map[uint16]uint8, 16)
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("1").Value()))] = 0x1
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("2").Value()))] = 0x2
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("3").Value()))] = 0x3
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("C").Value()))] = 0xC
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("4").Value()))] = 0x4
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("5").Value()))] = 0x5
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("6").Value()))] = 0x6
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("D").Value()))] = 0xD
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("7").Value()))] = 0x7
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("8").Value()))] = 0x8
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("9").Value()))] = 0x9
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("E").Value()))] = 0xE
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("A").Value()))] = 0xA
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("0").Value()))] = 0
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("B").Value()))] = 0xB
	keycodeMap[uint16(sdl.GetKeyFromName(keycfg.Section("").Key("F").Value()))] = 0xF
	keyboard.keycodeMap = keycodeMap

	// Create scancode map
	scancodeMap := make(map[uint16]uint8, 16)
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("1").Value()))] = 0x1
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("2").Value()))] = 0x2
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("3").Value()))] = 0x3
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("C").Value()))] = 0xC
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("4").Value()))] = 0x4
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("5").Value()))] = 0x5
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("6").Value()))] = 0x6
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("D").Value()))] = 0xD
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("7").Value()))] = 0x7
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("8").Value()))] = 0x8
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("9").Value()))] = 0x9
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("E").Value()))] = 0xE
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("A").Value()))] = 0xA
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("0").Value()))] = 0
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("B").Value()))] = 0xB
	scancodeMap[uint16(sdl.GetScancodeFromName(keycfg.Section("").Key("F").Value()))] = 0xF
	keyboard.scancodeMap = scancodeMap
	keyboard.scancodeReversed = reverseMap(scancodeMap)

	specialMap := make(map[string]uint16, 2)
	specialMap["QUIT"] = uint16(sdl.GetKeyFromName(keycfg.Section("").Key("QUIT").Value()))
	specialMap["PAUSE"] = uint16(sdl.GetKeyFromName(keycfg.Section("").Key("PAUSE").Value()))
	keyboard.specialMap = specialMap
}

func (keyboard *SDLKeyboard) waitForKeyPress() (uint8, bool) {
	var val uint8
	loop := true
	quit := false
	for loop {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				if t.Type == 768 {
					if val, ok := keyboard.keycodeMap[uint16(t.Keysym.Sym)]; ok {
						return val, quit
					}
				}
			}
		}
	}

	return val, quit

}

func (keyboard *SDLKeyboard) isKeyPressed(key uint8) bool {
	arr := sdl.GetKeyboardState()
	return arr[keyboard.scancodeReversed[key]] == 1
}

func (keyboard *SDLKeyboard) specialKeyPressed(paused bool) (bool, bool) {
	running := true
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyboardEvent:
			if t.Type == 768 {
				switch uint16(t.Keysym.Sym) {
				case keyboard.specialMap["PAUSE"]:
					paused = !paused
				case keyboard.specialMap["QUIT"]:
					running = false
				}
			}
		}
	}
	return paused, running
}
