package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/CapacitorSet/ja3-server/crypto/tls"
	"github.com/CapacitorSet/ja3-server/net/http"
	"github.com/go-redis/redis"
	"net"
	"os"
	"strconv"
)

var client *redis.Client

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	hash := md5.Sum([]byte(r.JA3Fingerprint))
	out := make([]byte, 32)
	hex.Encode(out, hash[:])
	if r.URL.Path == "/cached" {
		// Prevent results from being registered twice
		w.Header().Set("Cache-Control", "public,max-age=31556926,immutable")
		w.Header().Set("Expires", "Mon, 30 Dec 2019 08:00:00 GMT")
		w.Header().Set("Last-Modified", "Sun, 30 Dec 2018 08:00:00 GMT")
		w.WriteHeader(200)
		w.Write(out)

		_, err := client.HIncrBy("freqs", string(out), 1).Result()
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = client.HIncrBy("freqs", "total", 1).Result()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		num, err := client.HGet("freqs", string(out)).Result()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(200)
			w.Write([]byte("error"))
			return
		}
		numF, err := strconv.ParseFloat(num, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(200)
			w.Write([]byte("error"))
			return
		}
		total, err := client.HGet("freqs", "total").Result()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(200)
			w.Write([]byte("error"))
			return
		}
		totalF, err := strconv.ParseFloat(total, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(200)
			w.Write([]byte("error"))
			return
		}
		fmt.Fprintf(w, "%.2f", totalF/numF)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Syntax: %s redis_ip:redis_port path/to/certificate.pem path/to/key.pem\n", os.Args[0])
		return
	}
	client = redis.NewClient(&redis.Options{
		Addr:     os.Args[1],
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis up.")

	handler := http.HandlerFunc(handler)
	server := &http.Server{Addr: ":8443", Handler: handler}

	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	cert, err := tls.LoadX509KeyPair(os.Args[2], os.Args[3])
	if err != nil {
		panic(err)
	}
	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

	tlsListener := tls.NewListener(ln, &tlsConfig)
	fmt.Println("HTTP up.")
	err = server.Serve(tlsListener)
	if err != nil {
		panic(err)
	}

	ln.Close()
}
