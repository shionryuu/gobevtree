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

func (b *BlackBoard) GetValueAsInt(key int) (int, error) {
	if v, ok := b.board[key]; ok {
		if i, ok := v.(int); ok {
			return i, nil
		}
		return 0, ErrInvalidType
	}
	return 0, ErrInvalidKey
}

func (b *BlackBoard) SetValueAsInt(key int, v int) {
	b.board[key] = v
}

func (b *BlackBoard) GetValueAsFloat32(key int) (float32, error) {
	if v, ok := b.board[key]; ok {
		if i, ok := v.(float32); ok {
			return i, nil
		}
		return 0, ErrInvalidType
	}
	return 0, ErrInvalidKey
}

func (b *BlackBoard) SetValueAsFloat32(key int, v float32) {
	b.board[key] = v
}

func (b *BlackBoard) GetValueAsFloat64(key int) (float64, error) {
	if v, ok := b.board[key]; ok {
		if i, ok := v.(float64); ok {
			return i, nil
		}
		return 0, ErrInvalidType
	}
	return 0, ErrInvalidKey
}

func (b *BlackBoard) SetValueAsFloat64(key int, v float64) {
	b.board[key] = v
}

func (b *BlackBoard) GetValueAsString(key int) (string, error) {
	if v, ok := b.board[key]; ok {
		if i, ok := v.(string); ok {
			return i, nil
		}
		return "", ErrInvalidType
	}
	return "", ErrInvalidKey
}

func (b *BlackBoard) SetValueAsString(key int, v string) {
	b.board[key] = v
}

func (b *BlackBoard) GetValueAsInterface(key int) (interface{}, error) {
	if v, ok := b.board[key]; ok {
		return v, nil
	}
	return nil, ErrInvalidKey
}

func (b *BlackBoard) SetValueAsInterface(key int, v interface{}) {
	b.board[key] = v
}
