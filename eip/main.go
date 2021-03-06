package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/c-robinson/iplib"
)

// Extend the IP range by CIDR
// Usage: echo '1.2.3.4/24' | eip -s 32
// Usage: echo '1.2.3.4/24' | eip -p small

var (
	concurrency int
	sub         int
	port        string
	ports       []string
)

func main() {
	// cli arguments
	flag.IntVar(&concurrency, "c", 3, "concurrency ")
	flag.IntVar(&sub, "s", 32, "CIDR subnet (e.g: 24, 22)")
	flag.StringVar(&port, "p", "", "Append port after each IP (some predefined value: full, xlarge, large, small or f,x,l,s)")
	flag.Parse()

	if port != "" {
		ports = genPorts(port)
	}

	var wg sync.WaitGroup
	jobs := make(chan string, concurrency)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for job := range jobs {
			extendRange(job, sub)
		}
	}()

	sc := bufio.NewScanner(os.Stdin)
	go func() {
		for sc.Scan() {
			url := strings.TrimSpace(sc.Text())
			jobs <- url
		}
		close(jobs)
	}()
	wg.Wait()

}

func extendRange(rangeIP string, sub int) {
	_, ipna, err := iplib.ParseCIDR(rangeIP)
	if err != nil {
		ip := net.ParseIP(rangeIP)
		if ip != nil {
			if port == "" || sub != 32 {
				fmt.Println(ip)
			} else {
				for _, p := range ports {
					fmt.Printf("%s:%s\n", ip, p)
				}
			}
		}
		return
	}
	extendedIPs, err := ipna.Subnet(sub)
	if err != nil {
		return
	}
	for _, item := range extendedIPs {
		ip := item.String()
		if sub == 32 {
			ip = item.IP.String()
		}
		if port == "" || sub != 32 {
			fmt.Println(ip)
		} else {
			for _, p := range ports {
				fmt.Printf("%s:%s\n", ip, p)
			}
		}
	}
}

func genPorts(port string) []string {
	switch port {
	case "small":
		return []string{"80", "443", "8000", "8080", "8443"}
	case "s":
		return []string{"80", "443", "8000", "8080", "8443"}

	case "large":
		return []string{"80", "443", "81", "591", "2082", "2087", "2095", "2096", "3000", "8000", "8001", "8008", "8080", "8083", "8443", "8834", "8888"}
	case "l":
		return []string{"80", "443", "81", "591", "2082", "2087", "2095", "2096", "3000", "8000", "8001", "8008", "8080", "8083", "8443", "8834", "8888"}

	case "xlarge":
		return []string{"80", "443", "81", "300", "591", "593", "832", "981", "1010", "1311", "2082", "2087", "2095", "2096", "2480", "3000", "3128", "3333", "4243", "4567", "4711", "4712", "4993", "5000", "5104", "5108", "5800", "6543", "7000", "7396", "7474", "8000", "8001", "8008", "8014", "8042", "8069", "8080", "8081", "8083", "8088", "8090", "8091", "8118", "8123", "8172", "8222", "8243", "8280", "8281", "8333", "8443", "8500", "8834", "8880", "8888", "8983", "9000", "9043", "9060", "9080", "9090", "9091", "9200", "9443", "9800", "9981", "12443", "16080", "18091", "18092", "20720", "28017"}
	case "x":
		return []string{"80", "443", "81", "300", "591", "593", "832", "981", "1010", "1311", "2082", "2087", "2095", "2096", "2480", "3000", "3128", "3333", "4243", "4567", "4711", "4712", "4993", "5000", "5104", "5108", "5800", "6543", "7000", "7396", "7474", "8000", "8001", "8008", "8014", "8042", "8069", "8080", "8081", "8083", "8088", "8090", "8091", "8118", "8123", "8172", "8222", "8243", "8280", "8281", "8333", "8443", "8500", "8834", "8880", "8888", "8983", "9000", "9043", "9060", "9080", "9090", "9091", "9200", "9443", "9800", "9981", "12443", "16080", "18091", "18092", "20720", "28017"}
	case "full":
		var ports []string
		for i := 1; i <= 65535; i++ {
			ports = append(ports, fmt.Sprintf("%s", i))
		}
		return ports
	case "f":
		var ports []string
		for i := 1; i <= 65535; i++ {
			ports = append(ports, fmt.Sprintf("%d", i))
		}
		return ports
	default:
		return []string{"80", "443", "8000", "8080", "8443"}
	}

}
