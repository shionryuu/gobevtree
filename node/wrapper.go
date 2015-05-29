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

/*
 * Wrapper for Selector
 * https://groups.google.com/d/msg/golang-nuts/BKztgPqN87M/iUfZQIcNYfYJ
 */
type IBevSelector interface {
	IBevNode
}

type BevSelector struct {
	IBevSelector
}

func NewSelector(node IBevSelector) *BevSelector {
	return &BevSelector{node}
}

func (w *BevSelector) Evaluate(input interface{}) bool {
	nodePrecondition := w.IBevSelector.GetNodePrecondition()
	return (nodePrecondition == nil || nodePrecondition.ExternalCondition(input)) && w.IBevSelector.Evaluate(input)
}

/*
 * Wrapper for Terminal
 */
type IBevTerminal interface {
	IBevNode
	Enter(input interface{})
	Execute(input interface{}, output interface{}) BevRunningStatus
	Exit(input interface{}, exitStatus BevRunningStatus)
}

type BevTerminal struct {
	IBevTerminal
	nodeStatus TerminalNodeStaus
	needExit   bool
}

func NewTerminal(node IBevTerminal) *BevTerminal {
	return &BevTerminal{node, NodeReady, false}
}

func (w *BevTerminal) Evaluate(input interface{}) bool {
	nodePrecondition := w.IBevTerminal.GetNodePrecondition()
	return (nodePrecondition == nil || nodePrecondition.ExternalCondition(input)) && w.IBevTerminal.Evaluate(input)
}

func (node *BevTerminal) Transition(input interface{}) {
	if node.needExit {
		node.Exit(input, StateTransition)
	}

	node.SetActiveNode(nil)
	node.nodeStatus = NodeReady
	node.needExit = false
}

func (node *BevTerminal) Tick(input interface{}, output interface{}) BevRunningStatus {
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

/*
 * Wrapper used to reverse the evaluate result of node
 */
type BevReverse struct {
	IBevNode
}

func NewReverse(node IBevNode) *BevReverse {
	return &BevReverse{node}
}

func (w *BevReverse) Evaluate(input interface{}) bool {
	return !w.IBevNode.Evaluate(input)
}
