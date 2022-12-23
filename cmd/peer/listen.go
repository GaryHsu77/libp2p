package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func listenRun(cmd *cobra.Command, args []string) {

	argInit()

	if len(args) < 1 {
		logrus.Errorln("usage: lptunnel listen /ip4/$Relay_IP/tcp/$Relay_Port/p2p/$Relay_ID")
		return
	}

	var (
		err       error
		relayAddr = args[0]
	)

	// relay info
	relayMA, _ := multiaddr.NewMultiaddr(relayAddr)
	relayAddrInfo, _ := peer.AddrInfoFromP2pAddr(relayMA)
	logrus.Infoln(relayAddrInfo)

	// new listener
	listener := runPeer(relayAddrInfo)
	logrus.Infoln("peer is connected to relay")
	logrus.Infoln("I am", listener.ID())

	// set up a protocol handler on this peer
	listener.SetStreamHandler("/customprotocol", func(s network.Stream) {
		conn, err := net.Dial("tcp", forwardAddr)
		if err != nil {
			logrus.Errorf("unable to connect to %s, err:%v", forwardAddr, err)
			return
		}
		go runSession(s.ID(), s, conn)
	})

	// Hosts that want to have messages relayed on their behalf need to reserve a slot
	// with the circuit relay service host
	_, err = client.Reserve(context.Background(), listener, *relayAddrInfo)
	if err != nil {
		logrus.Errorln("caller failed to receive a relay reservation from relay. err:", err)
		return
	}
	logrus.Infoln("peer is reserved to relay")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
}
