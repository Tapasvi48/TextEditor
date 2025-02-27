package main

import (
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
	term.Restore(int(os.Stdin.Fd()), oldState)
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

func (e *Editor) draw() {
	// Clear screen
	fmt.Print("\033[2J\033[H")

	// Get text and manually split it correctly
	fullText := e.pt.GetText()

	// Split text by newlines, ensuring we get all lines
	bytes := fullText
	fmt.Print("> ")
	for _, b := range bytes {
		if b == 10 {
			fmt.Print("\r\n")
			fmt.Print("> ")
		} else {
			fmt.Print(string(b))
		}
	}

	// Print debug info with additional buffer info
	fmt.Print("\n") // Extra space
	fmt.Printf("DEBUG: Cursor(X:%d,Y:%d) |Cursor (%d)| Lines: %d | abs: %d\n|char at:%s",
		e.cursorPosition.xCorr,
		e.cursorPosition.yCorr,
		e.cursorPosition.absPoss,
		e.pt.GetCurrentLineLength(e.cursorPosition.yCorr),
		e.cursorPosition.absPoss,
		e.pt.CharAt(e.cursorPosition.absPoss-1))

	// Print hex representation of a few characters to debug potential issues
	// if len(fullText) > 0 {
	// 	fmt.Print("First few chars (hex): ")
	// 	for i := 0; i < min(10, len(fullText)); i++ {
	// 		fmt.Printf("%02x ", fullText[i])
	// 	}
	// 	fmt.Println()
	// }

	// Position cursor for editing
	fmt.Printf("\033[%d;%dH", e.cursorPosition.yCorr+1, e.cursorPosition.xCorr+3)
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (e *Editor) insertChar(ch string) {
	e.pt.InsertText(e.cursorPosition.absPoss, ch)
	e.cursorPosition.xCorr++
	e.cursorPosition.absPoss++
}

func (e *Editor) moveRight() {
	currLineLength := e.pt.GetCurrentLineLength(e.cursorPosition.yCorr)
	// if e.cursorPosition.xCorr == currLineLength {
	// 	if e.pt.NextLine(e.cursorPosition.yCorr, e.cursorPosition.absPoss) {
	// 		e.cursorPosition.yCorr++
	// 		e.cursorPosition.xCorr = 0
	// 	}
	// 	return

	// }
	if e.cursorPosition.xCorr == currLineLength-1 {
		if e.pt.GetLineCount() > e.cursorPosition.yCorr {
			e.cursorPosition.yCorr++
			e.cursorPosition.xCorr = 0
			return
		} else {
			return
		}

	}
	e.cursorPosition.absPoss++
	e.cursorPosition.xCorr++

}
func (e *Editor) moveUp() {

}
func (e *Editor) moveDown() {

}
func (e *Editor) moveLeft() {
	if e.cursorPosition.xCorr == 0 {
		if e.cursorPosition.yCorr > 0 {
			e.cursorPosition.xCorr = e.pt.GetCurrentLineLength(e.cursorPosition.yCorr-1) - 1
			e.cursorPosition.yCorr--
			return
		} else {
			return
		}

	}
	e.cursorPosition.xCorr--
	e.cursorPosition.absPoss--
}
func (e *Editor) newLine() {
	e.pt.InsertText(e.cursorPosition.absPoss, "\n")
	e.cursorPosition.yCorr++
	e.cursorPosition.xCorr = 0
	e.cursorPosition.absPoss++

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
	if n == 1 && buf[0] == 0x1b {
		return 0x1b
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

	needsRedraw := true // Initial draw

	for {
		if needsRedraw {
			editor.draw()
			needsRedraw = false
		}
		key := readKey()
		needsRedraw = true
		switch key {
		case 0x1b:
			restoreMode(oldState)
		case 0x0d:
			editor.newLine()
		case '↑':
			editor.moveUp()
		case '↓':
			editor.moveDown()
		case '→':
			editor.moveRight()
		case '←':
			editor.moveLeft()
		case 0x7f: // Backspace
			editor.backspace()
		default:
			if key >= 32 && key <= 126 {
				editor.insertChar(string(key))
			} else {
				needsRedraw = false
			}
		}

	}
}
