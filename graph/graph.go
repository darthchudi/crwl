// graph provides a data structure for storing
// visited URLs (nodes) and the links between them
package graph

import (
	"fmt"
	"sync"
)

type Node struct {
	// URL of the Node
	url string
}

type Graph struct {
	// nodes is our collection of visited URLs
	nodes map[string]*Node

	// edges represents the link between nodes
	edges map[Node][]*Node

	// mu is a mutex that protects our Graph for concurrent use
	mu sync.RWMutex
}

// NewGraph initializes a new Graph
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		edges: make(map[Node][]*Node),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(url string) {
	if g.HasNode(url) {
		return
	}

	node := &Node{url}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.nodes[url] = node
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(startURL, endURL string) error {
	g.mu.RLock()
	startNode := g.nodes[startURL]
	endNode := g.nodes[endURL]
	g.mu.RUnlock()

	if startNode == nil {
		return fmt.Errorf("failed to add edge, no node found for %v", startURL)
	}

	if endNode == nil {
		return fmt.Errorf("failed to add edge no node found for %v", endURL)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.edges[*startNode] = append(g.edges[*startNode], endNode)
	return nil
}

// HasNode checks if a node with a particular URL exists in the graph
func (g *Graph) HasNode(url string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.nodes[url]

	return exists
}
