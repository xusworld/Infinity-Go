package benchmark

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func batchReadWriteBenchmark(hashMaps []HashMap, goroutinesConf []int, N int, readPercentage float32, hint []string) {
	for _, goroutinesNum := range goroutinesConf {
		durations := make([]time.Duration, 0)

		for _, m := range hashMaps {
			durations = append(durations, readWriteBenchmark(m, goroutinesNum, N, readPercentage))
		}

		readTimes := cast.ToInt(cast.ToFloat32(N) * readPercentage)
		writeTimes := cast.ToInt(cast.ToFloat32(N) * (1.0 - readPercentage))
		fmt.Printf("Concurrency %d, ReadTimes %d, WirteTimes %d\n", goroutinesNum, readTimes, writeTimes)

		for idx, duration := range durations {
			fmt.Printf("TimeCost %s [%s]  \n", duration.String(), hint[idx])
		}
		fmt.Println()
		time.Sleep(time.Second)
	}
}

func readWriteBenchmark(m HashMap, goroutinesNum int, N int, readPercentage float32) time.Duration {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(goroutinesNum)

	for i := 0; i < goroutinesNum; i++ {
		if i%2 == 0 {
			// read
			go func() {
				defer wg.Done()
				for j := 0; j < int(float32(N)*readPercentage); j++ {
					_ = m.Get("MagicNumber")
				}
			}()
		} else {
			// write
			go func() {
				defer wg.Done()
				for j := 0; j < int(float32(N)*(1.0-readPercentage)); j++ {
					m.Set("Index "+strconv.Itoa(j), "MagicNumber")
				}
			}()
		}
	}

	wg.Wait()
	return time.Now().Sub(start)
}

func TestReadWrite(t *testing.T) {
	sm := &SyncMap{}
	rwm := NewRWMutexMap()
	mm := NewMutexMap()
	cm := NewConcurrentMap()

	// 并发读取数据的goroutines数量
	concurrency := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	// 单个goroutines中循环总的读写数量
	N := 10000
	batchReadWriteBenchmark([]HashMap{sm, rwm, mm, cm}, concurrency, N, 0.05,
		[]string{"Sync.Map", "RWMutexMap", "MutexMap", "ConcurrentMap"})
}
