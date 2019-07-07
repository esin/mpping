package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

func getCurrentTimeStamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func lookupIPAddr(poolAddr string) string {
	ips, err := net.LookupIP(poolAddr)
	if err != nil {
		fmt.Printf("Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		if strings.Index(ip.String(), ":") < 0 {
			return ip.String()
		}
	}
	return ""
}

type poolStruct struct {
	Address string
	Port    string
	Scheme  string
}

func checkForPoolAddr(urlArg string) (poolStruct, bool) {

	var newPool poolStruct

	strings.TrimLeft(urlArg, " ")

	if (!strings.HasPrefix(strings.ToLower(urlArg), "stratum://")) && (!strings.HasPrefix(strings.ToLower(urlArg), "stratum2://")) {
		urlArg = "stratum://" + urlArg
	}

	parsed, err := url.Parse(urlArg)
	if err != nil {
		return newPool, false
	}

	newPool.Address = parsed.Hostname()
	newPool.Port = parsed.Port()
	newPool.Scheme = parsed.Scheme

	return newPool, true
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("Catched %v", sig)
			os.Exit(0)
		}
	}()

	request := `{"id":1,"method":"mining.subscribe","params":["mpping-0.1","EthereumStratum/2.0.0"]}`
	flag.Parse()
	//poolAddrOpt := flag.String("pool", "abyss", "Pool address with port (for example: stratum.pool.com:3333")

	poolAddr := ""
	poolPort := ""

	for _, arg := range flag.Args() {
		newPool, ok := checkForPoolAddr(arg)
		if ok {
			poolAddr = newPool.Address
			poolPort = newPool.Port
		}
	}

	poolIPAddr := lookupIPAddr(poolAddr)

	fmt.Printf("MPPING %s:%s (%s:%s)\n", poolAddr, poolPort, poolIPAddr, poolPort)

	for {
		beforeConnect := getCurrentTimeStamp()
		poolConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%s", poolIPAddr, poolPort))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Fprintf(poolConnection, request+"\n")
		_, err = bufio.NewReader(poolConnection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		firstReply := getCurrentTimeStamp()

		fromUserToPool := (firstReply - beforeConnect) / 2

		fmt.Printf("From you to %s:%s %d msec\n", poolAddr, poolPort, fromUserToPool)
		poolConnection.Close()
		time.Sleep(1 * time.Second)
	}
}
