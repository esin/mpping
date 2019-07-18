package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

func getCurrentTimeStamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func lookupIPAddr(poolAddr string) (string, string) {
	ips, err := net.LookupIP(poolAddr)
	if err != nil {
		fmt.Printf("%s: %v\n", poolAddr, err)
		//os.Exit(1)
	}
	ipV4 := ""
	ipV6 := ""
	for _, ip := range ips {
		if govalidator.IsIPv4(ip.String()) && ipV4 == "" {
			ipV4 = ip.String()
		}
		if govalidator.IsIPv6(ip.String()) && ipV6 == "" {
			ipV6 = ip.String()
		}
	}
	return ipV4, ipV6
}

type poolStruct struct {
	PoolDomain           string `json:"pooldomain"`
	PoolPort             uint16 `json:"port"`
	PoolIPv4             string `json:"ipv4"`
	PoolIPv6             string `json:"ipv6"`
	PoolScheme           string
	TotalPacketsSent     uint64
	TotalPacketsReceived uint64
	TotalTimeMin         uint64
	TotalTimeMax         uint64
	TotalTime            uint64
}

func checkForPoolAddr(urlArg string) (poolStruct, bool) {

	var newPool poolStruct

	strings.TrimLeft(urlArg, " ")

	if strings.Index(strings.ToLower(urlArg), "://") < 0 {
		urlArg = "stratum://" + urlArg
	}

	parsed, err := url.Parse(urlArg)
	if err != nil {
		return newPool, false
	}

	newPool.PoolDomain = parsed.Hostname()
	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		//fmt.Printf("Crazy port: %v", parsed.Port())
		return newPool, false
	}
	newPool.PoolPort = uint16(port)
	newPool.PoolScheme = parsed.Scheme

	if govalidator.IsHost(newPool.PoolDomain) {
		newPool.PoolIPv4, newPool.PoolIPv6 = lookupIPAddr(newPool.PoolDomain)
	}

	if govalidator.IsIPv4(newPool.PoolDomain) {
		newPool.PoolIPv4 = newPool.PoolDomain
	}

	if govalidator.IsIPv6(newPool.PoolDomain) {
		newPool.PoolIPv6 = newPool.PoolDomain
	}

	if newPool.PoolIPv4 == "" && newPool.PoolIPv6 == "" {
		return newPool, false
	}

	return newPool, true
}

var poolList []poolStruct

func onStop() {

	for _, poolServer := range poolList {
		fmt.Printf("TIME total/min/max/avg\t\t%d ms/%d ms/%d ms/%d ms\n", poolServer.TotalTime, poolServer.TotalTimeMin, poolServer.TotalTimeMax, poolServer.TotalTime/poolServer.TotalPacketsReceived)
		fmt.Printf("PACKETS sent/received\t\t %d/%d\n", poolServer.TotalPacketsSent, poolServer.TotalPacketsReceived)
	}

	os.Exit(0)
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		onStop()
	}()

	request := `{"id":1,"method":"mining.subscribe","params":["mpping-0.1","EthereumStratum/2.0.0"]}`

	ipv4ProtoFlag := flag.Bool("4", true, "Ping only pool IP version 4")
	ipv6ProtoFlag := flag.Bool("6", false, "Ping only pool IP version 6")
	countPacketsFlag := flag.Int("count", 0, "Packets count. Inifinity by default")

	flag.Parse()
	poolAddr := ""
	poolPort := ""
	poolIPAddr := ""
	ipv4Proto := *ipv4ProtoFlag
	ipv6Proto := *ipv6ProtoFlag
	countPackets := *countPacketsFlag
	poolList = make([]poolStruct, 0)

	for _, arg := range flag.Args() {
		newPool, ok := checkForPoolAddr(arg)
		if ok {
			poolList = append(poolList, newPool)
		}
	}

	if len(poolList) == 0 {
		fmt.Println("MPPING - Mining Pool Ping tool, which counts time from you to first reply of mining pool")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(2)
	}

	//fmt.Printf("MPPING %s:%s (%s:%s)\n", poolAddr, poolPort, poolIPAddr, poolPort)
	fmt.Printf("%-37s%-37s\n", "POOLSERVER", "RTT msec")

	infinitLoop := true
	if countPackets > 0 {
		infinitLoop = false
	}

	for poolID, poolServer := range poolList {

		poolAddr = poolServer.PoolDomain
		poolPort = strconv.Itoa(int(poolServer.PoolPort))
		if ipv4Proto {
			poolIPAddr = poolServer.PoolIPv4
		}
		if ipv6Proto {
			poolIPAddr = fmt.Sprintf("[%s]", poolServer.PoolIPv6)
		}

		//var totalPacketsSent, totalPacketsRec uint64
		//var totalTimeMin, totalTimeMax, totalTime uint64

		for {
			if !infinitLoop {
				if countPackets == 0 {
					fmt.Println()
					onStop()
				}
				countPackets--
			}
			beforeConnect := getCurrentTimeStamp()
			//totalPacketsSent++
			poolList[poolID].TotalPacketsSent++
			poolConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%s", poolIPAddr, poolPort))
			if err != nil {
				fmt.Println(err)
				//os.Exit(1)
			}
			fmt.Fprintf(poolConnection, request+"\n")
			_, err = bufio.NewReader(poolConnection).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				//os.Exit(1)
			}

			firstReply := getCurrentTimeStamp()

			fromUserToPool := uint64((firstReply - beforeConnect))

			fmt.Printf("%-37s%-37s\n", fmt.Sprintf("%s:%s", poolAddr, poolPort), fmt.Sprintf("%d msec", fromUserToPool))

			poolConnection.Close()
			//totalTime = totalTime + fromUserToPool
			poolList[poolID].TotalTime += fromUserToPool
			if fromUserToPool < poolList[poolID].TotalTimeMin || poolList[poolID].TotalTimeMin == 0 {
				poolList[poolID].TotalTimeMin = fromUserToPool
			}

			if fromUserToPool > poolServer.TotalTimeMax || poolList[poolID].TotalTimeMax == 0 {
				poolList[poolID].TotalTimeMax = fromUserToPool
			}
			poolList[poolID].TotalPacketsReceived++
			time.Sleep(1 * time.Second)
		}
	}
	<-c

	os.Exit(0)
}
