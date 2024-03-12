package main

import (
	"flag"
	"log"

	"github.com/Jerry-se/libp2p-node/pkg/config"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	peerKeyPath := flag.String("peerkey", "", "the file path of peer key")
	flag.Parse()

	if *peerKeyPath == "" {
		log.Fatal("Please provide a filepath to save peer key")
	}

	privKey, pubKey, err := config.LoadPeerKey(*peerKeyPath)
	if err != nil {
		// log.Fatalf("Load peer key: %v", err)
		privKey, pubKey, err = config.GeneratePeerKey(*peerKeyPath)
		if err != nil {
			log.Fatalf("Generate peer key: %v", err)
		}
		log.Println("Generate peer key success")
	} else {
		log.Println("Load peer key success")
	}

	privkeyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		log.Fatalf("Marshal Private Key err: %v", err)
	}
	log.Println("Encode private key:", crypto.ConfigEncodeKey(privkeyBytes))

	pubkeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		log.Fatalf("Marshal Public Key err: %v", err)
	}
	log.Println("Encode public key:", crypto.ConfigEncodeKey(pubkeyBytes))

	id, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		log.Fatalf("Transform Peer ID err: %v", err)
	} else {
		log.Println("Transform Peer ID:", id)
	}
}
