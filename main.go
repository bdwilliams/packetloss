package main

import (
	"fmt"

	"github.com/go-ping/ping"
)

type PingResults struct {
	Sent       int
	Recv       int
	PingMin    float64
	PingAvg    float64
	PingMax    float64
	PingStdDev float64
	PacketLoss float64
}

func main() {
	hosts := map[string]string{
		"8.8.8.8":         "Google DNS",
		"1.1.1.1":         "Cloudflare",
		"208.67.222.222":  "OpenDNS",
		"69.162.81.155":   "Dallas, TX",
		"192.199.248.75":  "Denver, CO",
		"162.254.206.227": "Miami, FL",
		"209.142.68.29":   "Chicago, IL",
		"207.250.234.100": "Minneapolis, MN",
		"206.71.50.230":   "New York, NY",
		"65.49.22.66":     "San Francisco, CA",
		"23.81.0.59":      "Seattle, WA",
	}

	var pingResults []PingResults
	for h, l := range hosts {
		fmt.Printf("Starting ping test to %s (%s)\n", h, l)
		pinger, err := ping.NewPinger(h)
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}

		pinger.Count = 10

		pinger.OnRecv = func(pkt *ping.Packet) {
			// fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			var pingResult = PingResults{
				Sent:       stats.PacketsSent,
				Recv:       stats.PacketsRecv,
				PingMin:    stats.MinRtt.Seconds(),
				PingAvg:    stats.AvgRtt.Seconds(),
				PingMax:    stats.MaxRtt.Seconds(),
				PingStdDev: stats.StdDevRtt.Seconds(),
				PacketLoss: stats.PacketLoss,
			}

			if stats.PacketLoss > 0 {
				fmt.Println("\tPacket loss detected:", stats.PacketLoss)
			}

			pingResults = append(pingResults, pingResult)
		}

		err = pinger.Run()
		if err != nil {
			panic(err)
		}
	}

	// fmt.Println(pingResults)
}
