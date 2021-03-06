// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package opt

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/types"
)

//go:generate optgen -out operator.og.go ops ops/scalar.opt ops/relational.opt ops/enforcer.opt

// Operator describes the type of operation that a memo expression performs.
// Some operators are relational (join, select, project) and others are scalar
// (and, or, plus, variable).
type Operator uint16

// String returns the name of the operator as a string.
func (i Operator) String() string {
	if i >= Operator(len(opNames)-1) {
		return fmt.Sprintf("Operator(%d)", i)
	}

	return opNames[opIndexes[i]:opIndexes[i+1]]
}

// FuncOpDef defines the value of the Def private field of the Function
// operator. It provides the name and return type of the function, as well as a
// pointer to an already resolved builtin overload definition.
type FuncOpDef struct {
	Name     string
	Type     types.T
	Overload *tree.Builtin
}

func (f FuncOpDef) String() string {
	return f.Name
}

// ScanOpDef defines the value of the Def private field of the Scan operator.
type ScanOpDef struct {
	// Table identifies the table to scan. It is an index that can be passed to
	// the Metadata.Table method in order to fetch optbase.Table metadata.
	Table TableIndex

	// Cols specifies the set of columns that the scan operator projects. This
	// may be a subset of the columns that the table contains.
	Cols ColSet
}

// SetOpColMap defines the value of the ColMap private field of the set
// operators: Union, Intersect, Except, UnionAll, IntersectAll and ExceptAll.
// It matches columns from the left and right inputs of the operator
// with the output columns, since OutputCols are not ordered and may
// not correspond to each other.
//
// For example, consider the following query:
//   SELECT y, x FROM xy UNION SELECT b, a FROM ab
//
// Given:
//   col  index
//   x    1
//   y    2
//   a    3
//   b    4
//
// SetOpColMap will contain the following values:
//   Left:  [2, 1]
//   Right: [4, 3]
//   Out:   [5, 6]  <-- synthesized output columns
type SetOpColMap struct {
	Left  ColList
	Right ColList
	Out   ColList
}

// ComparisonOpReverseMap maps from an optimizer operator type to a semantic
// tree comparison operator type.
var ComparisonOpReverseMap = [...]tree.ComparisonOperator{
	EqOp:           tree.EQ,
	LtOp:           tree.LT,
	GtOp:           tree.GT,
	LeOp:           tree.LE,
	GeOp:           tree.GE,
	NeOp:           tree.NE,
	InOp:           tree.In,
	NotInOp:        tree.NotIn,
	LikeOp:         tree.Like,
	NotLikeOp:      tree.NotLike,
	ILikeOp:        tree.ILike,
	NotILikeOp:     tree.NotILike,
	SimilarToOp:    tree.SimilarTo,
	NotSimilarToOp: tree.NotSimilarTo,
	RegMatchOp:     tree.RegMatch,
	NotRegMatchOp:  tree.NotRegMatch,
	RegIMatchOp:    tree.RegIMatch,
	NotRegIMatchOp: tree.NotRegIMatch,
	IsOp:           tree.IsNotDistinctFrom,
	IsNotOp:        tree.IsDistinctFrom,
	ContainsOp:     tree.Contains,
}

// BinaryOpReverseMap maps from an optimizer operator type to a semantic tree
// binary operator type.
var BinaryOpReverseMap = [...]tree.BinaryOperator{
	BitandOp:        tree.Bitand,
	BitorOp:         tree.Bitor,
	BitxorOp:        tree.Bitxor,
	PlusOp:          tree.Plus,
	MinusOp:         tree.Minus,
	MultOp:          tree.Mult,
	DivOp:           tree.Div,
	FloorDivOp:      tree.FloorDiv,
	ModOp:           tree.Mod,
	PowOp:           tree.Pow,
	ConcatOp:        tree.Concat,
	LShiftOp:        tree.LShift,
	RShiftOp:        tree.RShift,
	FetchValOp:      tree.JSONFetchVal,
	FetchTextOp:     tree.JSONFetchText,
	FetchValPathOp:  tree.JSONFetchValPath,
	FetchTextPathOp: tree.JSONFetchTextPath,
}

// UnaryOpReverseMap maps from an optimizer operator type to a semantic tree
// unary operator type.
var UnaryOpReverseMap = [...]tree.UnaryOperator{
	UnaryMinusOp:      tree.UnaryMinus,
	UnaryComplementOp: tree.UnaryComplement,
}
