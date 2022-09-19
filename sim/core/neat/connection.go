package neat

import "fmt"

type Innovator struct {
	idx int
}

func (c *Innovator) Next() int {
	c.idx = c.idx + 1
	return c.idx - 1
}

type Connection struct {
	In     int
	Out    int
	Weight float64

	Expressed  bool
	Innovation int
}

func NewConnection(in int, out int, weight float64, expressed bool, innovation int) *Connection {
	return &Connection{
		In:     in,
		Out:    out,
		Weight: weight,

		Expressed:  expressed,
		Innovation: innovation,
	}
}

func (c *Connection) Disable() {
	c.Expressed = false
}

func (c *Connection) Copy() *Connection {
	return NewConnection(c.In, c.Out, c.Weight, c.Expressed, c.Innovation)
}

func (c *Connection) Print() {
	fmt.Println("CONNECTION(", c.Innovation, ") ", c.In, " -> ", c.Out, " [ ", c.Weight, " ] ", c.Expressed, " ")
}
