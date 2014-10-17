/*
 * Description: Behaviour tree in Go.
 * Copyright (c) 2014 ShionRyuu <shionryuu@outlook.com>.
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

package behavior_tree

import (
	"math/rand"
)

/*
 * BevRunningStatus
 */
const (
	StateTransition = -1
	StateExecuting  = iota
	StateFinish
)

/*
 * TerminalNodeStaus
 */
const (
	_ = iota
	NodeReady
	NodeRunning
	NodeFinish
)

/*
 * Other const variables
 */
const (
	ConstInfiniteLoop          int = -1
	ConstMaxChildNodeCnt       int = 16
	ConstInvalidChildNodeIndex int = 16
)

type BevRunningStatus int
type TerminalNodeStaus int

/*
 *
 */
type IBevNode interface {
	// public
	AddChildNode(childNode IBevNode) *BevNode
	SetNodePrecondition(nodePrecondition IBevNodePrecondition) *BevNode
	GetDebugName() string
	SetDebugName(debugName string) *BevNode
	GetLastActiveNode() IBevNode
	SetActiveNode(activeNode IBevNode)
	Evaluate(input interface{}) bool
	Transition(input interface{})
	Tick(input interface{}, output interface{}) BevRunningStatus
	// private
	doEvaluate(input interface{}) bool
	doTransition(input interface{})
	doTick(input interface{}, output interface{}) BevRunningStatus
	setParentNode(parentNode IBevNode)
	checkIndex(index int) bool
}

/*
 *
 */
type BevNode struct {
	nodePrecondition IBevNodePrecondition
	parentNode       IBevNode
	activeNode       IBevNode
	lastActiveNode   IBevNode
	childNodeCount   int
	debugName        string
	childNodeList    [ConstMaxChildNodeCnt]IBevNode
}

func NewBevNode(parentNode IBevNode, precondition IBevNodePrecondition) *BevNode {
	return &BevNode{nodePrecondition: precondition, parentNode: parentNode}
}

func (node *BevNode) AddChildNode(childNode IBevNode) *BevNode {
	if node.childNodeCount == ConstMaxChildNodeCnt {
		return node
	}

	node.childNodeList[node.childNodeCount] = childNode
	node.childNodeCount += 1
	return node
}

func (node *BevNode) SetNodePrecondition(nodePrecondition IBevNodePrecondition) *BevNode {
	if node.nodePrecondition != nodePrecondition {
		node.nodePrecondition = nodePrecondition
	}
	return node
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
	return (node.nodePrecondition == nil || node.nodePrecondition.ExternalCondition(input)) && node.doEvaluate(input)
}

func (node *BevNode) Transition(input interface{}) {
	node.doTransition(input)
}

func (node *BevNode) Tick(input interface{}, output interface{}) BevRunningStatus {
	return node.doTick(input, output)
}

func (node *BevNode) doEvaluate(input interface{}) bool {
	return true
}

func (node *BevNode) doTransition(input interface{}) {
}

func (node *BevNode) doTick(input interface{}, output interface{}) BevRunningStatus {
	return StateFinish
}

func (node *BevNode) setParentNode(parentNode IBevNode) {
	node.parentNode = parentNode
}

func (node *BevNode) checkIndex(index int) bool {
	return index >= 0 && index < node.childNodeCount
}

/*
 *
 */
type BevNodePrioritySelector struct {
	*BevNode
	currentSelectIndex int
	lastSelectIndex    int
}

func NewBevNodePrioritySelector(parentNode IBevNode, nodePrecondition IBevNodePrecondition) *BevNodePrioritySelector {
	return &BevNodePrioritySelector{NewBevNode(parentNode, nodePrecondition), ConstInvalidChildNodeIndex, ConstInvalidChildNodeIndex}
}

func (node *BevNodePrioritySelector) doEvaluate(input interface{}) bool {
	node.currentSelectIndex = ConstInvalidChildNodeIndex
	for i := 0; i < node.childNodeCount; i++ {
		if node.childNodeList[i].Evaluate(input) {
			node.currentSelectIndex = i
			return true
		}
	}
	return false
}

func (node *BevNodePrioritySelector) doTransition(input interface{}) {
	if node.checkIndex(node.lastSelectIndex) {
		node.childNodeList[node.lastSelectIndex].doTransition(input)
	}
	node.lastSelectIndex = ConstInvalidChildNodeIndex
}

func (node *BevNodePrioritySelector) doTick(input interface{}, output interface{}) BevRunningStatus {
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

/*
 *
 */
type BevNodeNonePrioritySelector struct {
	*BevNodePrioritySelector
}

func NewBevNodeNonePrioritySelector(parentNode IBevNode, nodePrecondition IBevNodePrecondition) *BevNodeNonePrioritySelector {
	return &BevNodeNonePrioritySelector{NewBevNodePrioritySelector(parentNode, nodePrecondition)}
}

func (node *BevNodeNonePrioritySelector) doEvaluate(input interface{}) bool {
	if node.checkIndex(node.currentSelectIndex) {
		curNode := node.childNodeList[node.currentSelectIndex]
		if curNode.Evaluate(input) {
			return true
		}
	}
	return node.Evaluate(input)
}

/*
 *
 */
type BevNodeRandomSelector struct {
	*BevNodePrioritySelector
}

func NewBevNodeRandomSelector(parentNode IBevNode, nodePrecondition IBevNodePrecondition) *BevNodeRandomSelector {
	return &BevNodeRandomSelector{NewBevNodePrioritySelector(parentNode, nodePrecondition)}
}

func (node *BevNodeRandomSelector) doEvaluate(input interface{}) bool {
	if node.childNodeCount >= 1 {
		randomIndex := rand.Intn(node.childNodeCount)
		if node.childNodeList[randomIndex].Evaluate(input) == true {
			node.currentSelectIndex = randomIndex
			return true
		}
	}
	return false
}

/*
 *
 */
type BevNodeSequence struct {
	*BevNode
	currentSelectIndex int
}

func NewBevNodeSequence(parentNode IBevNode, nodePrecondition IBevNodePrecondition) *BevNodeSequence {
	return &BevNodeSequence{NewBevNode(parentNode, nodePrecondition), ConstInvalidChildNodeIndex}
}

func (node *BevNodeSequence) doEvaluate(input interface{}) bool {
	Index := node.currentSelectIndex
	if !node.checkIndex(node.currentSelectIndex) && node.checkIndex(0) {
		Index = 0
	}
	if node.checkIndex(Index) {
		return node.childNodeList[Index].Evaluate(input)
	}
	return false
}

func (node *BevNodeSequence) doTransition(input interface{}) {
	if node.checkIndex(node.currentSelectIndex) {
		node.childNodeList[node.currentSelectIndex].Transition(input)
	}
	node.currentSelectIndex = ConstInvalidChildNodeIndex
}

func (node *BevNodeSequence) doTick(input interface{}, output interface{}) BevRunningStatus {
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

/*
 *
 */
type BevNodeParallel struct {
	*BevNode
}

func NewBevNodeParallel(parentNode IBevNode, nodePrecondition IBevNodePrecondition) *BevNodeParallel {
	return &BevNodeParallel{NewBevNode(parentNode, nodePrecondition)}
}

func (node *BevNodeParallel) doEvaluate(input interface{}) bool {
	for i := 0; i < node.childNodeCount; i++ {
		if !node.childNodeList[i].Evaluate(input) {
			return false
		}
	}
	return true
}

func (node *BevNodeParallel) doTransition(input interface{}) {
	for i := 0; i < node.childNodeCount; i++ {
		node.childNodeList[i].Transition(input)
	}
}

func (node *BevNodeParallel) doTick(input interface{}, output interface{}) BevRunningStatus {
	for i := 0; i < node.childNodeCount; i++ {
		if node.childNodeList[i].Tick(input, output) != StateFinish {
			return StateExecuting
		}
	}
	return StateFinish
}

/*
 *
 */
type BevNodeLoop struct {
	*BevNode
	loopCount    int
	currentCount int
}

func NewBevNodeLoop(parentNode IBevNode, nodePrecondition IBevNodePrecondition, totalLoopCount int) *BevNodeLoop {
	return &BevNodeLoop{NewBevNode(parentNode, nodePrecondition), totalLoopCount, 0}
}

func (node *BevNodeLoop) doEvaluate(input interface{}) bool {
	checkLoop := node.loopCount != ConstInfiniteLoop && node.currentCount > node.loopCount

	if !checkLoop {
		return false
	} else if node.checkIndex(0) {
		return node.childNodeList[0].Evaluate(input)
	} else {
		return false
	}
}

func (node *BevNodeLoop) doTransition(input interface{}) {
	if node.checkIndex(0) {
		node.childNodeList[0].Transition(input)
	}
}

func (node *BevNodeLoop) doTick(input interface{}, output interface{}) BevRunningStatus {
	if node.checkIndex(0) && node.childNodeList[0].Tick(input, output) == StateFinish {
		node.currentCount = node.currentCount + 1
	}

	if node.loopCount != ConstInfiniteLoop && node.currentCount > node.loopCount {
		return StateFinish
	} else {
		return StateExecuting
	}
}

/*
 *
 */
type BevNodeTerminal struct {
	*BevNode
	nodeStatus TerminalNodeStaus
	needExit   bool
}

func NewBevNodeTerminal(parentNode IBevNode, nodePrecondition IBevNodePrecondition) *BevNodeTerminal {
	return &BevNodeTerminal{NewBevNode(parentNode, nodePrecondition), NodeReady, false}
}

func (node *BevNodeTerminal) doTransition(input interface{}) {
	if node.needExit {
		node.doExit(input, StateTransition)
	}

	node.SetActiveNode(nil)
	node.nodeStatus = NodeReady
	node.needExit = false
}

func (node *BevNodeTerminal) doTick(input interface{}, output interface{}) BevRunningStatus {
	var bIsFinish BevRunningStatus = StateFinish

	if node.nodeStatus == NodeReady {
		node.Enter(input)
		node.needExit = true
		node.nodeStatus = NodeRunning
		node.SetActiveNode(node)
	}

	if node.nodeStatus == NodeRunning {
		bIsFinish = node.Execute(input, output)
		node.SetActiveNode(node)
		if bIsFinish == StateFinish || bIsFinish < 0 {
			node.nodeStatus = NodeFinish
		}
	}

	if node.nodeStatus == NodeFinish {
		if node.needExit {
			node.Exit(input, bIsFinish)
		}

		node.nodeStatus = NodeReady
		node.needExit = false
		node.SetActiveNode(nil)
	}

	return bIsFinish
}

func (node *BevNodeTerminal) Enter(input interface{}) {
	node.doEnter(input)
}

func (node *BevNodeTerminal) Execute(input interface{}, output interface{}) BevRunningStatus {
	return node.doExecute(input, output)
}

func (node *BevNodeTerminal) Exit(input interface{}, exitStatus BevRunningStatus) {
	node.doExit(input, exitStatus)
}

func (node *BevNodeTerminal) doEnter(input interface{}) {

}

func (node *BevNodeTerminal) doExecute(input interface{}, output interface{}) BevRunningStatus {
	return StateFinish
}

func (node *BevNodeTerminal) doExit(input interface{}, exitStatus BevRunningStatus) {

}
