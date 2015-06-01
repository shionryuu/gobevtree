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
type IPrecondition interface {
	ExternalCondition(input interface{}) bool
}

// always true precondition
type PreconditionTRUE struct {
}

func NewPreconditionTRUE() *PreconditionTRUE {
	return &PreconditionTRUE{}
}

func (Cond *PreconditionTRUE) ExternalCondition(input interface{}) bool {
	return true
}

// always false precondition
type PreconditionFALSE struct {
}

func NewPreconditionFALSE() *PreconditionFALSE {
	return &PreconditionFALSE{}
}

func (Cond *PreconditionFALSE) ExternalCondition(input interface{}) bool {
	return false
}

// return true if both preconditions return true
type PreconditionAND struct {
	first  IPrecondition
	second IPrecondition
}

func NewPreconditionAND(First IPrecondition, Second IPrecondition) *PreconditionAND {
	return &PreconditionAND{first: First, second: Second}
}

func (Cond *PreconditionAND) ExternalCondition(input interface{}) bool {
	return Cond.first.ExternalCondition(input) &&
		Cond.second.ExternalCondition(input)
}

// return true if one of the preconditions return true
type PreconditionOR struct {
	first  IPrecondition
	second IPrecondition
}

func NewPreconditionOR(First IPrecondition, Second IPrecondition) *PreconditionOR {
	return &PreconditionOR{first: First, second: Second}
}

func (Cond *PreconditionOR) ExternalCondition(input interface{}) bool {
	return Cond.first.ExternalCondition(input) ||
		Cond.second.ExternalCondition(input)
}
