package ir

import (
	"fmt"
)

const (
	jmpAlways = iota
	jmpIf
	jmpIfNot
)

type blockJump struct {
	typ  int
	cond Ref
	to   *Block
}

// Block is a basic block
type Block struct {
	id  int // basic block ida
	ops []op

	insts    []*inst
	jumpInst *inst
	spMoved  bool

	frameSize *int32

	jump *blockJump

	next *Block // next in the linked list

	instStart int32
	instEnd   int32
}

func (b *Block) String() string { return fmt.Sprintf("B%d", b.id) }

func (b *Block) addOp(op op) { b.ops = append(b.ops, op) }

// Comment adds an IR comment.
func (b *Block) Comment(s string) {
	b.addOp(&comment{s})
}

// Commentf adds an IR comment with a particular printing format.
func (b *Block) Commentf(s string, args ...interface{}) {
	b.Comment(fmt.Sprintf(s, args...))
}

// Arith append an arithmetic operation to the basic block
func (b *Block) Arith(dest Ref, x Ref, op string, y Ref) {
	b.addOp(&arithOp{dest, x, op, y})
}

// Assign appends an assignment operation to the basic block
func (b *Block) Assign(dest Ref, src Ref) {
	b.Arith(dest, nil, "", src)
}

// Zero appends zeroing operation to the basic block
func (b *Block) Zero(dest Ref) {
	b.addOp(&arithOp{dest, nil, "0", nil})
}

// Call appends a function call operation to the basic block
func (b *Block) Call(dests []Ref, f Ref, sig *FuncSig, args ...Ref) {
	argsCopy := make([]Ref, len(args))
	copy(argsCopy, args)
	b.addOp(&callOp{dests, f, sig, argsCopy})
}

// Jump sets the block always jump to the dest block at its end
func (b *Block) Jump(dest *Block) {
	if dest == b.next {
		b.jump = nil
	} else {
		b.jump = &blockJump{jmpAlways, nil, dest}
	}
}

// JumpIfNot sets the block to jump to its natural next when the
// condition is met, and jump to dest when the condition is not met
func (b *Block) JumpIfNot(cond Ref, dest *Block) {
	b.jump = &blockJump{jmpIfNot, cond, dest}
}

// JumpIf sets the block to jump to its natural next when the
// condition is not met, and jump to dest when the condition is met
func (b *Block) JumpIf(cond Ref, dest *Block) {
	b.jump = &blockJump{jmpIf, cond, dest}
}

func (b *Block) inst(i uint32) *inst {
	ret := &inst{inst: i}
	b.insts = append(b.insts, ret)
	return ret
}
