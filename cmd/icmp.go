package cmd

import (
	"bufio"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	DEFAULT_PRIVILEGED = runtime.GOOS == "windows"
)

var count int

// icmpCmd represents the icmp command
var icmpCmd = &cobra.Command{
	Use:   "icmp",
	Short: "icmp gfping network",
	Long:  `icmp gfping network`,
	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(runtime.NumCPU())

		var ipList []string

		if len(file) != 0 && (len(subnet) != 0 || len(singleip) != 0) || len(subnet) != 0 && (len(file) != 0 || len(singleip) != 0) || len(singleip) != 0 && (len(file) != 0 || len(subnet) != 0) {
			fmt.Println("can't specify both -f -g -i.")
			cmd.Help()
			return
		}

		// 检查输出文件
		if len(output) != 0 {
			if _, err := os.Stat(output); err != nil {
				f, err := os.Create(output)
				defer f.Close()
				if err != nil {
					fmt.Println("can't create output result file")
					return
				}
			}
			outputFile, err := os.Open(output)
			if err != nil {
				fmt.Println("can't open output result file")
				return
			}
			outputFile.Close()
		}

		if len(file) != 0 {
			fi, err := os.Open(file)
			if err != nil {
				fmt.Printf("can't open ip list file: %s\n", file)
				return
			}
			defer fi.Close()

			br := bufio.NewReader(fi)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				ipList = append(ipList, string(a))
			}
		}

		if len(subnet) != 0 {
			ips, err := subNetGet(subnet)
			if err != nil {
				panic(err)
			}
			ipList = ips
		}

		if len(singleip) != 0 {
			ipList = append(ipList, singleip)
		}

		var reachableIps []string
		var unreachableIps []string

		// var lock sync.Mutex
		var wg sync.WaitGroup

		buckets := make(chan bool, routinepool)
		for _, ip := range ipList {
			buckets <- true
			wg.Add(1)
			go func(ip string) {
				result, err := pingAlive(ip)
				// lock.Lock()
				// defer lock.Unlock()
				if err != nil {
					// unreachableIps = append(unreachableIps, ip)
					fmt.Printf("%s is unreachable\n", ip)
				} else {
					if result == true {
						// reachableIps = append(reachableIps, ip)
						fmt.Printf("%s is alive\n", ip)
					} else {
						time.Sleep(time.Microsecond * 10)
						secresult, err := pingAlive(ip)
						if err == nil && secresult == true {
							// reachableIps = append(reachableIps, ip)
							fmt.Printf("%s is alive\n", ip)
						} else {
							// unreachableIps = append(unreachableIps, ip)
							fmt.Printf("%s is unreachable\n", ip)
						}
					}
				}
				<-buckets
				wg.Done()
			}(ip)
		}
		wg.Wait()

		if len(output) != 0 {
			outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				fmt.Println("can't open output result file")
				return
			}
			defer outputFile.Close()

			for _, rip := range reachableIps {
				outputFile.WriteString(rip + ": ok\n")
			}
			for _, urip := range unreachableIps {
				outputFile.WriteString(urip + ": failed\n")
			}
			fmt.Printf("output result to file: %s \n", output)
		} else {
			// TODO
			// fmt.Printf("Reachable IP: %v\n", reachableIps)
			// fmt.Printf("UnReachable IP: %v\n", unreachableIps)
		}
	},
}

func subNetGet(subnet string) ([]string, error) {
	var ips []string
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return ips, err
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func pingAlive(ip string) (bool, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return false, err
	}
	pinger.Count = count
	pingTimeOut := time.Duration(timeout)
	pinger.Timeout = time.Duration(pingTimeOut * time.Millisecond)
	// TODO
	// pinger.SetPrivileged(true) // windows need
	pinger.Run()
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		return true, nil
	}
	return false, nil
}

func init() {
	rootCmd.AddCommand(icmpCmd)
	icmpCmd.Flags().IntVarP(&count, "count", "n", 3, "send N pings to each target")
}
