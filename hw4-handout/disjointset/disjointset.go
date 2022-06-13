package disjointset

// DisjointSet is the interface for the disjoint-set (or union-find) data
// structure.
// Do not change the definition of this interface.
type DisjointSet interface {
	// UnionSet(s, t) merges (unions) the sets containing s and t,
	// and returns the representative of the resulting merged set.
	UnionSet(int, int) int
	// FindSet(s) returns representative of the class that s belongs to.
	FindSet(int) int
}

// TODO: implement a type that satisfies the DisjointSet interface.
type SetDS struct {
	parent map[int]int
	rank map[int]int
}

func (DisjointS SetDS) UnionSet(union1 int, union2 int) int{
	root := DisjointS.FindSet(union1)
	root2 := DisjointS.FindSet(union2)
	
	if root == root2 {
		return root
	}
	if DisjointS.rank[root] < DisjointS.rank[root2] {
		DisjointS.parent[root] = root2
		return root2
	} else if DisjointS.rank[root] > DisjointS.rank[root2] {
		DisjointS.parent[root2] = root
		return root
	}
	DisjointS.parent[root2] = root
	DisjointS.rank[root]++
	return root
}

func (DisjointS SetDS) FindSet(value int) int{
	if val, ok := DisjointS.parent[value]; ok {
		DisjointS.parent[value] = DisjointS.FindSet(val)
		return DisjointS.parent[value]
	}
	return value
}

// NewDisjointSet creates a struct of a type that satisfies the DisjointSet interface.
func NewDisjointSet() DisjointSet {
	newDS := SetDS{}
	newDS.parent = make(map[int]int)
	newDS.rank = make(map[int]int)
	var newSet DisjointSet = newDS
	return newSet
}