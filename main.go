package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:",omitempty"`
	CreatedAt time.Time      `json:",omitempty"`
	UpdatedAt time.Time      `json:",omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:",omitempty"`
}

type Hosts struct {
	ID        uint
	IPAddress string
	Name      string
}

type Pings struct {
	Model
	ResultID  int `json:",omitempty"`
	Sequence  int
	IPAddress string
	Duration  time.Duration
	Date      time.Time
}
type PingResults struct {
	Model
	IPAddress  string
	Name       string
	Date       time.Time
	Sent       int
	Recv       int
	PingMin    time.Duration
	PingAvg    time.Duration
	PingMax    time.Duration
	PingStdDev time.Duration
	PacketLoss float64
	Pings      []Pings `json:"pings,omitempty" gorm:"foreignKey:ResultID;references:ID"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var pingResults []PingResults
	operatingsystem := runtime.GOOS
	timeout := flag.Duration("t", 5, "timeout in seconds")
	parse := flag.Bool("p", false, "parse json files")
	flag.Usage = func() {
		color.Green("Usage: %s [options] [file]", os.Args[0])
	}
	flag.Parse()

	f, err := os.Open("hosts.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	if *parse {
		Parse()
	} else {
		fmt.Printf("Packet Loss Tester\n")

		csvReader := csv.NewReader(f)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		hosts := []Hosts{}
		for _, d := range data {
			hosts = append(hosts, Hosts{
				IPAddress: d[0],
				Name:      d[1],
			})
		}

		const threshold = 75
		for _, d := range hosts {
			var pings []Pings
			fmt.Println("\n----------------------------------------------------------------")
			count := 1
			date := time.Now().Format("2006-01-02 15:04:05")
			pinger, err := ping.NewPinger(d.IPAddress)
			if operatingsystem == "windows" {
				pinger.SetPrivileged(true)
			}

			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}

			pinger.Timeout = *timeout
			pinger.Count = int(*timeout / time.Millisecond)

			color.White("Starting ping test (%v attempts) to %s (%s) - %s\n", (pinger.Count / 1000), d.IPAddress, d.Name, date)

			pinger.OnRecv = func(pkt *ping.Packet) {
				pings = append(pings, Pings{Sequence: count, IPAddress: d.IPAddress, Duration: pkt.Rtt, Date: time.Now()})
				if pkt.Rtt > time.Duration(threshold)*time.Millisecond {
					color.Red("\t%v - %d bytes from %s: icmp_seq=%d time=%v <-- high ping\n", count, pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
				} else {
					color.White("\t%v - %d bytes from %s: icmp_seq=%d time=%v\n", count, pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
				}
				count++
			}

			pinger.OnFinish = func(stats *ping.Statistics) {
				var pingResult = PingResults{
					Date:       time.Now(),
					IPAddress:  d.IPAddress,
					Name:       d.Name,
					Sent:       stats.PacketsSent,
					Recv:       stats.PacketsRecv,
					PingMin:    stats.MinRtt,
					PingAvg:    stats.AvgRtt,
					PingMax:    stats.MaxRtt,
					PingStdDev: stats.StdDevRtt,
					PacketLoss: stats.PacketLoss,
					Pings:      pings,
				}

				color.White("\n--- %s ping statistics ---\n", d.IPAddress)

				if stats.PacketLoss > 0 {
					color.Red("\t%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss detected\n",
						stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
				} else {
					color.White("\t%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
						stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
				}

				color.White("\tround-trip min/avg/max/stddev = %v/%v/%v/%v\n",
					stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

				pingResults = append(pingResults, pingResult)
			}

			err = pinger.Run()
			if err != nil {
				panic(err)
			}
		}

		color.White("\n----------------------------------------------------------------\n\n")

		b, err := json.Marshal(pingResults)
		if err != nil {
			fmt.Println(err)
			return
		}

		path := "output/"
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		write_file := path + "results_" + time.Now().Format("2006-01-02_15:04:05") + ".json"
		f, err = os.Create(write_file)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer f.Close()

		err = os.WriteFile(write_file, b, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
