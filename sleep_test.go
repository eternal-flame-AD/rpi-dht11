package dht11

import (
	"testing"
	"time"
)

func absInt64(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

func TestSleepAccuracy(t *testing.T) {
	start := time.Now()
	sleep(100 * time.Microsecond)
	if diff := time.Now().Sub(start).Nanoseconds() - 100000; absInt64(diff) > 1000 {
		t.Fatalf("sleep time is not accurate enough: diff = %d nanosec", diff)
	}
}
