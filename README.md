# P2P Tunneling

## Relay
* run  
  ```bash
  go run -mod=vendor ./cmd/peer relay
  ```
* output
  ```
  I am 12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
  Public Addresses:
          /ip6/::1/udp/4001/quic-v1/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip6/::1/udp/4001/quic/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip4/10.160.126.86/tcp/4001/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip4/10.160.126.86/udp/4001/quic-v1/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip4/127.0.0.1/udp/4001/quic-v1/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip4/10.160.126.86/udp/4001/quic/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
          /ip4/127.0.0.1/udp/4001/quic/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr
  ```


## Listener
* run  
  ```bash
  go run -mod=vendor ./cmd/peer listen /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr  
  ```
* output
  ```
  INFO[0000] {12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr: [/ip4/127.0.0.1/tcp/4001]}
  INFO[0000] peer is connected to relay
  INFO[0000] I am 12D3KooWFdraG2YDv4E9SyNzwJgMMd4dK2ZyKHLMbPKD94Roexuv
  INFO[0000] peer is reserved to relay
  ```

## Caller
* run  
  ```bash
  go run -mod=vendor ./cmd/peer call /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr 12D3KooWFdraG2YDv4E9SyNzwJgMMd4dK2ZyKHLMbPKD94Roexuv
  ```
* output
  ```
  INFO[0000] {12D3KooWFdraG2YDv4E9SyNzwJgMMd4dK2ZyKHLMbPKD94Roexuv: [/ip4/0.0.0.0/tcp/4001]}
  INFO[0000] {12D3KooWDzRUs2aoVz78Rbp4KVSsEAgNsz4bfa9byJYj4nGZsVZr: [/ip4/127.0.0.1/tcp/4001]}
  INFO[0000] peer is connected to relay
  INFO[0000] I am 12D3KooWPbNCpbiETz3GmqUhijJ2Z39zR44GoadxaCJdASTz1qkb
  INFO[0000] starting to connect the caller via the relay node
  INFO[0000] connected to listener: 12D3KooWFdraG2YDv4E9SyNzwJgMMd4dK2ZyKHLMbPKD94Roexuv
  INFO[0000] Proxy server start and listening on: 127.0.0.1:6601
  ```