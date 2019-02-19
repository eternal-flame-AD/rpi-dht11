package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"net/http"
	"time"

	wiringpi "github.com/eternal-flame-AD/go-wiringpi"
	dht11 "github.com/eternal-flame-AD/rpi-dht11"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var listen string
var pinNum int

func init() {
	prometheus.MustRegister(tempGauge)
	prometheus.MustRegister(humidGauge)
	l := flag.String("l", ":12800", "listen port")
	pin := flag.Uint("p", 1, "dht11 pin number")
	flag.Parse()
	listen = *l
	pinNum = int(*pin)
}

var tempGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "rpi_atmospheric",
	Subsystem: "dht11",
	Name:      "temp",
	Help:      "Temperature data collected from dht11 on RPI",
})

var humidGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "rpi_atmospheric",
	Subsystem: "dht11",
	Name:      "humidity",
	Help:      "Humidity data collected from dht11 on RPI",
})

func obtainDataOnce(gpio *wiringpi.GPIO) (h, t float64, err error) {
	for retry := 0; retry < 3; retry++ {
		h, t, err := dht11.Read(gpio, pinNum)
		if err == nil {
			return h, t, nil
		} else {
			log.Println(err)
			time.Sleep(1500 * time.Millisecond)
		}
	}
	return 0, 0, errors.New("too many fails, giving up")
}

var lastH float64
var lastT float64

func collect(gpio *wiringpi.GPIO, init bool) {
	h, t, err := obtainDataOnce(gpio)
	if err != nil {
		return
	}

	if init {
		lastH = h
		lastT = t
	}

	if math.Abs(h-lastH) > 10 || math.Abs(t-lastT) > 3 {
		newH, newT, err := obtainDataOnce(gpio)
		if err != nil {
			return
		}
		if math.Abs(newH-h) > 10 || math.Abs(newT-t) > 3 {
			return
		}
	}

	lastH = h
	lastT = t
	tempGauge.Set(t)
	humidGauge.Set(h)
}

func main() {
	gpio, err := wiringpi.Setup(wiringpi.WiringPiSetup)
	if err != nil {
		panic(err)
	}

	collect(gpio, true)
	go func() {
		ticker := time.NewTicker(12 * time.Second)
		for range ticker.C {
			collect(gpio, false)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
