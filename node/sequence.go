/*
 * Description: Behaviour tree in Go.
 * Copyright (c) 2014-2015 ShionRyuu <shionryuu@outlook.com>.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package node

import (
	p "github.com/ShionRyuu/gobevtree/precondition"
)

/*
 *
 */
type SequenceSelector struct {
	*BevNode
	currentSelectIndex int
}

func NewSequenceSelector(parentNode IBevNode, nodePrecondition p.IPrecondition) *SequenceSelector {
	return &SequenceSelector{NewBevNode(parentNode, nodePrecondition), ConstInvalidChildNodeIndex}
}

func (node *SequenceSelector) Evaluate(input interface{}) bool {
	Index := node.currentSelectIndex
	if !node.checkIndex(node.currentSelectIndex) && node.checkIndex(0) {
		Index = 0
	}
	if node.checkIndex(Index) {
		return node.childNodeList[Index].Evaluate(input)
	}
	return false
}

func (node *SequenceSelector) Transition(input interface{}) {
	if node.checkIndex(node.currentSelectIndex) {
		node.childNodeList[node.currentSelectIndex].Transition(input)
	}
	node.currentSelectIndex = ConstInvalidChildNodeIndex
}

func (node *SequenceSelector) Tick(input interface{}, output interface{}) BevRunningStatus {
	var bIsFinish BevRunningStatus = StateFinish

	if !node.checkIndex(node.currentSelectIndex) && node.checkIndex(0) {
		node.currentSelectIndex = 0
	}

	if node.checkIndex(node.currentSelectIndex) {
		bIsFinish = node.childNodeList[node.currentSelectIndex].Tick(input, output)
		if bIsFinish == StateFinish {
			node.currentSelectIndex += 1
			if node.currentSelectIndex >= node.childNodeCount {
				bIsFinish = StateFinish
			} else {
				bIsFinish = StateExecuting
			}
		}
	}

	if bIsFinish < 0 {
		node.currentSelectIndex = ConstInvalidChildNodeIndex
	}

	return bIsFinish
}
