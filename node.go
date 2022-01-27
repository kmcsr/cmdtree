
package cmdtree

import (
	io "io"
	strings "strings"
)

type (
	Executer func(*Context)(error)
	Requester func(*Context)(error)
	Suggester func(ctx *Context, cmd string)(suggestions []string)

	Node interface{
		Run(Executer)(Node)
		Executer()(Executer)
		Then(d Node)(Node)
		Nodes()([]Node)
		Request(Requester)(Node)
		OnSuggest(s Suggester)(Node)

		Parse(ctx *Context)(string, error)
		Suggest(ctx *Context, cmd string)(base string, suggestions []string)
	}
)

type source struct{
	io.Reader
	io.Writer
	attach interface{}
}

func NewSource(r io.Reader, w io.Writer, _attach ...interface{})(Source){
	var attach interface{} = nil
	if len(_attach) > 0 {
		attach = _attach[0]
	}
	return &source{
		Reader: r,
		Writer: w,
		attach: attach,
	}
}

func (s *source)AttachMent()(interface{}){
	return s.attach
}

type (
	Node0Ins interface{
		Node
		ParseI(ctx *Context, cmd string)(remain string, err error)
		SuggestI(ctx *Context, cmd string)(suggestions []string)
	}

	Node0 struct{
		Ins Node0Ins

		exec Executer
		then []Node
		requests []Requester
		suggester Suggester
	}
)

func (n *Node0)Executer()(Executer){
	return n.exec
}

func (n *Node0)Nodes()([]Node){
	return n.then
}

func (n *Node0)Run(e Executer)(Node){
	n.exec = e
	return n.Ins
}

func (n *Node0)Then(d Node)(Node){
	n.then = append(n.then, d)
	return n.Ins
}

func (n *Node0)Request(r Requester)(Node){
	n.requests = append(n.requests, r)
	return n.Ins
}

func (n *Node0)OnSuggest(s Suggester)(Node){
	n.suggester = s
	return n.Ins
}

func (n *Node0)Parse(ctx *Context)(remain string, err error){
	remain, err = n.Ins.ParseI(ctx, ctx.Remain())
	if err != nil { return }
	for _, r := range n.requests {
		err = r(ctx)
		if err != nil {
			return "", &ErrArgRequest{err}
		}
	}
	ctx.SetNexts(n.then...)
	return remain, nil
}

func (n *Node0)Suggest(ctx *Context, cmd string)(base string, suggestions []string){
	var (
		remain string
		err error
	)
	remain, err = n.Ins.ParseI(ctx, cmd)
	if len(remain) > 0 {
		for _, r := range n.requests {
			if err = r(ctx); err != nil { return }
		}
		b0 := strings.TrimSuffix(cmd, remain)
		remain = trimLeftStr(remain)
		for _, m := range n.then {
			bs, sg := m.Suggest(ctx.Clone(), remain)
			if len(sg) > 0 {
				return b0 + " " + bs, sg
			}
		}
	}else{
		if n.suggester != nil {
			return "", n.suggester(ctx, cmd)
		}
		return "", n.Ins.SuggestI(ctx, cmd)
	}
	return
}

func (n *Node0)SuggestI(*Context, string)([]string){
	return nil
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
	n.Ins = n
	return
}

func (n *LiteralNode)Names()([]string){
	return n.names
}

func (n *LiteralNode)ParseI(ctx *Context, cmd string)(remain string, err error){
	var s string
	s, remain = SplitNode(cmd)
	if !strInList(s, n.names) {
		return "", &ErrUnknownArg{cmd}
	}
	return
}

func (n *LiteralNode)SuggestI(ctx *Context, cmd string)(suggestions []string){
	for _, k := range n.names {
		if strings.HasPrefix(k, cmd) {
			suggestions = append(suggestions, k)
		}
	}
	return
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

func (n *RootNode)Suggest(ctx *Context, cmd string)(base string, suggestions []string){
	var (
		s string
		remain string
	)
	s, remain = SplitNode(cmd)
	if l, ok := n.literals[s]; ok && len(remain) > 0 {
		return l.Suggest(ctx.Clone(), cmd)
	}else{
		for k, _ := range n.literals {
			if strings.HasPrefix(k, s) {
				suggestions = append(suggestions, k)
			}
		}
	}
	return
}

func (n *RootNode)Parse(ctx *Context)(l *LiteralNode, err error){
	var s string
	s, _ = SplitNode(ctx.Remain())
	l, ok := n.literals[s]
	if !ok {
		return nil, &ErrUnknownCommand{s}
	}
	return l, nil
}
