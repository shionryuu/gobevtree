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
 * E_TerminalNodeStaus
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
	kInfiniteLoop                    int = -1
	k_BLimited_MaxChildNodeCnt       int = 16
	k_BLimited_InvalidChildNodeIndex int = 16
)

type BevRunningStatus int
type E_TerminalNodeStaus int

/*
 *
 */
type BevNode struct {
	nodePrecondition BevNodePrecondition
	parentNode       *BevNode // ? interface{}
	activeNode       *BevNode
	lastActiveNode   *BevNode
	childNodeCount   int
	debugName        string
	childNodeList    [k_BLimited_MaxChildNodeCnt]*BevNode
}

func NewBevNode(parentNode *BevNode, _o_NodeScript BevNodePrecondition) *BevNode {
	return &BevNode{nodePrecondition: _o_NodeScript, parentNode: parentNode}
}

func (node *BevNode) AddChildNode(childNode *BevNode) *BevNode {
	if node.childNodeCount == k_BLimited_MaxChildNodeCnt {
		return node
	}

	node.childNodeList[node.childNodeCount] = childNode
	node.childNodeCount += 1
	return node
}

func (node *BevNode) SetNodePrecondition(nodePrecondition BevNodePrecondition) *BevNode {
	if node.nodePrecondition != nodePrecondition {
		node.nodePrecondition = nodePrecondition
	}
	return node
}

func (node *BevNode) SetDebugName(_debugName string) *BevNode {
	node.debugName = _debugName
	return node
}

func (node *BevNode) GetLastActiveNode() *BevNode {
	return node.lastActiveNode
}

func (node *BevNode) setActiveNode(activeNode *BevNode) {
	node.lastActiveNode = node.activeNode
	node.activeNode = activeNode

	if node.parentNode != nil {
		node.parentNode.setActiveNode(activeNode)
	}
}

func (node *BevNode) GetDebugName() string {
	return node.debugName
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

func (node *BevNode) setParentNode(parentNode *BevNode) {
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

func NewBevNodePrioritySelector(parentNode *BevNode, nodePrecondition BevNodePrecondition) *BevNodePrioritySelector {
	return &BevNodePrioritySelector{NewBevNode(parentNode, nodePrecondition), k_BLimited_InvalidChildNodeIndex, k_BLimited_InvalidChildNodeIndex}
}

func (node *BevNodePrioritySelector) doEvaluate(input interface{}) bool {
	node.currentSelectIndex = k_BLimited_InvalidChildNodeIndex
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
	node.lastSelectIndex = k_BLimited_InvalidChildNodeIndex
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
			node.lastSelectIndex = k_BLimited_InvalidChildNodeIndex
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

func NewBevNodeNonePrioritySelector(parentNode *BevNode, nodePrecondition BevNodePrecondition) *BevNodeNonePrioritySelector {
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

func NewBevNodeRandomSelector(parentNode *BevNode, nodePrecondition BevNodePrecondition) *BevNodeRandomSelector {
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

func NewBevNodeSequence(parentNode *BevNode, nodePrecondition BevNodePrecondition) *BevNodeSequence {
	return &BevNodeSequence{NewBevNode(parentNode, nodePrecondition), k_BLimited_InvalidChildNodeIndex}
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
	node.currentSelectIndex = k_BLimited_InvalidChildNodeIndex
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
		node.currentSelectIndex = k_BLimited_InvalidChildNodeIndex
	}

	return bIsFinish
}

/*
 *
 */
type BevNodeParallel struct {
	*BevNode
}

func NewBevNodeParallel(parentNode *BevNode, nodePrecondition BevNodePrecondition) *BevNodeParallel {
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

func NewBevNodeLoop(parentNode *BevNode, nodePrecondition BevNodePrecondition, totalLoopCount int) *BevNodeLoop {
	return &BevNodeLoop{NewBevNode(parentNode, nodePrecondition), totalLoopCount, 0}
}

func (node *BevNodeLoop) doEvaluate(input interface{}) bool {
	checkLoop := node.loopCount != kInfiniteLoop && node.currentCount > node.loopCount

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

	if node.loopCount != kInfiniteLoop && node.currentCount > node.loopCount {
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
	nodeStatus E_TerminalNodeStaus
	needExit   bool
}

func NewBevNodeTerminal(parentNode *BevNode, nodePrecondition BevNodePrecondition) *BevNodeTerminal {
	return &BevNodeTerminal{NewBevNode(parentNode, nodePrecondition), NodeReady, false}
}

func (node *BevNodeTerminal) doTransition(input interface{}) {
	if node.needExit {
		node.doExit(input, StateTransition)
	}

	node.setActiveNode(nil)
	node.nodeStatus = NodeReady
	node.needExit = false
}

func (node *BevNodeTerminal) doTick(input interface{}, output interface{}) BevRunningStatus {
	var bIsFinish BevRunningStatus = StateFinish

	if node.nodeStatus == NodeReady {
		node.Enter(input)
		node.needExit = true
		node.nodeStatus = NodeRunning
		// node.setActiveNode(node)
	}

	if node.nodeStatus == NodeRunning {
		bIsFinish = node.Execute(input, output)
		// node.setActiveNode(node)
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
		node.setActiveNode(nil)
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
