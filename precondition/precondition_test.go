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

package precondition

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
		cond := NewPreconditionTRUE()
		Convey("It should always be true", func() {
			So(cond.ExternalCondition(1), ShouldEqual, true)
		})
	})
	Convey("Alway false precondition", t, func() {
		cond := NewPreconditionFALSE()
		Convey("It should always be false", func() {
			So(cond.ExternalCondition(1), ShouldEqual, false)
		})
	})
}

func TestAndOrCond(t *testing.T) {
	trueCond := NewPreconditionTRUE()
	falseCond := NewPreconditionFALSE()

	Convey("And precondition", t, func() {
		Convey("Return false if one of the cond return false", func() {
			falseAndFalseCond := NewPreconditionAND(falseCond, falseCond)
			So(falseAndFalseCond.ExternalCondition(1), ShouldEqual, false)

			trueAndFalseCond := NewPreconditionAND(trueCond, falseCond)
			So(trueAndFalseCond.ExternalCondition(1), ShouldEqual, false)

			falseAndTrueCond := NewPreconditionAND(falseCond, trueCond)
			So(falseAndTrueCond.ExternalCondition(1), ShouldEqual, false)
		})

		Convey("Return true if both of the cond return true", func() {
			trueAndTrueCond := NewPreconditionAND(trueCond, trueCond)
			So(trueAndTrueCond.ExternalCondition(1), ShouldEqual, true)
		})
	})

	Convey("Or precondition", t, func() {
		Convey("Return true if one of the cond return true", func() {
			trueOrTrueCond := NewPreconditionOR(trueCond, trueCond)
			So(trueOrTrueCond.ExternalCondition(1), ShouldEqual, true)

			trueOrFalseCond := NewPreconditionOR(trueCond, falseCond)
			So(trueOrFalseCond.ExternalCondition(1), ShouldEqual, true)

			falseOrTrueCond := NewPreconditionOR(falseCond, trueCond)
			So(falseOrTrueCond.ExternalCondition(1), ShouldEqual, true)
		})

		Convey("Return false if both of the cond return false", func() {
			falseOrFalseCond := NewPreconditionOR(falseCond, falseCond)
			So(falseOrFalseCond.ExternalCondition(1), ShouldEqual, false)
		})
	})
}
