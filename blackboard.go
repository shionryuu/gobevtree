/*
 * Description: Behaviour tree in Go.
 * Copyright (c) 2014 ShionRyuu <shionryuu@outlook.com>.
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

package behavior_tree

import (
	"errors"
)

var (
	ErrInvalidKey  = errors.New("Invalid Key")
	ErrInvalidType = errors.New("Invalid Type")
)

type IBlackBoard interface{}

type BlackBoard struct {
	board map[int]interface{}
}

func NewBlackboard() *BlackBoard {
	mBoard := make(map[int]interface{}, 1000)
	return &BlackBoard{board: mBoard}
}

/*
 * GetValueAsBool, SetValueAsBool
 * GetValueAsInt, SetValueAsInt
 * GetValueAsFloat, SetValueAsFloat
 * GetValueAsString, SetValueAsString
 * GetValuesAsInterface, SetValuesAsInterface
 */
func (b *BlackBoard) GetValueAsBool(key int) (bool, error) {
	if v, ok := b.board[key]; ok {
		if i, ok := v.(bool); ok {
			return i, nil
		}
		return false, ErrInvalidType
	}
	return false, ErrInvalidKey
}

func (b *BlackBoard) SetValueAsBool(key int, v bool) {
	b.board[key] = v
}

func (b *BlackBoard) GetValuesAsInterface(key int) (interface{}, error) {
	if v, ok := b.board[key]; ok {
		return v, nil
	}
	return nil, ErrInvalidKey
}

func (b *BlackBoard) SetValuesAsInterface(key int, v interface{}) {
	b.board[key] = v
}
