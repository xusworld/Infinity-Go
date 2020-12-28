package benchmark

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func batchWriteOnlyBenchmark(hashMaps []HashMap, goroutinesConf []int, writeTimes int, hint []string) {
	for _, goroutinesNum := range goroutinesConf {
		durations := make([]time.Duration, 0)

		for _, m := range hashMaps {
			durations = append(durations, writeOnlyBenchmark(m, goroutinesNum, writeTimes))
		}

		fmt.Printf("Concurrency %d, WirteTimes %d\n", goroutinesNum, writeTimes)
		for idx, duration := range durations {
			fmt.Printf("TimeCost %s [%s]  \n", duration.String(), hint[idx])
		}
		fmt.Println()
		time.Sleep(time.Second)
	}
}

func writeOnlyBenchmark(m HashMap, goroutinesNum int, writeTimes int) time.Duration {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(goroutinesNum)

	for i := 0; i < goroutinesNum; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < writeTimes; j++ {
				m.Set("Index "+strconv.Itoa(j), "MagicNumber")
			}
		}()
	}
	wg.Wait()

	return time.Now().Sub(start)
}



func TestWriteOnly(t *testing.T) {
	sm := &SyncMap{}
	rwm := NewRWMutexMap()
	mm := NewMutexMap()
	cm := NewConcurrentMap()

	// 并发读取数据的goroutines数量
	concurrency := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
	// 单个goroutines中循环读取数据的次数
	writeTimes := 10000
	batchWriteOnlyBenchmark([]HashMap{sm, rwm, mm, cm}, concurrency, writeTimes,
		[]string{"Sync.Map", "RWMutexMap", "MutexMap", "ConcurrentMap"})

}

func TestWriteOnlyV1(t *testing.T) {
	sm := &SyncMap{}
	rwm := NewRWMutexMap()
	mm := NewMutexMap()
	cm := NewConcurrentMap()

	// 并发读取数据的goroutines数量
	concurrency := []int{1000}
	// 单个goroutines中循环读取数据的次数
	writeTimes := 10000
	batchWriteOnlyBenchmark([]HashMap{sm, rwm, mm, cm}, concurrency, writeTimes,
		[]string{"Sync.Map", "RWMutexMap", "MutexMap", "ConcurrentMap"})

}