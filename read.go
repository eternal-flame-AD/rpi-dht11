package dht11

import (
	"errors"
	"time"

	wiringpi "github.com/eternal-flame-AD/go-wiringpi"
)

// nanoseconds for high interval threshold
const highThres = 50000

func verifyChecksum(data []byte) bool {
	sum := 0
	for i := 0; i < len(data)-1; i++ {
		sum += int(data[i])
		sum &= 0xff
	}
	return sum == int(data[len(data)-1])
}

func bitToBytes(data []bool) []byte {
	if len(data)%8 != 0 {
		panic("length is not a multiple of 8")
	}
	res := make([]byte, len(data)/8)
	for i := 0; i < len(res); i++ {
		for j := 0; j < 8; j++ {
			d := 0
			if data[i*8+j] {
				d = 1
			}
			res[i] += byte(d << (7 - uint(j)))
		}
	}
	return res
}

func expect(req func() bool, timeout time.Duration) bool {
	start := time.Now()
	for time.Now().Sub(start) < timeout {
		if req() {
			return true
		}
	}
	return false
}

func readBit(gpio *wiringpi.GPIO, pin int) (bool, error) {
	if !expect(func() bool {
		return gpio.DigitalRead(pin) == wiringpi.Low
	}, 100*time.Millisecond) {
		return false, errors.New("no response from dht11 while waiting for bit low start")
	}

	// start read
	if !expect(func() bool {
		return gpio.DigitalRead(pin) == wiringpi.High
	}, 100*time.Millisecond) {
		return false, errors.New("no response from dht11 while waiting for bit high")
	}

	start := time.Now()

	if !expect(func() bool {
		return gpio.DigitalRead(pin) == wiringpi.Low
	}, 1000*time.Millisecond) {
		return false, errors.New("no response from dht11 while waiting for bit end")
	}
	diff := time.Now().Sub(start).Nanoseconds()

	return diff > highThres, nil
}

// Read reads dht11
func Read(gpio *wiringpi.GPIO, pin int) (h, t float64, err error) {
	// start signal
	gpio.PinMode(pin, wiringpi.Output)
	gpio.DigitalWrite(pin, wiringpi.Low)
	sleep(20 * time.Millisecond)
	gpio.DigitalWrite(pin, wiringpi.High)
	sleep(40 * time.Microsecond)

	// wait for reply
	gpio.PinMode(pin, wiringpi.Input)
	gpio.Pull(pin, wiringpi.PullUp)
	if !expect(func() bool {
		return gpio.DigitalRead(pin) == wiringpi.Low
	}, 100*time.Millisecond) {
		return 0, 0, errors.New("no response from dht11")
	}
	if !expect(func() bool {
		return gpio.DigitalRead(pin) == wiringpi.High
	}, 100*time.Millisecond) {
		return 0, 0, errors.New("no response from dht11")
	}

	var datas [40]bool
	for i := 0; i < 40; i++ {
		data, err := readBit(gpio, pin)
		if err != nil {
			return 0, 0, err
		}
		datas[i] = data
	}

	dataBytes := bitToBytes(datas[:])
	if checksumPass := verifyChecksum(dataBytes); !checksumPass {
		return 0, 0, errors.New("checksum mismatch")
	}
	return calcH(dataBytes[0], dataBytes[1]), calcT(dataBytes[2], dataBytes[3]), nil
}
