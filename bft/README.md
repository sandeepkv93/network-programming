## Byzantine Fault Tolerance (BFT)

Byzantine Fault Tolerance enables distributed systems to reach consensus even when some nodes are faulty or malicious (Byzantine failures).

## What is Byzantine Fault Tolerance?

The Byzantine Generals Problem describes the challenge of achieving consensus when some participants may be unreliable or malicious.

**Byzantine Failure**: Node exhibits arbitrary behavior
- Sends conflicting information
- Remains silent
- Sends corrupted data
- Acts maliciously

**PBFT (Practical Byzantine Fault Tolerance)**:
Most well-known BFT algorithm for asynchronous systems.

**Three Phases**:

**1. Pre-Prepare** (Primary):
- Primary receives client request
- Assigns sequence number
- Broadcasts Pre-Prepare message

**2. Prepare** (All Replicas):
- Replicas validate Pre-Prepare
- Broadcast Prepare message
- Wait for 2f+1 Prepare messages (prepared)

**3. Commit** (All Replicas):
- Broadcast Commit message
- Wait for 2f+1 Commit messages (committed)
- Execute request

**Requirements**:
- At least 3f+1 nodes to tolerate f Byzantine nodes
- Quorum: 2f+1 agreeing nodes
- Digital signatures for authentication

**Properties**:
- **Safety**: All honest nodes agree
- **Liveness**: System makes progress
- **Byzantine Resilience**: Works with < n/3 faulty nodes

## Comparison with Other Consensus

**Crash Fault Tolerance** (Paxos, Raft):
- Simpler, faster
- Assumes nodes fail by stopping (crash)
- Requires n/2+1 nodes (f+1 for f failures)

**Byzantine Fault Tolerance**:
- More complex, slower
- Handles malicious behavior
- Requires 3f+1 nodes (2f+1 for f failures)

## Applications

- Blockchain (some variants)
- Hyperledger Fabric
- Distributed databases
- Aerospace systems
- Financial systems

## Further Reading

- [Practical Byzantine Fault Tolerance - Castro & Liskov](http://pmg.csail.mit.edu/papers/osdi99.pdf)
- [Byzantine Generals Problem - Lamport](https://lamport.azurewebsites.net/pubs/byz.pdf)
- [Byzantine Fault - Wikipedia](https://en.wikipedia.org/wiki/Byzantine_fault)
