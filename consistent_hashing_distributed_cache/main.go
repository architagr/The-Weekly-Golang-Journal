package main

import (
	cacheserver "consistent_hashing/cache_server"
	"hash/fnv"
	"log"
	"math/rand"
	"time"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Seed randomness for reproducibility
	rand.Seed(time.Now().UnixNano())

	// 🚀 Step 1: Initialize caching server with 5 nodes
	cache := cacheserver.InitCachingServer(fnv.New64a, 5)

	// 📝 Step 2: Put some data into the cache
	sampleData := map[string]string{
		"user:101": "Alice",
		"user:102": "Bob",
		"user:103": "Charlie",
		"user:104": "Daisy",
		"user:105": "Eve",
	}

	log.Println("\n===== Putting data into cache =====")
	for k, v := range sampleData {
		err := cache.Put(k, v)
		if err != nil {
			log.Printf("❌ Failed to put key %s: %v", k, err)
		} else {
			log.Printf("✅ Stored key '%s' with value '%s'", k, v)
		}
	}

	// 🕒 Wait a bit to allow simulated node failures
	log.Println("\n🕒 Waiting 12 seconds to allow node failures (simulated)...")
	time.Sleep(12 * time.Second)

	// 🔍 Step 3: Try retrieving keys (some may fail or trigger retry logic)
	log.Println("\n===== Getting data from cache =====")
	keysToGet := []string{"user:101", "user:102", "user:999", "user:104", "user:105"}
	for _, k := range keysToGet {
		val, err := cache.Get(k)
		if err != nil {
			log.Printf("❌ Failed to get key %s: %v", k, err)
		} else {
			log.Printf("✅ Retrieved key '%s' with value '%v'", k, val)
		}
	}

	// ❌ Step 4: Try deleting a key
	log.Println("\n===== Deleting a key from cache =====")
	err := cache.Delete("user:103")
	if err != nil {
		log.Printf("❌ Failed to delete key 'user:103': %v", err)
	} else {
		log.Println("✅ Deleted key 'user:103'")
	}

	// 🔁 Step 5: Try getting it again
	log.Println("\n===== Trying to get deleted key =====")
	val, err := cache.Get("user:103")
	if err != nil {
		log.Printf("Expected failure: %v", err)
	} else {
		log.Printf("Unexpectedly retrieved: %v", val)
	}

	log.Println("\n🏁 Done with demo")
}
