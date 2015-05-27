package behavior_tree

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type A struct {
	*BevNodeTerminal
	v int
}

func NewA(parentNode IBevNode, cond IBevNodePrecondition, v int) *A {
	return &A{
		NewBevNodeTerminal(parentNode, cond),
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
		selector := NewSelector(NewBevNodeParallel(nil, NewBevNodePreconditionTRUE()))
		selector.AddChildNode(NewTerminal(NewA(selector, NewBevNodePreconditionTRUE(), 1)))
		selector.AddChildNode(NewTerminal(NewA(selector, NewBevNodePreconditionTRUE(), 10)))
		output := 0
		if selector.Evaluate(nil) {
			selector.Tick(nil, &output)
		}
		So(output, ShouldEqual, 10)
	})

	Convey("Parallel will fail if one of its children fails", t, func() {
		selector := NewSelector(NewBevNodeParallel(nil, NewBevNodePreconditionTRUE()))
		selector.AddChildNode(NewTerminal(NewA(selector, NewBevNodePreconditionTRUE(), 1)))
		selector.AddChildNode(NewTerminal(NewA(selector, NewBevNodePreconditionFALSE(), 10)))
		output := 0
		if selector.Evaluate(nil) {
			selector.Tick(nil, &output)
		}
		So(output, ShouldEqual, 0)
	})
}
