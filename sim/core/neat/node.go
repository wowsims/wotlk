package neat

import "fmt"

type NodeKind uint32

const (
	NodeKind_Input NodeKind = iota
	NodeKind_Hidden
	NodeKind_Output
)

type Node struct {
	Id   int
	Kind NodeKind
}

func NewNode(kind NodeKind, id int) *Node {
	return &Node{
		Id:   id,
		Kind: kind,
	}
}

func (n *Node) Copy() *Node {
	return NewNode(n.Kind, n.Id)
}

func (n *Node) IsInput() bool {
	return n.Kind == NodeKind_Input
}

func (n *Node) IsHidden() bool {
	return n.Kind == NodeKind_Hidden
}

func (n *Node) IsOutput() bool {
	return n.Kind == NodeKind_Output
}

func (n *Node) Print() {
	kindStr := " INPUT"
	if n.IsHidden() {
		kindStr = " HIDDEN"
	} else if n.IsOutput() {
		kindStr = " OUTPUT"
	}
	fmt.Println("NODE ", n.Id, kindStr)
}
