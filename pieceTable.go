package main

type Node struct {
	bufferIndex int
	startIndex  int
	length      int
	lineStart   []int
}

type PieceTable struct {
	buffers []string
	nodes   []Node
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
func (p *PieceTable) GetText() (outputString string) {
	for _, node := range p.nodes {
		startOffset := 0
		if len(node.lineStart) == 0 {
			outputString += p.buffers[node.bufferIndex]
		}
		for _, lines := range node.lineStart {
			endOffset := lines
			outputString += p.buffers[node.bufferIndex][startOffset:endOffset]
			startOffset = endOffset + 1
		}
	}
	return outputString
}

func (p *PieceTable) InsertText(position int, line int, text string) {
	if len(text) == 0 {
		return
	}

	p.buffers = append(p.buffers, text)
	bufferIndex := len(p.buffers) - 1

	offset := 0
	var newNodes []Node

	for i, node := range p.nodes {
		if offset+node.length > position {
			// Splitting the current node
			before := Node{node.bufferIndex, node.startIndex, position - offset, append([]int{}, node.lineStart...)}
			insert := Node{bufferIndex, 0, len(text), []int{line}}
			after := Node{node.bufferIndex, node.startIndex + before.length, node.length - before.length, append([]int{}, node.lineStart...)}

			if before.length > 0 {
				newNodes = append(newNodes, before)
			}
			newNodes = append(newNodes, insert)
			if after.length > 0 {
				newNodes = append(newNodes, after)
			}

			// Append the rest of the existing nodes
			newNodes = append(newNodes, p.nodes[i+1:]...)
			p.nodes = newNodes
			return
		}
		offset += node.length
		newNodes = append(newNodes, node)
	}

	// If position is at the end, simply append the new character
	newNode := Node{bufferIndex, 0, len(text), []int{len(text)}}
	p.nodes = append(p.nodes, newNode)
}
