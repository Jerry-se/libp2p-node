package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/Jerry-se/libp2p-node/pkg/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/ipfs/go-log/v2"

	"github.com/multiformats/go-multiaddr"
)

var (
	DefaultBootstrapPeers = convertPeers([]string{
		"/ip4/122.99.183.54/tcp/7001/p2p/12D3KooWSpgWzEXE5GNjY6hgdAhuuBLe4d3ocqWDnVLdCa8U3cig",
		// "/ip4/122.99.183.54/tcp/8511/p2p/12D3KooWJaNdwbwsvESYZmeEhTwkG1KirKxXpE6CdPx2Z9VqM3Rt",
		// "/ip4/122.99.183.54/tcp/8515/p2p/12D3KooWNp5pyEAtXs52RqssBkrCGojyNLWiZysjr8EpXKK4rpyp",
		"/ip4/82.157.50.32/tcp/7001/p2p/12D3KooWFrTcDtocZWEvEAk2X4poyn13LzT3G7JMBRoPD73YPAoB",
	})
)

var logger = log.Logger("rendezvous")

func convertPeers(peers []string) []peer.AddrInfo {
	pinfos := make([]peer.AddrInfo, len(peers))
	for i, addr := range peers {
		maddr := multiaddr.StringCast(addr)
		p, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			logger.Fatalln(err)
		}
		pinfos[i] = *p
	}
	return pinfos
}

func handleStream(stream network.Stream) {
	logger.Info("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func main() {
	log.SetAllLoggers(log.LevelInfo)
	log.SetLogLevel("rendezvous", "debug")
	help := flag.Bool("h", false, "Display Help")
	listenF := flag.Int("l", 6000, "listening port waiting for incoming connections")
	peerKeyPath := flag.String("peerkey", "", "the file path of peer key")
	pskString := flag.String("psk", "", "Pre-Shared Key")
	rendezvousString := flag.String("rendezvous", "meet me here",
		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	protocolPrefix := flag.String("protocol", "", "the prefix attached to all DHT protocols")
	flag.Parse()

	if *help {
		fmt.Println("This program demonstrates a simple p2p chat application using libp2p")
		fmt.Println()
		fmt.Println("Usage: Run './chat in two different terminals. Let them connect to the bootstrap nodes, announce themselves and connect to the peers")
		flag.PrintDefaults()
		return
	}

	if *peerKeyPath == "" {
		logger.Fatal("Please provide a filepath to save peer key")
	}
	peerKey, _, err := config.LoadPeerKey(*peerKeyPath)
	if err != nil {
		// log.Fatalf("Load peer key: %v", err)
		peerKey, _, err = config.GeneratePeerKey(*peerKeyPath)
		if err != nil {
			logger.Fatalf("Generate peer key: %v", err)
		}
	} else {
		logger.Info("Load peer key success")
	}

	ctx := context.Background()
	// var kademliaDHT *dht.IpfsDHT

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *listenF),
		),
		libp2p.Identity(peerKey),
		libp2p.DefaultPrivateTransports,
		// libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		// libp2p.ProtocolVersion("ipfs/0.1.0"),
		libp2p.NATPortMap(),
		// libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		// 	dhtOpts := []dht.Option{
		// 		dht.Mode(dht.ModeAuto),
		// 		dht.BootstrapPeers(DefaultBootstrapPeers...),
		// 	}
		// 	if *protocolPrefix != "" {
		// 		dhtOpts = append(dhtOpts, dht.ProtocolPrefix(protocol.ID(*protocolPrefix)))
		// 	}
		// 	kademliaDHT, err = dht.New(ctx, h, dhtOpts...)
		// 	return kademliaDHT, err
		// }),
		libp2p.ForceReachabilityPrivate(),
		libp2p.EnableAutoRelayWithStaticRelays(
			DefaultBootstrapPeers,
			autorelay.WithNumRelays(1),
			autorelay.WithMinCandidates(1),
		),
		libp2p.EnableHolePunching(),
	}

	if *pskString != "" {
		psk, err := hex.DecodeString(*pskString)
		if err != nil {
			logger.Fatalf("Decoding PSK: %v", err)
		}
		logger.Infof("Pre-Shared Key %v", psk)
		opts = append(opts, libp2p.PrivateNetwork(psk))
	}

	// libp2p.New constructs a new libp2p Host. Other options can be added
	// here.
	host, err := libp2p.New(opts...)
	if err != nil {
		panic(err)
	}
	defer host.Close()

	logger.Info("Host created. We are:", host.ID())
	logger.Info(host.Addrs())

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler("/chat/1.0.0", handleStream)

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	dhtOpts := []dht.Option{
		dht.Mode(dht.ModeAuto),
		dht.BootstrapPeers(DefaultBootstrapPeers...),
	}
	if *protocolPrefix != "" {
		dhtOpts = append(dhtOpts, dht.ProtocolPrefix(protocol.ID(*protocolPrefix)))
	}
	kademliaDHT, err := dht.New(ctx, host, dhtOpts...)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerInfo := range DefaultBootstrapPeers {
		wg.Add(1)
		go func(peerInfo peer.AddrInfo) {
			defer wg.Done()
			host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)
			if err := host.Connect(ctx, peerInfo); err != nil {
				logger.Warning(err)
			} else {
				logger.Info("Connection established with bootstrap node:", peerInfo)
			}
		}(peerInfo)
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	logger.Info("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, *rendezvousString)
	logger.Debug("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	logger.Debug("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, *rendezvousString)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		logger.Debug("Found peer:", peer)

		logger.Debug("Connecting to:", peer)
		stream, err := host.NewStream(ctx, peer.ID, "/chat/1.0.0")

		if err != nil {
			logger.Warning("Connection failed:", err)
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}

		logger.Info("Connected to:", peer)
	}

	select {}
}
