package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Node struct {
	Name     string
	Parent   *Node
	Children []*Node
}

func (n *Node) Orbit(parent *Node) error {
	if n == nil {
		return nil
	}

	if n.Parent != nil {
		return fmt.Errorf("Node %q cannot orbit multiple nodes", n.Name)
	}

	parent.Children = append(parent.Children, n)
	n.Parent = parent

	return nil
}

func (n *Node) Distance() int {
	if n == nil {
		return -1
	}
	return n.Parent.Distance() + 1
}

func (n *Node) Path() []*Node {
	if n == nil {
		return []*Node{}
	}
	return append(n.Parent.Path(), n)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (n *Node) DistanceTo(other *Node) (common *Node, nDist int, otherDist int) {
	nPath := n.Path()
	oPath := other.Path()

	nLen := len(nPath)
	oLen := len(oPath)
	l := min(nLen, oLen)

	nPath, oPath = nPath[:l], oPath[:l]

	for ; l > 0; l-- {
		i := l - 1
		if nPath[i] == oPath[i] {
			return nPath[i], nLen - l, oLen - l
		}
	}

	return nil, -1, -1
}

type NodeMap map[string]*Node

func (m NodeMap) Exists(name string) bool {
	_, ok := m[name]
	return ok
}

func (m NodeMap) Get(name string) *Node {
	if m.Exists(name) {
		return m[name]
	}

	n := &Node{Name: name}
	m[name] = n
	return n
}

func loadMap(filename string) (NodeMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	nodeMap := NodeMap{}

	for _, orbit := range strings.Split(string(data), "\n") {
		if orbit == "" {
			continue
		}

		nodes := strings.Split(orbit, ")")
		if len(nodes) != 2 {
			return nil, fmt.Errorf("Invalid number of nodes: %d", len(nodes))
		}

		parent := nodeMap.Get(nodes[0])
		child := nodeMap.Get(nodes[1])

		if err := child.Orbit(parent); err != nil {
			return nil, err
		}
	}

	return nodeMap, nil
}

func main() {
	nodeMap, err := loadMap("data.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	total := 0

	for _, node := range nodeMap {
		total += node.Distance()
	}

	fmt.Println(total)

	you, santa := nodeMap["YOU"], nodeMap["SAN"]

	_, dYou, dSan := you.Parent.DistanceTo(santa.Parent)
	fmt.Println(dYou + dSan)
}
