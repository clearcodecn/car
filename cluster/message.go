package cluster

type Msg int

func (n *Node) GlobalHandlers(h ...HandlerFunc) {
	n.globalHandlers = append(n.globalHandlers, h...)
}

func (n *Node) OnMessage(msg Msg, h ...HandlerFunc) {
	if _, ok := n.handlers[msg]; !ok {
		n.handlers[msg] = make([]HandlerFunc, 0)
	}
	n.handlers[msg] = append(n.handlers[msg], n.globalHandlers...)
	n.handlers[msg] = append(n.handlers[msg], h...)
}
