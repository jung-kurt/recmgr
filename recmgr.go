package recmgr

import (
	"github.com/google/btree"
)

type lessFncType func(a, b interface{}) bool

// VisitFncType defines the application function that will be called by the
// record manager library when traversing a record collection.
type VisitFncType func(recPtr interface{}) bool

// IndexType associates a record collection with a record comparison function.
// An instance of this type is used to access records in index order.
type IndexType struct {
	bt   *btree.BTree
	less lessFncType
}

// GrpType aggregates a number of btree indexes. Adding and deleting records is
// done with instances of this type and operates on all indexes.
type GrpType struct {
	list []IndexType
}

type btreeRecType struct {
	less   lessFncType
	recPtr interface{}
}

func (rec btreeRecType) Less(than btree.Item) bool {
	return rec.less(rec.recPtr, than.(btreeRecType).recPtr)
}

// Index adds an index to the record manager instance.
//
// degree specifies the underlying btree degree.
//
// less is a function that should return true if the record pointed to be aPtr
// sorts less than the record pointed to be bPtr.
//
// This function should be called, once for each index, before adding any
// records to the manager, that is, before calling ReplaceOrInsert()
func (rm *GrpType) Index(degree int, less func(aPtr, bPtr interface{}) bool) (idx IndexType) {
	idx = IndexType{bt: btree.New(degree), less: less}
	rm.list = append(rm.list, idx)
	return
}

// Delete removes all keys of the record pointed to by recPtr from the record
// manager. The referenced record is not modified. All fields involved in each
// of the indexes must be specified. A convenient way to do this is to pass
// Delete() the result of Get() since only one key has to be fully specified.
// This method returns zero if recPtr is nil
func (rm *GrpType) Delete(recPtr interface{}) (count int) {
	if recPtr != nil {
		for _, index := range rm.list {
			if nil != index.bt.Delete(btreeRecType{less: index.less, recPtr: recPtr}) {
				count++
			}
		}
	}
	return
}

// ReplaceOrInsert inserts, or replaces if already present, a record reference
// to each of the underlying btree indexes. No action is taken if recPtr is
// nil.
func (rm *GrpType) ReplaceOrInsert(recPtr interface{}) {
	if recPtr != nil {
		for _, index := range rm.list {
			index.bt.ReplaceOrInsert(btreeRecType{less: index.less, recPtr: recPtr})
		}
	}
}

func (idx IndexType) rec(recPtr interface{}) btreeRecType {
	return btreeRecType{less: idx.less, recPtr: recPtr}
}

func visitWrap(fnc VisitFncType) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		return fnc(item.(btreeRecType).recPtr)
	}
}

// Ascend calls fnc for every value in the collection, that is, [first, last],
// in the order specified by idx. The traversal is terminated if fnc returns
// false.
func (idx IndexType) Ascend(fnc VisitFncType) {
	idx.bt.Ascend(visitWrap(fnc))
}

// AscendLessThan calls fnc for every value in the collection less than the
// value pointed to by ltPtr, that is, [first, ltPtr), in the order specified
// by idx. The traversal is terminated if fnc returns false. All key fields
// needed by idx must be assigned in the value pointed to by ltPtr.
func (idx IndexType) AscendLessThan(ltPtr interface{}, fnc VisitFncType) {
	idx.bt.AscendLessThan(idx.rec(ltPtr), visitWrap(fnc))
}

// AscendGreaterOrEqual calls fnc for every value in the collection greater
// than or equal to the value pointed to by gePtr, that is, [gePtr, last], in
// the order specified by idx. The traversal is terminated if fnc returns
// false. All key fields needed by idx must be assigned in the value pointed to
// by gePtr.
func (idx IndexType) AscendGreaterOrEqual(gePtr interface{}, fnc VisitFncType) {
	idx.bt.AscendGreaterOrEqual(idx.rec(gePtr), visitWrap(fnc))
}

// AscendRange calls fnc for every value in the collection greater than or
// equal to the value pointed to by gePtr and less than the value pointed to by
// ltPtr, that is, [gePtr, lePtr), in the order specified by idx. The traversal
// is terminated if fnc returns false. All key fields needed by idx must be
// assigned in the values pointed to by gePtr and ltPtr.
func (idx IndexType) AscendRange(gePtr, ltPtr interface{}, fnc VisitFncType) {
	idx.bt.AscendRange(idx.rec(gePtr), idx.rec(ltPtr), visitWrap(fnc))
}

// Get retrieves the value associated with the key specified by keyPtr. All key
// fields needed by idx must be assigned in the value pointed to by keyPtr. If
// the value cannot be located, nil is returned.
func (idx IndexType) Get(keyPtr interface{}) (recPtr interface{}) {
	item := idx.bt.Get(idx.rec(keyPtr))
	if item != nil {
		recPtr = item.(btreeRecType).recPtr
	}
	return
}

// Has returns true if the value associated with the key specified by keyPtr is
// present, false otherwise. All key fields needed by idx must be assigned in
// the value pointed to by keyPtr.
func (idx IndexType) Has(keyPtr interface{}) bool {
	return idx.bt.Has(idx.rec(keyPtr))
}

// Len returns the number of items in the btree associated with idx.
func (idx IndexType) Len() int {
	return idx.bt.Len()
}
