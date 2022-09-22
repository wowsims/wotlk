package neat

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/wowsims/wotlk/sim/core"
)

type Genome struct {
	Connections map[int]*Connection
	Nodes       map[int]*Node

	NumInputs  int
	NumOutputs int
}

func NewGenome() *Genome {
	g := &Genome{}
	g.Connections = make(map[int]*Connection)
	g.Nodes = make(map[int]*Node)
	return g
}

func NewGenomeFromFile(path string) *Genome {
	file, err := os.Open(path)
	if err == nil {
		g := NewGenome()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var nodeKey int = 0
		for scanner.Scan() {
			var nodes []*Node = make([]*Node, 0)
			var nodeDecl bool = false
			var nodeKind NodeKind
			var connectionDecl bool = false
			var connectionDataIdx int = 0
			var connectionData [4]float64
			var connectionExpressed bool = false
			for _, token := range strings.Fields(scanner.Text()) {
				if strings.Contains(token, "in") {
					nodeDecl = true
					connectionDecl = false
					nodeKind = NodeKind_Input
				} else if strings.Contains(token, "hidden") {
					nodeDecl = true
					connectionDecl = false
					nodeKind = NodeKind_Hidden
				} else if strings.Contains(token, "out") {
					nodeDecl = true
					connectionDecl = false
					nodeKind = NodeKind_Output
				} else if strings.Contains(token, "connection") {
					connectionDecl = true
					connectionDataIdx = 0
					nodeDecl = false
				} else if strings.Contains(token, "t") {
					connectionExpressed = true
				} else if strings.Contains(token, "f") {
					connectionExpressed = false
				} else if val, err := strconv.Atoi(token); err == nil {
					if nodeDecl {
						switch nodeKind {
						case NodeKind_Input:
							for i := 0; i < val; i++ {
								n := NewNode(NodeKind_Input, i+nodeKey)
								nodes = append(nodes, n)
								g.AddNode(n)
							}
							g.NumInputs = val
							nodeKey += val
						case NodeKind_Hidden:
							for i := 0; i < val; i++ {
								n := NewNode(NodeKind_Hidden, i+nodeKey)
								nodes = append(nodes, n)
								g.AddNode(n)
							}
							nodeKey += val
						case NodeKind_Output:
							for i := 0; i < val; i++ {
								n := NewNode(NodeKind_Output, i+nodeKey)
								nodes = append(nodes, n)
								g.AddNode(n)
							}
							g.NumOutputs = val
							nodeKey += val
						}

						nodeDecl = false
					} else if connectionDecl {
						connectionData[connectionDataIdx] = float64(val)
						connectionDataIdx++
					}
				} else if val, err := strconv.ParseFloat(token, 32); err == nil {
					if connectionDecl {
						connectionData[connectionDataIdx] = val
						connectionDataIdx++
					}
				}
			}

			if connectionDecl {
				g.AddConnection(NewConnection(int(connectionData[0]), int(connectionData[1]), connectionData[2], connectionExpressed, int(connectionData[3])))
			}
		}

		return g
	}

	return nil
}

func (g *Genome) AddConnection(c *Connection) {
	g.Connections[c.Innovation] = c
}

func (g *Genome) AddNode(n *Node) {
	g.Nodes[n.Id] = n
}

func (g *Genome) Mutation() {
	for _, c := range g.Connections {
		if rand.Float64() <= 0.9 {
			c.Weight *= rand.Float64()*4.0 - 2.0
		} else {
			c.Weight = rand.Float64()*4.0 - 2.0
		}
	}
}

func (g *Genome) AddConnectionMut(innov *Innovator) {
	ln := len(g.Nodes)
	n1 := g.Nodes[rand.Intn(ln)]
	n2 := g.Nodes[rand.Intn(ln)]

	reversed := (n1.IsHidden() && n2.IsInput() || n1.IsOutput() && n2.IsHidden() || n1.IsOutput() && n2.IsInput())

	for _, c := range g.Connections {
		if (c.In == n1.Id && c.Out == n2.Id) || (c.In == n2.Id && c.Out == n1.Id) {
			return
		}
	}

	a := core.TernaryInt(reversed, n2.Id, n1.Id)
	b := core.TernaryInt(reversed, n1.Id, n2.Id)
	g.AddConnection(NewConnection(a, b, rand.Float64()*2.0-1.0, true, innov.Next()))
}

func (g *Genome) AddNodeMut(innov *Innovator) {
	lc := len(g.Connections)
	c := g.Connections[rand.Intn(lc)]

	n1 := g.Nodes[c.In]
	n2 := g.Nodes[c.Out]

	c.Disable()

	nn := NewNode(NodeKind_Hidden, len(g.Nodes))
	inToNew := NewConnection(n1.Id, nn.Id, 1.0, true, innov.Next())
	newToOut := NewConnection(nn.Id, n2.Id, c.Weight, true, innov.Next())

	g.AddNode(nn)
	g.AddConnection(inToNew)
	g.AddConnection(newToOut)
}

func Cross(g1 *Genome, g2 *Genome) *Genome {
	child := NewGenome()

	for _, g1n := range g1.Nodes {
		child.AddNode(g1n.Copy())
	}

	for _, g1c := range g1.Connections {
		_, exists := g2.Connections[g1c.Innovation]
		if exists {
			childCon := (*Connection)(nil)
			if rand.Intn(2) == 1 {
				childCon = g1c.Copy()
			} else {
				childCon = g2.Connections[g1c.Innovation].Copy()
			}
			child.AddConnection(childCon)
		} else {
			child.AddConnection(g1c.Copy())
		}
	}

	return child
}

func (g *Genome) SortedNodeKeys() []int {
	keys := make([]int, 0, len(g.Nodes))
	i := 0
	for k := range g.Nodes {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func (g *Genome) SortedConnectionKeys() []int {
	keys := make([]int, 0, len(g.Connections))
	i := 0
	for k := range g.Connections {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinIntSlice(slice []int) int {
	if len(slice) == 0 {
		panic("No elements in slice!")
	}

	result := slice[0]
	for _, v := range slice {
		if v <= result {
			result = v
		}
	}
	return result
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxIntSlice(slice []int) int {
	if len(slice) == 0 {
		panic("No elements in slice!")
	}

	result := slice[0]
	for _, v := range slice {
		if v >= result {
			result = v
		}
	}
	return result
}

func CountMatchingGenes(g1 *Genome, g2 *Genome) int {
	matching := 0

	nk1 := g1.SortedNodeKeys()
	nk2 := g2.SortedNodeKeys()
	indices := MaxInt(MaxIntSlice(nk1), MaxIntSlice(nk2))

	for i := 0; i <= indices; i++ {
		_, e1 := g1.Nodes[i]
		_, e2 := g2.Nodes[i]
		if e1 && e2 {
			matching++
		}
	}

	ck1 := g1.SortedConnectionKeys()
	ck2 := g2.SortedConnectionKeys()
	indices = MaxInt(MaxIntSlice(ck1), MaxIntSlice(ck2))

	for i := 0; i <= indices; i++ {
		_, e1 := g1.Connections[i]
		_, e2 := g2.Connections[i]
		if e1 && e2 {
			matching++
		}
	}

	return matching
}

func CountExcessGenes(g1 *Genome, g2 *Genome) int {
	return 0
}

func CountDisjointGenes(g1 *Genome, g2 *Genome) int {
	return 0
}

func AverageWeightDiff(g1 *Genome, g2 *Genome) float64 {
	weightDiff := 0.0
	matching := 0

	ck1 := g1.SortedConnectionKeys()
	ck2 := g2.SortedConnectionKeys()
	indices := MaxInt(MaxIntSlice(ck1), MaxIntSlice(ck2))

	for i := 0; i <= indices; i++ {
		c1, e1 := g1.Connections[i]
		c2, e2 := g2.Connections[i]
		if e1 && e2 {
			matching++
			weightDiff += math.Abs(c1.Weight - c2.Weight)
		}
	}

	return weightDiff / float64(matching)
}

func CompatibilityDistance(g1 *Genome, g2 *Genome, c1 float64, c2 float64, c3 float64) float64 {
	excessGenes := CountExcessGenes(g1, g2)
	disjointGenes := CountDisjointGenes(g1, g2)
	averageWeightDiff := AverageWeightDiff(g1, g2)
	return float64(excessGenes)*c1 + float64(disjointGenes)*c2 + averageWeightDiff*c3
}

func Activation(value float64) float64 {
	return 1.0 / (1.0 + math.Exp(-value))
}

func (g *Genome) Evaluate(inputs []float64) (int, []float64) {
	for i, in := range inputs {
		g.Nodes[i].Output = in
	}

	for _, n := range g.Nodes {
		if n.IsInput() {
			continue
		}

		s := 0.0
		for _, c := range g.Connections {
			if c.Out == n.Id && c.Expressed {
				s += c.Weight * g.Nodes[c.In].Output
			}
		}
		n.Output = Activation(s)
	}

	outIndicesIdx := 0
	outIndices := make([]int, g.NumOutputs)
	for _, n := range g.Nodes {
		if n.IsOutput() {
			outIndices[outIndicesIdx] = n.Id
			outIndicesIdx++
		}
	}

	sort.Ints(outIndices)

	outMax := 0.0
	outMaxIdx := 0
	outIndicesIdx = 0
	out := make([]float64, g.NumOutputs)
	for _, i := range outIndices {
		out[outIndicesIdx] = g.Nodes[i].Output

		if out[outIndicesIdx] >= outMax {
			outMax = out[outIndicesIdx]
			outMaxIdx = outIndicesIdx
		}

		outIndicesIdx++
	}

	return outMaxIdx, out
}

func (g *Genome) Print() {
	for _, n := range g.Nodes {
		n.Print()
	}

	for _, c := range g.Connections {
		c.Print()
	}
}
