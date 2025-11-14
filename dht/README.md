## Distributed Hash Table (DHT)

A Distributed Hash Table is a decentralized distributed system that provides a lookup service similar to a hash table, where key-value pairs are stored across a distributed network of nodes.

## What is a DHT?

DHTs distribute data across many nodes in a network, allowing efficient lookup without central coordination.

**Popular DHT Algorithms**:
- **Chord**: Circular identifier space with finger tables
- **Kademlia**: XOR metric-based routing (used in BitTorrent)
- **Pastry**: Prefix-based routing
- **CAN**: Coordinate-based routing in d-dimensional space

**Key Properties**:
- **Scalability**: Supports millions of nodes
- **Fault Tolerance**: Continues working with node failures
- **Decentralization**: No single point of failure
- **Efficiency**: O(log N) lookup complexity

## How Chord Works

This implementation is inspired by Chord:

1. **Identifier Space**: Nodes and keys hashed to m-bit identifier circle
2. **Successor**: Each node knows its successor on the circle
3. **Finger Table**: Shortcuts for O(log N) lookups
4. **Routing**: Forward request to closest preceding node

**Finger Table**:
- finger[i] = successor of (n + 2^i) mod 2^m
- Provides exponential distance coverage
- Enables logarithmic-time lookups

**Key Operations**:
- **Put(key, value)**: Store key-value pair
- **Get(key)**: Retrieve value by key
- **Join**: Add node to network
- **Leave**: Remove node (with data transfer)

## Further Reading

- [Chord: A Scalable Peer-to-peer Lookup Protocol](https://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf)
- [Kademlia](https://en.wikipedia.org/wiki/Kademlia)
- [Distributed Hash Table - Wikipedia](https://en.wikipedia.org/wiki/Distributed_hash_table)
