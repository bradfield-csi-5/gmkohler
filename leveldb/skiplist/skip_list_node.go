package skiplist

import (
	"bytes"
	"errors"
	"leveldb"
)

type Node interface {
	// CompareKey returns an int whose meaning is like that of a Compare() function:
	// -1 if less than, 0 if equal, 1 if greater than.
	CompareKey(key leveldb.Key) int
	// Key returns the key of a node
	Key() leveldb.Key
	// Value returns the value of a node
	Value() leveldb.Value
	// SetValue sets the value of a node
	SetValue(value leveldb.Value) error
	// ForwardNodeAtLevel seems like it should be private; can we just delete the interface and use nil for NilNode?
	ForwardNodeAtLevel(lvl level) Node
	// SetForwardNodeAtLevel seems like it should be private; can we just delete the interface and use nil for NilNode?
	SetForwardNodeAtLevel(lvl level, node Node) error
	// Next is a convenience function for ForwardNodeAtLevel(1).
	Next() Node
}

// newValueNode does not do validation of the values being inserted, hence private to this package
func newValueNode(key leveldb.Key, value leveldb.Value) *valueNode {
	forwardNodes := forwardList{}
	for j := range forwardNodes {
		// consider a singleton
		forwardNodes[j] = NilNode
	}
	return &valueNode{
		key:          key,
		value:        value,
		forwardNodes: forwardNodes,
	}
}

// newHeaderNode() provides a Node that conforms to the description
// of the "header" of a SkipList. Namely, its forwardList is full of NIL
// nodes, and its key is less than any other valid key.
func newHeaderNode() *valueNode {
	/**
	 * An element NIL is allocated and given a key greater than any legal key.
	 * All levels of all skip lists are terminated with NIL. A new list is
	 * initialized so that the level of the list is equal to 1 and all forward
	 * pointers of the list’s header point to NIL.
	 */
	return newValueNode(nil, nil)
}

type valueNode struct {
	key          leveldb.Key
	value        leveldb.Value
	forwardNodes forwardList
}

func (vn *valueNode) SetForwardNodeAtLevel(lvl level, node Node) error {
	// TODO bounds checking
	vn.forwardNodes[lvl-1] = node
	return nil
}

func (vn *valueNode) CompareKey(k leveldb.Key) int {
	return bytes.Compare(vn.key, k)
}
func (vn *valueNode) Next() Node           { return vn.ForwardNodeAtLevel(1) }
func (vn *valueNode) Key() leveldb.Key     { return vn.key }
func (vn *valueNode) Value() leveldb.Value { return vn.value }

func (vn *valueNode) ForwardNodeAtLevel(lvl level) Node {
	// TODO: bounds checking
	return vn.forwardNodes[lvl-1]
}

func (vn *valueNode) setForwardNodeAtLevel(lvl level, node Node) error {
	// TODO bounds checking
	vn.forwardNodes[lvl-1] = node
	return nil
}

func (vn *valueNode) SetValue(value leveldb.Value) error {
	if len(value) == 0 {
		return errors.New("setting empty value is not allowed")
	}
	vn.value = value
	return nil
}

var NilNode Node = &nilNode{}

type nilNode struct{}

// CompareKey always considers nilNode's value to be greater than another Node's value.
func (nn *nilNode) CompareKey(leveldb.Key) int { return 1 }

// SetForwardNodeAtLevel should not be called on nilNode
func (nn *nilNode) SetForwardNodeAtLevel(level, Node) error {
	panic("nilNode should not be asked to set a new node")
}

// ForwardNodeAtLevel should not be called on nilNode
func (nn *nilNode) ForwardNodeAtLevel(level) Node {
	panic("nilNode should not be asked for a ForwardNode")
}

// Next should not be called on nilNode
func (nn *nilNode) Next() Node {
	panic("should not ask for next node of nilNode")
}

// Key should not be called on nilNode
func (nn *nilNode) Key() leveldb.Key {
	panic("should not ask for key of nilNode")
}

// Value should not be called on nilNode
func (nn *nilNode) Value() leveldb.Value {
	panic("should not ask for value of nilNode")
}

// SetValue should not be called on nilNode
func (nn *nilNode) SetValue(leveldb.Value) error {
	panic("should not ask to set value of nilNode")
}

// forwardList is a type to hang get/set methods off to help avoid level-starts-at-one errors
type forwardList [maxLevel]Node

// getLevel handles translating the level to an array index to avoid off-by-one issues
func (f *forwardList) getLevel(lvl level) Node { return f[lvl-1] }

// setLevel handles translating the level to an array index to avoid off-by-one issues
func (f *forwardList) setLevel(lvl level, node Node) { f[lvl-1] = node }
