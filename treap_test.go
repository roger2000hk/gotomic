
package gotomic

import (
	"testing"
	"runtime"
	"reflect"
	"math/rand"
	"fmt"
)

type s string
func (self s) Compare(t Thing) int {
	if other, ok := t.(s); ok {
		return compStrings(string(self), string(other))
	}
	panic(fmt.Errorf("%#v can only compare to other s's, not %#v of type %T", self, t, t))
}

func fiddleTreap(t *testing.T, treap *Treap, x string, do, done chan bool) {
	<- do
	n := int(10000 + rand.Int31() % 1000)
	vals := make([]s, n)
	for i := 0; i < n; i++ {
		v := s(fmt.Sprint(rand.Int63(), ".", i, ".", x))
		vals[i] = v
		_, ok := treap.Put(v, v)
		if ok {
			fmt.Printf("1 %v should not contain %v\n", treap.Describe(), v)
			t.Fatalf("1 %v should not contain %v\n", treap.Describe(), v)
		}
 		value, ok := treap.Get(v)
		if !ok {
			fmt.Printf("2 %v should contain %v\n", treap.Describe(), v)
			t.Fatalf("2 %v should contain %v\n", treap.Describe(), v)
		}
		if v.Compare(value) != 0 {
			fmt.Printf("3 %v should contain %v\n", treap.Describe(), v)
			t.Fatalf("3 %v should contain %v\n", treap.Describe(), v)
		}
	} 
	for i := 0; i < n; i++ {
		v := vals[i]
		old, ok := treap.Delete(v)
		if !ok {
			fmt.Printf("4 %v should contain %v\n", treap.Describe(), v)
			t.Fatalf("4 %v should contain %v\n", treap.Describe(), v)
		}
		if old != v {
			fmt.Printf("5 %v should contain %v\n", treap.Describe(), v)
			t.Fatalf("5 %v should contain %v\n", treap.Describe(), v)
		}
		_, ok = treap.Get(v)
		if ok {
			fmt.Printf("6 %v should not contain %v\n", treap.Describe(), v)
			t.Fatalf("6 %v should not contain %v\n", treap.Describe(), v)
		}
	}
	done <- true
}

func TestTreapConc(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	treap := NewTreap()
	for i := 9; i >= 0; i-- {
		v := s(fmt.Sprint(i))
		treap.Put(v, v)
	}
	assertTreapSlice(t, treap, []Comparable{s("0"), s("1"), s("2"), s("3"), s("4"), s("5"), s("6"), s("7"), s("8"), s("9")}, []Thing{s("0"), s("1"), s("2"), s("3"), s("4"), s("5"), s("6"), s("7"), s("8"), s("9")})
	do := make(chan bool)
	done := make(chan bool)
	for i := 0; i < runtime.NumCPU(); i++ {
		go fiddleTreap(t, treap, fmt.Sprint("fiddler-", i, "-"), do, done)
	}
	close(do)
	for i := 0; i < runtime.NumCPU(); i++ {
		<-done
	}
	assertTreapSlice(t, treap, []Comparable{s("0"), s("1"), s("2"), s("3"), s("4"), s("5"), s("6"), s("7"), s("8"), s("9")}, []Thing{s("0"), s("1"), s("2"), s("3"), s("4"), s("5"), s("6"), s("7"), s("8"), s("9")})
}

func TestTreapPreviousNext(t *testing.T) {
	treap := NewTreap()
	for i := 9; i >= 0; i-- {
		treap.Put(c(i), fmt.Sprint(i))
	}
	assertTreapSlice(t, treap, []Comparable{c(0), c(1), c(2), c(3), c(4), c(5), c(6), c(7), c(8), c(9)}, []Thing{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}) 
	k, v, ok := treap.Next(c(4))
	if !ok {
		t.Error("should have something after 4")
	}
	if k != c(5) {
		t.Error("5 should be after 4")
	}
	if v != "5" {
		t.Error("5 should be after 4")
	}
	k, v, ok = treap.Previous(c(7))
	if !ok {
		t.Error("should have something before 7")
	}
	if k != c(6) {
		t.Error("6 should be before 7")
	}
	if v != "6" {
		t.Error("6 should be before 7")
	}
	k, v, ok = treap.Previous(c(0))
	if ok {
		t.Error("should not have anything before 0")
	}
	k, v, ok = treap.Next(c(9))
	if ok {
		t.Error("should not have anything after 9")
	}
}

func TestTreapPutGetDelete(t *testing.T) {
	treap := NewTreap()
	_, ok := treap.Get(c(3))
	if ok {
		t.Error("should not contain 3")
	}
	treap.Put(c(3), 44)
	v, ok := treap.Get(c(3))
	if !ok {
		t.Error("should contain 3")
	}
	if v != 44 {
		t.Error("should be 44")
	}
	v, ok = treap.Delete(c(3))
	if !ok {
		t.Error("should contain 3")
	}
	if v != 44 {
		t.Error("should be 44")
	}
	v, ok = treap.Get(c(3))
	if ok {
		t.Error("should not contain 3")
	}
	v, ok = treap.Delete(c(3))
	if v == 44 {
		t.Error("should not be 44")
	}
	if ok {
		t.Error("should not contain 3")
	}
}

func assertTreapSlice(t *testing.T, treap *Treap, keys []Comparable, values []Thing) {
	found_keys, found_values := treap.ToSlice()
	if !reflect.DeepEqual(keys, found_keys) {
		t.Errorf("%v.ToSlice keys should be %#v but was %#v", treap, keys, found_keys)
	}
	if !reflect.DeepEqual(values, found_values) {
		t.Errorf("%v.ToSlice values should be %#v but was %#v", treap, values, found_values)
	}
}

func TestTreapToSlice(t *testing.T) {
	treap := NewTreap()
	treap.Put(c(4), "4")
	treap.Put(c(6), "6")
	treap.Put(c(1), "1")
	treap.Put(c(8), "8")
	treap.Put(c(5), "5")
	assertTreapSlice(t, treap, []Comparable{c(1), c(4), c(5), c(6), c(8)}, []Thing{"1", "4", "5", "6", "8"})
}

func TestTreapMin(t *testing.T) {
	treap := NewTreap()
	k, v, ok := treap.Min()
	if ok {
		t.Error("should not have min value")
	}
	treap.Put(c(3), "3")
	k, v, ok = treap.Min()
	if !ok {
		t.Error("should have min value")
	}
	if k != c(3) {
		t.Error("min should be 3")
	}
	if v != "3" {
		t.Error("min should be 3")
	}
	treap.Put(c(2), "2")
	k, v, ok = treap.Min()
	if !ok {
		t.Error("should have min value")
	}
	if k != c(2) {
		t.Error("min should be 2")
	}
	if v != "2" {
		t.Errorf("min should be 2, not %#v", v)
	}
	treap.Put(c(4), "4")
	k, v, ok = treap.Min()
	if !ok {
		t.Error("should have min value")
	}
	if k != c(2) {
		t.Error("min should be 2")
	}
	if v != "2" {
		t.Error("min should be 2")
	}
}

func TestTreapMax(t *testing.T) {
	treap := NewTreap()
	k, v, ok := treap.Max()
	if ok {
		t.Error("should not have max value")
	}
	treap.Put(c(3), "3")
	k, v, ok = treap.Max()
	if !ok {
		t.Error("should have max value")
	}
	if k != c(3) {
		t.Error("max should be 3")
	}
	if v != "3" {
		t.Error("max should be 3")
	}
	treap.Put(c(2), "2")
	k, v, ok = treap.Max()
	if !ok {
		t.Error("should have max value")
	}
	if k != c(3) {
		t.Error("max should be 3")
	}
	if v != "3" {
		t.Errorf("max should be 3, not %#v", v)
	}
	treap.Put(c(4), "4")
	k, v, ok = treap.Max()
	if !ok {
		t.Error("should have max value")
	}
	if k != c(4) {
		t.Error("max should be 4")
	}
	if v != "4" {
		t.Error("max should be 4")
	}
}