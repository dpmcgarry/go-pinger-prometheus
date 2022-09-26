package cmd

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func pingDest(host string) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		panic(err)
	}

	pinger.Count = 10
	pinger.Timeout = time.Duration(10) * time.Second

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Println("Onrecv")
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Println("OnDuprecv")
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Println("OnFin")
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
}

// dopingCmd represents the doping command
var dopingCmd = &cobra.Command{
	Use:   "doping",
	Short: "Executes pings",
	Long:  `Execuates pings`,
	Run: func(cmd *cobra.Command, args []string) {
		destinations := viper.GetStringSlice("destinations")
		for _, dest := range destinations {
			fmt.Println(dest)
			pingDest(dest)
		}
	},
}

func init() {
	rootCmd.AddCommand(dopingCmd)

}
