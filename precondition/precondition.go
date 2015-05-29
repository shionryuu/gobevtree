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

//
type IBevNodePrecondition interface {
	ExternalCondition(input interface{}) bool
}

// always true precondition
type BevNodePreconditionTRUE struct {
}

func NewBevNodePreconditionTRUE() *BevNodePreconditionTRUE {
	return &BevNodePreconditionTRUE{}
}

func (Cond *BevNodePreconditionTRUE) ExternalCondition(input interface{}) bool {
	return true
}

// always false precondition
type BevNodePreconditionFALSE struct {
}

func NewBevNodePreconditionFALSE() *BevNodePreconditionFALSE {
	return &BevNodePreconditionFALSE{}
}

func (Cond *BevNodePreconditionFALSE) ExternalCondition(input interface{}) bool {
	return false
}

// return true if both preconditions return true
type BevNodePreconditionAND struct {
	first  IBevNodePrecondition
	second IBevNodePrecondition
}

func NewBevNodePreconditionAND(First IBevNodePrecondition, Second IBevNodePrecondition) *BevNodePreconditionAND {
	return &BevNodePreconditionAND{first: First, second: Second}
}

func (Cond *BevNodePreconditionAND) ExternalCondition(input interface{}) bool {
	return Cond.first.ExternalCondition(input) &&
		Cond.second.ExternalCondition(input)
}

// return true if one of the preconditions return true
type BevNodePreconditionOR struct {
	first  IBevNodePrecondition
	second IBevNodePrecondition
}

func NewBevNodePreconditionOR(First IBevNodePrecondition, Second IBevNodePrecondition) *BevNodePreconditionOR {
	return &BevNodePreconditionOR{first: First, second: Second}
}

func (Cond *BevNodePreconditionOR) ExternalCondition(input interface{}) bool {
	return Cond.first.ExternalCondition(input) ||
		Cond.second.ExternalCondition(input)
}
