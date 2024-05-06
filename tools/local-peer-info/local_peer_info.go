package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/libp2p/go-libp2p/core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

var LOCAL_PEER_ENDPOINT = "http://localhost:7801/api/v0/id"

// Borrowed from ipfs code to parse the results of the command `ipfs id`
type IdOutput struct {
	ID              string
	PublicKey       string
	Addresses       []string
	AgentVersion    string
	ProtocolVersion string
}

// quick and dirty function to get the local ipfs daemons address for bootstrapping
func getLocalPeerInfo() []peer.AddrInfo {
	resp, err := http.PostForm(LOCAL_PEER_ENDPOINT, nil)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var js IdOutput
	err = json.Unmarshal(body, &js)
	if err != nil {
		log.Fatalln(err)
	}
	pid, err := peer.Decode(js.ID)
	if err != nil {
		log.Fatalln("Decode Peer ID ", err)
	}
	for _, addr := range js.Addresses {
		// For some reason, possibly NAT traversal, we need to grab the loopback ip address
		if addr[0:8] == "/ip4/127" {
			return convertPeers([]string{addr}, pid)
		}
	}
	log.Fatalln(err)
	return make([]peer.AddrInfo, 1) // not reachable, but keeps the compiler happy
}

func convertPeers(peers []string, pid peer.ID) []peer.AddrInfo {
	pinfos := make([]peer.AddrInfo, len(peers))
	for i, addr := range peers {
		maddr := ma.StringCast(addr)
		// p, err := peer.AddrInfoFromP2pAddr(maddr)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// pinfos[i] = *p
		pinfos[i] = peer.AddrInfo{
			ID:    pid,
			Addrs: []ma.Multiaddr{maddr},
		}
	}
	return pinfos
}

func main() {
	bootstrapPeers := getLocalPeerInfo()
	fmt.Println("get Local Peer Info:", bootstrapPeers)
}
