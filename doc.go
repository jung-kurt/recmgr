/*
 * Copyright (c) 2015 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
Package recmgr provides a thin, goroutine-safe wrapper around Google's btree
package. It facilitates the use of multiple indexes to manage an in-memory
collection of records.

This package operates on pointers to values, typically structs that can be
indexed in multiple ways.

The methods in this package correspond to the methods of the same name in the
btree package. Because multiple indexes are processed as a group, some methods
are not supported, for example DeleteMin() and DeleteMax(). Similarly, some
method semantics are different, for example Delete() returns the number of
removed keys rather than the deleted item. Additionally, variations of the
traversal methods are available that return a slice of record pointers.

All methods in this package are safe for concurrent goroutine use.

License

recmgr is copyrighted by Kurt Jung and is released under the MIT License.

Installation

To install the package on your system, run

	go get github.com/jung-kurt/recmgr

Later, to receive updates, run

	go get -u github.com/jung-kurt/recmgr

Quick Start

The following Go code demonstrates the creation and operation of a record manager instance.

	var grp recmgr.GrpType
	type person struct {
		name string
		num  int
	}
	personList := []person{
		{"Athos", 1},
		{"Porthos", 2},
		{"Aramis", 3}}
	idxName := grp.Index(4, func(a, b interface{}) bool {
		return a.(*person).name < b.(*person).name
	})
	idxNum := grp.Index(4, func(a, b interface{}) bool {
		return a.(*person).num < b.(*person).num
	})
	for j := range personList {
		grp.ReplaceOrInsert(&personList[j])
	}
	print := func(recPtr interface{}) bool {
		p := recPtr.(*person)
		fmt.Printf("    %-8s %2d\n", p.name, p.num)
		return true
	}
	fmt.Println("Name order")
	idxName.Ascend(print)
	fmt.Println("Number order")
	idxNum.Ascend(print)
	// Output:
	// Name order
	//     Aramis    3
	//     Athos     1
	//     Porthos   2
	// Number order
	//     Athos     1
	//     Porthos   2
	//     Aramis    3

Limitations

The records managed by this package are referenced by pointers so they should
remain accessible for the duration of the recmgr instance. Notice in the range
loop shown above that the record address passed to ReplaceOrInsert() points
into the array that underlies the personList slice, not the address of the
range expression's ephemeral second value.

Within the managed collection of records, if any key field (that is, any struct
field that is used in the less function passed to Index()) is modified, it is
advised to delete the record before modification and add it again afterward to
keep the underlying btrees consistent. Non-key fields can be changed with
impunity.
*/
package recmgr
