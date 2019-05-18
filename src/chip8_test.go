package main

import "testing"
import "fmt"

func returnVM(rombytes []byte) VM {
	vm := VM{}
	vm.init(rombytes, "on", "on", 1300, 60, 1)
	display := SDLDisplay{}
	display.init(int32(8), uint32(0xFFFFFFFF), uint32(0x00000000))
	keyboard := SDLKeyboard{}
	keyboard.generateKeymaps()
	vm.loop(&display, &keyboard)
	return vm

}

func Test00E0(t *testing.T) {
	rombytes := []byte{0x60, 0x08, 0xA0, 0x55, 0xD0, 0x05, 0x00, 0xE0}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x055) {
		t.Errorf("RESET screen instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x055))
	}
	if vm.V != [16]uint8{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.screen[0] != [8]uint8{0, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[1] != [8]uint8{0, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[2] != [8]uint8{0, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[3] != [8]uint8{0, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[4] != [8]uint8{0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected screen incorrect, got: %s\nfont memory: %s",
			fmt.Sprint(vm.screen),
			fmt.Sprint(vm.memory[0x050:0x200]))
	}
}

func Test00EE(t *testing.T) {
	rombytes := []byte{0xA0, 0x05,
		0x60, 0x03,
		0x22, 0x0A,
		0x61, 0x02,
		0x12, 0x0E,
		0xF0, 0x1E,
		0x00, 0xEE,
		0x00, 0x00}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x008) {
		t.Errorf("DRAW instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x008))
	}
	if vm.V != [16]uint8{3, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x210) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x210))
	}
}

func Test1nnn(t *testing.T) {
	rombytes := []byte{0x13, 0x10}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("JMP instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.pc != uint16(0x0310) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x0310))
	}
}

func Test2nnn(t *testing.T) {
	rombytes := []byte{0x23, 0xe6}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("CALL instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.pc != uint16(0x03e6) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x03e6))
	}
	if vm.sp != uint16(1) {
		t.Errorf("Expected stack pointer incorrect, got: %d, want: %d.", vm.pc, uint16(1))
	}
	if vm.stack != [16]uint16{0x200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected stack incorrect, got: %s", fmt.Sprint(vm.stack))
	}
}

func Test3xkk(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x36, 0x04}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x204+2) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x204+2))
	}
}

func Test3xkk_noskip(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x36, 0x01}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x204) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x204))
	}
}

func Test4xkk(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x46, 0x01}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x204+2) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x204+2))
	}
}

func Test4xkk_noskip(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x46, 0x04}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x204) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x204))
	}
}

func Test5xy0(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x67, 0x04, 0x56, 0x70}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0x04, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x206+2) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x206+2))
	}
}

func Test5xy0_noskip(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x67, 0x05, 0x56, 0x70}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0x05, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x206) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x206))
	}
}

func Test6xkk(t *testing.T) {
	// LOAD Vx
	rombytes := []byte{0x66, 0x04}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("LOAD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test7xkk(t *testing.T) {
	// ADD Vx
	rombytes := []byte{0x66, 0x04, 0x76, 0x02}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("ADD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x06, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test7xkk_overflow(t *testing.T) {
	// ADD Vx
	rombytes := []byte{0x66, 0x04, 0x76, 0xFF}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("ADD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x03, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy0(t *testing.T) {
	// LD Vx, Vy
	rombytes := []byte{0x60, 0x55, 0x81, 0x00}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("LD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x55, 0x55, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy1(t *testing.T) {
	// OR Vx, Vy
	rombytes := []byte{0x60, 0xFF, 0x61, 0x0F, 0x81, 0x01}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("OR instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0xFF, 0xFF, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy2(t *testing.T) {
	// AND Vx, Vy
	rombytes := []byte{0x60, 0xFF, 0x61, 0x0F, 0x80, 0x12}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("AND instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x0F, 0x0F, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy3(t *testing.T) {
	// XOR Vx, Vy
	rombytes := []byte{0x60, 0xFF, 0x61, 0x0F, 0x80, 0x13}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("XOR instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0xF0, 0x0F, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy4_nocarry(t *testing.T) {
	// ADD Vx, Vy
	rombytes := []byte{0x60, 0x02, 0x61, 0x08, 0x80, 0x14}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("ADD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x0A, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy4_carry(t *testing.T) {
	// ADD Vx, Vy
	rombytes := []byte{0x60, 0xFF, 0x61, 0x08, 0x80, 0x14}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("ADD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x07, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy5_noborrow(t *testing.T) {
	// SUB Vx, Vy
	rombytes := []byte{0x60, 0x0A, 0x61, 0x08, 0x80, 0x15}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SUB instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x02, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy5_borrow(t *testing.T) {
	// SUB Vx, Vy
	rombytes := []byte{0x60, 0x0A, 0x61, 0x0F, 0x80, 0x15}
	// TODO: Output 5 or 251 depends on overflow handling
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SUB instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0xFB, 0x0F, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy6_0(t *testing.T) {
	// SUB Vx, Vy
	rombytes := []byte{0x60, 0xF0, 0x61, 0x08, 0x80, 0x16}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SHR instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x78, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy6_1(t *testing.T) {
	// SUB Vx, Vy
	rombytes := []byte{0x60, 0x0F, 0x61, 0x08, 0x80, 0x16}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SHR instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x07, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy7_noborrow(t *testing.T) {
	// SUBN Vx, Vy
	rombytes := []byte{0x60, 0x08, 0x61, 0x0F, 0x80, 0x17}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SUBN instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x07, 0x0F, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xy7_borrow(t *testing.T) {
	// SUBN Vx, Vy
	rombytes := []byte{0x60, 0x0C, 0x61, 0x0A, 0x80, 0x17}
	// TODO: Output 2 or 254 depends on overflow handling
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SUBN instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0xFE, 0x0A, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xyE_0(t *testing.T) {
	// SHL Vx, Vy
	rombytes := []byte{0x60, 0x0F, 0x61, 0x08, 0x80, 0x1E}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SHL instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0x1E, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test8xyE_1(t *testing.T) {
	// SHL Vx, Vy
	rombytes := []byte{0x60, 0xFF, 0x61, 0x08, 0x80, 0x1E}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SHL instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0xFE, 0x08, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func Test9xy0(t *testing.T) {
	// inverse of 5xy0
	// SNE Vx Vy
	rombytes := []byte{0x66, 0x04, 0x67, 0x04, 0x96, 0x70}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SNE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0x04, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x206) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x206))
	}
}

func Test9xy0_noskip(t *testing.T) {
	rombytes := []byte{0x66, 0x04, 0x67, 0x05, 0x96, 0x70}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("SNE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0x04, 0x05, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.pc != uint16(0x206+2) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x206+2))
	}
}

func TestAnnn(t *testing.T) {
	rombytes := []byte{0xA0, 0x05}
	vm := returnVM(rombytes)

	if vm.I != uint16(5) {
		t.Errorf("Load instruction incorrect, got: %d, want: %d.", vm.I, uint16(5))
	}
	if vm.pc != uint16(0x200+2) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x200+2))
	}
}

func TestBnnn(t *testing.T) {
	// JP V0, addr
	rombytes := []byte{0x60, 0x04, 0xB2, 0x66}
	vm := returnVM(rombytes)

	if vm.I != uint16(0) {
		t.Errorf("Load instruction incorrect, got: %d, want: %d.", vm.I, uint16(0))
	}
	if vm.pc != uint16(0x26A) {
		t.Errorf("Expected program counter incorrect, got: %d, want: %d.", vm.pc, uint16(0x26A))
	}
}

func TestDxyn(t *testing.T) {
	rombytes := []byte{0x60, 0x00, 0xA0, 0x50, 0xD0, 0x05, 0x00, 0x00}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x050) {
		t.Errorf("DRAW instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x050))
	}
	if vm.V != [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.screen[0] != [8]uint8{0xF0, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[1] != [8]uint8{0x90, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[2] != [8]uint8{0x90, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[3] != [8]uint8{0x90, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[4] != [8]uint8{0xF0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected screen incorrect, got: %s\nfont memory: %s",
			fmt.Sprint(vm.screen),
			fmt.Sprint(vm.memory[0x050:0x200]))
	}
}

func TestDxyn_allchars(t *testing.T) {
	rombytes := []byte{0x60, 0x00, 0x61, 0x04, 0xA0, 0x50,
		0xD0, 0x05,
		0x63, 0x01, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x09, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x0E, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x13, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x18, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x1D, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x22, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x27, 0xF3, 0x29, 0xD1, 0x05,
		0x73, 0x01, 0x61, 0x2C, 0xF3, 0x29, 0xD1, 0x05,
		0x62, 0x06,
		0x73, 0x01, 0x61, 0x00, 0xF3, 0x29, 0xD1, 0x25,
		0x73, 0x01, 0x61, 0x05, 0xF3, 0x29, 0xD1, 0x25,
		0x73, 0x01, 0x61, 0x0A, 0xF3, 0x29, 0xD1, 0x25,
		0x73, 0x01, 0x61, 0x0F, 0xF3, 0x29, 0xD1, 0x25,
		0x73, 0x01, 0x61, 0x14, 0xF3, 0x29, 0xD1, 0x25,
		0x73, 0x01, 0x61, 0x19, 0xF3, 0x29, 0xD1, 0x25}

	vm := returnVM(rombytes)

	if vm.V != [16]uint8{0, 0x19, 0x06, 0x0F, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.screen[0] != [8]uint8{242, 123, 210, 247, 189, 239, 0, 0} ||
		vm.screen[1] != [8]uint8{150, 8, 82, 132, 5, 41, 0, 0} ||
		vm.screen[2] != [8]uint8{146, 123, 222, 247, 137, 239, 0, 0} ||
		vm.screen[3] != [8]uint8{146, 64, 66, 20, 145, 33, 0, 0} ||
		vm.screen[4] != [8]uint8{247, 123, 194, 247, 145, 239, 0, 0} ||
		vm.screen[5] != [8]uint8{0, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[6] != [8]uint8{247, 61, 207, 120, 0, 0, 0, 0} ||
		vm.screen[7] != [8]uint8{148, 161, 40, 64, 0, 0, 0, 0} ||
		vm.screen[8] != [8]uint8{247, 33, 47, 120, 0, 0, 0, 0} ||
		vm.screen[9] != [8]uint8{148, 161, 40, 64, 0, 0, 0, 0} ||
		vm.screen[10] != [8]uint8{151, 61, 207, 64, 0, 0, 0, 0} {
		t.Errorf("Expected screen incorrect, got: %s\nfont memory: %s",
			fmt.Sprint(vm.screen),
			fmt.Sprint(vm.memory[0x050:0x200]))
	}
}

func TestDxyn_overlap(t *testing.T) {
	rombytes := []byte{0x60, 0x00, 0x61, 0x04, 0x63, 0x00,
		0x61, 0x00, 0xF3, 0x29, 0xD1, 0x05, 0x73, 0x01,
		0x61, 0x00, 0xF3, 0x29, 0xD1, 0x05}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x055) {
		t.Errorf("DRAW instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x050))
	}
	if vm.V != [16]uint8{0, 0, 0, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.screen[0] != [8]uint8{208, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[1] != [8]uint8{240, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[2] != [8]uint8{176, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[3] != [8]uint8{176, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[4] != [8]uint8{128, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected screen incorrect, got: %s\nfont memory: %s",
			fmt.Sprint(vm.screen),
			fmt.Sprint(vm.memory[0x050:0x200]))
	}
}

func TestDxyn_overlap2(t *testing.T) {
	rombytes := []byte{0x60, 0x00, 0x63, 0x0F,
		0x61, 0x00, 0xF3, 0x29, 0xD1, 0x05, 0x63, 0x00,
		0x61, 0x02, 0xF3, 0x29, 0xD1, 0x05}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x050) {
		t.Errorf("DRAW instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x050))
	}
	if vm.V != [16]uint8{0, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
	if vm.screen[0] != [8]uint8{204, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[1] != [8]uint8{164, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[2] != [8]uint8{212, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[3] != [8]uint8{164, 0, 0, 0, 0, 0, 0, 0} ||
		vm.screen[4] != [8]uint8{188, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected screen incorrect, got: %s\nfont memory: %s",
			fmt.Sprint(vm.screen),
			fmt.Sprint(vm.memory[0x050:0x200]))
	}
}

func TestFx1E(t *testing.T) {
	// I = I + Vx
	rombytes := []byte{0x60, 0x02, 0xA0, 0x08, 0xF0, 0x1E}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x0A) {
		t.Errorf("ADD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x0A))
	}
	if vm.V != [16]uint8{0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func TestFx29(t *testing.T) {
	// Set I to sprite location
	rombytes := []byte{0x60, 0x03, 0xA0, 0x08, 0xF0, 0x29}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x05F) {
		t.Errorf("SPRITE instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x05F))
	}
	if vm.V != [16]uint8{0x03, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}
}

func TestFx33(t *testing.T) {
	// BCD of Vx in I, I+1, I+2
	rombytes := []byte{0x60, 0xFF, 0xA5, 0x10, 0xF0, 0x33}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x510) {
		t.Errorf("LOAD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x0510))
	}
	if vm.V != [16]uint8{0xFF, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}

	if vm.memory[0x0510] != 2 || vm.memory[0x0511] != 5 || vm.memory[0x0512] != 5 {

		t.Errorf("Expected memory incorrect, got: %s", fmt.Sprint(vm.memory[0x0510:0x0513]))
	}
}

func TestFx55(t *testing.T) {
	// Store registers V0 through Vx in memory starting at location I.
	rombytes := []byte{0x60, 0xDE, 0x61, 0xAD, 0x62, 0xBE, 0x63, 0xEF, 0xA5, 0x10, 0xF3, 0x55}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x510) {
		t.Errorf("LOAD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x0510))
	}
	if vm.V != [16]uint8{0xDE, 0xAD, 0xBE, 0xEF, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}

	if vm.memory[0x0510] != 0xDE || vm.memory[0x0511] != 0xAD || vm.memory[0x0512] != 0xBE || vm.memory[0x0513] != 0xEF {

		t.Errorf("Expected memory incorrect, got: %s", fmt.Sprint(vm.memory[0x0510:0x0514]))
	}
}

func TestFx65(t *testing.T) {
	// Store registers V0 through Vx in memory starting at location I.
	rombytes := []byte{0xF0, 0x29, 0xF4, 0x65}
	vm := returnVM(rombytes)

	if vm.I != uint16(0x050) {
		t.Errorf("LOAD instruction incorrect, got: %d, want: %d.", vm.I, uint16(0x050))
	}
	if vm.V != [16]uint8{0xF0, 0x90, 0x90, 0x90, 0xF0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Expected V incorrect, got: %s", fmt.Sprint(vm.V))
	}

}
