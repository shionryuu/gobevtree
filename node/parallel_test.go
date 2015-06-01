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
	. "github.com/ShionRyuu/gobevtree/precondition"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type A struct {
	*TerminalNode
	v int
}

func NewA(parentNode IBevNode, cond IPrecondition, v int) *A {
	return &A{
		NewTerminalNode(parentNode, cond),
		v,
	}
}

func (node *A) Enter(input interface{}) {
}

func (node *A) Execute(input interface{}, output interface{}) BevRunningStatus {
	if o, ok := output.(*int); ok {
		*o = node.v
	}
	return StateFinish
}

func (node *A) Exit(input interface{}, exitStatus BevRunningStatus) {
}

func TestParallel(t *testing.T) {
	Convey("Parallel will succ if all of its children succs", t, func() {
		selector := NewSelector(NewParallelSelector(nil, NewPreconditionTRUE()))
		selector.AddChildNode(NewTerminal(NewA(selector, NewPreconditionTRUE(), 1)))
		selector.AddChildNode(NewTerminal(NewA(selector, NewPreconditionTRUE(), 10)))
		output := 0
		if selector.Evaluate(nil) {
			selector.Tick(nil, &output)
		}
		So(output, ShouldEqual, 10)
	})

	Convey("Parallel will fail if one of its children fails", t, func() {
		selector := NewSelector(NewParallelSelector(nil, NewPreconditionTRUE()))
		selector.AddChildNode(NewTerminal(NewA(selector, NewPreconditionTRUE(), 1)))
		selector.AddChildNode(NewTerminal(NewA(selector, NewPreconditionFALSE(), 10)))
		output := 0
		if selector.Evaluate(nil) {
			selector.Tick(nil, &output)
		}
		So(output, ShouldEqual, 0)
	})
}
