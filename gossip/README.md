## Gossip Protocol

Gossip protocol is a peer-to-peer communication protocol inspired by how rumors spread in social networks. It's used for information dissemination in distributed systems.

## What is Gossip Protocol?

Also known as "epidemic protocol", gossip protocol spreads information by having each node periodically exchange data with a random subset of peers.

**Characteristics**:
- **Decentralized**: No central coordinator
- **Scalable**: Logarithmic message complexity
- **Fault-tolerant**: Continues working with node failures
- **Eventually consistent**: All nodes converge to same state

**Applications**:
- Cluster membership (Cassandra, Consul)
- Failure detection
- Database replication
- Blockchain consensus (some variants)

## How It Works

1. **Periodic Selection**: Every T seconds, select k random peers
2. **Information Exchange**: Send local state to selected peers
3. **Merge**: Peers merge received information with local state
4. **Propagation**: Process repeats, spreading information exponentially

**Parameters**:
- **Fanout**: Number of peers to gossip to (typically 3-7)
- **Interval**: Time between gossip rounds (typically 1-10 seconds)
- **TTL**: Message time-to-live (prevents infinite loops)

**Conflict Resolution**:
- Last-write-wins (timestamp)
- Vector clocks
- CRDTs (Conflict-free Replicated Data Types)

## Further Reading

- [Gossip Protocol - Wikipedia](https://en.wikipedia.org/wiki/Gossip_protocol)
- [Epidemic Algorithms for Replicated Database Maintenance](https://dl.acm.org/doi/10.1145/41840.41841)
