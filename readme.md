# WeakRef

## What, why?

A garbage collector is a complex mechanism. As long there is a reference to an object, it will not be cleared by the
garbage collector
But, if you have a recursive reference, it will never be cleared.

# Use case 1
Example; If you walk through all the nodes of a xml structure, it can be very handy to store the reference to the
parent, in each child.
However, once you're done with the xml structure, and leave the function, the garbage collector will never clear the
data. As the child depend on the parent, and the parent got a list of children. Creating a memory leak

To overcome this issue, a weak ref is one of these ways to fix it. Referring weakly in the child to the parent, will not
count as a reference in the eyes of the garbage collector, and will be cleared once the user discards of the xml
structure.

# Use case 2
It can be used to store things in memory as long there is space for it. Garbage collectors are notoriously lazy don't clean up until its absolutely needed.
This can be used in the programmers advantage too, as it makes it possible to store something in memory if it fits.
Think of some file buffer that you need all the time, but also don't explicitly want to store in memory so your code still works on machines with low memory,
it doesn't become a memory hog. Think of a http server that reads files from disk all the time. It makes sense to reuse the buffer of the file, if it isnt deleted yet. 

While you could do this with a cache that has a TTL, this enforces your server to have a specific memory usage over time, and it can run out. by using weakref's,
it can never run out as the garbage collector has space to delete   



## Weakrefs in other languages

Weakrefs isn't an uncommon solution, and is native in many other languages. Go sadly doesn't have this support natively

- [C#](https://docs.microsoft.com/en-us/dotnet/api/system.weakreference?view=net-6.0)
- [java](https://docs.oracle.com/javase/8/docs/api/java/lang/ref/WeakReference.html)
- [javascript](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WeakRef)
- [php](https://www.php.net/manual/en/class.weakreference.php)
- [python](https://docs.python.org/3/library/weakref.html)
- [rust](https://doc.rust-lang.org/std/rc/struct.Weak.html)

## How to use

````go
package main

import (
	"fmt"
	"github.com/jaenster/weakref"
	"runtime"
	"time"
)

type Point struct {
	X int
	Y int
}

func main() {

	point := &Point{5, 5}

	wr := weakref.New[Point](point)
	
	if wr.IsAlive() {
		fmt.Printf("Point still exists")
		fmt.Sprintln(wr.Deref())
	}

	// Remove the only reference to the point
	point = nil

	go (func() {
		// GC is lazy, sometimes it doesn't clear up right away.
		// Call it over and over to ensure point gets removed
		for {
			runtime.Gosched() // Run other goroutines
			runtime.GC()      // make sure the gc is called all the time
		}
	})()
	
	// Will be cleared in GC
	time.Sleep(750)
	if wr.IsAlive() == false {
		fmt.Printf("Point is removed from memory")
	}
}
````

## Credits

Token a lot of inspiration and the literal test from
[ivanrad/go-weakref](https://github.com/ivanrad/go-weakref).

While this wasn't production ready due to the lack of generics in go, this should be properly functional as the compiler
will be type aware and don't have runtime dependencies on interfaces
