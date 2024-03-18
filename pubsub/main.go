package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

var (
	DefaultBootstrapPeers = convertPeers([]string{
		"/ip4/122.99.183.54/tcp/7001/p2p/12D3KooWSpgWzEXE5GNjY6hgdAhuuBLe4d3ocqWDnVLdCa8U3cig",
		"/ip4/122.99.183.54/tcp/8511/p2p/12D3KooWJaNdwbwsvESYZmeEhTwkG1KirKxXpE6CdPx2Z9VqM3Rt",
		"/ip4/82.157.50.32/tcp/7001/p2p/12D3KooWFrTcDtocZWEvEAk2X4poyn13LzT3G7JMBRoPD73YPAoB",
	})
	topicNameFlag  = flag.String("topicName", "applesauce", "name of topic to join")
	protocolPrefix = flag.String("protocol", "", "the prefix attached to all DHT protocols")
)

func convertPeers(peers []string) []multiaddr.Multiaddr {
	maddrs := make([]multiaddr.Multiaddr, len(peers))
	for i, peer := range peers {
		maddr, err := multiaddr.NewMultiaddr(peer)
		if err != nil {
			panic(err)
		}
		maddrs[i] = maddr
	}
	return maddrs
}

func main() {
	pskString := flag.String("psk", "", "Pre-Shared Key")
	flag.Parse()
	ctx := context.Background()

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	}
	if *pskString != "" {
		psk, err := hex.DecodeString(*pskString)
		if err != nil {
			panic(err)
		}
		fmt.Println("Pre-Shared Key ", psk)
		opts = append(opts, libp2p.PrivateNetwork(psk))
	}
	h, err := libp2p.New(opts...)
	if err != nil {
		panic(err)
	}
	defer h.Close()
	go discoverPeers(ctx, h)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	topic, err := ps.Join(*topicNameFlag)
	if err != nil {
		panic(err)
	}
	go streamConsoleTo(ctx, topic)

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	printMessagesFrom(ctx, sub)
}

func initDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	dhtOpts := []dht.Option{
		dht.Mode(dht.ModeClient),
	}
	if *protocolPrefix != "" {
		dhtOpts = append(dhtOpts, dht.ProtocolPrefix(protocol.ID(*protocolPrefix)))
	}

	kademliaDHT, err := dht.New(ctx, h, dhtOpts...)
	if err != nil {
		panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			} else {
				fmt.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

func discoverPeers(ctx context.Context, h host.Host) {
	kademliaDHT := initDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, *topicNameFlag)

	// Look for others who have announced and attempt to connect to them
	// anyConnected := false
	// for !anyConnected {
	// 	fmt.Println("Searching for peers...")
	// 	peerChan, err := routingDiscovery.FindPeers(ctx, *topicNameFlag)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for peer := range peerChan {
	// 		if peer.ID == h.ID() {
	// 			continue // No self connection
	// 		}
	// 		err := h.Connect(ctx, peer)
	// 		if err != nil {
	// 			fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
	// 		} else {
	// 			fmt.Println("Connected to:", peer.ID)
	// 			anyConnected = true
	// 		}
	// 	}
	// }
	fmt.Println("Peer discovery complete")
}

func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
	}
}
