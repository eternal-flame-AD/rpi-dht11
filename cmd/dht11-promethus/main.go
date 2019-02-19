package main

import (
	"flag"
	"log"
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

func collect(gpio *wiringpi.GPIO) {
	for retry := 0; retry < 3; retry++ {
		h, t, err := dht11.Read(gpio, pinNum)
		if err == nil {
			tempGauge.Set(t)
			humidGauge.Set(h)
			break
		} else {
			log.Println(err)
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	gpio, err := wiringpi.Setup(wiringpi.WiringPiSetup)
	if err != nil {
		panic(err)
	}

	collect(gpio)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			collect(gpio)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
