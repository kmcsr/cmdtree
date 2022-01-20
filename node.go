
package cmdtree

import (
	io "io"
)

type (
	Source interface{
		io.Reader
		io.Writer
	}
	ArgMap map[string]interface{}
	Executer func(Source, ArgMap)(error)

	Node interface{
		Run(Executer)(Node)
		Then(d Node)(Node)

		Parse(args ArgMap, cmd string)(Executer, error)
	}
)

type Node0 struct{
	ins interface{
		Node
		ParseI(args ArgMap, cmd string)(l string, err error)
	}

	exec Executer
	then []Node
}

func (n *Node0)Run(e Executer)(Node){
	n.exec = e
	return n.ins
}

func (n *Node0)Then(d Node)(Node){
	n.then = append(n.then, d)
	return n.ins
}

func (n *Node0)Parse(args ArgMap, cmd string)(exec Executer, err error){
	cmd, err = n.ins.ParseI(args, cmd)
	if err != nil { return }
	if cmd = trimLeftStr(cmd); len(cmd) == 0 {
		if n.exec == nil {
			return nil, ErrUnknownArg
		}
		return n.exec, nil
	}
	if len(n.then) == 0 {
		return nil, ErrUnknownArg
	}
	if len(n.then) == 1 {
		return n.then[0].Parse(args, cmd)
	}
	for _, t := range n.then {
		exec, err = t.Parse(args, cmd)
		if err == nil {
			return
		}
	}
	return nil, ErrUnknownArg
}

type LiteralNode struct{
	Node0
	names []string
}

func Literal(names ...string)(n *LiteralNode){
	if len(names) == 0 {
		panic("Literal node at least must have one name")
	}
	n = &LiteralNode{
		names: names,
	}
	n.ins = n
	return
}

func (n *LiteralNode)Names()([]string){
	return n.names
}

func (n *LiteralNode)ParseI(_ ArgMap, cmd string)(_ string, err error){
	var s string
	s, cmd = splitNode(cmd)
	if !strInList(s, n.names) {
		return "", ErrUnknownArg
	}
	return cmd, nil
}

type RootNode struct{
	literals map[string]*LiteralNode
}

func NewRoot()(*RootNode){
	return &RootNode{
		literals: make(map[string]*LiteralNode),
	}
}

func (n *RootNode)Then(d Node)(*RootNode){
	l := d.(*LiteralNode)
	for _, s := range l.Names() {
		n.literals[s] = l
	}
	return n
}

func (n *RootNode)Parse(args ArgMap, cmd string)(exec Executer, err error){
	var s string
	s, cmd = splitNode(cmd)
	l, ok := n.literals[s]
	if !ok {
		return nil, ErrUnknownArg
	}
	if cmd = trimLeftStr(cmd); len(cmd) == 0 {
		if l.exec == nil {
			return nil, ErrUnknownArg
		}
		return l.exec, nil
	}
	if len(l.then) == 0 {
		return nil, ErrUnknownArg
	}
	if len(l.then) == 1 {
		return l.then[0].Parse(args, cmd)
	}
	for _, t := range l.then {
		exec, err = t.Parse(args, cmd)
		if err == nil {
			return
		}
	}
	return nil, ErrUnknownArg
}

func (n *RootNode)Execute(src Source, cmd string)(err error){
	var (
		args ArgMap = make(ArgMap)
		exec Executer
	)
	exec, err = n.Parse(args, cmd)
	if err != nil { return }
	return exec(src, args)
}
