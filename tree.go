package gobreak

import (
	"strings"

	linq "github.com/ahmetalpbalkan/go-linq"
)

type TreePath interface {
	GetMc() string
	GetUri() string
	GetQz() int
}

type TreeNode struct {
	Mc    string
	Uri   string
	M     T
	Nodes []*TreeNode
}

func newTreeNode(m TreePath) *TreeNode {
	return &TreeNode{m.GetMc(), m.GetUri(), m, make([]*TreeNode, 0)}
}

func BuildTree(src T) []*TreeNode {
	root := &TreeNode{"", "", nil, make([]*TreeNode, 0)}
	buildTreeNodes(src, root, "")
	return root.Nodes
}

func buildTreeNodes(src T, r *TreeNode, prefix string) {
	results, _ := queryChildren(linq.From(src), prefix)
	for _, it := range results {
		child := newTreeNode(it.(TreePath))
		r.Nodes = append(r.Nodes, child)
		buildTreeNodes(src, child, it.(TreePath).GetUri()+".")
	}
}

func queryChildren(q linq.Query, prefix string) ([]linq.T, error) {
	return q.Where(func(s linq.T) (bool, error) {
		last := strings.TrimPrefix(s.(TreePath).GetUri(), prefix)
		return strings.HasPrefix(s.(TreePath).GetUri(), prefix) && !strings.Contains(last, "."), nil
	}).OrderBy(func(a, b linq.T) bool {
		return a.(TreePath).GetQz() > b.(TreePath).GetQz()
	}).Results()
}
