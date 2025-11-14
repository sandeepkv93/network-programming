## Raft Consensus Algorithm

Raft is a consensus algorithm designed to be easy to understand. It's equivalent to Paxos in fault-tolerance and performance but structured differently.

## What is Raft?

Raft was designed as a more understandable alternative to Paxos, emphasizing clarity while maintaining the same guarantees.

**Node States**:
- **Follower**: Passive, responds to leaders and candidates
- **Candidate**: Actively seeking votes to become leader
- **Leader**: Handles all client requests, replicates log

**Key Components**:

**1. Leader Election**:
- Nodes start as followers
- If election timeout expires without heartbeat, become candidate
- Request votes from peers
- Majority vote → become leader
- Leader sends periodic heartbeats

**2. Log Replication**:
- Leader accepts client requests
- Appends to local log
- Replicates to followers via AppendEntries RPC
- Once majority replicate, commits entry
- Followers apply committed entries

**3. Safety**:
- **Election Safety**: At most one leader per term
- **Leader Append-Only**: Leader never overwrites log
- **Log Matching**: If two logs contain same entry, all preceding entries are identical
- **Leader Completeness**: Committed entry present in all future leaders
- **State Machine Safety**: If node applies entry, all nodes apply same entry at that index

**Terms**:
- Logical clock dividing time
- Each term starts with election
- At most one leader per term

## Advantages over Paxos

- Easier to understand and implement
- Stronger leadership (simplifies operation)
- Clear separation of concerns
- Membership changes handled explicitly

## Applications

- etcd (Kubernetes)
- Consul
- CockroachDB
- TiKV

## Further Reading

- [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)
- [Raft Visualization](https://raft.github.io/)
- [Raft - Wikipedia](https://en.wikipedia.org/wiki/Raft_(algorithm))
