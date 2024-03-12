package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Jerry-se/libp2p-node/pkg/config"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"

	// ds "github.com/ipfs/go-datastore"
	// dsync "github.com/ipfs/go-datastore/sync"
	golog "github.com/ipfs/go-log/v2"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

var DefaultBootstrapPeers = convertPeers([]string{
	"/ip4/122.99.183.54/tcp/7001/p2p/12D3KooWSpgWzEXE5GNjY6hgdAhuuBLe4d3ocqWDnVLdCa8U3cig",
	"/ip4/122.99.183.54/tcp/8511/p2p/12D3KooWJaNdwbwsvESYZmeEhTwkG1KirKxXpE6CdPx2Z9VqM3Rt",
	"/ip4/122.99.183.54/tcp/8515/p2p/12D3KooWNp5pyEAtXs52RqssBkrCGojyNLWiZysjr8EpXKK4rpyp",
	"/ip4/82.157.50.32/tcp/7001/p2p/12D3KooWFrTcDtocZWEvEAk2X4poyn13LzT3G7JMBRoPD73YPAoB",
})

func convertPeers(peers []string) []multiaddr.Multiaddr {
	maddrs := make([]multiaddr.Multiaddr, len(peers))
	for i, peer := range peers {
		maddr, err := multiaddr.NewMultiaddr(peer)
		if err != nil {
			log.Fatalln(err)
		}
		maddrs[i] = maddr
	}
	return maddrs
}

// doEcho reads a line of data from a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s\n", str)
	_, err = s.Write([]byte(str))
	return err
}

func main() {
	golog.SetAllLoggers(golog.LevelInfo)

	listenF := flag.Int("l", 6000, "listening port waiting for incoming connections")
	pskString := flag.String("psk", "", "Pre-Shared Key")
	peerKeyPath := flag.String("peerkey", "", "the file path of peer key")
	ping := flag.Bool("ping", false, "whether to enable ipfs ping")
	protocolPrefix := flag.String("protocol", "", "the prefix attached to all DHT protocols")
	flag.Parse()

	if *peerKeyPath == "" {
		log.Fatal("Please provide a filepath to save peer key")
	}
	peerKey, _, err := config.LoadPeerKey(*peerKeyPath)
	if err != nil {
		// log.Fatalf("Load peer key: %v", err)
		peerKey, _, err = config.GeneratePeerKey(*peerKeyPath)
		if err != nil {
			log.Fatalf("Generate peer key: %v", err)
		}
	} else {
		log.Println("Load peer key success")
	}

	ctx := context.Background()
	var kadDHT *dht.IpfsDHT

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		log.Fatalf("Create connection manager: %v", err)
	}

	opts := []libp2p.Option{
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *listenF),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", *listenF),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1/webtransport", *listenF),
			// "/ip4/0.0.0.0/tcp/0",
			// "/ip4/0.0.0.0/udp/0/quic-v1",
			// "/ip4/0.0.0.0/udp/0/quic-v1/webtransport",
			// "/ip6/::/tcp/0",
			// "/ip6/::/udp/0/quic-v1",
			// "/ip6/::/udp/0/quic-v1/webtransport",
		),
		// Use the keypair we generated
		libp2p.Identity(peerKey),
		libp2p.Ping(*ping),
		// libp2p.DefaultPrivateTransports,
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			dhtOpts := []dht.Option{
				dht.Mode(dht.ModeServer),
				dht.BucketSize(20),
			}
			if *protocolPrefix != "" {
				dhtOpts = append(dhtOpts, dht.ProtocolPrefix(protocol.ID(*protocolPrefix)))
			}
			kadDHT, err = dht.New(ctx, h, dhtOpts...)
			return kadDHT, err
		}),
		libp2p.ProtocolVersion("ipfs/0.1.0"),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.DisableRelay(),
		libp2p.EnableNATService(),
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableRelayService(relay.WithResources(relay.DefaultResources())),
		// libp2p.EnableHolePunching(),
	}

	if *pskString != "" {
		psk, err := hex.DecodeString(*pskString)
		if err != nil {
			log.Fatalf("Decoding PSK: %v", err)
		}
		log.Println("Pre-Shared Key ", psk)
		opts = append(opts, libp2p.PrivateNetwork(psk))
	}

	node, err := libp2p.New(opts...)
	if err != nil {
		log.Fatalf("Create libp2p host: %v", err)
	}

	// _, err = relay.New(node, relay.WithResources(relay.DefaultResources()))
	// if err != nil {
	// 	log.Fatalf("Failed to instantiate the relay service: %v", err)
	// }

	log.Println("Listen addresses:", node.Addrs())
	log.Println("Node id:", node.ID())

	// Set a stream handler on host A. /chat/1.0.0 is
	// a user-defined protocol name.
	node.SetStreamHandler("/chat/1.0.0", func(s network.Stream) {
		log.Println("Got a new stream!")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	// print the node's PeerInfo in multiaddr format
	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	log.Println("libp2p node address:", addrs) // addrs[0]

	// dstore := dsync.MutexWrap(ds.NewMapDatastore())
	// kadDHT := dht.NewDHT(ctx, node, dstore)

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			node.Peerstore().AddAddrs(peerinfo.ID, peerinfo.Addrs, peerstore.PermanentAddrTTL)
			if err := node.Connect(ctx, *peerinfo); err != nil {
				log.Println("Connect bootstrap node", *peerinfo, err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	err = kadDHT.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("Bootstrap the host: %v", err)
	}

	log.Println("listening for connections")
	select {} // hang forever
}
