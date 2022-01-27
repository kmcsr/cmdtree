
package cmdtree

import (
	io "io"
	sync "sync"
)

type (
	Source interface{
		io.Reader
		io.Writer
		AttachMent()(interface{})
	}
	ArgMap map[string]interface{}

	Context struct{
		Source
		command string

		args ArgMap
		remain string
		nexts []Node
		exec Executer
		err error

		origin *Context
		count int
		ctlk chan struct{}
		ctlkc func()
		done chan struct{}
	}
)

func (m ArgMap)Clone()(c ArgMap){
	c = make(ArgMap, len(m))
	for k, v := range m {
		c[k] = v
	}
	return
}

func NewContext(src Source)(ctx *Context){
	return &Context{
		Source: src,
	}
}

func (ctx *Context)Args()(ArgMap){
	return ctx.args
}

func (ctx *Context)Arg(key string)(interface{}){
	return ctx.ArgDefault(key, nil)
}

func (ctx *Context)ArgDefault(key string, def interface{})(interface{}){
	v, ok := ctx.args[key]
	if !ok { return def }
	return v
}

func (ctx *Context)Command()(string){
	return ctx.command
}

func (ctx *Context)Remain()(string){
	return ctx.remain
}

func (ctx *Context)SetNexts(n ...Node){
	ctx.nexts = n
}

func (ctx *Context)Nexts()(n []Node){
	return ctx.nexts
}

func (ctx *Context)Origin()(*Context){
	if ctx.origin == nil {
		return ctx
	}
	return ctx.origin.Origin()
}

func (ctx *Context)Clone()(*Context){
	return &Context{
		Source: ctx.Source,
		command: ctx.command,
		args: ctx.args.Clone(),
		remain: ctx.remain,
		origin: ctx,
	}
}

func (ctx *Context)Begin(command string){
	if ctx.origin != nil {
		panic("cloned context")
	}
	ctx.args = make(ArgMap)
	ctx.command = command
	ctx.remain = command
	ctx.exec = nil
	ctx.count = 0
	ctx.ctlk = make(chan struct{})
	var ctlkc sync.Once
	ctx.ctlkc = func(){
		ctlkc.Do(func(){ close(ctx.ctlk) })
	}
	ctx.done = make(chan struct{})
}

func (ctx *Context)gop(n Node){
	o := ctx.Origin()
	defer func(){
		o.count--
		if o.count == 0 {
			o.ctlkc()
		}
	}()
	select{
	case <-o.done: return
	default:
	}

	var (
		remain string
		err error
	)
	remain, err = n.Parse(ctx)
	select{
	case <-o.done: return
	default:
	}
	if err != nil {
		if _, ok := err.(*ErrUnknownArg); !ok {
			o.err = err
		}
		return
	}
	if ctx.remain = trimLeftStr(remain); len(ctx.remain) == 0 {
		if n.Executer() != nil {
			o.args = ctx.args
			o.remain = ""
			o.exec = n.Executer()
			close(o.done)
		}
	}else{
		for _, n := range ctx.nexts {
			o.count++
			go ctx.Clone().gop(n)
		}
	}
}

func (ctx *Context)Suggest(root *RootNode, command string)(base string, suggestions []string){
	ctx.Begin(command)
	base, suggestions = root.Suggest(ctx, command)
	suggestions = delDuplication(suggestions)
	return
}

func (ctx *Context)Execute(root *RootNode, command string)(err error){
	if len(command) == 0 { return nil }
	ctx.Begin(command)
	var n Node
	n, err = root.Parse(ctx)
	if err != nil { return }
	ctx.count++
	ctx.gop(n)
	select{
	case <-ctx.done:
	case <-ctx.ctlk:
	}
	if ctx.exec != nil {
		return ctx.exec(ctx)
	}
	if ctx.err != nil {
		return ctx.err
	}
	return &ErrUnknownArg{ctx.remain}
}
