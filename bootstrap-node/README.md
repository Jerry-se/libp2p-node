# libp2p-bootstrap-node

libp2p bootstrap node

## go

代码参考: <https://github.com/libp2p/go-libp2p>

## rust

代码参考: <https://github.com/libp2p/rust-libp2p/tree/master/misc/server>

```bash
# 编译
cargo build
cargo build --release
# 运行
cargo run -- --config ./config.json --enable-kademlia --enable-autonat
./libp2p-server --config ./config.json --enable-kademlia --enable-autonat
```

打开 tracing log 需要设置环境变量 `RUST_LOG=info`

## 说明

1. 部署引导节点的服务器必须要有公网 IP。

## 参考资料

- [JS dht bootstrapper](https://github.com/libp2p/js-libp2p-amino-dht-bootstrapper)
- [Rust libp2p Server](https://github.com/libp2p/rust-libp2p/tree/master/misc/server)
- [IPFS Bitswap Protocol](https://specs.ipfs.tech/bitswap-protocol/)
- [go-libp2p-relay-daemon](https://github.com/libp2p/go-libp2p-relay-daemon)
- [go Punchr dcutr](https://github.com/libp2p/punchr)
- [Hydra Booster - A DHT Indexer node & Peer Router](https://github.com/libp2p/hydra-booster)
- [go libp2p daemon](https://github.com/libp2p/go-libp2p-daemon)

当前应该有两种引导节点的实现，一种是 Rust libp2p Server，另一种是 Kubo (Go 语言的实现) 。

以前应该是可以通过 JS 搭建的，后来要使用 CID，研究了新的 amino protocol。

- [A Rusty Bootstrapper](https://blog.ipfs.tech/2023-rust-libp2p-based-ipfs-bootstrap-node/)
- [Hole punching in libp2p - Overcoming Firewalls](https://blog.ipfs.tech/2022-01-20-libp2p-hole-punching/)
- [How to setup ipfs relay](https://discuss.ipfs.tech/t/how-to-setup-v1-relay-in-the-new-config/14019?page=2)

