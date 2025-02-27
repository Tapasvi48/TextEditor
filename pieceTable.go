package main

type Node struct {
	bufferIndex int
	startIndex  int
	length      int
	lineStart   []int
	//will contain bufferIndex of line end maybe
}

type PieceTable struct {
	buffers   []string
	nodes     []Node
	nodeIndex []int
}

func NewPieceTable(initialText string) *PieceTable {
	return &PieceTable{
		buffers: []string{initialText},
		nodes:   []Node{},
	}

}

// return the string at line n
func (pt *PieceTable) GoToLine(n int) (bufferIndex int, startOffset int, endOffset int) {
	lineCount := 0

	for _, node := range pt.nodes {
		numLines := len(node.lineStart)

		if lineCount+numLines > n {
			startOffset := node.startIndex
			if n-lineCount > 0 {
				startOffset += node.lineStart[n-lineCount]
			}

			endOffset := node.startIndex + node.length
			if n-lineCount+1 < numLines {
				endOffset = node.startIndex + node.lineStart[n-lineCount+1]
			}

			return node.bufferIndex, startOffset, endOffset
		}

		lineCount += numLines
	}

	return -1, -1, -1
}

func (p *PieceTable) GetText() []byte {
	var result []byte

	for _, node := range p.nodes {
		text := p.buffers[node.bufferIndex]
		nodeText := text[node.startIndex : node.startIndex+node.length]
		result = append(result, nodeText...)
	}

	// Replace all occurrences of "/n" with actual newlines
	return result
}

func (p *PieceTable) GetLineStart(line int) {

}

func (p *PieceTable) GetCurrentLineLength(line int) int {
	if line < 0 {
		return 0
	}
	currLine := 0
	currlength := 0
	prevLength := 0
	//line 1
	// curr 0
	//0+1-1
	if line == 0 {
		for _, node := range p.nodes {
			if len(node.lineStart) != 0 {
				currlength += node.lineStart[0] - node.startIndex
				return currlength + 1
			}
			currlength += node.length
		}
		return currlength + 1
	}
	//1
	//curr Line=0

	for indx, node := range p.nodes {
		if currLine+len(node.lineStart) >= line {
			//traverse till next line start keep adding the lenght of node
			if len(node.lineStart) == 1 {
				currlength += node.length - (node.lineStart[0] - node.startIndex) - 1
				for j := indx + 1; j < len(p.nodes); j++ {
					if len(p.nodes[j].lineStart) != 0 {
						currlength += p.nodes[j].lineStart[0] - p.nodes[j].startIndex
						return currlength + 1
					}
					currlength += p.nodes[j].length
				}
				return currlength + 1
			}
			//3
			//4 [0,5,1,4]
			currlength += node.lineStart[line-currLine] - node.lineStart[line-currLine-1]
			return currlength + 1
		}
		currLine += len(node.lineStart)
		prevLength += node.length
	}

	// If the requested line is out of range
	return -1
}
func (p *PieceTable) CharAt(position int) string {
	offset := 0
	//t
	//abs ->1
	if position == -1 {
		return ""
	}
	for _, node := range p.nodes {
		if offset+node.length > position {
			return string(p.buffers[node.bufferIndex][node.startIndex+position-offset])
		}
		offset += node.length
	}
	return ""
}

func (p *PieceTable) GetLineCount() int {
	lineCount := 0

	for _, node := range p.nodes {
		lineCount += len(node.lineStart)
	}
	return lineCount
}
func (p *PieceTable) GetAllNodes() []Node {
	return p.nodes
}
func (p *PieceTable) NextLine(line int, position int) bool {
	totalLines := p.GetLineCount()
	if line >= totalLines-1 {
		return false
	}

	nextLineLength := p.GetCurrentLineLength(line + 1)

	if position <= nextLineLength {
		return true
	} else {

		return true
	}

}

func (p *PieceTable) InsertText(position int, text string) {
	if len(text) == 0 {
		return
	}

	p.buffers = append(p.buffers, text)
	bufferIndex := len(p.buffers) - 1

	offset := 0
	var newNodes []Node

	for i, node := range p.nodes {
		if offset+node.length > position {
			insertOffset := position - offset
			// Split existing node if inserting in the middle
			beforeText := Node{
				bufferIndex: node.bufferIndex,
				startIndex:  node.startIndex,
				length:      insertOffset,
				lineStart:   extractLineStart(node.startIndex, node.lineStart, insertOffset),
			}

			afterText := Node{
				bufferIndex: node.bufferIndex,
				startIndex:  node.startIndex + insertOffset,
				length:      node.length - insertOffset,
				lineStart:   updateAfterTextLines(node.lineStart, insertOffset),
			}

			// Extract new line start positions
			newLineStart := []int{}
			for idx, ch := range text {
				if ch == '\n' {
					newLineStart = append(newLineStart, idx+1)
				}
			}

			insertedNode := Node{
				bufferIndex: bufferIndex,
				startIndex:  0,
				length:      len(text),
				lineStart:   newLineStart,
			}

			if beforeText.length > 0 {
				newNodes = append(newNodes, beforeText)
			}
			newNodes = append(newNodes, insertedNode)
			if afterText.length > 0 {
				newNodes = append(newNodes, afterText)
			}
			p.nodes = append(newNodes, p.nodes[i+1:]...)
			return
		}

		offset += node.length
		newNodes = append(newNodes, node)
	}

	// If position is at the end
	newLineStart := []int{}
	for idx, ch := range text {
		if ch == '\n' {
			newLineStart = append(newLineStart, idx)
		}
	}
	newNode := Node{
		bufferIndex: bufferIndex,
		startIndex:  0,
		length:      len(text),
		//not to insert the char
		lineStart: newLineStart,
	}
	p.nodes = append(p.nodes, newNode)
}

// Helper function to extract line start positions before the insert point
func extractLineStart(startIndex int, lineStart []int, insertOffset int) []int {
	newLineStart := []int{}
	for _, pos := range lineStart {
		if pos < insertOffset {
			newLineStart = append(newLineStart, pos)
		} else {
			break
		}
	}
	return newLineStart
}

func updateAfterTextLines(lineStart []int, insertOffset int) []int {
	newLineStart := []int{}
	for _, pos := range lineStart {
		if pos > insertOffset {
			newLineStart = append(newLineStart, pos-insertOffset)
		}
	}
	return newLineStart
}
