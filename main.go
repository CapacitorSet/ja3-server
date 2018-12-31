package main

import (
	"net"
	"github.com/CapacitorSet/ja3-server/crypto/tls"
	"github.com/CapacitorSet/ja3-server/net/http"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(".")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Prevent results from being registered twice
	w.Header().Set("Cache-Control", "public,max-age=31556926,immutable")
	w.Header().Set("Expires", "Mon, 30 Dec 2019 08:00:00 GMT")
	w.Header().Set("Last-Modified", "Sun, 30 Dec 2018 08:00:00 GMT")
	w.WriteHeader(200)

	hash := md5.Sum([]byte(r.JA3Fingerprint))
	out := make([]byte, 32)
	hex.Encode(out, hash[:])
	w.Write(out)
}

func main() {
	handler := http.HandlerFunc(handler)
	server := &http.Server{Addr: ":8443", Handler: handler}

	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	cert, err := tls.LoadX509KeyPair("fullchain.pem", "privkey.pem")
	if err != nil {
		panic(err)
	}
	tlsConfig := tls.Config{Certificates:[]tls.Certificate{cert}}

	tlsListener := tls.NewListener(ln, &tlsConfig)
	err = server.Serve(tlsListener)
	if err != nil {
    	panic(err)
    }

	ln.Close()
}
