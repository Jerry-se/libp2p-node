# go-libp2p-tutorial

go-libp2p-tutorial

## build

```bash
$ go build -o libp2p-node
```

## run and test

In one terminal window, let’s start a listening node (i.e. don’t pass any command line arguments):

```bash
$ ./libp2p-node
libp2p node address: /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL
```

In another terminal window, let’s run a second node but pass the address of the first node, and we should see some ping responses logged:

```bash
$ ./libp2p-node /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL
libp2p node address: /ip4/127.0.0.1/tcp/61846/ipfs/QmVyKLTLswap3VYbpBATsgNpi6JdwSwsZALPxEnEbEndup
sending 5 ping messages to /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL
pinged /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL in 431.231µs
pinged /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL in 164.94µs
pinged /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL in 220.544µs
pinged /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL in 208.761µs
pinged /ip4/127.0.0.1/tcp/61790/ipfs/QmZKjsGJ6ukXVRXVEcExx9GhiyWoJC97onYpzBwCHPWqpL in 201.37µs
```

Success! Our two peers are now communicating using go-libp2p! Sure, they can only say “ping”, but it’s a start!

## reference

Getting started with go-libp2p: <https://docs.libp2p.io/tutorials/getting-started/go/>
