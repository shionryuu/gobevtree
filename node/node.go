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
 * BevNode
 */
type BevNode struct {
	nodePrecondition p.IPrecondition
	parentNode       IBevNode
	activeNode       IBevNode
	lastActiveNode   IBevNode
	childNodeCount   int
	debugName        string
	childNodeList    [ConstMaxChildNodeCnt]IBevNode
}

func NewBevNode(parentNode IBevNode, nodePrecondition p.IPrecondition) *BevNode {
	return &BevNode{nodePrecondition: nodePrecondition, parentNode: parentNode}
}

func (node *BevNode) AddChildNode(childNode IBevNode) *BevNode {
	if node.childNodeCount == ConstMaxChildNodeCnt {
		return node
	}

	node.childNodeList[node.childNodeCount] = childNode
	node.childNodeCount += 1
	return node
}

func (node *BevNode) SetNodePrecondition(nodePrecondition p.IPrecondition) *BevNode {
	node.nodePrecondition = nodePrecondition
	return node
}

func (node *BevNode) GetNodePrecondition() p.IPrecondition {
	return node.nodePrecondition
}

func (node *BevNode) GetDebugName() string {
	return node.debugName
}

func (node *BevNode) SetDebugName(debugName string) *BevNode {
	node.debugName = debugName
	return node
}

func (node *BevNode) GetLastActiveNode() IBevNode {
	return node.lastActiveNode
}

func (node *BevNode) SetActiveNode(activeNode IBevNode) {
	node.lastActiveNode = node.activeNode
	node.activeNode = activeNode

	if node.parentNode != nil {
		node.parentNode.SetActiveNode(activeNode)
	}
}

func (node *BevNode) Evaluate(input interface{}) bool {
	return true
}

func (node *BevNode) Transition(input interface{}) {
}

func (node *BevNode) Tick(input interface{}, output interface{}) BevRunningStatus {
	return StateFinish
}

func (node *BevNode) checkIndex(index int) bool {
	return index >= 0 && index < node.childNodeCount
}

func (node *BevNode) getChildNodeCount() int {
	return node.childNodeCount
}

func (node *BevNode) Reconstruct(parentNode IBevNode) {
	node.parentNode = parentNode
	node.childNodeCount = len(node.childNodeList)
	for _, childNode := range node.childNodeList {
		childNode.Reconstruct(node)
	}
}
