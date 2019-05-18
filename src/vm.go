package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Display interface {
	clearDisplay()
	updateDisplay()
	drawPixel(x int32, y int32)
}

type Keyboard interface {
	waitForKeyPress() (uint8, bool)
	isKeyPressed(key uint8) bool                // argument is 0-F key value
	specialKeyPressed(paused bool) (bool, bool) // paused, running
}

// VM : Class for virtual machine - holds all memory and registers
type VM struct {
	romlength              uint16
	pc, I, opcode, sp      uint16
	V                      [16]uint8
	memory                 [4096]uint8
	screen                 [32][8]uint8 // bitmap 64x32
	delayTimer, soundTimer uint8
	stack                  [16]uint16
	drawflag               bool
	wrapX                  string
	wrapY                  string
	clockSpeed             uint16
	timerSpeed             uint16
	screenBuffer           uint8
}

func (vm *VM) printState() {
	fmt.Printf("PC: 0x%x\n", vm.pc)
	fmt.Printf("I: 0x%x\n", vm.I)
	fmt.Printf("Opcode: 0x%x\n", vm.opcode)
	fmt.Println("Memory:")
	fmt.Println(vm.memory)
	fmt.Println("V:")
	fmt.Println(vm.V)
	fmt.Println("Screen:")
	fmt.Println(vm.screen)
	fmt.Printf("DT: %d\n", vm.delayTimer)
	fmt.Printf("ST: %d\n", vm.soundTimer)
	fmt.Printf("SP: %d\n", vm.sp)
	fmt.Println("Stack:")
	fmt.Println(vm.stack)
}

func (vm *VM) initialiseFont() {
	//0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	//0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	//0 0x050
	vm.memory[0x050] = 0xF0
	vm.memory[0x051] = 0x90
	vm.memory[0x052] = 0x90
	vm.memory[0x053] = 0x90
	vm.memory[0x054] = 0xF0
	// 1
	vm.memory[0x055] = 0x20
	vm.memory[0x056] = 0x60
	vm.memory[0x057] = 0x20
	vm.memory[0x058] = 0x20
	vm.memory[0x059] = 0x70
	// 2
	vm.memory[0x05A] = 0xF0
	vm.memory[0x05B] = 0x10
	vm.memory[0x05C] = 0xF0
	vm.memory[0x05D] = 0x80
	vm.memory[0x05E] = 0xF0
	// 3
	vm.memory[0x05F] = 0xF0
	vm.memory[0x060] = 0x10
	vm.memory[0x061] = 0xF0
	vm.memory[0x062] = 0x10
	vm.memory[0x063] = 0xF0
	// 4
	vm.memory[0x064] = 0x90
	vm.memory[0x065] = 0x90
	vm.memory[0x066] = 0xF0
	vm.memory[0x067] = 0x10
	vm.memory[0x068] = 0x10
	// 5
	vm.memory[0x069] = 0xF0
	vm.memory[0x06A] = 0x80
	vm.memory[0x06B] = 0xF0
	vm.memory[0x06C] = 0x10
	vm.memory[0x06D] = 0xF0
	// 6
	vm.memory[0x06E] = 0xF0
	vm.memory[0x06F] = 0x80
	vm.memory[0x070] = 0xF0
	vm.memory[0x071] = 0x90
	vm.memory[0x072] = 0xF0
	// 7
	vm.memory[0x073] = 0xF0
	vm.memory[0x074] = 0x10
	vm.memory[0x075] = 0x20
	vm.memory[0x076] = 0x40
	vm.memory[0x077] = 0x40
	// 8
	vm.memory[0x078] = 0xF0
	vm.memory[0x079] = 0x90
	vm.memory[0x07A] = 0xF0
	vm.memory[0x07B] = 0x90
	vm.memory[0x07C] = 0xF0
	// 9
	vm.memory[0x07D] = 0xF0
	vm.memory[0x07E] = 0x90
	vm.memory[0x07F] = 0xF0
	vm.memory[0x080] = 0x10
	vm.memory[0x081] = 0xF0
	// A
	vm.memory[0x082] = 0xF0
	vm.memory[0x083] = 0x90
	vm.memory[0x084] = 0xF0
	vm.memory[0x085] = 0x90
	vm.memory[0x086] = 0x90
	// B
	vm.memory[0x087] = 0xE0
	vm.memory[0x088] = 0x90
	vm.memory[0x089] = 0xE0
	vm.memory[0x08A] = 0x90
	vm.memory[0x08B] = 0xE0
	// C
	vm.memory[0x08C] = 0xF0
	vm.memory[0x08D] = 0x80
	vm.memory[0x08E] = 0x80
	vm.memory[0x08F] = 0x80
	vm.memory[0x090] = 0xF0
	// D
	vm.memory[0x091] = 0xE0
	vm.memory[0x092] = 0x90
	vm.memory[0x093] = 0x90
	vm.memory[0x094] = 0x90
	vm.memory[0x095] = 0xE0
	// E
	vm.memory[0x096] = 0xF0
	vm.memory[0x097] = 0x80
	vm.memory[0x098] = 0xF0
	vm.memory[0x099] = 0x80
	vm.memory[0x09A] = 0xF0
	// F
	vm.memory[0x09B] = 0xF0
	vm.memory[0x09C] = 0x80
	vm.memory[0x09D] = 0xF0
	vm.memory[0x09E] = 0x80
	vm.memory[0x09F] = 0x80
}

func (vm *VM) loadROM(rombytes []byte) {
	vm.romlength = uint16(len(rombytes))
	for i, byt := range rombytes {
		vm.memory[0x200+i] = byt
	}
}

func (vm *VM) init(rombytes []byte, wrapX string, wrapY string, clockSpeed int, timerSpeed int, screenBuffer int) {
	vm.initialiseFont()
	vm.loadROM(rombytes)
	vm.pc = 0x200
	vm.drawflag = false
	vm.wrapX = wrapX
	vm.wrapY = wrapY
	vm.clockSpeed = uint16(clockSpeed)
	vm.timerSpeed = uint16(timerSpeed)
	vm.screenBuffer = uint8(screenBuffer)

}

func (vm *VM) parseOpcode(keyboard Keyboard) bool {
	var running bool
	vm.opcode = uint16(vm.memory[vm.pc])<<8 | uint16(vm.memory[vm.pc+1]) // big-endian
	vm.drawflag = false
	switch vm.opcode & 0xF000 {
	case 0x0000:
		switch vm.opcode & 0x00FF {
		case 0x00E0:
			// clear vm.screen
			// fmt.Printf("Clear vm.screen % x, % d\n", vm.opcode, vm.pc)
			for yp := 0; yp < 32; yp++ {
				for xb := 0; xb < 8; xb++ {
					vm.screen[yp][xb] = 0
				}
			}
			vm.drawflag = true
			vm.pc += 2
		case 0x00EE:
			// 00EE - RET
			// Return from a subroutine.
			if vm.sp <= 0 {
				log.Fatal(fmt.Errorf("vm.stack pointer below 0"))
			}
			// fmt.Printf("RET pc: %x, new pc: %x, opcode: %x\n", vm.pc, vm.stack[vm.sp]+2, vm.opcode)
			vm.sp--
			vm.pc = vm.stack[vm.sp] + 2
			// fmt.Printf("RET pc: %d\n", vm.pc)
		default:
			// fmt.Printf("SYS vm.opcode ignored: % x, % d\n", vm.opcode, vm.pc)
			vm.pc += 2
		}
	case 0x1000:
		// 1nnn - JP addr
		// Jump to location nnn.
		vm.pc = 0x0FFF & vm.opcode
		if vm.pc < 0x200 || vm.pc > 0xFFF {
			log.Fatal(fmt.Errorf("illegal JMP instruction - PC: %x, opcode: %x", vm.pc, vm.opcode))
		}
		// fmt.Printf("JMP to : % x, % x\n", vm.pc, vm.opcode)
		// note endless jumps used as halt

	case 0x2000:
		// 2nnn - CALL addr
		// Call subroutine at nnn.
		vm.stack[vm.sp] = vm.pc
		vm.sp++
		// fmt.Printf("CALL pc: %x, new pc: %x, opcode: %x\n", vm.pc, (0x0FFF & vm.opcode), vm.opcode)
		vm.pc = 0x0FFF & vm.opcode
		if vm.pc < 0x200 || vm.pc > 0xFFF {
			log.Fatal(fmt.Errorf("illegal JMP instruction - PC: %x, opcode: %x", vm.pc, vm.opcode))
		}
		// fmt.Printf("CALL stack: %s\n", fmt.Sprint(vm.stack))

	case 0x3000:
		// 3xkk - SE vm.Vx, byte
		// Skip next instruction if vm.Vx = kk.
		if vm.V[0x0F00&vm.opcode>>8] == uint8(0x00FF&vm.opcode) {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x4000:
		// 4xkk - SNE vm.Vx, byte
		// Skip next instruction if vm.Vx != kk.
		if vm.V[0x0F00&vm.opcode>>8] != uint8(0x00FF&vm.opcode) {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x5000:
		// 5xy0 - SE vm.Vx, vm.Vy
		// Skip next instruction if vm.Vx = vm.Vy.
		if vm.V[0x0F00&vm.opcode>>8] == vm.V[uint8(0x00F0&vm.opcode)>>4] {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0x6000:
		// 6xkk - LD vm.Vx, byte
		// Set vm.Vx = kk.
		vm.V[0x0F00&vm.opcode>>8] = uint8(0x00FF & vm.opcode)
		vm.pc += 2

	case 0x7000:
		// 7xkk - ADD vm.Vx, byte
		// Set vm.Vx = vm.Vx + kk.
		vm.V[0x0F00&vm.opcode>>8] += uint8(0x0FF & vm.opcode)
		vm.pc += 2

	case 0x8000:
		switch vm.opcode & 0x000F {
		case 0x0000:
			//8xy0 - LD vm.Vx, vm.Vy
			//Set vm.Vx = vm.Vy.
			vm.V[vm.opcode&0x0F00>>8] = vm.V[vm.opcode&0x00F0>>4]
			vm.pc += 2
		case 0x0001:
			// 8xy1 - OR vm.Vx, vm.Vy
			// Set vm.Vx = vm.Vx OR vm.Vy.
			vm.V[vm.opcode&0x0F00>>8] |= vm.V[vm.opcode&0x00F0>>4]
			vm.pc += 2
		case 0x0002:
			// 8xy2 - AND vm.Vx, vm.Vy
			// Set vm.Vx = vm.Vx AND vm.Vy.
			vm.V[vm.opcode&0x0F00>>8] &= vm.V[vm.opcode&0x00F0>>4]
			vm.pc += 2
		case 0x0003:
			// 8xy3 - XOR vm.Vx, vm.Vy
			// Set vm.Vx = vm.Vx XOR vm.Vy.
			vm.V[vm.opcode&0x0F00>>8] ^= vm.V[vm.opcode&0x00F0>>4]
			vm.pc += 2
		case 0x0004:
			// 8xy4 - ADD vm.Vx, vm.Vy
			// Set vm.Vx = vm.Vx + vm.Vy, set vm.VF = carry.
			tempvar := uint16(vm.V[vm.opcode&0x0F00>>8]) + uint16(vm.V[vm.opcode&0x00F0>>4])
			vm.V[vm.opcode&0x0F00>>8] = uint8(0x00FF & tempvar)
			if 0xFF00&tempvar != 0 {
				vm.V[0xF] = 1
			} else {
				vm.V[0xF] = 0
			}
			vm.pc += 2
		case 0x0005:
			// 8xy5 - SUB vm.Vx, vm.Vy
			// Set vm.Vx = vm.Vx - vm.Vy, set vm.VF = NOT borrow.
			if vm.V[vm.opcode&0x0F00>>8] > vm.V[vm.opcode&0x00F0>>4] {
				vm.V[0xF] = 1
			} else {
				vm.V[0xF] = 0
			}
			vm.pc += 2
			vm.V[vm.opcode&0x0F00>>8] -= vm.V[vm.opcode&0x00F0>>4]
		case 0x0006:
			// 8xy6 - SHR vm.Vx {, vm.Vy}
			// Set vm.Vx = vm.Vx SHR 1.
			if vm.V[vm.opcode&0x0F00>>8]&1 == 1 {
				vm.V[0xF] = 1
			} else {
				vm.V[0xF] = 0
			}
			vm.V[vm.opcode&0x0F00>>8] >>= 1
			vm.pc += 2
		case 0x0007:
			// 8xy7 - SUBN vm.Vx, vm.Vy
			// Set vm.Vx = vm.Vy - vm.Vx, set vm.VF = NOT borrow.
			if vm.V[vm.opcode&0x00F0>>4] > vm.V[vm.opcode&0x0F00>>8] {
				vm.V[0xF] = 1
			} else {
				vm.V[0xF] = 0
			}
			vm.pc += 2
			vm.V[vm.opcode&0x0F00>>8] = vm.V[vm.opcode&0x00F0>>4] - vm.V[vm.opcode&0x0F00>>8]
		case 0x000E:
			// 8xyE - SHL vm.Vx {, vm.Vy}
			// Set vm.Vx = vm.Vx SHL 1.
			if vm.V[vm.opcode&0x0F00>>8]&128 == 128 {
				vm.V[0xF] = 1
			} else {
				vm.V[0xF] = 0
			}
			vm.V[vm.opcode&0x0F00>>8] <<= 1
			vm.pc += 2
		default:
			fmt.Printf("Bad vm.opcode: % x\n", vm.opcode)
		}

	case 0x9000:
		// 9xy0 - SNE vm.Vx, vm.Vy
		// Skip next instruction if vm.Vx != vm.Vy.
		if vm.V[0x0F00&vm.opcode>>8] != vm.V[0x00F0&vm.opcode>>4] {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0xA000:
		// Annn - LD vm.I, addr
		// Set vm.I = nnn.
		vm.I = vm.opcode & 0x0FFF
		vm.pc += 2
		// fmt.Printf("% x\n", vm.I)

	case 0xB000:
		// Bnnn - JP vm.V0, addr
		// Jump to location nnn + vm.V0.
		vm.pc = 0x0FFF&vm.opcode + uint16(vm.V[0])

	case 0xC000:
		// Cxkk - RND vm.Vx, byte
		// Set vm.Vx = random byte AND kk.
		vm.V[0x0F00&vm.opcode>>8] = uint8(rand.Intn(256)) & uint8(0x00FF&vm.opcode)
		vm.pc += 2

	case 0xD000:
		// Dxyn - DRW vm.Vx, vm.Vy, nibble
		// Display n-byte sprite starting at vm.memory location vm.I at (vm.Vx, vm.Vy), set vm.VF = collision.
		// fmt.Printf("Draw sprite + collide: % x, %d\n", vm.opcode, vm.pc)

		vm.drawflag = true
		n := 0x000F & vm.opcode
		x := vm.V[0x0F00&vm.opcode>>8]
		y := vm.V[0x00F0&vm.opcode>>4]
		vm.V[0xF] = 0
		if x > 63 {
			switch vm.wrapX {
			case "on":
				x = 0 + x%64
			case "off":
				vm.drawflag = false
			case "error":
				log.Fatal(fmt.Errorf("illegal DRAW instruction X - PC: 0x%x, opcode: 0x%x",
					vm.pc, vm.opcode))
			}
		}
		if y > 31 {
			switch vm.wrapY {
			case "on":
				y = 0 + y%32
			case "off":
				vm.drawflag = false
			case "error":
				log.Fatal(fmt.Errorf("illegal DRAW instruction Y - PC: 0x%x, opcode: 0x%x",
					vm.pc, vm.opcode))
			}
		}
		vm.pc += 2
		if vm.drawflag {
			for i := uint16(0); i < n; i++ {
				sprite := vm.memory[vm.I+i]
				if x%8 == 0 {
					if vm.screen[y][x/8]&sprite > 0 {
						vm.V[0xF] = 1
					}
					vm.screen[y][x/8] = sprite ^ vm.screen[y][x/8]

				} else {
					// first part
					s := sprite >> (x % 8)
					if vm.screen[y][x/8]&s > 0 {
						vm.V[0xF] = 1
					}
					vm.screen[y][x/8] = s ^ vm.screen[y][x/8]
					// second part - handle wrap

					s = sprite & uint8(math.Pow(2, float64(x%8))-1) << (8 - x%8)
					newx := x/8 + 1
					if newx > 7 {
						newx = 0

					}
					if vm.wrapX == "on" || newx > 0 {
						if vm.screen[y][newx]&s > 0 {
							vm.V[0xF] = 1
						}
						vm.screen[y][newx] = s ^ vm.screen[y][newx]
					}

				}
				y++
				if y > 31 {
					if vm.wrapY == "on" {
						y = 0
					} else {
						break
					}
				}
			}
		}
		// The interpreter reads n bytes from vm.memory, starting at the address stored in vm.I.
		// These bytes are then displayed as sprites on vm.screen at coordinates (vm.Vx, vm.Vy).
		// Sprites are XORed onto the existing vm.screen.
		// If this causes any pixels to be erased, vm.VF is set to 1, otherwise it is set to 0.
		// If the sprite is positioned so part of it is outside the coordinates of the display,
		// it wraps around to the opposite side of the vm.screen.

	case 0xE000:
		switch vm.opcode & 0x00FF {
		case 0x009E:
			// Ex9E - SKP vm.Vx
			// Skip next instruction if key with the value of vm.Vx is pressed.
			if keyboard.isKeyPressed(vm.V[0x0F00&vm.opcode>>8]) {
				vm.pc += 4
			} else {
				vm.pc += 2
			}

		case 0x00A1:
			// ExA1 - SKNP vm.Vx
			// Skip next instruction if key with the value of vm.Vx is not pressed.
			if !keyboard.isKeyPressed(vm.V[0x0F00&vm.opcode>>8]) {
				vm.pc += 4
			} else {
				vm.pc += 2
			}

		default:
			fmt.Printf("Bad vm.opcode: % x\n", vm.opcode)
		}

	case 0xF000:
		switch vm.opcode & 0x00FF {
		case 0x0007:
			// Fx07 - LD vm.Vx, DT
			// Set vm.Vx = delay timer value.
			vm.V[0x0F00&vm.opcode>>8] = vm.delayTimer
			vm.pc += 2

		case 0x000A:
			// Fx0A - LD vm.Vx, K
			// Wait for a key press, store the value of the key in vm.Vx.
			vm.V[0x0F00&vm.opcode>>8], running = keyboard.waitForKeyPress()
			if !running {
				return false
			}
			vm.pc += 2

		case 0x0015:
			// Fx15 - LD DT, vm.Vx
			// Set delay timer = vm.Vx.
			vm.delayTimer = vm.V[0x0F00&vm.opcode>>8]
			vm.pc += 2

		case 0x0018:
			// Fx18 - LD ST, vm.Vx
			// Set sound timer = vm.Vx.
			vm.soundTimer = vm.V[0x0F00&vm.opcode>>8]
			vm.pc += 2

		case 0x001E:
			// Fx1E - ADD vm.I, vm.Vx
			// Set vm.I = vm.I + vm.Vx.
			vm.I += uint16(vm.V[0x0F00&vm.opcode>>8])
			vm.pc += 2

		case 0x0029:
			// Fx29 - LD F, vm.Vx
			// Set vm.I = location of sprite for digit vm.Vx.
			// The value of vm.I is set to the location for the hexadecimal sprite corresponding to the value of vm.Vx.

			vm.I = 0x050 + 5*uint16(vm.V[0x0F00&vm.opcode>>8])
			vm.pc += 2

		case 0x0033:
			// Fx33 - LD B, vm.Vx
			// Store BCD representation of vm.Vx in vm.memory locations vm.I, vm.I+1, and vm.I+2.
			// The interpreter takes the decimal value of vm.Vx,
			// and places the hundreds digit in vm.memory at location in vm.I,
			// the tens digit at location vm.I+1,
			// and the ones digit at location vm.I+2.
			vm.memory[vm.I] = vm.V[0x0F00&vm.opcode>>8] / 100
			vm.memory[vm.I+1] = (vm.V[0x0F00&vm.opcode>>8] - vm.memory[vm.I]*100) / 10
			vm.memory[vm.I+2] = vm.V[0x0F00&vm.opcode>>8] - vm.memory[vm.I]*100 - vm.memory[vm.I+1]*10

			vm.pc += 2

		case 0x0055:
			// Fx55 - LD [vm.I], vm.Vx
			// Store registers vm.V0 through vm.Vx in vm.memory starting at location vm.I.
			var i uint16
			for i = 0; i <= 0x0F00&vm.opcode>>8; i++ {
				vm.memory[vm.I+i] = vm.V[i]
			}
			vm.pc += 2

		case 0x0065:
			// Fx65 - LD vm.Vx, [vm.I]
			// Read registers vm.V0 through vm.Vx from vm.memory starting at location vm.I.
			var i uint16
			for i = 0; i <= 0x0F00&vm.opcode>>8; i++ {
				vm.V[i] = vm.memory[vm.I+i]
			}
			vm.pc += 2

		default:
			fmt.Printf("Bad vm.opcode: % x\n", vm.opcode)
		}

	default:
		fmt.Printf("Bad vm.opcode: % x\n", vm.opcode)
	}
	return true

}

func (vm *VM) loop(display Display, keyboard Keyboard) {
	var timecount uint8
	var running = true
	var andscreen [32][8]uint8 // bitmap 64x32

	paused := false
	bell := []byte{7}
	screenarray := make([][32][8]uint8, vm.screenBuffer)

	// main loop
	for running {
		time.Sleep(time.Duration(1E6/uint32(vm.clockSpeed)) * time.Microsecond)
		running = vm.parseOpcode(keyboard)

		// Do not run SDL code in test
		if !strings.HasSuffix(os.Args[0], ".test") {
		input:
			paused, running = keyboard.specialKeyPressed(paused)

			if paused && running {
				goto input
			}

			// display tick
			if vm.drawflag {
				display.clearDisplay()
				for yp := 0; yp < 32; yp++ {
					for xb := 0; xb < 8; xb++ {
						for z := uint8(0); z < vm.screenBuffer; z++ {
							andscreen[yp][xb] = screenarray[z][yp][xb] | vm.screen[yp][xb]
						}
						for xp := 0; xp < 8; xp++ {
							if andscreen[yp][xb]&uint8(math.Pow(2, float64(7-xp)))>>uint8(7-xp) == 1 {
								display.drawPixel(int32(8)*int32(xb)+int32(xp), int32(yp))
							}
						}
					}
				}
				display.updateDisplay()
				screenarray[0] = vm.screen
				for z := vm.screenBuffer - 1; z > 0; z-- {
					screenarray[z] = screenarray[z-1]
				}

			} // drawflag end
		} // unit test ignore end

		if timecount >= uint8(vm.clockSpeed/vm.timerSpeed) { //timer start
			timecount = 0
			if vm.delayTimer > 0 {
				vm.delayTimer--
			}
			if vm.soundTimer > 0 {
				os.Stdout.Write(bell) // TODO: continuous tone
				vm.soundTimer--
			}
		} // timer end
		timecount++
		if vm.pc-0x200 >= vm.romlength {
			running = false
		}
	}
}
