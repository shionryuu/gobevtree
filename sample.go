package main

import (
	"fmt"
	_ "math/rand"
	"time"

	btboard "github.com/ShionRyuu/gobevtree/blackboard"
	btnode "github.com/ShionRyuu/gobevtree/node"
	btcond "github.com/ShionRyuu/gobevtree/precondition"
)

//print action
type TestTerNode struct {
	*btnode.TerminalNode
	data string
}

func (this *TestTerNode) Enter(input interface{}) {
	fmt.Println("enter node ", this.data)
}

func (this *TestTerNode) Execute(input interface{}, output interface{}) btnode.BevRunningStatus {
	fmt.Println("Execute node ", this.data)
	return btnode.StateFinish
}

func (this *TestTerNode) Exit(input interface{}, exitStatus btnode.BevRunningStatus) {
	fmt.Println("Exit node", this.data)
}

//wait action
const (
	delayTimeFrame = 1
	frame          = 0
)

type WaitActNode struct {
	*btnode.TerminalNode
	waitTime int //等多久
	useTime  int //目前等待时间
}

func NewWaitActNode(waitTime int, parent btnode.IBevNode) *WaitActNode {
	return &WaitActNode{btnode.NewTerminalNode(parent, nil), waitTime, 0}
}

func (this *WaitActNode) Enter(input interface{}) {
	fmt.Println("enter wait ", this.waitTime)
	this.useTime = 0
}

func (this *WaitActNode) Execute(input interface{}, output interface{}) btnode.BevRunningStatus {
	fmt.Println("Execute wait ", this.useTime, "/", this.waitTime)
	if this.useTime >= this.waitTime {
		return btnode.StateFinish
	}
	this.useTime += delayTimeFrame
	return btnode.StateExecuting
}

func (this *WaitActNode) Exit(input interface{}, exitStatus btnode.BevRunningStatus) {
	fmt.Println("Exit wait", this.waitTime)
	this.useTime = 0
}

//通过blackboard比较int条件
type PreconditionLessInt struct {
	first  int //变量索引
	second int //变量索引
}

func NewPreconditionLessInt(FirstIndex int, SecIndex int) *PreconditionLessInt {
	return &PreconditionLessInt{first: FirstIndex, second: SecIndex}
}

func (Cond *PreconditionLessInt) ExternalCondition(input interface{}) bool {
	board := input.(*btboard.BlackBoard)
	a, erra := board.GetValueAsInt(Cond.first)
	if erra != nil {
		fmt.Errorf(erra.Error())
	}

	b, errb := board.GetValueAsInt(Cond.second)
	if errb != nil {
		fmt.Errorf(errb.Error())
	}

	return a < b
}

func main() {
	fmt.Println("begin")

	testSequenceSelector()
	testParallelSelector()
	testPrioritySelector()
	testRandomSelector()
	testSimple()
	fmt.Println("end")
}

func testSequenceSelector() {
	fmt.Println("SequenceSelector===========>")
	inboard := btboard.NewBlackboard()
	outboard := btboard.NewBlackboard()

	tree := btnode.NewSequenceSelector(nil, nil)
	tree.SetDebugName("seq")
	node1 := &TestTerNode{btnode.NewTerminalNode(nil, nil), "node1"}
	node2 := &TestTerNode{btnode.NewTerminalNode(nil, nil), "node2"}
	wrap1 := btnode.NewTerminal(node1)
	wrap1.SetDebugName("w1")
	wrap2 := btnode.NewTerminal(node2)
	wrap2.SetDebugName("w2")
	tree.AddChildNode(wrap1)
	tree.AddChildNode(wrap2)
	renderTree(tree, 2, inboard, outboard, 0)

}
func testParallelSelector() {
	fmt.Println("ParallelSelector===========>")
	inboard := btboard.NewBlackboard()
	outboard := btboard.NewBlackboard()

	tree := btnode.NewParallelSelector(nil, nil)
	node1 := &TestTerNode{btnode.NewTerminalNode(nil, nil), "node1"}
	node2 := &TestTerNode{btnode.NewTerminalNode(nil, nil), "node2"}
	wrap1 := btnode.NewTerminal(node1)
	wrap2 := btnode.NewTerminal(node2)
	tree.AddChildNode(wrap1)
	tree.AddChildNode(wrap2)
	renderTree(tree, 2, inboard, outboard, 0)

}

func testPrioritySelector() {
	fmt.Println("PrioritySelector===========>")
	inboard := btboard.NewBlackboard()
	outboard := btboard.NewBlackboard()
	indexA := 1
	indexB := 2
	inboard.SetValueAsInt(indexA, 33) //设置1号变量
	inboard.SetValueAsInt(indexB, 10) //设置2号变量

	tree := btnode.NewPrioritySelector(nil, nil)
	node1 := &TestTerNode{btnode.NewTerminalNode(nil, NewPreconditionLessInt(indexA, indexB)), "node1"}
	node2 := &TestTerNode{btnode.NewTerminalNode(nil, btcond.NewPreconditionTRUE()), "node2"}
	wrap1 := btnode.NewTerminal(node1)
	wrap2 := btnode.NewTerminal(node2)
	tree.AddChildNode(wrap1)
	tree.AddChildNode(wrap2)
	renderTree(tree, 2, inboard, outboard, 0)

}

func testRandomSelector() {
	fmt.Println("RandomSelector===========>")
	inboard := btboard.NewBlackboard()
	outboard := btboard.NewBlackboard()

	tree := btnode.NewRandomSelector(nil, nil)
	node1 := &TestTerNode{btnode.NewTerminalNode(nil, nil), "node1"}
	node2 := &TestTerNode{btnode.NewTerminalNode(nil, nil), "node2"}
	wrap1 := btnode.NewTerminal(node1)
	wrap2 := btnode.NewTerminal(node2)
	tree.AddChildNode(wrap1)
	tree.AddChildNode(wrap2)
	renderTree(tree, 10, inboard, outboard, 0)

}
func renderTree(tree btnode.IBevNode, count int, inboard *btboard.BlackBoard, outboard *btboard.BlackBoard, delayTime int) {
	btnode.PrintbevTree(tree, 0)
	for i := 0; i < count; i++ {
		if tree.Evaluate(inboard) {
			tree.Tick(inboard, outboard)
		} // else {
		//	tree.Transition(inboard)
		//}
		if delayTime > 0 {
			time.Sleep(time.Duration(delayTime) * time.Second)
		}

	}

}

func testSimple() {
	/*				   selector（a<b）
	*				 /          \
	*           seq(selector)    rand（selector）
	*       /   |    \           /      \
	* say(11) wait(1s) say(12)  say(21) say(22)
	*
	 */
	//注意：只有经过wrapper封装，调用NewSelector，NewTerminal才会在执行precondition条件判断
	//结果：在seq下执行10次，再到右侧ran执行10次
	fmt.Println("testSimple===========>")
	inboard := btboard.NewBlackboard()
	outboard := btboard.NewBlackboard()
	indexA := 1
	indexB := 2

	inboard.SetValueAsInt(indexA, 0)  //设置1号变量
	inboard.SetValueAsInt(indexB, 10) //设置2号变量

	tree := btnode.NewPrioritySelector(nil, nil)
	cond := NewPreconditionLessInt(indexA, indexB)
	seq := btnode.NewSelector(btnode.NewSequenceSelector(nil, cond))
	randn := btnode.NewSelector(btnode.NewRandomSelector(nil, btcond.NewPreconditionTRUE()))
	tree.AddChildNode(seq)
	tree.AddChildNode(randn)

	node11 := &TestTerNode{btnode.NewTerminalNode(seq, nil), "node11"}
	node12 := &TestTerNode{btnode.NewTerminalNode(seq, nil), "node12"}
	node21 := &TestTerNode{btnode.NewTerminalNode(randn, nil), "node21"}
	node22 := &TestTerNode{btnode.NewTerminalNode(randn, nil), "node22"}
	waitAct := NewWaitActNode(5, seq)

	wrap11 := btnode.NewTerminal(node11)
	wrapWait := btnode.NewTerminal(waitAct)
	wrap12 := btnode.NewTerminal(node12)
	wrap21 := btnode.NewTerminal(node21)
	wrap22 := btnode.NewTerminal(node22)
	seq.AddChildNode(wrap11)
	seq.AddChildNode(wrapWait)
	seq.AddChildNode(wrap12)

	randn.AddChildNode(wrap21)
	randn.AddChildNode(wrap22)
	btnode.PrintbevTree(tree, 0)
	//renderTree
	for i := 0; i < 20; i++ {
		//one frame
		if tree.Evaluate(inboard) {
			tree.Tick(inboard, outboard)
		}
		fmt.Println("frame:", i)
		time.Sleep(time.Duration(delayTimeFrame) * time.Second)
		inboard.SetValueAsInt(indexA, i) //设置1号变量
	}

}
