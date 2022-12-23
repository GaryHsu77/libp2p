package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/swarm"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func callRun(cmd *cobra.Command, args []string) {

	argInit()

	if len(args) < 2 {
		logrus.Errorln("usage: lp call /ip4/$Relay_IP/tcp/$Relay_Port/p2p/$Relay_ID $Dest_Peer_ID")
		return
	}

	var (
		err          error
		relayAddr    = args[0]
		listenerAddr = args[1]
	)

	// dest peer info
	listenerMA, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/4001/p2p/%s", listenerAddr))
	listenerAddrInfo, _ := peer.AddrInfoFromP2pAddr(listenerMA)
	logrus.Infoln(listenerAddrInfo)

	// relay info
	relayMA, _ := multiaddr.NewMultiaddr(relayAddr)
	relayAddrInfo, _ := peer.AddrInfoFromP2pAddr(relayMA)
	logrus.Infoln(relayAddrInfo)

	// new caller
	caller := runPeer(relayAddrInfo)
	logrus.Infoln("peer is connected to relay")
	logrus.Infoln("I am", caller.ID())

	// Now create a new address for unreachable2 that specifies to communicate via
	// relay1 using a circuit relay
	relayaddr, err := multiaddr.NewMultiaddr("/p2p/" + relayAddrInfo.ID.String() + "/p2p-circuit/p2p/" + listenerAddrInfo.ID.String())
	if err != nil {
		logrus.Errorln(err)
		return
	}

	// Since we just tried and failed to dial, the dialer system will, by default
	// prevent us from redialing again so quickly. Since we know what we're doing, we
	// can use this ugly hack (it's on our TODO list to make it a little cleaner)
	// to tell the dialer "no, its okay, let's try this again"
	caller.Network().(*swarm.Swarm).Backoff().Clear(listenerAddrInfo.ID)

	logrus.Infoln("starting to connect the caller via the relay node")

	// Open a connection to the previously unreachable host via the relay address
	listenerRelayinfo := peer.AddrInfo{
		ID:    listenerAddrInfo.ID,
		Addrs: []multiaddr.Multiaddr{relayaddr},
	}
	if err := caller.Connect(context.Background(), listenerRelayinfo); err != nil {
		logrus.Errorf("Unexpected error here. Failed to connect unreachable1 and unreachable2: %v", err)
		return
	}
	logrus.Infoln("connected to listener:", listenerRelayinfo.ID.String())

	// bind local server
	proxyAddr := "127.0.0.1:6601"
	proxy, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	logrus.Infoln("Proxy server start and listening on:", proxyAddr)

	// run local server
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	go func() {
		sessionID := 1
		for ctx.Err() == nil {
			conn, err := proxy.Accept()
			if ctx.Err() != nil {
				return
			}
			if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			if err != nil {
				logrus.Warnln("failed to accept client connection, err:", err)
				return
			}

			s, err := caller.NewStream(network.WithUseTransient(context.Background(), "customprotocol"), listenerRelayinfo.ID, "/customprotocol")
			if err != nil {
				logrus.Errorln("failed to connect the listener via relay, err:", err)
				continue
			}
			go runSession(s.ID(), s, conn)

			if sessionID++; sessionID > 999999999 {
				sessionID = 1
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
}
