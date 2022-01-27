
package cmdtree

import (
	math "math"
	strconv "strconv"
	strings "strings"
)

type AttrNode0 struct{
	Node0
	key string
}

func NewAttrNode0(ins Node0Ins, key string)(*AttrNode0){
	return &AttrNode0{
		Node0: Node0{
			Ins: ins,
		},
		key: key,
	}
}

func (a *AttrNode0)Key()(string){
	return a.key
}

type IntegerNode struct{
	AttrNode0
	min, max int64
}

func Integer(key string, mm ...int64)(n *IntegerNode){
	n = new(IntegerNode)
	n.Ins = n
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

func (a *IntegerNode)ParseI(ctx *Context, cmd string)(remain string, err error){
	var (
		s string
		n int64
	)
	s, remain = SplitNode(cmd)
	n, err = strconv.ParseInt(s, 10, 64)
	if err != nil { return }
	if n < a.min || n > a.max {
		return "", &ErrArgOutOfRange{a.min, a.max}
	}
	ctx.Args()[a.key] = n
	return
}

type FloatNode struct{
	AttrNode0
	min, max float64
}

func Float(key string, mm ...float64)(n *FloatNode){
	n = new(FloatNode)
	n.Ins = n
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

func (a *FloatNode)ParseI(ctx *Context, cmd string)(remain string, err error){
	var (
		s string
		n float64
	)
	s, remain = SplitNode(cmd)
	n, err = strconv.ParseFloat(s, 64)
	if err != nil { return }
	if n < a.min || n > a.max {
		return "", &ErrArgOutOfRange{a.min, a.max}
	}
	ctx.Args()[a.key] = n
	return
}

type BooleanNode struct{
	AttrNode0
}

func Boolean(key string)(n *BooleanNode){
	n = new(BooleanNode)
	n.Ins = n
	n.key = key
	return
}

func (a *BooleanNode)ParseI(ctx *Context, cmd string)(remain string, err error){
	var s string
	s, remain = SplitNode(cmd)
	switch s {
	case "T", "TRUE", "True", "true":
		ctx.Args()[a.key] = true
	case "F", "FALSE", "False", "false":
		ctx.Args()[a.key] = false
	default:
		return "", &ErrUnknownArg{cmd}
	}
	return
}

func (a *BooleanNode)SuggestI(ctx *Context, cmd string)(suggestions []string){
	if len(cmd) == 0 {
		return []string{
			"T", "TRUE", "True", "true",
			"F", "FALSE", "False", "false",
		}
	}
	suggestions = make([]string, 0, 8)
	if cmd == "T" {
		suggestions = append(suggestions, "T")
	}
	if strings.HasPrefix("TRUE", cmd) {
		suggestions = append(suggestions, "TRUE")
	}
	if strings.HasPrefix("True", cmd) {
		suggestions = append(suggestions, "True")
	}
	if strings.HasPrefix("true", cmd) {
		suggestions = append(suggestions, "true")
	}
	if len(suggestions) == 0 {
		if cmd == "F" {
			suggestions = append(suggestions, "F")
		}
		if strings.HasPrefix("FALSE", cmd) {
			suggestions = append(suggestions, "FALSE")
		}
		if strings.HasPrefix("False", cmd) {
			suggestions = append(suggestions, "False")
		}
		if strings.HasPrefix("false", cmd) {
			suggestions = append(suggestions, "false")
		}
	}
	return
}

type TextNode struct{
	AttrNode0
	min, max int
}

func Text(key string, mm ...int)(n *TextNode){
	n = new(TextNode)
	n.Ins = n
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

func (a *TextNode)ParseI(ctx *Context, cmd string)(remain string, err error){
	var s string
	s, remain = SplitNode(cmd)
	if (a.min > 0 && len(s) < a.min) || (a.max > 0 && len(s) > a.max) {
		return "", &ErrUnknownArg{cmd}
	}
	ctx.Args()[a.key] = s
	return remain, nil
}

type GreedyTextNode struct{
	AttrNode0
	min, max int
}

func GreedyText(key string, mm ...int)(n *GreedyTextNode){
	n = new(GreedyTextNode)
	n.Ins = n
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

func (a *GreedyTextNode)ParseI(ctx *Context, cmd string)(_ string, err error){
	s := cmd
	if (a.min > 0 && len(s) < a.min) || (a.max > 0 && len(s) > a.max) {
		return "", &ErrUnknownArg{s}
	}
	ctx.Args()[a.key] = s
	return "", nil
}
