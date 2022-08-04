package weakref

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type WeakReferable any

type WeakRef[T WeakReferable] struct {
	target uintptr
}

func New[T WeakReferable](ref *T) *WeakRef[T] {
	i := ^uintptr(unsafe.Pointer(&ref))
	weakRef := &WeakRef[T]{i}
	runtime.SetFinalizer(ref, func(_ *T) {
		atomic.StoreUintptr(&weakRef.target, uintptr(0))
	})
	return weakRef
}

func (w *WeakRef[T]) IsAlive() bool {
	return atomic.LoadUintptr(&w.target) != 0
}

func (w *WeakRef[T]) Deref() (v *T) {
	target := atomic.LoadUintptr(&w.target)
	if target != 0 {
		i := (*[1]uintptr)(unsafe.Pointer(&v))
		i[0] = ^target
	}
	return
}
