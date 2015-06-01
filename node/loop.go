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
type LoopSelector struct {
	*BevNode
	loopCount    int
	currentCount int
}

func NewLoopSelector(parentNode IBevNode, nodePrecondition p.IPrecondition, totalLoopCount int) *LoopSelector {
	return &LoopSelector{NewBevNode(parentNode, nodePrecondition), totalLoopCount, 0}
}

func (node *LoopSelector) Evaluate(input interface{}) bool {
	checkLoop := node.loopCount != ConstInfiniteLoop && node.currentCount > node.loopCount

	if !checkLoop {
		return false
	} else if node.checkIndex(0) {
		return node.childNodeList[0].Evaluate(input)
	} else {
		return false
	}
}

func (node *LoopSelector) Transition(input interface{}) {
	if node.checkIndex(0) {
		node.childNodeList[0].Transition(input)
	}
}

func (node *LoopSelector) Tick(input interface{}, output interface{}) BevRunningStatus {
	if node.checkIndex(0) && node.childNodeList[0].Tick(input, output) == StateFinish {
		node.currentCount = node.currentCount + 1
	}

	if node.loopCount != ConstInfiniteLoop && node.currentCount > node.loopCount {
		return StateFinish
	} else {
		return StateExecuting
	}
}
