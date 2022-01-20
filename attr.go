
package cmdtree

import (
	math "math"
	strconv "strconv"
)

type AttrNode0 struct{
	Node0
	key string
}

type IntegerNode struct{
	AttrNode0
	min, max int64
}

func Integer(key string, mm ...int64)(n *IntegerNode){
	n = new(IntegerNode)
	n.ins = n
	n.key = key
	if len(mm) > 0 {
		n.min = mm[0]
	}else{
		n.min = math.MinInt64
	}
	if len(mm) > 1 {
		n.max = mm[1]
	}else{
		n.max = math.MaxInt64
	}
	return
}

func (a *IntegerNode)AtMin(m int64)(*IntegerNode){
	a.min = m
	return a
}

func (a *IntegerNode)AtMax(m int64)(*IntegerNode){
	a.max = m
	return a
}

func (a *IntegerNode)ParseI(args ArgMap, cmd string)(_ string, err error){
	var (
		s string
		n int64
	)
	s, cmd = splitNode(cmd)
	n, err = strconv.ParseInt(s, 10, 64)
	if err != nil { return }
	if n < a.min || n > a.max {
		return "", &ErrArgOutOfRange{a.min, a.max}
	}
	args[a.key] = n
	return cmd, nil
}

type FloatNode struct{
	AttrNode0
	min, max float64
}

func Float(key string, mm ...float64)(n *FloatNode){
	n = new(FloatNode)
	n.ins = n
	n.key = key
	if len(mm) > 0 {
		n.min = mm[0]
	}else{
		n.min = -math.MaxFloat64
	}
	if len(mm) > 1 {
		n.max = mm[1]
	}else{
		n.min = math.MaxFloat64
	}
	return
}

func (a *FloatNode)AtMin(m float64)(*FloatNode){
	a.min = m
	return a
}

func (a *FloatNode)AtMax(m float64)(*FloatNode){
	a.max = m
	return a
}

func (a *FloatNode)ParseI(args ArgMap, cmd string)(_ string, err error){
	var (
		s string
		n float64
	)
	s, cmd = splitNode(cmd)
	n, err = strconv.ParseFloat(s, 64)
	if err != nil { return }
	if n < a.min || n > a.max {
		return "", &ErrArgOutOfRange{a.min, a.max}
	}
	args[a.key] = n
	return cmd, nil
}

type BooleanNode struct{
	AttrNode0
}

func Boolean(key string)(n *BooleanNode){
	n = new(BooleanNode)
	n.ins = n
	n.key = key
	return
}

func (a *BooleanNode)ParseI(args ArgMap, cmd string)(_ string, err error){
	var s string
	s, cmd = splitNode(cmd)
	switch s {
	case "T", "True", "TRUE", "true":
		args[a.key] = true
	case "F", "False", "FALSE", "false":
		args[a.key] = false
	default:
		return "", ErrUnknownArg
	}
	return cmd, nil
}

type TextNode struct{
	AttrNode0
	min, max int
}

func Text(key string, mm ...int)(n *TextNode){
	n = new(TextNode)
	n.ins = n
	n.key = key
	if len(mm) > 0 {
		n.min = mm[0]
	}
	if len(mm) > 1 {
		n.max = mm[1]
	}
	return
}

func (a *TextNode)MinLen(m int)(*TextNode){
	a.min = m
	return a
}

func (a *TextNode)MaxLen(m int)(*TextNode){
	a.max = m
	return a
}

func (a *TextNode)ParseI(args ArgMap, cmd string)(_ string, err error){
	var s string
	s, cmd = splitNode(cmd)
	if (a.min > 0 && len(s) < a.min) || (a.max > 0 && len(s) > a.max) {
		return "", ErrUnknownArg
	}
	args[a.key] = s
	return cmd, nil
}

type GreedyTextNode struct{
	AttrNode0
	min, max int
}

func GreedyText(key string, mm ...int)(n *GreedyTextNode){
	n = new(GreedyTextNode)
	n.ins = n
	n.key = key
	if len(mm) > 0 {
		n.min = mm[0]
	}
	if len(mm) > 1 {
		n.max = mm[1]
	}
	return
}

func (a *GreedyTextNode)MinLen(m int)(*GreedyTextNode){
	a.min = m
	return a
}

func (a *GreedyTextNode)MaxLen(m int)(*GreedyTextNode){
	a.max = m
	return a
}

func (a *GreedyTextNode)ParseI(args ArgMap, cmd string)(_ string, err error){
	if (a.min > 0 && len(cmd) < a.min) || (a.max > 0 && len(cmd) > a.max) {
		return "", ErrUnknownArg
	}
	args[a.key] = cmd
	return "", nil
}
