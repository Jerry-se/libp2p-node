package config

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func TestGeneratePreSharedKey(t *testing.T) {
	pskString, err := GeneratePreSharedKey()
	if err != nil {
		t.Fatal("Failed to generate Pre-Shared Key")
	} else {
		t.Logf("Generate Pre-Shared Key: %s", pskString)
	}
}

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../../peer.json")
	if err != nil {
		t.Fatalf("Load json configuration failed %v", err)
	} else {
		t.Logf("Load json configuration info %s", config)
	}
}

func TestPeerKey(t *testing.T) {
	priv, pub, err := LoadPeerKey("../../peer.key")
	if err != nil {
		t.Log("Peer.key not existed, Generating ...")
		priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
		if err != nil {
			t.Fatalf("Failed to generate ed25519 key %v", err)
		} else {
			SavePeerKey("../../peer.key", priv)
			t.Logf("Generate ed25529 pub key %v", pub)
			id, err := peer.IDFromPublicKey(pub)
			if err != nil {
				t.Fatalf("Generate Peer ID err: %v", err)
			} else {
				t.Logf("Generate Peer ID %v", id)
			}
		}
	} else {
		t.Logf("Read ed25529 pub key %v", pub)
		t.Logf("Get Pub from priv %v", priv.GetPublic())
		id, err := peer.IDFromPublicKey(pub)
		if err != nil {
			t.Fatalf("Transform Peer ID err: %v", err)
		} else {
			t.Logf("Transform Peer ID %v", id)
		}
	}
}
