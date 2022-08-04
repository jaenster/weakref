package weakref

import (
	"log"
	"runtime"
	"testing"
)

type Point3D struct {
	x, y, z int
}

func makeWeakReference() *WeakRef[Point3D] {
	point := Point3D{1, 2, 3}
	weakRef := New[Point3D](&point)
	return weakRef
}

func TestWeakRef(t *testing.T) {
	weakRef := makeWeakReference()
	if weakRef == nil {
		t.Fatal("unexpected: weakRef is nil")
	}
	if weakRef.Deref() == nil {
		t.Fatal("unexpected: weakRef.Deref() is nil")
	}
	t.Log(weakRef.Deref())
	if weakRef.IsAlive() == false {
		t.Fatal("unexpected: weakRef.IsAlive() is false")
	}

	// 10 loops to make sure its gone (taken from example code)
	for i := 1; i < 10; i++ {
		runtime.Gosched()
		runtime.GC()
	}

	if weakRef == nil {
		log.Fatal("unexpected: weakRef is nil after GC")
	}
	if weakRef.Deref() != nil {
		log.Fatal("unexpected: weakRef.Deref() is not nil after GC")
	}
	if weakRef.IsAlive() == true {
		log.Fatal("unexpected: weakRef.IsAlive() is true")
	}
}
