package deps

import (
	"errors"
)

// Semantically, this is Map<Node, Set<Node>>.
type depmap map[interface{}]map[interface{}]struct{}

type Graph struct {
	// Maintain dependency relationships in both directions.
	// `dependencies` tracks child -> parent, and `dependents` tracks parents -> children.
	dependencies, dependents depmap
}

func New() *Graph {
	return &Graph{
		dependencies: make(depmap),
		dependents:   make(depmap),
	}
}

func (g *Graph) DependOn(node, dep interface{}) error {
	if node == dep {
		return errors.New("self-referential dependencies not allowed")
	}
	if g.DependsOn(dep, node) {
		return errors.New("circular dependencies not allowed")
	}

	updateSet(g.dependencies, node, func(nodes map[interface{}]struct{}) {
		nodes[dep] = struct{}{}
	})
	updateSet(g.dependents, dep, func(nodes map[interface{}]struct{}) {
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
	allNodes := make(map[interface{}]struct{}, approxSize)
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
	for node := range g.dependents {
		if _, ok := g.dependencies[node]; !ok {
			out = append(out, node)
		}
	}
	return out
}

// TopoSorted returns a slice of all of the graph nodes in topological sort order. That is,
// if `B` depends on `A`, then `A` is guaranteed to come before `B` in the sorted output.
// The graph is guaranteed to be cycle-free because cycles are detected while building the
// graph. This implements Kahn's algorithm for topological sorting.
func (g *Graph) TopoSorted() []interface{} {
	out := []interface{}{}

	// Copy dependency information to be mutated below.
	// dependencies: child -> parents
	dependencies := make(depmap)
	for k, v := range g.dependencies {
		dependencies[k] = v
	}
	// dependents: parent -> children
	dependents := make(depmap)
	for k, v := range g.dependents {
		dependents[k] = v
	}

	// Keep track of the set of nodes whose dependencies have already been met.
	dependenciesMet := make(map[interface{}]struct{})
	markDependenciesMet := func(node interface{}) {
		// Mark node as not depending on anything already in `out`.
		dependenciesMet[node] = struct{}{}
		// Remove the record of everything that depended on `node`.
		dependsOnNode := dependents[node]
		for dependingNode := range dependsOnNode {
			delete(dependencies[dependingNode], node)
			if len(dependencies[dependingNode]) == 0 {
				delete(dependencies, dependingNode)
			}
		}
	}
	for _, leafNode := range g.Leaves() {
		markDependenciesMet(leafNode)
	}

	for len(dependenciesMet) > 0 {
		elem := setPop(dependenciesMet)
		out = append(out, elem)
		for dependentNode := range dependents[elem] {
			if _, hasDependents := dependencies[dependentNode]; !hasDependents {
				markDependenciesMet(dependentNode)
			}
		}
	}

	return out
}

func (g *Graph) transitiveDependencies(node interface{}) map[interface{}]struct{} {
	return g.buildTransitive(node, g.immediateDependencies)
}

func (g *Graph) immediateDependencies(node interface{}) map[interface{}]struct{} {
	if deps, ok := g.dependencies[node]; ok {
		return deps
	}
	return nil
}

func (g *Graph) transitiveDependents(node interface{}) map[interface{}]struct{} {
	return g.buildTransitive(node, g.immediateDependents)
}

func (g *Graph) immediateDependents(node interface{}) map[interface{}]struct{} {
	if deps, ok := g.dependents[node]; ok {
		return deps
	}
	return nil
}

// buildTransitive starts at `root` and continues calling `nextFn` to keep discovering more nodes until
// the graph cannot produce any more. It returns the set of all discovered nodes.
func (g *Graph) buildTransitive(root interface{}, nextFn func(interface{}) map[interface{}]struct{}) map[interface{}]struct{} {
	out := make(map[interface{}]struct{})
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

type updateFn = func(nodes map[interface{}]struct{})

func updateSet(ds depmap, node interface{}, fn updateFn) {
	nodeSet, ok := ds[node]
	if !ok {
		nodeSet = make(map[interface{}]struct{})
		ds[node] = nodeSet
	}
	fn(nodeSet)
}

func setPop(set map[interface{}]struct{}) interface{} {
	for elem := range set {
		defer delete(set, elem)
		return elem
	}
	return nil
}
