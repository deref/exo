package deps

import (
	"errors"
	"fmt"
	"strings"
)

type Node interface {
	ID() string
}

type nodeset map[string]Node

// The value is a simple slice of strings. This assumes that the number of dependencies for any given node
// is small enough that a linear search is usually going to be at least as efficient as a lookup in a hash
// table, and it has less memory overhead.
type depmap map[string][]string

type Graph struct {
	// Maintain dependency relationships in both directions.
	// `dependencies` tracks child -> parents, and `dependents` tracks parent -> children.
	dependencies, dependents depmap
	nodes                    nodeset
}

func New() *Graph {
	return &Graph{
		dependencies: make(depmap),
		dependents:   make(depmap),
		nodes:        make(nodeset),
	}
}

func (g *Graph) Dump() string {
	var out strings.Builder
	out.WriteString("Nodes:\n")
	for node := range g.nodes {
		fmt.Fprintf(&out, "\t%v\n", node)
	}

	out.WriteString("Dependencies:\n")
	for node, deps := range g.dependencies {
		fmt.Fprintf(&out, "\t%v <-", node)
		for _, dep := range deps {
			fmt.Fprintf(&out, " %v", dep)
		}
		out.WriteByte('\n')
	}

	out.WriteString("Dependents:\n")
	for node, deps := range g.dependents {
		fmt.Fprintf(&out, "\t%v ->", node)
		for _, dep := range deps {
			fmt.Fprintf(&out, " %v", dep)
		}
		out.WriteByte('\n')
	}

	return out.String()
}

// AddNode will register some node with the graph even if that node does not depend on
// anything else.
func (g *Graph) AddNode(node Node) {
	g.nodes[node.ID()] = node
}

func (g *Graph) DependOn(node, dep Node) error {
	if node.ID() == dep.ID() {
		return errors.New("self-referential dependencies not allowed")
	}

	if err := g.AddEdge(node.ID(), dep.ID()); err != nil {
		return err
	}

	g.AddNode(node)
	g.AddNode(dep)

	return nil
}

// AddEdge registers a dependency between the node identified by parentID and the node identified by childID.
// This assumes that the nodes will be added separately via AddNode() or DependOn().
func (g *Graph) AddEdge(parentID, childID string) error {
	if g.dependsOn(childID, parentID) {
		return errors.New("circular dependencies not allowed")
	}

	updateNodeIDs(g.dependencies, parentID, func(ids []string) []string {
		return addID(ids, childID)
	})
	updateNodeIDs(g.dependents, childID, func(ids []string) []string {
		return addID(ids, parentID)
	})

	return nil
}

func (g *Graph) DependsOn(node, dep Node) bool {
	return g.dependsOn(node.ID(), dep.ID())
}

func (g *Graph) dependsOn(parentID, childID string) bool {
	tds := g.Dependencies(parentID)
	_, ok := tds[childID]
	return ok
}

func (g *Graph) HasDependent(node, dep Node) bool {
	return g.hasDependent(node.ID(), dep.ID())
}

func (g *Graph) hasDependent(childID, parendID string) bool {
	tds := g.Dependents(childID)
	_, ok := tds[parendID]
	return ok
}

func (g *Graph) Nodes() []Node {
	out := make([]Node, 0, len(g.nodes))
	for _, node := range g.nodes {
		out = append(out, node)
	}
	return out
}

func (g *Graph) Leaves() []Node {
	out := make([]Node, 0)
	for id, node := range g.nodes {
		ids, ok := g.dependencies[id]
		if !ok {
			out = append(out, node)
		} else {
			// Additionally, if no dependencies exist in the graph, consider this a leaf.
			var foundReference bool
			for _, referencedID := range ids {
				_, foundReference = g.nodes[referencedID]
				if foundReference {
					break
				}
			}
			if !foundReference {
				out = append(out, node)
			}
		}
	}
	return out
}

// TopoSortedLayers returns a slice of all of the graph nodes in topological sort order. That is,
// if `B` depends on `A`, then `A` is guaranteed to come before `B` in the sorted output.
// The graph is guaranteed to be cycle-free because cycles are detected while building the
// graph. Additionally, the output is grouped into "layers", which are guaranteed to not have
// any dependencies within each layer. This is useful, e.g. when building an execution plan for
// some DAG, in which case each element within each layer could be executed in parallel. If you
// do not need this layered property, use `Graph.TopoSorted()`, which flattens all elements.
func (g *Graph) TopoSortedLayers() [][]Node {
	out := [][]Node{}

	shrinkingGraph := g.clone()
	for {
		leaves := shrinkingGraph.Leaves()
		if len(leaves) == 0 {
			break
		}

		out = append(out, leaves)
		for _, leafNode := range leaves {

			dependents := shrinkingGraph.dependents[leafNode.ID()]

			for _, dependentID := range dependents {
				// Should be safe because every relationship is bidirectional.
				dependencies := shrinkingGraph.dependencies[dependentID]
				if len(dependencies) == 1 {
					// The only dependent _must_ be `leafNode`, so we can delete the `dep` entry entirely.
					delete(shrinkingGraph.dependencies, dependentID)
				} else {
					shrinkingGraph.dependencies[dependentID] = removeID(dependencies, leafNode.ID())
				}
			}
			delete(shrinkingGraph.dependents, leafNode.ID())
		}

		nextLeaves := shrinkingGraph.Leaves()
		// nodes must be removed after the next iteration's leaves have been evaluated so that we do not
		// delete the last layer's elements before the last iteration.
		for _, leafNode := range leaves {
			delete(shrinkingGraph.nodes, leafNode.ID())
		}
		leaves = nextLeaves
	}

	return out
}

// TopoSorted returns all the nodes in the graph is topological sort order.
// See also `Graph.TopoSortedLayers()`.
func (g *Graph) TopoSorted() []Node {
	nodeCount := 0
	layers := g.TopoSortedLayers()
	for _, layer := range layers {
		nodeCount += len(layer)
	}

	allNodes := make([]Node, 0, nodeCount)
	for _, layer := range layers {
		for _, node := range layer {
			allNodes = append(allNodes, node)
		}
	}

	return allNodes
}

func (g *Graph) Dependencies(id string) nodeset {
	return g.buildTransitive(id, g.immediateDependencies)
}

func (g *Graph) immediateDependencies(node Node) nodeset {
	if depIDs, ok := g.dependencies[node.ID()]; ok {
		out := make(nodeset)
		for _, nodeID := range depIDs {
			if node, ok := g.nodes[nodeID]; ok {
				out[nodeID] = node
			}
		}
		return out
	}
	return nil
}

func (g *Graph) Dependents(id string) nodeset {
	return g.buildTransitive(id, g.immediateDependents)
}

func (g *Graph) immediateDependents(node Node) nodeset {
	if depIDs, ok := g.dependents[node.ID()]; ok {
		out := make(nodeset)
		for _, nodeID := range depIDs {
			if node, ok := g.nodes[nodeID]; ok {
				out[nodeID] = node
			}
		}
		return out
	}
	return nil
}

func (g *Graph) clone() *Graph {
	return &Graph{
		dependencies: copyDepmap(g.dependencies),
		dependents:   copyDepmap(g.dependents),
		nodes:        copyNodeset(g.nodes),
	}
}

// buildTransitive starts at `root` and continues calling `nextFn` to keep discovering more nodes until
// the graph cannot produce any more. It returns the set of all discovered nodes.
func (g *Graph) buildTransitive(rootID string, nextFn func(Node) nodeset) nodeset {
	root, ok := g.nodes[rootID]
	if !ok {
		return nil
	}

	out := make(nodeset)
	searchNext := []Node{root}
	for len(searchNext) > 0 {
		// List of new nodes from this layer of the dependency graph. This is
		// assigned to `searchNext` at the end of the outer "discovery" loop.
		discovered := []Node{}
		for _, node := range searchNext {
			// For each node to discover, find the next nodes.
			for id, nextNode := range nextFn(node) {
				// If we have not seen the node before, add it to the output as well
				// as the list of nodes to traverse in the next iteration.
				if _, ok := out[id]; !ok {
					out[id] = nextNode
					discovered = append(discovered, nextNode)
				}
			}
		}
		searchNext = discovered
	}
	return out
}

func copyNodeset(s nodeset) nodeset {
	out := make(nodeset, len(s))
	for k, v := range s {
		out[k] = v
	}
	return out
}

func copyDepmap(m depmap) depmap {
	out := make(depmap, len(m))
	for k, v := range m {
		out[k] = make([]string, len(v))
		copy(out[k], v)
	}
	return out
}

type updateFn = func(nodeIDs []string) []string

func updateNodeIDs(ds depmap, id string, fn updateFn) {
	nodeIDs, ok := ds[id]
	if !ok {
		nodeIDs = []string{}
		ds[id] = nodeIDs
	}
	ds[id] = fn(nodeIDs)
}

func addID(ids []string, id string) []string {
	for _, idElem := range ids {
		if id == idElem {
			return ids
		}
	}
	return append(ids, id)
}

func removeID(ids []string, id string) []string {
	idx := -1
	for i, idElem := range ids {
		if id == idElem {
			idx = i
			break
		}
	}
	if idx == -1 {
		// Not found.
		return ids
	}

	ids[idx] = ids[len(ids)-1]
	return ids[:len(ids)-1]
}
