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

/*
 * BevRunningStatus
 */
const (
	k_BRS_ERROR_Transition = -1
	k_BRS_Executing        = iota
	k_BRS_Finish
)

/*
 * E_TerminalNodeStaus
 */
const (
	_ = iota
	k_TNS_Ready
	k_TNS_Running
	k_TNS_Finish
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
	mo_NodePrecondition BevNodePrecondition
	mo_ParentNode       *BevNode
	mo_ActiveNode       *BevNode
	mo_LastActiveNode   *BevNode
	mul_ChildNodeCount  int
	mz_DebugName        string
	mao_ChildNodeList   [k_BLimited_MaxChildNodeCnt]*BevNode
}

func NewBevNode(_o_ParentNode *BevNode, _o_NodeScript BevNodePrecondition) *BevNode {
	return &BevNode{mo_NodePrecondition: _o_NodeScript, mo_ParentNode: _o_ParentNode}
}

func (node *BevNode) AddChildNode(_o_ChildNode *BevNode) *BevNode {
	if node.mul_ChildNodeCount == k_BLimited_MaxChildNodeCnt {
		return node
	}

	node.mao_ChildNodeList[node.mul_ChildNodeCount] = _o_ChildNode
	node.mul_ChildNodeCount += 1
	return node
}

func (node *BevNode) SetNodePrecondition(_o_NodePrecondition BevNodePrecondition) *BevNode {
	if node.mo_NodePrecondition != _o_NodePrecondition {
		node.mo_NodePrecondition = _o_NodePrecondition
	}
	return node
}

func (node *BevNode) SetDebugName(_debugName string) *BevNode {
	node.mz_DebugName = _debugName
	return node
}

func (node *BevNode) GetLastActiveNode() *BevNode {
	return node.mo_LastActiveNode
}

func (node *BevNode) SetActiveNode(_o_Node *BevNode) {
	node.mo_LastActiveNode = node.mo_ActiveNode
	node.mo_ActiveNode = _o_Node

	if node.mo_ParentNode != nil {
		node.mo_ParentNode.SetActiveNode(_o_Node)
	}
}

func (node *BevNode) GetDebugName() string {
	return node.mz_DebugName
}

func (node *BevNode) Evaluate(input interface{}) bool {
	return (node.mo_NodePrecondition == nil || node.mo_NodePrecondition.ExternalCondition(input)) && node.doEvaluate(input)
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
	return k_BRS_Finish
}

func (node *BevNode) setParentNode(_o_ParentNode *BevNode) {
	node.mo_ParentNode = _o_ParentNode
}

func (node *BevNode) checkIndex(_ui_Index int) bool {
	return _ui_Index >= 0 && _ui_Index < node.mul_ChildNodeCount
}

/*
 *
 */
type BevNodePrioritySelector struct {
	*BevNode
	mui_CurrentSelectIndex int
	mui_LastSelectIndex    int
}

func NewBevNodePrioritySelector(_o_ParentNode *BevNode, _o_NodePrecondition BevNodePrecondition) *BevNodePrioritySelector {
	return &BevNodePrioritySelector{NewBevNode(_o_ParentNode, _o_NodePrecondition), k_BLimited_InvalidChildNodeIndex, k_BLimited_InvalidChildNodeIndex}
}

func (node *BevNodePrioritySelector) doEvaluate(input interface{}) bool {
	return true
}

func (node *BevNodePrioritySelector) doTransition(input interface{}) {

}

func (node *BevNodePrioritySelector) doTick(input interface{}, output interface{}) BevRunningStatus {
	return k_BRS_Finish
}

/*
 *
 */
type BevNodeNonePrioritySelector struct {
	*BevNodePrioritySelector
}

func NewBevNodeNonePrioritySelector(_o_ParentNode *BevNode, _o_NodePrecondition BevNodePrecondition) *BevNodeNonePrioritySelector {
	return &BevNodeNonePrioritySelector{NewBevNodePrioritySelector(_o_ParentNode, _o_NodePrecondition)}
}

func (node *BevNodeNonePrioritySelector) doEvaluate(input interface{}) bool {
	return true
}

/*
 *
 */
type BevNodeSequence struct {
	*BevNode
	mui_CurrentNodeIndex int
}

func NewBevNodeSequence(_o_ParentNode *BevNode, _o_NodePrecondition BevNodePrecondition) *BevNodeSequence {
	return &BevNodeSequence{NewBevNode(_o_ParentNode, _o_NodePrecondition), k_BLimited_InvalidChildNodeIndex}
}

func (node *BevNodeSequence) doEvaluate(input interface{}) bool {
	return true
}

func (node *BevNodeSequence) doTransition(input interface{}) {

}

func (node *BevNodeSequence) doTick(input interface{}, output interface{}) BevRunningStatus {
	return k_BRS_Finish
}

/*
 *
 */
type BevNodeParallel struct {
	node                *BevNode
	mab_ChildNodeStatus [k_BLimited_MaxChildNodeCnt]BevRunningStatus
}

func NewBevNodeParallel(_o_ParentNode *BevNode, _o_NodePrecondition BevNodePrecondition) *BevNodeParallel {
	return &BevNodeParallel{node: NewBevNode(_o_ParentNode, _o_NodePrecondition)}
}

func (node *BevNodeParallel) doEvaluate(input interface{}) bool {
	return true
}

func (node *BevNodeParallel) doTransition(input interface{}) {

}

func (node *BevNodeParallel) doTick(input interface{}, output interface{}) BevRunningStatus {
	return k_BRS_Finish
}

/*
 *
 */
type BevNodeLoop struct {
	*BevNode
	mi_LoopCount    int
	mi_CurrentCount int
}

func NewBevNodeLoop(_o_ParentNode *BevNode, _o_NodePrecondition BevNodePrecondition, _i_LoopCnt int) *BevNodeLoop {
	return &BevNodeLoop{NewBevNode(_o_ParentNode, _o_NodePrecondition), _i_LoopCnt, 0}
}

func (node *BevNodeLoop) doEvaluate(input interface{}) bool {
	return true
}

func (node *BevNodeLoop) doTransition(input interface{}) {

}

func (node *BevNodeLoop) doTick(input interface{}, output interface{}) BevRunningStatus {
	return k_BRS_Finish
}

/*
 *
 */
type BevNodeTerminal struct {
	*BevNode
	me_Status   E_TerminalNodeStaus
	mb_NeedExit bool
}

func NewBevNodeTerminal(_o_ParentNode *BevNode, _o_NodePrecondition BevNodePrecondition) *BevNodeTerminal {
	return &BevNodeTerminal{NewBevNode(_o_ParentNode, _o_NodePrecondition), k_TNS_Ready, false}
}

func (node *BevNodeTerminal) doTransition(input interface{}) {

}

func (node *BevNodeTerminal) doTick(input interface{}, output interface{}) BevRunningStatus {
	return k_BRS_Finish
}

func doEnter(input interface{}) {

}

func doExecute(input interface{}, output interface{}) BevRunningStatus {
	return k_BRS_Finish
}

func doExit(input interface{}, _ui_ExitID BevRunningStatus) {

}
