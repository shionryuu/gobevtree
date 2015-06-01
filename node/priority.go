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
type PrioritySelector struct {
	*BevNode
	currentSelectIndex int
	lastSelectIndex    int
}

func NewPrioritySelector(parentNode IBevNode, nodePrecondition p.IPrecondition) *PrioritySelector {
	return &PrioritySelector{NewBevNode(parentNode, nodePrecondition), ConstInvalidChildNodeIndex, ConstInvalidChildNodeIndex}
}

func (node *PrioritySelector) Evaluate(input interface{}) bool {
	node.currentSelectIndex = ConstInvalidChildNodeIndex
	for i := 0; i < node.childNodeCount; i++ {
		if node.childNodeList[i].Evaluate(input) {
			node.currentSelectIndex = i
			return true
		}
	}
	return false
}

func (node *PrioritySelector) Transition(input interface{}) {
	if node.checkIndex(node.lastSelectIndex) {
		node.childNodeList[node.lastSelectIndex].Transition(input)
	}
	node.lastSelectIndex = ConstInvalidChildNodeIndex
}

func (node *PrioritySelector) Tick(input interface{}, output interface{}) BevRunningStatus {
	var isFinish BevRunningStatus = StateFinish

	if node.checkIndex(node.currentSelectIndex) {
		if node.lastSelectIndex != node.currentSelectIndex {
			if node.checkIndex(node.lastSelectIndex) {
				lastNode := node.childNodeList[node.lastSelectIndex]
				lastNode.Transition(input)
			}
			node.lastSelectIndex = node.currentSelectIndex
		}
	}

	if node.checkIndex(node.lastSelectIndex) {
		curNode := node.childNodeList[node.lastSelectIndex]
		isFinish = curNode.Tick(input, output)
		if isFinish == StateFinish {
			node.lastSelectIndex = ConstInvalidChildNodeIndex
		}
	}

	return isFinish
}
