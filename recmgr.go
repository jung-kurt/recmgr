package recmgr

import (
	"github.com/google/btree"
	"sync"
)

type lessFncType func(a, b interface{}) bool

// VisitFncType defines the application function that will be called by the
// record manager library when traversing a record collection. The boolean
// false value is returned to terminate the traversal.
type VisitFncType func(recPtr interface{}) bool

// IndexType associates a record collection with a record comparison function.
// A value of this type is used to access records in index order. Its methods
// are safe for concurrent goroutine use.
type IndexType struct {
	bt       *btree.BTree
	less     lessFncType
	mutexPtr *sync.RWMutex
}

// GrpType aggregates a number of btree indexes. Adding and deleting records is
// done with instances of this type and operates on all indexes. Values of this
// type require no special initialization before use. Its methods are safe for
// concurrent goroutine use.
type GrpType struct {
	list  []IndexType
	mutex sync.RWMutex
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
// less is a function that should return true if the record pointed to by aPtr
// sorts less than the record pointed to by bPtr.
//
// This function should be called, once for each index, before adding any
// records to the manager, that is, before calling ReplaceOrInsert().
func (grp *GrpType) Index(degree int, less func(aPtr, bPtr interface{}) bool) (idx IndexType) {
	grp.mutex.Lock()
	idx = IndexType{bt: btree.New(degree), less: less, mutexPtr: &grp.mutex}
	grp.list = append(grp.list, idx)
	grp.mutex.Unlock()
	return
}

// Delete removes all keys of the record pointed to by recPtr from the record
// manager. The referenced record is not modified. All fields involved in each
// of the indexes must be specified. A convenient way to do this is to pass
// Delete() the result of Get() since only one key has to be fully specified.
// This method returns zero if recPtr is nil.
func (grp *GrpType) Delete(recPtr interface{}) (count int) {
	if recPtr != nil {
		grp.mutex.Lock()
		for _, index := range grp.list {
			if nil != index.bt.Delete(btreeRecType{less: index.less, recPtr: recPtr}) {
				count++
			}
		}
		grp.mutex.Unlock()
	}
	return
}

// Has returns the number of indexes that contain a key for the record
// specified by recPtr. Each key field used in an index must be assigned. Zero
// is returned if recPtr is nil.
func (grp *GrpType) Has(recPtr interface{}) (count int) {
	if recPtr != nil {
		grp.mutex.Lock()
		for _, index := range grp.list {
			if index.bt.Has(btreeRecType{less: index.less, recPtr: recPtr}) {
				count++
			}
		}
		grp.mutex.Unlock()
	}
	return
}

// ReplaceOrInsert inserts, or replaces if already present, a record reference
// to each of the underlying btree indexes. No action is taken if recPtr is
// nil.
func (grp *GrpType) ReplaceOrInsert(recPtr interface{}) {
	if recPtr != nil {
		grp.mutex.Lock()
		for _, index := range grp.list {
			index.bt.ReplaceOrInsert(btreeRecType{less: index.less, recPtr: recPtr})
		}
		grp.mutex.Unlock()
	}
}

func (idx IndexType) rec(recPtr interface{}) btreeRecType {
	return btreeRecType{less: idx.less, recPtr: recPtr}
}

func listGen(listPtr *[]interface{}) func(item btree.Item) bool {
	return func(item btree.Item) bool {
		*listPtr = append(*listPtr, item.(btreeRecType).recPtr)
		return true
	}
}

func (idx IndexType) traverse(list []interface{}, fnc VisitFncType) {
	for _, recPtr := range list {
		if !fnc(recPtr) {
			return
		}
	}
}

// AscendList returns a slice of record pointers for every value in the
// collection in the order associated with idx.
func (idx IndexType) AscendList() (list []interface{}) {
	(*idx.mutexPtr).RLock()
	idx.bt.Ascend(listGen(&list))
	(*idx.mutexPtr).RUnlock()
	return
}

// Ascend calls fnc for every value in the collection, that is, [first, last],
// in the order associated with idx. The traversal is terminated if fnc returns
// false.
func (idx IndexType) Ascend(fnc VisitFncType) {
	idx.traverse(idx.AscendList(), fnc)
}

// AscendLessThanList returns a slice with every value in the collection less
// than the value pointed to by ltPtr, that is, [first, ltPtr), in the order
// associated with idx. All key fields needed by idx must be assigned in the
// value pointed to by ltPtr.
func (idx IndexType) AscendLessThanList(ltPtr interface{}) (list []interface{}) {
	(*idx.mutexPtr).RLock()
	idx.bt.AscendLessThan(idx.rec(ltPtr), listGen(&list))
	(*idx.mutexPtr).RUnlock()
	return
}

// AscendLessThan calls fnc for every value in the collection less than the
// value pointed to by ltPtr, that is, [first, ltPtr), in the order determined
// by idx. The traversal is terminated if fnc returns false. All key fields
// needed by idx must be assigned in the value pointed to by ltPtr.
func (idx IndexType) AscendLessThan(ltPtr interface{}, fnc VisitFncType) {
	idx.traverse(idx.AscendLessThanList(ltPtr), fnc)
}

// AscendGreaterOrEqualList returns a slice with every value in the collection
// greater than or equal to the value pointed to by gePtr, that is, [gePtr,
// last], in the order associated with idx. All key fields needed by idx must
// be assigned in the value pointed to by gePtr.
func (idx IndexType) AscendGreaterOrEqualList(gePtr interface{}) (list []interface{}) {
	(*idx.mutexPtr).RLock()
	idx.bt.AscendGreaterOrEqual(idx.rec(gePtr), listGen(&list))
	(*idx.mutexPtr).RUnlock()
	return
}

// AscendGreaterOrEqual calls fnc for every value in the collection greater
// than or equal to the value pointed to by gePtr, that is, [gePtr, last], in
// the order associated with idx. The traversal is terminated if fnc returns
// false. All key fields needed by idx must be assigned in the value pointed to
// by gePtr.
func (idx IndexType) AscendGreaterOrEqual(gePtr interface{}, fnc VisitFncType) {
	idx.traverse(idx.AscendGreaterOrEqualList(gePtr), fnc)
}

// AscendRangeList returns a slice with every value in the collection greater
// than or equal to the value pointed to by gePtr and less than the value
// pointed to by ltPtr, that is, [gePtr, lePtr), in the order specified by idx.
// All key fields needed by idx must be assigned in the values pointed to by
// gePtr and ltPtr.
func (idx IndexType) AscendRangeList(gePtr, ltPtr interface{}) (list []interface{}) {
	(*idx.mutexPtr).RLock()
	idx.bt.AscendRange(idx.rec(gePtr), idx.rec(ltPtr), listGen(&list))
	(*idx.mutexPtr).RUnlock()
	return
}

// AscendRange calls fnc for every value in the collection greater than or
// equal to the value pointed to by gePtr and less than the value pointed to by
// ltPtr, that is, [gePtr, lePtr), in the order specified by idx. The traversal
// is terminated if fnc returns false. All key fields needed by idx must be
// assigned in the values pointed to by gePtr and ltPtr.
func (idx IndexType) AscendRange(gePtr, ltPtr interface{}, fnc VisitFncType) {
	idx.traverse(idx.AscendRangeList(gePtr, ltPtr), fnc)
}

// Get retrieves the value associated with the key specified by keyPtr. All key
// fields needed by idx must be assigned in the value pointed to by keyPtr. If
// the value cannot be located, nil is returned.
func (idx IndexType) Get(keyPtr interface{}) (recPtr interface{}) {
	(*idx.mutexPtr).RLock()
	item := idx.bt.Get(idx.rec(keyPtr))
	(*idx.mutexPtr).RUnlock()
	if item != nil {
		recPtr = item.(btreeRecType).recPtr
	}
	return
}

// Has returns true if the value associated with the key specified by keyPtr is
// present, false otherwise. All key fields needed by idx must be assigned in
// the value pointed to by keyPtr.
func (idx IndexType) Has(keyPtr interface{}) (ok bool) {
	(*idx.mutexPtr).RLock()
	ok = idx.bt.Has(idx.rec(keyPtr))
	(*idx.mutexPtr).RUnlock()
	return
}

// Len returns the number of items in the btree associated with idx.
func (idx IndexType) Len() (count int) {
	(*idx.mutexPtr).RLock()
	count = idx.bt.Len()
	(*idx.mutexPtr).RUnlock()
	return
}
