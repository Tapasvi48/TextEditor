package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/term"
)

type CursorPosition struct {
	xCorr int
	yCorr int
}

type Editor struct {
	textContent    [][]rune
	cursorPosition *CursorPosition
}

func NewEditor() *Editor {
	return &Editor{
		textContent: [][]rune{}, // Start with one empty line
		cursorPosition: &CursorPosition{
			xCorr: 0,
			yCorr: 0,
		},
	}
}

func (editor *Editor) InsertChar(char rune) {
	y := editor.cursorPosition.yCorr
	if len(editor.textContent) <= y {
		editor.textContent = append(editor.textContent, []rune{})
	}
	editor.textContent[y] = append(editor.textContent[y], char)
	// Move cursor forward
	editor.cursorPosition.xCorr++

}

func (editor *Editor) NewLine() {
	editor.cursorPosition.yCorr++
	editor.cursorPosition.xCorr = 0
	editor.textContent = append(editor.textContent, []rune{})

}

func (editor *Editor) PrevLine() {
	//merge this line with prev line
	//in case of no line
	y := editor.cursorPosition.yCorr
	prevY := y - 1
	prevLen := len(editor.textContent[prevY])
	if len(editor.textContent[y]) == 0 {
		editor.cursorPosition.yCorr = prevY
		editor.cursorPosition.xCorr = prevLen
		return
	}

	editor.textContent[prevY] = append(editor.textContent[prevY], editor.textContent[y]...)
	editor.textContent = append(editor.textContent[:y], editor.textContent[y+1:]...)
	editor.cursorPosition.yCorr = prevY
	editor.cursorPosition.xCorr = prevLen

}

func (editor *Editor) Backspace() {

	y := editor.cursorPosition.yCorr
	x := editor.cursorPosition.xCorr
	if x == 0 && y > 0 {
		editor.PrevLine()
		return
	}

	if y < 0 || y >= len(editor.textContent) || x <= 0 || x > len(editor.textContent[y]) {
		return
	}
	editor.textContent[y] = append(editor.textContent[y][:x-1], editor.textContent[y][x:]...)
	editor.cursorPosition.xCorr--

}

func (editor *Editor) PrintScreen() {
	cmd := exec.Command("clear")

	cmd.Stdout = os.Stdout
	cmd.Run()

	for _, line := range editor.textContent {
		fmt.Print(string(line) + "\r\n")
	}

	fmt.Printf("\033[%d;%dH", editor.cursorPosition.yCorr+1, editor.cursorPosition.xCorr+1)

	// Move cursor to the correct position

}

func readKey() rune {
	var buf [1]byte
	syscall.Syscall(syscall.SYS_READ, 0, uintptr(unsafe.Pointer(&buf)), 1)

	return rune(buf[0])
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	editor := NewEditor()

	fmt.Println("Simple Text Editor (Press ESC to Exit)")

	for {
		editor.PrintScreen()
		key := readKey()

		switch key {
		case 27: // ESC to exit
			fmt.Println("\nExiting...")
			return
		case '\n', '\r': // Enter key
			editor.NewLine()
		case 8, 127:
			editor.Backspace()
		default:
			editor.InsertChar(key)
		}
	}
}
