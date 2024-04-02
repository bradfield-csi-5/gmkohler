package leveldb

import "bytes"

type skipListNode interface {
	// CompareKey returns an int whose meaning is like that of a Compare() function:
	// -1 if less than, 0 if equal, 1 if greater than.
	CompareKey(key Key) int
	// Value returns the value of a node
	Value() Value
	// SetValue sets the value of a node
	SetValue(value Value)
	// ForwardNodeAtLevel retrieves the "forward" node at a particular level.
	ForwardNodeAtLevel(lvl level) skipListNode
	// SetForwardNodeAtLevel sets the "forward" node at a particular level.
	SetForwardNodeAtLevel(lvl level, node skipListNode) error
}

type forwardList [maxLevel]skipListNode

// GetLevel handles translating the level to an array index to avoid off-by-one issues
func (f *forwardList) GetLevel(lvl level) skipListNode { return f[lvl-1] }

// SetLevel handles translating the level to an array index to avoid off-by-one issues
func (f *forwardList) SetLevel(lvl level, node skipListNode) { f[lvl-1] = node }

type valueNode struct {
	key          Key
	value        Value
	forwardNodes forwardList
}

func (vn *valueNode) SetForwardNodeAtLevel(lvl level, node skipListNode) error {
	// TODO bounds checking
	vn.forwardNodes[lvl-1] = node
	return nil
}

func (vn *valueNode) CompareKey(k Key) int {
	return bytes.Compare(vn.key, k)
}

func (vn *valueNode) ForwardNodeAtLevel(lvl level) skipListNode {
	// TODO: bounds checking
	return vn.forwardNodes[lvl-1]
}

func (vn *valueNode) Value() Value { return vn.value }

func (vn *valueNode) SetValue(value Value) { vn.value = value }

func newValueNode(key Key, value Value) skipListNode {
	forwardNodes := forwardList{}
	for j := range forwardNodes {
		// consider a singleton
		forwardNodes[j] = &nilNode{}
	}
	return &valueNode{
		key:          key,
		value:        value,
		forwardNodes: forwardNodes,
	}
}

// newHeaderNode() provides a skipListNode that conforms to the description
// of the "header" of a SkipList. Namely, its forwardList is full of NIL
// nodes, and its key is less than any other valid key.
func newHeaderNode() skipListNode {
	/**
	 * An element NIL is allocated and given a key greater than any legal key.
	 * All levels of all skip lists are terminated with NIL. A new list is
	 * initialized so that the level of the list is equal to 1 and all forward
	 * pointers of the list’s header point to NIL.
	 */
	return newValueNode(nil, nil)
}

type nilNode struct{}

func (nn *nilNode) SetForwardNodeAtLevel(level, skipListNode) error {
	panic("should not set forward node on nilNode")
}

func (nn *nilNode) ForwardNodeAtLevel(level) skipListNode {
	panic("should not ask forward node of nilNode")
}

func (nn *nilNode) CompareKey(Key) int { return 1 }

func (nn *nilNode) Value() Value {
	panic("should not ask for value of nilNode")
}

func (nn *nilNode) SetValue(Value) {
	panic("should not ask to set value of nilNode")
}
