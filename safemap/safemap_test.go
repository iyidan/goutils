package safemap

import (
	"fmt"
	"testing"
)

func TestSafeMap(t *testing.T) {
	m := New()
	for i := 0; i < 100; i++ {
		go func(i int) {
			m.Set(i, i)
			fmt.Printf("Exist m[%d]=%v\n", i, m.Exist(i))
			fmt.Printf("Get m[%d]=%d\n", i, m.Get(i))
			fmt.Printf("Del m[%d]\n", i)
			m.Del(i)
			fmt.Printf("Exist m[%d]=%v\n", i, m.Exist(i))
			fmt.Printf("Get m[%d]=%d\n", i, m.Get(i))
		}(i)
	}
}

func TestSafeMap_Len(t *testing.T) {
	m := New()
	if m.Len() != 0 {
		t.Fatal("m.len != 0")
	}
	m.Add("test", "test")
	m.Add("test1", "k")
	if m.Len() != 2 {
		t.Fatal("m.len != 2")
	}
	m.Del("test")
	if m.Len() != 1 {
		t.Fatal("m.len != 1")
	}
}

func TestSafeMap_Add(t *testing.T) {
	m := New()

	if m.Add("test", "test") != true {
		t.Fatal(`m.Add("test", "test") != true`)
	}

	if m.Add("test", "test2") != false {
		t.Fatal(`m.Add("test", "test2") != false`)
	}

	m.Del("test")

	if m.Add("test", "test3") != true {
		t.Fatal(`m.Add("test", "test3") != true`)
	}
}

func TestSafeMap_CasSet(t *testing.T) {
	m := New()

	if !m.CasSet("test", 1, 0) {
		t.Fatal("casset failed")
	}

	if !m.CasSet("test", 2, 1) {
		t.Fatal("casset failed")
	}

	if m.Get("test").(int) != 2 {
		t.Fatal("casset failed")
	}

	if m.CasSet("test", 3, 1) {
		t.Fatal("casset failed")
	}
}

func TestSafeMap_CasMultiSet(t *testing.T) {
	m := New()
	m.Set("t1", 1)
	m.Set("t2", 2)

	old := map[interface{}]interface{}{"t1": 1, "t2": 2}
	now := map[interface{}]interface{}{"t1": 3, "t2": 4}

	if !m.CasMultiSet(now, old) {
		t.Fatal(`!m.CaseMultiSet(now, old)`)
	}
	if m.Get("t1").(int) != 3 {
		t.Fatal(`m.Get("t1").(int) != 3`)
	}
	if m.Get("t2").(int) != 4 {
		t.Fatal(`m.Get("t2").(int) != 4`)
	}

	if m.CasMultiSet(now, old) {
		t.Fatal(`m.CaseMultiSet(now, old)`)
	}
}

func TestSafeMap_GetAll(t *testing.T) {
	m := New()

	m.Add("k", "v")
	m.Add("k2", "v2")
	m.Add("k3", "v3")

	cmap := m.GetAll()

	t.Log(cmap)

	if len(cmap) != 3 {
		t.Fatal("len(cmap) != 3")
	}
}

func TestSafeMap_DelAll(t *testing.T) {
	m := New()

	m.Add("k", "v")
	m.Add("k2", "v2")
	m.Add("k3", "v3")

	m.DelAll()

	cmap := m.GetAll()
	t.Log(cmap)

	if len(cmap) != 0 {
		t.Fatal("len(cmap) != 0")
	}
}
