package skiplist

// translating https://www.epaperpress.com/sortsearch/download/skiplist.pdf

import (
	"errors"
	"leveldb"
	"math/rand/v2"
)

const (
	maxLevel level   = 16
	p        float32 = 0.5
)

// level should range from 1 to maxLevel.  Methods exist on the appropriate types
// to translate this to zero-based indexing when accessing arrays.
type level uint8

// SkipList provides a storage mechanism.
// A level `i` node has `i` forward pointers, indexed 1 through `i`.
// We do not need to store the level of a node in the node.
// Levels are capped at some appropriate level MaxLevel
// The level of a list is maximum level currently in the list, or 1 if list is empty
type SkipList struct {
	header     Node
	level      level
	numEntries uint64
}

// NewSkipList builds a SkipList with the appropriate state.
// An element NIL is allocated and given a key greater than any legal key.
// All levels of all skip lists are terminated with NIL.
// A new list is initialized so that the level of the list is equal to 1 and
// all forward pointers of the header
func NewSkipList() *SkipList {
	return &SkipList{
		header: newHeaderNode(),
		level:  1, // should this be one or zero?
	}
}

// Search searches for an element by traversing forward pointers that do not
// overshoot the node containing the element being searched for.  When
// no more progress can be made at the current level of forward pointers,
// the search moves down to the next level.  When we can make no more
// progress at level 1, we must be immediately in front of the node that
// contains the desired element (if it is in the list).
func (sl *SkipList) Search(searchKey leveldb.Key) (leveldb.Value, error) {
	// consider µ-optimization:
	//if sl.deletedKeys.contains(searchKey) {
	//	return nil, leveldb.NewNotFoundError(searchKey)
	//}
	currentNode, err := sl.TraverseUntil(searchKey, nil)
	if err != nil {
		return nil, err
	}

	currentNode = currentNode.Next()
	if currentNode.CompareKey(searchKey) == 0 {
		return currentNode.Value(), nil
	} else {
		return nil, leveldb.NewNotFoundError(searchKey)
	}
}

// Insert follows a "search and splice" approach.  A vector _update_
// is maintained so that when the search is complete (and we are ready to
// perform the splice), `update[i]` contains a pointer to the rightmost node of
// level i or higher that is to the left o the location of insertion/deletion.
//
// If an insertion generates a node with a level greater than the previous
// maximum  level of the list, we update hte maximum level of the list and
// initialize portions of the update vector.  After each deletion, we check
// if we have deleted the maximum element of the list and if so, decrease the
// maximum level of the list.
func (sl *SkipList) Insert(searchKey leveldb.Key, newValue leveldb.Value) error {
	if len(searchKey) == 0 {
		return errors.New("cannot insert blank key")
	}

	var (
		lastNodeTraversedPerLevel = forwardList{}
		currentNode               Node
		err                       error
	)

	currentNode, err = sl.TraverseUntil(searchKey, lastNodeTraversedPerLevel.setLevel)
	if err != nil {
		return err
	}
	currentNode = currentNode.Next()
	if currentNode.CompareKey(searchKey) == 0 {
		if err = currentNode.SetValue(newValue); err != nil {
			return err
		}
		return nil
	}

	insertionLevel := randomLevel()
	if insertionLevel > sl.level {
		for lvl := sl.level + 1; lvl <= insertionLevel; lvl++ {
			lastNodeTraversedPerLevel.setLevel(lvl, sl.header)
		}
		sl.level = insertionLevel
	}
	newNode := newValueNode(searchKey, newValue)
	for lvl := level(1); lvl <= insertionLevel; lvl++ {
		nodeToUpdate := lastNodeTraversedPerLevel.getLevel(lvl)
		if nodeToUpdate == nil {
			continue
		}
		if err := newNode.SetForwardNodeAtLevel(
			lvl,
			nodeToUpdate.ForwardNodeAtLevel(lvl),
		); err != nil {
			return err
		}
		if err := nodeToUpdate.SetForwardNodeAtLevel(
			lvl,
			newNode,
		); err != nil {
			return err
		}
	}

	sl.numEntries++

	return nil
}

func (sl *SkipList) Delete(searchKey leveldb.Key) error {
	var (
		nodesToUpdate = forwardList{}
		currentNode   Node
		err           error
	)
	currentNode, err = sl.TraverseUntil(searchKey, nodesToUpdate.setLevel)
	if err != nil {
		return err
	}
	currentNode = currentNode.ForwardNodeAtLevel(1)
	if currentNode.CompareKey(searchKey) == 0 {
		for j := range sl.level {
			lvl := j + 1
			if nodesToUpdate.getLevel(lvl).ForwardNodeAtLevel(lvl) != currentNode {
				break
			}
			err = nodesToUpdate.getLevel(lvl).SetForwardNodeAtLevel(lvl, currentNode.ForwardNodeAtLevel(lvl))
			if err != nil {
				return err // panic because incomplete operation?
			}
			for sl.level > 1 && sl.header.ForwardNodeAtLevel(sl.level) == NilNode {
				sl.level--
			}
		}
		sl.numEntries--

	} /*else {
		return leveldb.NewNotFoundError(searchKey)
	}*/
	return nil
}

func (sl *SkipList) Scan(start, limit leveldb.Key) ([]leveldb.Value, error) {
	nearestNode, err := sl.TraverseUntil(start, nil)
	if err != nil {
		return nil, err
	}
	firstNode := nearestNode.ForwardNodeAtLevel(1)
	var values []leveldb.Value
	for currentNode := firstNode; currentNode.CompareKey(limit) <= 0; currentNode = currentNode.ForwardNodeAtLevel(1) {
		values = append(values, currentNode.Value())
	}
	return values, nil
}

func (sl *SkipList) Reset() error {
	sl.level = 1
	sl.numEntries = 0
	sl.header = newHeaderNode() // forget all data
	return nil
}

// TraverseUntil returns a node that is in the spot where the first-level forward entry would be the key, or would be
// greater than the key.  In other words, it returns the node just before the desired key, whether or not that key
// exists.
func (sl *SkipList) TraverseUntil(key leveldb.Key, listener func(level, Node)) (Node, error) {
	var currentNode = sl.header
	for currentLevel := sl.level; currentLevel > 0; currentLevel-- {
		for currentNode.ForwardNodeAtLevel(currentLevel).CompareKey(key) < 0 {
			currentNode = currentNode.ForwardNodeAtLevel(currentLevel)
		}
		if listener != nil {
			listener(currentLevel, currentNode)
		}
	}
	return currentNode, nil
}

// randomLevel() selects a level between 1 and the maxLevel.  Note that this is "base 1" and should
// be decremented by 1 when accessing arrays of nodes at each level.
func randomLevel() level {
	var lvl level = 1 // consider "lvl" like "length" and handle off-by-one in iteration/access
	for rand.Float32() < p && lvl < maxLevel {
		lvl++
	}
	return lvl
}

// type keySet map[string]struct{}
//
// func (ks keySet) remove(key leveldb.Key) {
// 	delete(ks, string(key))
// }
// func (ks keySet) add(key leveldb.Key) {
// 	ks[string(key)] = struct{}{}
// }
