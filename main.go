package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func enableRawMode() (*term.State, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return oldState, nil
}

func restoreMode(oldState *term.State) {
	_ = term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Println("\nExiting...")
	os.Exit(0)
}

type CursorPosition struct {
	xCorr   int
	yCorr   int
	absPoss int
}

type Editor struct {
	pt             *PieceTable
	cursorPosition *CursorPosition
}

func NewEditor() *Editor {
	return &Editor{
		pt: NewPieceTable(""),
		cursorPosition: &CursorPosition{
			xCorr:   0,
			yCorr:   0,
			absPoss: 0,
		},
	}

}

func (e *Editor) readKey() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	return char
}
func (e *Editor) draw() {
	fmt.Print("\x1b[2J") // Clear screen
	fmt.Print("\x1b[H")  // Move cursor home
	text := e.pt.GetText()
	lines := bytes.Split([]byte(text), []byte("\n"))
	for i := 0; i < len(lines); i++ {
		if i == 0 {
			fmt.Print("> ")
		} else {
			fmt.Print("  ")
		}
		fmt.Println(string(lines[i]))
	}
	fmt.Printf("\x1b[%d;%dH", e.cursorPosition.yCorr+1, e.cursorPosition.xCorr+2)

}
func (e *Editor) insertChar(ch string) {
	e.pt.InsertText(e.cursorPosition.absPoss, e.cursorPosition.yCorr, ch)
	e.cursorPosition.xCorr++
	e.cursorPosition.absPoss++

}
func (e *Editor) newLine() {

}
func (e *Editor) moveRight() {

}
func (e *Editor) moveLeft() {

}
func (e *Editor) backspace() {

}
func readKey() rune {
	var buf [3]byte
	n, err := syscall.Read(0, buf[:])
	if err != nil {
		return 0
	}

	if n == 1 {
		return rune(buf[0])
	}

	if n == 3 && buf[0] == 0x1b && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			return '↑'
		case 'B':
			return '↓'
		case 'C':
			return '→'
		case 'D':
			return '←'
		}
	}

	return 0
}
func main() {
	oldState, err := enableRawMode()

	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer restoreMode(oldState)

	fmt.Println("Simple Text Editor (Press ESC twice to Exit)")
	editor := NewEditor()
	escPressed := false
	needsRedraw := true // Initial draw

	for {
		if needsRedraw {
			editor.draw()
			needsRedraw = false
		}

		key := readKey()
		needsRedraw = true // Assume we need to redraw by default

		switch key {
		case 0x1b: // Escape
			if escPressed {
				restoreMode(oldState)
			}
			escPressed = true
			needsRedraw = false // No visual change for ESC alone
		case '↑':

		case '↓':

		case '→':
			editor.moveRight()
		case '←':
			editor.moveLeft()
		case 0x7f: // Backspace
			editor.backspace()
		case 0x0d: // Enter
			editor.insertChar("\n")
		default:
			if key >= 32 && key <= 126 {
				editor.insertChar(string(key))
			} else {
				needsRedraw = false
			}
		}

		escPressed = false
	}
}
