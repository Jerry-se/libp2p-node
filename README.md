# libp2p-node

libp2p node

## 架构

libp2p + IPFS + Stable Diffusion

- go-libp2p: <https://github.com/libp2p/go-libp2p>

开发环境:
- Windows: Go 1.21.6 (at least 1.19)

对外提供接口: Gin + Gorilla WebSocket

## 功能
1. 节点发现 Discovery
2. 节点连接 Connect
3. NAT 穿越 NAT traversal
4. 私有 DHT 网络

1. 要有 Route 才能使用 PeerID 来连接，否则只能使用 `/ip4/127.0.0.1/tcp/53724/p2p/<peer-id>`。
2. 要支持 TCP/QUIC 协议，如果要支持从浏览器发起连接，需要支持 WebSocket 协议或者 WebRTC 协议。
3. 要支持 Relay ，或者公网服务器上的节点自带 Relay 功能，否则家用电脑无法跨域 NAT 和防火墙。

## 编译和打包

编译:

1. go mod tidy
2. go build -o go-libp2p-tutorial.exe go-libp2p-tutorial\main.go

go build -o libp2p-address.exe app1\main.go
go test -v -timeout 30s -count=1 -run TestPeerKey github.com/Jerry-se/libp2p-node/pkg/config

## 参考资料

- [Run a go-libp2p node](https://docs.libp2p.io/guides/getting-started/go/)
