package deps

import (
	"bytes"
	"errors"
	"fmt"
)

type nodeset map[interface{}]struct{}

// Semantically, this is Map<Node, Set<Node>>.
type depmap map[interface{}]nodeset

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
	var out bytes.Buffer
	out.WriteString("Nodes:\n")
	for node := range g.dependencies {
		fmt.Fprintf(&out, "\t%v\n", node)
	}

	out.WriteString("Dependencies:\n")
	for node, deps := range g.dependencies {
		fmt.Fprintf(&out, "\t%v <-", node)
		for dep := range deps {
			fmt.Fprintf(&out, " %v", dep)
		}
		out.WriteByte('\n')
	}

	out.WriteString("Dependents:\n")
	for node, deps := range g.dependents {
		fmt.Fprintf(&out, "\t%v ->", node)
		for dep := range deps {
			fmt.Fprintf(&out, " %v", dep)
		}
		out.WriteByte('\n')
	}

	return out.String()
}

func (g *Graph) DependOn(node, dep interface{}) error {
	if node == dep {
		return errors.New("self-referential dependencies not allowed")
	}
	if g.DependsOn(dep, node) {
		return errors.New("circular dependencies not allowed")
	}

	g.nodes[node] = struct{}{}
	g.nodes[dep] = struct{}{}

	updateSet(g.dependencies, node, func(nodes nodeset) {
		nodes[dep] = struct{}{}
	})
	updateSet(g.dependents, dep, func(nodes nodeset) {
		nodes[node] = struct{}{}
	})

	return nil
}

func (g *Graph) DependsOn(node, dep interface{}) bool {
	tds := g.transitiveDependencies(node)
	_, ok := tds[dep]
	return ok
}

func (g *Graph) HasDependent(node, dep interface{}) bool {
	tds := g.transitiveDependents(node)
	_, ok := tds[dep]
	return ok
}

func (g *Graph) Nodes() []interface{} {
	approxSize := len(g.dependencies) + len(g.dependents)/2
	allNodes := make(nodeset, approxSize)
	var nodeCount int
	for node := range g.dependencies {
		nodeCount++
		allNodes[node] = struct{}{}
	}
	for node := range g.dependents {
		if _, exists := allNodes[node]; !exists {
			nodeCount++
			allNodes[node] = struct{}{}
		}
	}

	distinctNodes := make([]interface{}, nodeCount)
	for node := range allNodes {
		distinctNodes[nodeCount-1] = node
		nodeCount--
	}

	return distinctNodes
}

func (g *Graph) Leaves() []interface{} {
	out := make([]interface{}, 0)
	for node := range g.nodes {
		if _, ok := g.dependencies[node]; !ok {
			out = append(out, node)
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
// do not need this layered property, use `Graph.TopoSorted()`, which flattens all elements
func (g *Graph) TopoSortedLayers() [][]interface{} {
	out := [][]interface{}{}

	shrinkingGraph := g.clone()
	for {
		leaves := shrinkingGraph.Leaves()
		if len(leaves) == 0 {
			break
		}

		out = append(out, leaves)
		for _, leafNode := range leaves {

			dependents := shrinkingGraph.dependents[leafNode]

			for dependent := range dependents {
				// Should be safe because every relationship is bidirectional.
				dependencies := shrinkingGraph.dependencies[dependent]
				if len(dependencies) == 1 {
					// The only dependent _must_ be `leafNode`, so we can delete the `dep` entry entirely.
					delete(shrinkingGraph.dependencies, dependent)
				} else {
					delete(dependencies, leafNode)
				}
			}
			delete(shrinkingGraph.dependents, leafNode)
		}

		nextLeaves := shrinkingGraph.Leaves()
		// nodes must be removed after the next iteration's leaves have been evaluated so that we do not
		// delete the last layer's elements before the last iteration.
		for _, leafNode := range leaves {
			delete(shrinkingGraph.nodes, leafNode)
		}
		leaves = nextLeaves
	}

	return out
}

// TopoSorted returns all the nodes in the graph is topological sort order.
// See also `Graph.TopoSortedLayers()`.
func (g *Graph) TopoSorted() []interface{} {
	nodeCount := 0
	layers := g.TopoSortedLayers()
	for _, layer := range layers {
		nodeCount += len(layer)
	}

	allNodes := make([]interface{}, 0, nodeCount)
	for _, layer := range layers {
		for _, node := range layer {
			allNodes = append(allNodes, node)
		}
	}

	return allNodes
}

func (g *Graph) transitiveDependencies(node interface{}) nodeset {
	return g.buildTransitive(node, g.immediateDependencies)
}

func (g *Graph) immediateDependencies(node interface{}) nodeset {
	if deps, ok := g.dependencies[node]; ok {
		return deps
	}
	return nil
}

func (g *Graph) transitiveDependents(node interface{}) nodeset {
	return g.buildTransitive(node, g.immediateDependents)
}

func (g *Graph) immediateDependents(node interface{}) nodeset {
	if deps, ok := g.dependents[node]; ok {
		return deps
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
func (g *Graph) buildTransitive(root interface{}, nextFn func(interface{}) nodeset) nodeset {
	out := make(nodeset)
	searchNext := []interface{}{root}
	for len(searchNext) > 0 {
		// List of new nodes from this layer of the dependency graph. This is
		// assigned to `searchNext` at the end of the outer "discovery" loop.
		discovered := []interface{}{}
		for _, node := range searchNext {
			// For each node to discover, find the next nodes.
			for nextNode := range nextFn(node) {
				// If we have not seen the node before, add it to the output as well
				// as the list of nodes to traverse in the next iteration.
				if _, ok := out[nextNode]; !ok {
					out[nextNode] = struct{}{}
					discovered = append(discovered, nextNode)
				}
			}
		}
		searchNext = discovered
	}
	return out
}

func copyDepmap(m depmap) depmap {
	out := make(depmap, len(m))
	for k, v := range m {
		out[k] = copyNodeset(v)
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

type updateFn = func(nodes nodeset)

func updateSet(ds depmap, node interface{}, fn updateFn) {
	nodeSet, ok := ds[node]
	if !ok {
		nodeSet = make(nodeset)
		ds[node] = nodeSet
	}
	fn(nodeSet)
}
