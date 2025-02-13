package store

type URLHeap struct {
	Nodes []*URLNode
	By    func(i, j *URLNode) bool // Sort by func
}

func (h *URLHeap) Len() int           { return len(h.Nodes) }
func (h *URLHeap) Less(i, j int) bool { return h.By(h.Nodes[i], h.Nodes[j]) }
func (h *URLHeap) Swap(i, j int)      { h.Nodes[i], h.Nodes[j] = h.Nodes[j], h.Nodes[i] }

func (h *URLHeap) Push(x interface{}) {
	h.Nodes = append(h.Nodes, x.(*URLNode))
}

func (h *URLHeap) Pop() interface{} {
	old := h.Nodes
	n := len(old)
	x := old[n-1]
	h.Nodes = old[0 : n-1]
	return x
}

func findIndex(h *URLHeap, node *URLNode) int {
	for i, n := range h.Nodes {
		if n == node {
			return i
		}
	}
	return -1
}
