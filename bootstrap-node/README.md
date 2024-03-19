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
