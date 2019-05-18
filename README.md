# chip8go

A simple CHIP-8 interpreter/emulator in Go

![Tetris with default colours](screens/tetris.png)

![Breakout with different colours](screens/breakout.png)

## Installation

### Install dependencies

```bash
go get github.com/veandco/go-sdl2/sdl
go get gopkg.in/ini.v1
go get github.com/vharitonsky/iniflags
```

### Build

```bash
go build -o chip8go src/*
```

## Usage

Run a rom:

```bash
./chip8go ./path/to/rom.ch8
```

### Command-line Arguments

Full list of options:

```
  -allowMissingConfig
    	Don't terminate the app if the ini file cannot be read.
  -allowUnknownFlags
    	Don't terminate the app if ini file contains unknown flags.
  -bg string
    	Colour for background (active pixels) as hexadecimal string (default: 0x00000000)
  -clock-speed int
    	Approximate cycle speed in Hz (default: 1300)
  -config string
    	Path to ini config for using in go flags. May be relative to the current executable path.
  -configUpdateInterval duration
    	Update interval for re-reading config file set via -config flag. Zero disables config file re-reading.
  -debug
    	Produce output for debugging (default: False)
  -dumpflags
    	Dumps values for all flags defined in the app into stdout in ini-compatible syntax and terminates the app.
  -fg string
    	Colour for foreground (active pixels) as hexadecimal string (default: 0xFFFFFFFF)
  -scaling-factor int
    	Scaling factor for pixels (sets screen size) (default: 8)
  -screen-buffer int
    	Number of frames to merge for output to prevent flickering (default: 1)
  -timer-speed int
    	Approximate timer speed in Hz (default: 60)
  -wrapX string
    	Wrap screen horizontally: on, off, error (default "on")
  -wrapY string
    	Wrap screen vertically: on, off, error (default "on")
```

Note you can use -config to pass the above arguments in a .ini.

#### Key mapping

The key mapping can be set in keys.ini, the default mapping is:

```
1 = 1
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
```

Where pause and quit are special emulator keys.

The [SDL names for the keys](https://wiki.libsdl.org/SDL_Keycode) should be used for assignment, these usually correspond to the normal key label.

# CHIP-8 Description

[CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) is an interpreted programming language, created for home computers in the 1970s, so that one ROM could be played on many different systems using the CHIP-8 interpreter and virtual machine.

Note that the opcodes are for the virtual machine, and were never executed as native machine code.

## Virtual Machine specifications

### Clock rate

There is no specific clock rate for the CHIP-8, and the best rate may vary depending on the ROM.

In practice, a value around 700Hz will usually perform well.

#### Timers

The CHIP-8 has two timers, the delay timer and the sound timer.

The delay timer when set to a value greater than 0, will decrement by 1 at a rate of 60Hz until reaching 0.

The sound timer when set to a value greater than 0, will decrement by 1 at a rate of 60Hz until reaching 0, emitting a continuous tone while it is above 0.

The timers should always tick at 60Hz.

In practice these are implemented as unsigned 8-bit integers.

### Memory
The virtual machine has 4096 bits of memory (i.e. 0x000 to 0xFFF inclusive).

It has 16 8-bit registers, called the __V__ registers. The last of these, V[0xF] is used to hold special values by some of the opcodes, such as the carry flag for addition or the collision flag for draw operations. 

It has one 16-bit register called the __I__ register, mainly used for pointing to the 12-bit memory addresses.

There is also a stack of up to 16 16-bit values and a stack pointer to track the position on the stack. This is used to store return addresses for function calls (note there is no stack frame aside from the return address, since there are no static or local variables, etc.).

The ROM is placed in memory starting at 0x200 (to a maximum of 0xFFF). 

#### Font

The built-in font, providing sprites for the characters 0-F should be stored from 0x050 to 0x09F inclusive.

Each sprite is 5 bytes long, and so "occupies" 8x5 pixels on the screen (but pixels are XORed on the current screen when drawn).

The font is defined as follows:

```
+---------------------+-----------------+--------------------+-------------------+  
| Symbol  | Address   | Sprite          | Binary             | Hex               |  
+=====================+=================+====================+===================+  
| 0       | 0x050     | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 1       | 0x055     |   *             | 00100000           | 0x20              |  
|         |           |  **             | 01100000           | 0x60              |  
|         |           |   *             | 00100000           | 0x20              |  
|         |           |   *             | 00100000           | 0x20              |  
|         |           |  ***            | 01110000           | 0x70              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 2       | 0x05A     | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 3       | 0x05F     | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 4       | 0x064     | *  *            | 10010000           | 0x90              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           |    *            | 00010000           | 0x10              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 5       | 0x069     | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 6       | 0x06E     | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 7       | 0x073     | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           |   *             | 00100000           | 0x20              |  
|         |           |  *              | 01000000           | 0x40              |  
|         |           |  *              | 01000000           | 0x40              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 8       | 0x078     | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 9       | 0x07D     | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           |    *            | 00010000           | 0x10              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 10      | 0x082     | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | *  *            | 10010000           | 0x90              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 11      | 0x087     | ***             | 11100000           | 0xE0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ***             | 11100000           | 0xE0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ***             | 11100000           | 0xE0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 12      | 0x08C     | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 13      | 0x091     | ***             | 11100000           | 0xE0              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | *  *            | 10010000           | 0x90              |  
|         |           | ***             | 11100000           | 0xE0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 14      | 0x096     | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
+---------+-----------+-----------------+--------------------+-------------------+  
| 15      | 0x09B     | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | ****            | 11110000           | 0xF0              |  
|         |           | *               | 10000000           | 0x80              |  
|         |           | *               | 10000000           | 0x80              |  
+---------+-----------+-----------------+--------------------+-------------------+  
```

#### Memory Map
```
+---------------+= 0xFFF (4095) End of Chip-8 RAM
|               |
|               |
|               |
|               |
|               |
| 0x200 to 0xFFF|
|     Chip-8    |
| Program / Data|
|     Space     |
|               |
|               |
|               |
+- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
|               |
|               |
|               |
+---------------+= 0x200 (512) Start of most Chip-8 programs
| 0x000 to 0x1FF|
| Reserved for  |
|  interpreter  |  0x050 (80) to 0x09F (159) inclusive for built-in font
+---------------+= 0x000 (0) Start of Chip-8 RAM
```
### Display

The display is not part of the memory and should be handled separately.

It is monochrome and has a 64x32 resolution. This can easily be represented using 32 \* 8 8-bit integers and managed with bitwise operations.

It is indexed with (0,0) in upper left-most corner, and (63,31) in the bottom right corner.

```
(0,0)	(63,0)
(0,31)	(63,31)
```

#### Wrapping

X-wrapping should be enabled, that is part of a sprite that exceeds the screen width should wrap around, and a draw instruction to draw off-screen should normally also wrap around the screen horizontally.

Y-wrapping (vertically) is rarely used, although it is enabled by default in this emulator.

### Keyboard

The keyboard contains the characters 0-F in the following arrangement:

```
1	2	3	C
4	5	6	D
7	8	9	E
A	0	B	F
```

## Opcodes

There is no official assembly language for the CHIP-8, only the documented opcodes.

Instructions are two bytes long, with the most significant byte first.

The opcode documentation was taken from: http://devernay.free.fr/hacks/chip8/C8TECH10.HTM

### 0nnn - SYS addr
Jump to a machine code routine at nnn.

This instruction is only used on the old computers on which Chip-8 was originally implemented. It is ignored by modern interpreters.


### 00E0 - CLS
Clear the display.


### 00EE - RET
Return from a subroutine.

The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.


### 1nnn - JP addr
Jump to location nnn.

The interpreter sets the program counter to nnn.


### 2nnn - CALL addr
Call subroutine at nnn.

The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.


### 3xkk - SE Vx, byte
Skip next instruction if Vx = kk.

The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.


### 4xkk - SNE Vx, byte
Skip next instruction if Vx != kk.

The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.


### 5xy0 - SE Vx, Vy
Skip next instruction if Vx = Vy.

The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.


### 6xkk - LD Vx, byte
Set Vx = kk.

The interpreter puts the value kk into register Vx.


### 7xkk - ADD Vx, byte
Set Vx = Vx + kk.

Adds the value kk to the value of register Vx, then stores the result in Vx. 

### 8xy0 - LD Vx, Vy
Set Vx = Vy.

Stores the value of register Vy in register Vx.


### 8xy1 - OR Vx, Vy
Set Vx = Vx OR Vy.

Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values, and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0. 


### 8xy2 - AND Vx, Vy
Set Vx = Vx AND Vy.

Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0. 


### 8xy3 - XOR Vx, Vy
Set Vx = Vx XOR Vy.

Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0. 


### 8xy4 - ADD Vx, Vy
Set Vx = Vx + Vy, set VF = carry.

The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.


### 8xy5 - SUB Vx, Vy
Set Vx = Vx - Vy, set VF = NOT borrow.

If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.


### 8xy6 - SHR Vx {, Vy}
Set Vx = Vx SHR 1.

If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.


### 8xy7 - SUBN Vx, Vy
Set Vx = Vy - Vx, set VF = NOT borrow.

If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.


### 8xyE - SHL Vx {, Vy}
Set Vx = Vx SHL 1.

If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.


### 9xy0 - SNE Vx, Vy
Skip next instruction if Vx != Vy.

The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.


### Annn - LD I, addr
Set I = nnn.

The value of register I is set to nnn.


### Bnnn - JP V0, addr
Jump to location nnn + V0.

The program counter is set to nnn plus the value of V0.


### Cxkk - RND Vx, byte
Set Vx = random byte AND kk.

The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.


### Dxyn - DRW Vx, Vy, nibble
Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.

The interpreter reads n bytes from memory, starting at the address stored in I. 

These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing screen. 

If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.

If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen.


### Ex9E - SKP Vx
Skip next instruction if key with the value of Vx is pressed.

Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.


### ExA1 - SKNP Vx
Skip next instruction if key with the value of Vx is not pressed.

Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.


### Fx07 - LD Vx, DT
Set Vx = delay timer value.

The value of DT is placed into Vx.


### Fx0A - LD Vx, K
Wait for a key press, store the value of the key in Vx.

All execution stops until a key is pressed, then the value of that key is stored in Vx.


### Fx15 - LD DT, Vx
Set delay timer = Vx.

DT is set equal to the value of Vx.


### Fx18 - LD ST, Vx
Set sound timer = Vx.

ST is set equal to the value of Vx.


### Fx1E - ADD I, Vx
Set I = I + Vx.

The values of I and Vx are added, and the results are stored in I.


### Fx29 - LD F, Vx
Set I = location of sprite for digit Vx.

The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.


### Fx33 - LD B, Vx
Store BCD representation of Vx in memory locations I, I+1, and I+2.

The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.


### Fx55 - LD [I], Vx
Store registers V0 through Vx in memory starting at location I.

The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.


### Fx65 - LD Vx, [I]
Read registers V0 through Vx from memory starting at location I.

The interpreter reads values from memory starting at location I into registers V0 through Vx.

## Common development issues

If you're planning to build your own CHIP-8 emulator, here are some bugs I encountered to watch out for:

### SDL Rect() arguments

Watch out that SDL's Rect constructor takes the upper left (X, Y) co-ordinates, and then the width and height of the rectangle in pixels. 

For a while I was passing the upper left and bottom right co-ordinates (instead of the width and height), which are still valid integers and can be hard to debug.

I recommend using a linter like govet that forces you to name the parameters, that made the error clear.

### Fx29 opcode

Note that the Fx29 opcode to get the memory address for a built-in font sprite, does not get it for the value x in the opcode, but rather for the value in V[x].

Again, x and V[x] have the same type and size and so it can be difficult to debug this error. Write careful unit tests.

### Screen buffering

I implemented a screen buffer - that is, the last X screens (1 by default) are ORed together to avoid flickering in ROMs that update the screen by first clearing sprites in place (i.e. re-drawing over old position) and then drawing the new ones - this is very common in practice.

I recommend having this as an option, otherwise a lot of ROMs will suffer from flickering sprites.

### Screen wrapping

I implemented screen wrapping in both directions, even when the draw instruction is called for a starting point off-screen. Some ROMs depend on the latter behaviour, such as the 1979 version of Breakout by Carmelo Cortez. 

So make sure you can handle drawing off-screen without out-of-bounds errors.

### V[0xF] Draw collision flag

Remember that sprites can take up many rows, i.e. the 5x8 sprites for the built-in font. The V[0xF] collision flag should be set if **any** of the rows of the sprite results in a collision, not just the first or last. 

Write careful unit tests to detect this error if you accidentally reset V[0xF] at each row.

### Buggy ROMs

Some ROMs are quite buggy themselves, so be sure to test with other emulators when building your own.

For example, due to how the 1979 version of Breakout handles collisions, sometimes the ball can pass through the bricks. The same can happen in the Pong (1 player) ROM with the paddles.

Some ROMs might require a higher clock speed, or need Y-wrapping, etc.

### SDL threading

On OS X at least, it seems the SDL rendering must be done on the main thread. Due to this issue I had to disable the SDL code in the VM loop for the unit tests to run correctly.

I am not sure if this is an OS X/golang specific issue.

## Future development

If there's time, it'd be great to add the following features in the future:

* Check correct directories for config and keys files (i.e. XDG config directories on Linux, etc.)
* Package chip8go for the AUR
* Write a curses frontend so it can be run in the terminal too
* Add Super CHIP-8 support for ROMs that use the additional opcodes and higher resolution
* Refactor font code to set the font in one line
* Refactor VM code to reduce complexity
* Use enums for the command-line option types (not strings)
* Fix SDL pixel format - to use a monochrome multiplexed format rather than drawing to an RGB surface
* Write an assembler and add pseudo-instructions for common operations- JEQ, JNE, etc.
* Write a sprite creator to easily generate the hex for sprites
* Write a disassembler to convert ROMs to created assembly language and try to annotate common logic (loops, etc.)
* Write a working ROM using the assembler
* Add network play (with shared controls and screen) to play ROMs that support two players
* Write a debugger/cheat option to be able to monitor and edit the state of the virtual machine in play.
