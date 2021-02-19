package model

import (
	"testing"
	"unsafe"
)

type xxx struct {
	ID     string
	Medias []*Media
}

func TestUnsafe(t *testing.T) {

	list := make([]*Media, 0)
	list = append(list, &Media{
		Business:   "XXX",
		BusinessId: "1",
	})

	olist := make([]*xxx, 0)
	olist = append(olist, &xxx{
		ID: "1",
	})

	list2 := *(*[]*UnsafeMedia)(unsafe.Pointer(&olist))

	for _, v := range list2 {
		for _, vm := range list {
			if v.ID == vm.BusinessId {
				v.Medias = append(v.Medias, vm)
				// log.Fatal(list2)

			}
		}
	}
	for _, v := range list2 {
		t.Logf("%+v", v.Medias[0])
	}
}

type aaa struct {
	a int
}

func TestXXX(t *testing.T) {
	var a = make([]aaa, 0)
	a = append(a, aaa{
		a: 1,
	})
	a = append(a, aaa{
		a: 2,
	})
	for k, v := range a {
		a[k].a = v.a * 2
	}
	for _, v := range a {
		t.Log(v.a)
	}
}
