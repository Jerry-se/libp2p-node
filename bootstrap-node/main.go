package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/Jerry-se/libp2p-node/pkg/config"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	golog "github.com/ipfs/go-log/v2"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

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
	pskString := flag.String("psk",
		"b8d65dc757e9ebfc0b5e7aafd4ebfdb7409273e4d0a6ac9ba716a6e15d603d92",
		"Pre-Shared Key")
	peerKeyPath := flag.String("peerkey", "", "the file path of peer key")
	flag.Parse()

	psk, err := hex.DecodeString(*pskString)
	if err != nil {
		log.Fatalf("Decoding PSK: %v", err)
	}
	log.Println("Pre-Shared Key ", psk)

	var peerKey crypto.PrivKey
	if *peerKeyPath == "" {
		log.Fatal("Please provide a filepath to save peer key")
	} else {
		peerKey, _, err = config.LoadPeerKey(*peerKeyPath)
		if err != nil {
			// log.Fatalf("Load peer key: %v", err)
			peerKey, _, err = config.GeneratePeerKey(*peerKeyPath)
			if err != nil {
				log.Fatalf("Generate peer key: %v", err)
			}
		} else {
			log.Println("Load peer key success")
		}
	}

	node, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *listenF)),
		libp2p.Identity(peerKey),
		libp2p.PrivateNetwork(psk),
		// libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	)
	if err != nil {
		log.Fatalf("Create libp2p host: %v", err)
	}

	ctx := context.Background()

	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	kadDHT := dht.NewDHT(ctx, node, dstore)

	err = kadDHT.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("Bootstrap the host: %v", err)
	}

	fmt.Println("Listen addresses:", node.Addrs())
	fmt.Println("Node id:", node.ID())

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	node.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("Got a new stream!")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	log.Println("listening for connections")
	select {} // hang forever
}
