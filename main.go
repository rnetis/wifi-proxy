package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/cssivision/reverseproxy"
	socks "github.com/nwtgck/go-socks"
)

func handleSocksServer(ln net.Listener) {
	conf := &socks.Config{}
	socksServer, err := socks.New(conf)
	if err != nil {
		log.Println(err)
	}

	socksServer.Serve(ln)
}

func handleHTTPServer(ln net.Listener) {
	http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := reverseproxy.NewReverseProxy(r.URL)
		proxy.ServeHTTP(w, r)
	}))
}

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:7979")
	if err != nil {
		log.Println("Failed to create listener: ", err)
	}
	defer listener.Close()

	//Socks Server
	go handleSocksServer(listener)
	//HTTP and HTTPS Server
	go handleHTTPServer(listener)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("Error Getting Addresses: ", err)
		return
	}
	fmt.Println("...Servers Started...")
	for i, v := range addrs {
		addr := strings.Split(v.String(), "/")
		fmt.Printf("%d) %s:7979\n", i+1, addr[0])
	}

	select {}

}
