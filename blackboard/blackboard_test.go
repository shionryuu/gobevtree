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

package blackboard

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBoolValue(t *testing.T) {
	blackboard := NewBlackboard()
	Convey("Get unknown value should return error", t, func() {
		value, err := blackboard.GetValueAsBool(1)
		So(value, ShouldEqual, false)
		So(err, ShouldNotEqual, nil)
	})
	Convey("You should get what you set as bool", t, func() {
		blackboard.SetValueAsBool(1, true)
		value, err := blackboard.GetValueAsBool(1)
		So(value, ShouldEqual, true)
		So(err, ShouldEqual, nil)
	})
}

func TestIntValue(t *testing.T) {
	blackboard := NewBlackboard()
	Convey("Get unknown value should return error", t, func() {
		value, err := blackboard.GetValueAsInt(1)
		So(value, ShouldEqual, 0)
		So(err, ShouldNotEqual, nil)
	})
	Convey("You should get what you set as int", t, func() {
		blackboard.SetValueAsInt(1, 1)
		value, err := blackboard.GetValueAsInt(1)
		So(value, ShouldEqual, 1)
		So(err, ShouldEqual, nil)
	})
}

func TestFloat32Value(t *testing.T) {
	blackboard := NewBlackboard()
	Convey("Get unknown value should return error", t, func() {
		value, err := blackboard.GetValueAsFloat32(1)
		So(value, ShouldEqual, 0)
		So(err, ShouldNotEqual, nil)
	})
	Convey("You should get what you set as float32", t, func() {
		blackboard.SetValueAsFloat32(1, 1)
		value, err := blackboard.GetValueAsFloat32(1)
		So(value, ShouldEqual, 1)
		So(err, ShouldEqual, nil)
	})
}

func TestFloat64Value(t *testing.T) {
	blackboard := NewBlackboard()
	Convey("Get unknown value should return error", t, func() {
		value, err := blackboard.GetValueAsFloat64(1)
		So(value, ShouldEqual, 0)
		So(err, ShouldNotEqual, nil)
	})
	Convey("You should get what you set as float64", t, func() {
		blackboard.SetValueAsFloat64(1, 1)
		value, err := blackboard.GetValueAsFloat64(1)
		So(value, ShouldEqual, 1)
		So(err, ShouldEqual, nil)
	})
}

func TestStringValue(t *testing.T) {
	blackboard := NewBlackboard()
	Convey("Get unknown value should return error", t, func() {
		value, err := blackboard.GetValueAsString(1)
		So(value, ShouldEqual, "")
		So(err, ShouldNotEqual, nil)
	})
	Convey("You should get what you set as string", t, func() {
		blackboard.SetValueAsString(1, "true")
		value, err := blackboard.GetValueAsString(1)
		So(value, ShouldEqual, "true")
		So(err, ShouldEqual, nil)
	})
}

func TestInterfaceValue(t *testing.T) {
	blackboard := NewBlackboard()
	Convey("Get unknown value should return error", t, func() {
		value, err := blackboard.GetValueAsInterface(1)
		So(value, ShouldEqual, nil)
		So(err, ShouldNotEqual, nil)
	})
	Convey("You should get what you set as interface", t, func() {
		blackboard.SetValueAsInterface(1, blackboard)
		value, err := blackboard.GetValueAsInterface(1)
		So(value, ShouldEqual, blackboard)
		So(err, ShouldEqual, nil)
	})
}
