// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/hashicorp/vault-pki-cieps-example/pkg/route"
)

func main() {
	serverCertFlag := flag.String("server-cert", "server.crt", "Path to the server certificate file")
	serverKeyFlag := flag.String("server-key", "server.key", "Path to the server key file corresponding to the given certificate file")
	listenAddrFlag := flag.String("listen", ":443", "TCP listen address in host:port format")
	flag.Parse()

	mux := http.NewServeMux()
	route.RegisterHandlers(mux)

	err := http.ListenAndServeTLS(*listenAddrFlag, *serverCertFlag, *serverKeyFlag, mux)
	if err != nil {
		log.Fatalf("ListenAndServe(%v, %v, %v): %v", *listenAddrFlag, *serverCertFlag, *serverKeyFlag, err)
	}
}
