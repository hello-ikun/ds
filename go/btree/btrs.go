package btrs

import (
	"fmt"
)

type BTree struct {
	Root *Node
	size int
	m    int
}
type Node struct {
	Parent   *Node
	Entries  []*Entry
	Children []*Node
}

// Entry 存储的实体数据
type Entry struct {
	Key int
	Val interface{}
}

func NewTree(order int) *BTree {
	if order < 3 {
		order = 3
	}
	return &BTree{m: order}
}
func (t *BTree) Put(key int, val interface{}) {
	entry := &Entry{Key: key, Val: val}
	if t.Root == nil {
		t.Root = &Node{Entries: []*Entry{entry}}
		t.size++
		return
	}
	if t.insert(t.Root, entry) {
		t.size++
	}
}
func (t *BTree) insert(node *Node, entry *Entry) bool {
	if t.isLeaf(node) {
		return t.insertIntoLeaf(node, entry)
	}
	return t.insertIntoInternal(node, entry)
}
func (t *BTree) insertIntoLeaf(node *Node, entry *Entry) bool {
	idx, ok := t.search(node, entry.Key)
	if ok {
		node.Entries[idx] = entry
		return false
	}
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[idx+1:], node.Entries[idx:])
	node.Entries[idx] = entry
	t.split(node)
	return true
}
func (t *BTree) split(node *Node) {
	if !t.shouldSplit(node) {
		return
	}
	fmt.Println("node:", node, "start split")
	if node == t.Root {
		t.splitRoot()
		return
	}
	t.splitNonRoot(node)
}
func (t *BTree) splitRoot() {
	mid := t.middle()

	left := &Node{
		Entries: append([]*Entry(nil), t.Root.Entries[:mid]...),
	}
	right := &Node{
		Entries: append([]*Entry(nil), t.Root.Entries[mid+1:]...),
	}
	if !t.isLeaf(t.Root) {
		left.Children = append([]*Node(nil), t.Root.Children[:mid+1]...)
		right.Children = append([]*Node(nil), t.Root.Children[mid+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}
	newRoot := &Node{
		Entries:  []*Entry{t.Root.Entries[mid]},
		Children: []*Node{left, right},
	}
	left.Parent = newRoot
	right.Parent = newRoot
	t.Root = newRoot
}
func (t *BTree) splitNonRoot(node *Node) {
	middle := t.middle()
	parent := node.Parent

	left := &Node{Entries: append([]*Entry(nil), node.Entries[:middle]...), Parent: parent}
	right := &Node{Entries: append([]*Entry(nil), node.Entries[middle+1:]...), Parent: parent}

	if !t.isLeaf(node) {
		left.Children = append([]*Node(nil), node.Children[:middle+1]...)
		right.Children = append([]*Node(nil), node.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := t.search(parent, node.Entries[middle].Key)

	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	parent.Children[insertPosition] = left

	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right

	t.split(parent)
}
func setParent(nodes []*Node, parent *Node) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

// 节点个数大于最大个数的话
func (t *BTree) shouldSplit(node *Node) bool {
	return len(node.Entries) > t.maxEntries()
}
func (t *BTree) Empty() bool {
	return t.size == 0
}
func (t *BTree) middle() int {
	return (t.m - 1) / 2
}
func (t *BTree) maxChildren() int {
	return t.m
}
func (t *BTree) minChildren() int {
	return (t.m + 1) / 2
}
func (t *BTree) maxEntries() int {
	return t.maxChildren() - 1
}
func (t *BTree) minEntries() int {
	return t.minChildren() - 1
}
func (t *BTree) isLeaf(node *Node) bool {
	return len(node.Children) == 0
}
func (t *BTree) insertIntoInternal(node *Node, entry *Entry) bool {
	idx, ok := t.search(node, entry.Key)
	if ok {
		node.Entries[idx] = entry
		return false
	}
	return t.insert(node.Children[idx], entry)
}

func (t *BTree) search(node *Node, key int) (int, bool) {
	left := 0
	right := len(node.Entries) - 1
	ans := -1

	for left <= right {
		mid := (left + right) / 2
		if node.Entries[mid].Key == key {
			return mid, true
		} else if node.Entries[mid].Key < key {
			ans = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return ans + 1, false
}
func (t *BTree) searchRecur(node *Node, key int) (st *Node, idx int, ok bool) {
	if t.Empty() {
		return nil, -1, false
	}
	st = node
	for {
		idx, ok = t.search(node, key)
		if ok {
			return node, idx, true
		}
		if t.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[idx]
	}
}
func (t *BTree) Get(key int) (interface{}, bool) {
	node, idx, ok := t.searchRecur(t.Root, key)
	if ok {
		return node.Entries[idx].Val, true
	}
	return nil, false
}
func (t *BTree) Remove(key int) {
	node, idx, ok := t.searchRecur(t.Root, key)
	if ok {
		t.delete(node, idx)
		t.size--
	}
}
func (t *BTree) delete(node *Node, idx int) {
	if t.isLeaf(node) {
		delKey := node.Entries[idx].Key
		t.deleteEntry(node, idx)
		t.rebalance(node, delKey)
		if len(t.Root.Entries) == 0 {
			t.Root = nil
		}
		return
	}
	leftLargestNode := t.right(node.Children[idx])
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[idx] = leftLargestNode.Entries[leftLargestEntryIndex]
	deleteKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	t.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	t.rebalance(leftLargestNode, deleteKey)
}

func (t *BTree) deleteEntry(node *Node, idx int) {
	copy(node.Entries[idx:], node.Entries[idx+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}
func (t *BTree) rebalance(node *Node, deleteKey int) {
	if node == nil || len(node.Entries) >= t.minEntries() {
		return
	}
	leftSibling, leftSiblingIndex := t.leftSibling(node, deleteKey)
	if leftSibling != nil && len(leftSibling.Entries) > t.minEntries() {
		node.Entries = append(
			[]*Entry{node.Parent.Entries[leftSiblingIndex]},
			node.Entries...,
		)
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		t.deleteEntry(leftSibling, len(leftSibling.Entries)-1)
		if !t.isLeaf(leftSibling) {
			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*Node{leftSiblingRightMostChild}, node.Children...)
			t.deleteEntry(leftSibling, len(leftSibling.Children)-1)
		}
		return
	}
	rightSibling, rightSiblingIndex := t.rightSibling(node, deleteKey)
	if rightSibling != nil && len(rightSibling.Entries) > t.minEntries() {
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		t.deleteEntry(rightSibling, 0)
		if !t.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			t.deleteEntry(rightSibling, 0)
		}
		return
	}
	if rightSibling != nil {
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deleteKey = node.Parent.Entries[rightSiblingIndex-1].Key
		t.deleteEntry(node.Parent, rightSiblingIndex-1)
		t.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		t.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSibling != nil {
		entries := append([]*Entry(nil), leftSibling.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])

		node.Entries = append(entries, node.Entries...)
		deleteKey = node.Parent.Entries[leftSiblingIndex].Key
		t.deleteEntry(node.Parent, leftSiblingIndex)
		t.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		t.deleteChild(node.Parent, leftSiblingIndex)
	}
	if node.Parent == t.Root && len(t.Root.Entries) == 0 {
		t.Root = node
		node.Parent = nil
		return
	}
	t.rebalance(node.Parent, deleteKey)
}
func (t *BTree) left(node *Node) *Node {
	if t.Empty() {
		return nil
	}
	cur := node
	for {
		if t.isLeaf(cur) {
			return cur
		}
		cur = cur.Children[0]
	}
}
func (t *BTree) right(node *Node) *Node {
	if t.Empty() {
		return nil
	}
	cur := node
	for {
		if t.isLeaf(cur) {
			return cur
		}
		cur = cur.Children[len(cur.Children)-1]
	}
}
func (t *BTree) leftSibling(node *Node, key int) (*Node, int) {
	if node.Parent != nil {
		idx, _ := t.search(node.Parent, key)
		idx--
		if idx >= 0 && idx < len(node.Parent.Children) {
			return node.Parent.Children[idx], idx
		}
	}
	return nil, -1
}
func (t *BTree) rightSibling(node *Node, key int) (*Node, int) {
	if node.Parent != nil {
		idx, _ := t.search(node.Parent, key)
		idx++
		if idx < len(node.Parent.Children) {
			return node.Parent.Children[idx], idx
		}
	}
	return nil, -1
}
func (t *BTree) deleteChild(node *Node, idx int) {
	if idx >= len(node.Children) {
		return
	}
	copy(node.Children[idx:], node.Children[idx+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}
func (t *BTree) appendChildren(fromNode, toNode *Node) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}
func (t *BTree) prependChildren(fromNode, toNode *Node) {
	children := append([]*Node(nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

type Logger struct{}

func (l *Logger) tree(node *Node, childName string, dsc func(*Node) string, depth int, prefix string) {
	if depth == 0 {
		fmt.Printf("+--%s\n", dsc(node))
		depth++
	}

	childCount := len(node.Children)
	for idx, child := range node.Children {
		isLastChild := idx == childCount-1

		newPrefix := prefix
		if isLastChild {
			newPrefix += "    "
		} else {
			newPrefix += " |   "
		}

		fmt.Println(prefix, childPrefix(isLastChild), dsc(child))

		l.tree(child, childName, dsc, depth+1, newPrefix)
	}
}

func childPrefix(isLastChild bool) string {
	if isLastChild {
		return "└-- "
	}
	return "|-- "
}
func printTree(node *Node) {
	fmt.Println("\n************BTree*************")

	dsc := func(node *Node) string {
		s := ""
		for _, keyword := range node.Entries {
			s += fmt.Sprintf("%v,", keyword.Key)
		}
		s = s[:len(s)-1]
		return s
	}
	logger := &Logger{}
	logger.tree(node, "child_nodes", dsc, 0, "    ")

	fmt.Println("******************************")
}
