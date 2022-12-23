package main

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

const maxPacketSize = 1024
const forwardAddr = "10.123.12.138:8443"

func runPeer(relayAddrInfo *peer.AddrInfo) host.Host {

	peer, err := libp2p.New(
		libp2p.NoListenAddrs,
		libp2p.EnableRelay(),
	)
	if err != nil {
		logrus.Errorf("Failed to create unreachable1: %v", err)
		return nil
	}

	// Connect both unreachable1 and unreachable2 to relay1
	if err := peer.Connect(context.Background(), *relayAddrInfo); err != nil {
		logrus.Errorf("Failed to connect unreachable1 and relay1: %v", err)
		return nil
	}

	return peer
}

func runSession(id string, stream network.Stream, conn net.Conn) {

	var (
		wg          sync.WaitGroup
		log         = logrus.New().WithField("session", id)
		ctx, cancel = context.WithCancel(context.TODO())
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(stream, conn)
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(conn, stream)
		cancel()
	}()

	log.Infoln("started")
	<-ctx.Done()

	// exit
	stream.Close()
	conn.Close()
	wg.Wait()
	log.Infoln("stopped")
}
