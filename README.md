# rpi-dht11

Library for reading DHT11 sensor from raspberry pi GPIO

## Usage

Read programatically
```golang
gpio, err := wiringpi.Setup(wiringpi.WiringPiSetup)
if err != nil {
    panic(err)
}

h, t, err := dht11.Read(gpio, 1)
if err != nil {
    panic(err)
}
fmt.Printf("T=%.1fdegC H=%.1f%%\n", t, h)
```

Command line util
```bash
$ go get github.com/eternal-flame-AD/rpi-dht11/cmd/dht11
$ GPIO_DHT11=1 dht11 # subsitute gpio index to wiringpi index
T=19.0degC H=53.0%
```