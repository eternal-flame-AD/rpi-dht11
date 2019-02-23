package dht11

import "time"

func sleep(t time.Duration) {
	start := time.Now()
	if t > 200*time.Microsecond {
		time.Sleep(t)
	}
	for time.Now().Sub(start) < t {
	}
}
