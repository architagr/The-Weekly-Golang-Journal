package main

import (
	"fmt"
	"sync"
	"sync_pool/unidirectional"
	"time"
)

type Cache struct {
	Id int
}

func main() {
	basicExample()
	fmt.Println("------------")
	unidirectional.Run()
}

func basicExample() {
	countOfCacheInstances := 0
	objPool := sync.Pool{
		New: func() any {
			// just to add a artificial delay
			time.Sleep(500 * time.Millisecond)
			// this is just to check how many time has this code been called
			// which indicates how many instances of Cache are created
			countOfCacheInstances++

			fmt.Println("cache object created")
			return &Cache{
				// assign id to each object
				// to know if this same instance is been reused
				Id: countOfCacheInstances,
			}
		},
	}

	// Try to get an object instance from the pool
	// if no avaibale instance then it will call the New

	obj := objPool.Get().(*Cache) // -- 1️⃣
	// do some taks using obj
	obj2 := objPool.Get().(*Cache) // - 2️⃣
	// do some taks using obj2
	objPool.Put(obj)               // -- 3️⃣
	objPool.Put(obj2)              // -- 4️⃣
	obj3 := objPool.Get().(*Cache) // -- 5️⃣
	// do some taks using obj3
	obj4 := objPool.Get().(*Cache) // -- 6️⃣
	// do some taks using obj4

	fmt.Println(obj)
	fmt.Println(obj2)
	fmt.Println(obj3)
	fmt.Println(obj4)

	fmt.Println("number of cache obje created", countOfCacheInstances)
}
