package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	wiringpi "github.com/eternal-flame-AD/go-wiringpi"
	dht11 "github.com/eternal-flame-AD/rpi-dht11"
)

var pin int

func init() {
	pinStr := os.Getenv("GPIO_DHT11")
	pinNum, err := strconv.Atoi(pinStr)
	if err != nil {
		panic(err)
	}
	pin = pinNum
}

func main() {
	gpio, err := wiringpi.Setup(wiringpi.WiringPiSetup)
	if err != nil {
		panic(err)
	}

	for retry := 0; retry < 3; retry++ {
		h, t, err := dht11.Read(gpio, 1)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("T=%.1fdegC H=%.1f%%\n", t, h)
		break
	}
}
