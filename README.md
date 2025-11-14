## Network Programming

### Table of Contents

| No. | Problem Statement & Explanation                      | Code                                                                                               |
| --- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| 1   | [UDP Client & Server](./udp/README.md)               | [server.go](./udp/server.go), [client.go](./udp/client.go), [stress_test.go](./udp/stress_test.go) |
| 2   | [Ping](./ping/README.md)                             | [ping.go](./ping/ping.go)                                                                          |
| 3   | [TCP Client & Server](./tcp/README.md)               | [server.go](./tcp/server.go), [client.go](./tcp/client.go)                                         |
| 4   | [HTTP Client & Server](./http/README.md)             | [server.go](./http/server.go), [client.go](./http/client.go)                                       |
| 5   | [Echo Client & Server](./echo/README.md)             | [server.go](./echo/server.go), [client.go](./echo/client.go)                                       |
| 6   | [DNS Server](./dns/README.md)                        | [server.go](./dns/server.go), [client.go](./dns/client.go)                                         |
| 7   | [FTP Client & Server](./ftp/README.md)               | [server.go](./ftp/server.go), [client.go](./ftp/client.go)                                         |
| 8   | [Port Scanner](./portscanner/README.md)              | [scanner.go](./portscanner/scanner.go)                                                             |
| 9   | [Proxy Server](./proxy/README.md)                    | [server.go](./proxy/server.go)                                                                     |
| 10  | [Mail Client & Server](./mail/README.md)             | [server.go](./mail/server.go), [client.go](./mail/client.go)                                       |
| 11  | [SMTP Client & Server](./smtp/README.md)             | [server.go](./smtp/server.go), [client.go](./smtp/client.go)                                       |
| 12  | [Load Balancer](./loadbalancer/README.md)            | [loadbalancer.go](./loadbalancer/loadbalancer.go)                                                  |
| 13  | [Simple Chat App](./chat/README.md)                  | [server.go](./chat/server.go), [client.go](./chat/client.go)                                       |
| 14  | [Telnet Client & Server](./telnet/README.md)         | [server.go](./telnet/server.go), [client.go](./telnet/client.go)                                   |
| 15  | [SSH Client & Server](./ssh/README.md)               | [server.go](./ssh/server.go), [client.go](./ssh/client.go)                                         |
| 16  | [File Transfer over network](./filetransfer/README.md) | [server.go](./filetransfer/server.go), [client.go](./filetransfer/client.go)                     |
| 17  | [Ifconfig implementation](./ifconfig/README.md)      | [ifconfig.go](./ifconfig/ifconfig.go)                                                              |
| 18  | [Traceroute implementation](./traceroute/README.md)  | [traceroute.go](./traceroute/traceroute.go)                                                        |
| 19  | [Netstat implementation](./netstat/README.md)        | [netstat.go](./netstat/netstat.go)                                                                 |
| 20  | [nslookup implementation](./nslookup/README.md)      | [nslookup.go](./nslookup/nslookup.go)                                                              |
| 21  | [ARP implementation](./arp/README.md)                | [arp.go](./arp/arp.go)                                                                             |
| 22  | [DHCP implementation](./dhcp/README.md)              | [server.go](./dhcp/server.go), [client.go](./dhcp/client.go)                                       |
| 23  | [IP Scanner](./ipscanner/README.md)                  | [scanner.go](./ipscanner/scanner.go)                                                                |
| 24  | [WebSocket Client & Server](./websocket/README.md)   | [server.go](./websocket/server.go), [client.go](./websocket/client.go)                             |
| 25  | [WebRTC Client & Server](./webrtc/README.md)         | [server.go](./webrtc/server.go), [client.go](./webrtc/client.go)                                   |
| 26  | [VPN Client & Server](./vpn/README.md)               | [server.go](./vpn/server.go), [client.go](./vpn/client.go)                                         |
| 27  | [Remote Execution](./remoteexec/README.md)           | [server.go](./remoteexec/server.go), [client.go](./remoteexec/client.go)                           |
| 28  | [Remote Login](./remotelogin/README.md)              | [server.go](./remotelogin/server.go), [client.go](./remotelogin/client.go)                         |
| 29  | [Remote Procedure Call](./rpc/README.md)             | [server.go](./rpc/server.go), [client.go](./rpc/client.go)                                         |
| 30  | [Tunneling](./tunneling/README.md)                   | [server.go](./tunneling/server.go), [client.go](./tunneling/client.go)                             |
| 31  | [Heartbeat Server](./heartbeat/README.md)            | [server.go](./heartbeat/server.go), [client.go](./heartbeat/client.go)                             |
| 32  | [Rate Limiter](./ratelimiter/README.md)              | [limiter.go](./ratelimiter/limiter.go)                                                              |
| 33  | [Web Crawler](./webcrawler/README.md)                | [crawler.go](./webcrawler/crawler.go)                                                               |
| 34  | [Packet Sniffer](./packetsniffer/README.md)          | [sniffer.go](./packetsniffer/sniffer.go)                                                            |
| 35  | [Port Forwarding](./portforwarding/README.md)        | [forwarder.go](./portforwarding/forwarder.go)                                                       |
| 36  | [Content Delivery Network](./cdn/README.md)          | [server.go](./cdn/server.go), [origin.go](./cdn/origin.go)                                          |
| 37  | [HTTPS Client & Server](./https/README.md)           | [server.go](./https/server.go), [client.go](./https/client.go)                                      |
| 38  | [FTPS Client & Server](./ftps/README.md)             | [server.go](./ftps/server.go), [client.go](./ftps/client.go)                                        |
| 39  | [SFTP Client & Server](./sftp/README.md)             | [server.go](./sftp/server.go), [client.go](./sftp/client.go)                                        |
| 40  | [Voice over IP](./voip/README.md)                    | [server.go](./voip/server.go), [client.go](./voip/client.go)                                        |
| 41  | [Video over IP](./videoip/README.md)                 | [server.go](./videoip/server.go), [client.go](./videoip/client.go)                                  |
| 42  | [Video Streaming](./videostreaming/README.md)        | [server.go](./videostreaming/server.go), [client.go](./videostreaming/client.go)                    |
| 43  | [Video Conferencing](./videoconference/README.md)    | [server.go](./videoconference/server.go)                                                             |
| 44  | [IP Spoofing](./ipspoofing/README.md)                | [spoof.go](./ipspoofing/spoof.go)                                                                   |
| 45  | [Firewall](./firewall/README.md)                     | [firewall.go](./firewall/firewall.go)                                                               |
| 46  | [IP Tables](./iptables/README.md)                    | [iptables.go](./iptables/iptables.go)                                                               |
| 47  | [OAuth 2.0 Client & Server](./oauth/README.md)       | [server.go](./oauth/server.go), [client.go](./oauth/client.go)                                      |
| 48  | [Two Factor Authentication](./twofa/README.md)       | [twofa.go](./twofa/twofa.go)                                                                         |
| 49  | [Gossip Protocol](./gossip/README.md)                | [protocol.go](./gossip/protocol.go)                                                                 |
| 50  | [Distributed Hash Table](./dht/README.md)            | [dht.go](./dht/dht.go)                                                                               |
| 51  | [Paxos](./paxos/README.md)                           | [paxos.go](./paxos/paxos.go)                                                                         |
| 52  | [Raft](./raft/README.md)                             | [raft.go](./raft/raft.go)                                                                            |
| 53  | [Byzantine Fault Tolerance](./bft/README.md)         | [bft.go](./bft/bft.go)                                                                               |
| 54  | [Consensus](./consensus/README.md)                   | [consensus.go](./consensus/consensus.go)                                                             |
