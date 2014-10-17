package behavior_tree

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPrecondition(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given some integer with a starting value", t, func() {
		x := 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}

func TestTrueFalseCond(t *testing.T) {
	Convey("Alway true precondition", t, func() {
		cond := NewBevNodePreconditionTRUE()
		Convey("It should always be true", func() {
			So(cond.ExternalCondition(1), ShouldEqual, true)
		})
	})
	Convey("Alway false precondition", t, func() {
		cond := NewBevNodePreconditionFALSE()
		Convey("It should always be false", func() {
			So(cond.ExternalCondition(1), ShouldEqual, false)
		})
	})
}

func TestAndOrCond(t *testing.T) {
	trueCond := NewBevNodePreconditionTRUE()
	falseCond := NewBevNodePreconditionFALSE()

	Convey("And precondition", t, func() {
		Convey("Return false if one of the cond return false", func() {
			falseAndFalseCond := NewBevNodePreconditionAND(falseCond, falseCond)
			So(falseAndFalseCond.ExternalCondition(1), ShouldEqual, false)

			trueAndFalseCond := NewBevNodePreconditionAND(trueCond, falseCond)
			So(trueAndFalseCond.ExternalCondition(1), ShouldEqual, false)

			falseAndTrueCond := NewBevNodePreconditionAND(falseCond, trueCond)
			So(falseAndTrueCond.ExternalCondition(1), ShouldEqual, false)
		})

		Convey("Return true if both of the cond return true", func() {
			trueAndTrueCond := NewBevNodePreconditionAND(trueCond, trueCond)
			So(trueAndTrueCond.ExternalCondition(1), ShouldEqual, true)
		})
	})

	Convey("Or precondition", t, func() {
		Convey("Return true if one of the cond return true", func() {
			trueOrTrueCond := NewBevNodePreconditionOR(trueCond, trueCond)
			So(trueOrTrueCond.ExternalCondition(1), ShouldEqual, true)

			trueOrFalseCond := NewBevNodePreconditionOR(trueCond, falseCond)
			So(trueOrFalseCond.ExternalCondition(1), ShouldEqual, true)

			falseOrTrueCond := NewBevNodePreconditionOR(falseCond, trueCond)
			So(falseOrTrueCond.ExternalCondition(1), ShouldEqual, true)
		})

		Convey("Return false if both of the cond return false", func() {
			falseOrFalseCond := NewBevNodePreconditionOR(falseCond, falseCond)
			So(falseOrFalseCond.ExternalCondition(1), ShouldEqual, false)
		})
	})
}
