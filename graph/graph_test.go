package graph

import (
	"fmt"
	"testing"
)

func TestAddNode(t *testing.T) {
	g := NewGraph()

	url := "https://example.com"

	g.AddNode(url)

	if !g.HasNode(url) {
		t.Fatalf("expected graph to have node %v", url)
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph()

	startURL := "https://example.com"
	neighborURL := "https://example.com/savings"

	g.AddNode(startURL)
	g.AddNode(neighborURL)

	err := g.AddEdge(startURL, neighborURL)

	if err != nil {
		t.Fatalf("add edge error: %v", err)
	}
}

func TestAddEdgeError(t *testing.T) {
	g := NewGraph()

	startURL := "https://example.com"
	neighborURL := "https://example.com/cards"

	g.AddNode(startURL)

	err := g.AddEdge(startURL, neighborURL)

	if err == nil {
		t.Fatalf("expected add edge operation to fail")
	}

	got := err.Error()
	expected := fmt.Sprintf("failed to add edge no node found for %v", neighborURL)

	if got != expected {
		t.Fatalf("expected error message %v, got %v", expected, got)
	}
}
