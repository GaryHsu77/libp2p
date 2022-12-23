package main

import (
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

func relayRun(cmd *cobra.Command, args []string) {

	argInit()

	var opts []libp2p.Option
	opts = append(opts, libp2p.ListenAddrs(multiaddr.StringCast("/ip4/0.0.0.0/udp/4001/quic")))
	opts = append(opts, libp2p.ListenAddrs(multiaddr.StringCast("/ip6/::/udp/4001/quic")))
	opts = append(opts, libp2p.ListenAddrs(multiaddr.StringCast("/ip4/127.0.0.1/tcp/4001")))
	opts = append(opts, libp2p.ListenAddrs(multiaddr.StringCast("/ip4/0.0.0.0/tcp/4001")))

	// Create a host to act as a middleman to relay messages on our behalf
	host, err := libp2p.New(
		opts...,
	)
	if err != nil {
		log.Printf("Failed to create relay1: %v", err)
		return
	}

	relayrc := relay.WithLimit(nil)
	// Configure the host to offer the ciruit relay service.
	_, err = relay.New(host, relayrc)
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
		return
	}

	fmt.Printf("I am %s\n", host.ID())
	fmt.Printf("Public Addresses:\n")
	for _, addr := range host.Addrs() {
		fmt.Printf("\t%s/p2p/%s\n", addr, host.ID())
	}

	c := make(chan bool)
	<-c
}
