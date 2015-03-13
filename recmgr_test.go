package recmgr_test

import (
	"fmt"
	"github.com/jung-kurt/recmgr"
	"strings"
)

type RecType struct {
	name string
	num  int
}

func (rec RecType) String() string {
	numStr := fmt.Sprintf("%d", rec.num)
	return fmt.Sprintf("%s%s%s", rec.name, strings.Repeat(".", 18-len(numStr)-len(rec.name)), numStr)
}

func Example_basic() {
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
	fmt.Printf("Name order (%d)\n", idxName.Len())
	idxName.Ascend(print)
	fmt.Printf("Number order (%d)\n", idxNum.Len())
	idxNum.Ascend(print)
	// Output:
	// Name order (3)
	//     Aramis    3
	//     Athos     1
	//     Porthos   2
	// Number order (3)
	//     Athos     1
	//     Porthos   2
	//     Aramis    3
}

func Example() {
	var rm recmgr.GrpType
	recList := []RecType{
		{name: "Brahms", num: 1833},
		{name: "Bach", num: 1685},
		{name: "Palestrina", num: 1525},
		{name: "Mozart", num: 1756},
		{name: "Schubert", num: 1797},
	}
	indent := func(fmtStr string, args ...interface{}) {
		fmt.Printf("    %s\n", fmt.Sprintf(fmtStr, args...))
	}
	print := func(recPtr interface{}) bool {
		if recPtr != nil {
			indent("%s", *recPtr.(*RecType))
		}
		return true
	}
	idxName := rm.Index(4, func(a, b interface{}) bool {
		return a.(*RecType).name < b.(*RecType).name
	})
	idxNum := rm.Index(4, func(a, b interface{}) bool {
		return a.(*RecType).num < b.(*RecType).num
	})
	for j := range recList {
		rm.ReplaceOrInsert(&recList[j])
	}
	fmt.Println("Name order")
	idxName.Ascend(print)
	fmt.Println("Number order")
	idxNum.Ascend(print)
	fmt.Println("Delete \"Palestrina\"")
	indent("keys deleted: %d", rm.Delete(idxName.Get(&RecType{name: "Palestrina"})))
	fmt.Println("Name order")
	idxName.Ascend(print)
	fmt.Println("Number order")
	idxNum.Ascend(print)
	fmt.Println("Number < 1797")
	idxNum.AscendLessThan(&RecType{num: 1797}, print)
	fmt.Println("Name >= \"Mozart\"")
	idxName.AscendGreaterOrEqual(&RecType{name: "Mozart"}, print)
	fmt.Println("1685 <= Number < 1797")
	idxNum.AscendRange(&RecType{num: 1685}, &RecType{num: 1797}, print)
	fmt.Println("Get 1756")
	print(idxNum.Get(&RecType{num: 1756}))
	fmt.Println("Get 1800")
	print(idxNum.Get(&RecType{num: 1800}))
	fmt.Println("Get \"Schubert\"")
	print(idxName.Get(&RecType{name: "Schubert"}))
	fmt.Println("Get \"Beethoven\"")
	print(idxName.Get(&RecType{name: "Beethoven"}))
	fmt.Println("Has 1756")
	indent("%v", idxNum.Has(&RecType{num: 1756}))
	fmt.Println("Has 1770")
	indent("%v", idxNum.Has(&RecType{num: 1770}))
	fmt.Println("Terminated")
	idxName.Ascend(func(recPtr interface{}) bool { return false })
	// Output:
	// Name order
	//     Bach..........1685
	//     Brahms........1833
	//     Mozart........1756
	//     Palestrina....1525
	//     Schubert......1797
	// Number order
	//     Palestrina....1525
	//     Bach..........1685
	//     Mozart........1756
	//     Schubert......1797
	//     Brahms........1833
	// Delete "Palestrina"
	//     keys deleted: 2
	// Name order
	//     Bach..........1685
	//     Brahms........1833
	//     Mozart........1756
	//     Schubert......1797
	// Number order
	//     Bach..........1685
	//     Mozart........1756
	//     Schubert......1797
	//     Brahms........1833
	// Number < 1797
	//     Bach..........1685
	//     Mozart........1756
	// Name >= "Mozart"
	//     Mozart........1756
	//     Schubert......1797
	// 1685 <= Number < 1797
	//     Bach..........1685
	//     Mozart........1756
	// Get 1756
	//     Mozart........1756
	// Get 1800
	// Get "Schubert"
	//     Schubert......1797
	// Get "Beethoven"
	// Has 1756
	//     true
	// Has 1770
	//     false
	// Terminated
}
