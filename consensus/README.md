## Consensus

Consensus is the process of achieving agreement among distributed processes or systems. It's fundamental to distributed computing.

## What is Consensus?

Consensus algorithms allow multiple nodes to agree on a single value, even in the presence of failures.

**Requirements** (FLP Impossibility addresses these):
- **Agreement**: All correct nodes decide on the same value
- **Validity**: Decided value was proposed by some node
- **Termination**: All correct nodes eventually decide
- **Integrity**: Nodes decide at most once

**FLP Impossibility**: In asynchronous systems with even one crash failure, no deterministic consensus algorithm can guarantee termination.

**Solutions**:
- Randomization (eventual consensus)
- Partial synchrony assumptions
- Failure detectors
- Weakened requirements

## Common Consensus Algorithms

**1. Majority/Voting**:
- Simplest form
- Requires majority agreement
- Fast but not fault-tolerant to Byzantine failures

**2. Quorum-Based**:
- Requires Q votes (Q configurable)
- Flexible quorum sizes
- Used in many distributed databases

**3. Two-Phase Commit (2PC)**:
- Prepare phase
- Commit phase
- Blocking (coordinator failure blocks system)

**4. Three-Phase Commit (3PC)**:
- Non-blocking variant of 2PC
- Adds pre-commit phase
- More complex

**5. Paxos**:
- Proven safe, tolerates crashes
- Complex to understand and implement

**6. Raft**:
- Understandable alternative to Paxos
- Leader-based
- Widely adopted

**7. PBFT**:
- Byzantine fault tolerant
- Requires 3f+1 nodes
- More expensive but handles malicious nodes

## Use Cases

**Distributed Databases**:
- Transaction commits
- Data replication
- Shard coordination

**Blockchain**:
- Block validation
- State transitions
- Network agreement

**Distributed Systems**:
- Leader election
- Configuration management
- Service discovery

**Cloud Computing**:
- Resource allocation
- Load balancing
- Failover coordination

## This Implementation

Provides several consensus variants:

**SimpleConsensus**: Basic majority voting
**QuorumConsensus**: Configurable quorum
**TimedConsensus**: With timeout support

## Further Reading

- [Consensus - Wikipedia](https://en.wikipedia.org/wiki/Consensus_(computer_science))
- [FLP Impossibility](https://en.wikipedia.org/wiki/Consensus_(computer_science)#Impossibility_results)
- [Byzantine Fault Tolerance](https://en.wikipedia.org/wiki/Byzantine_fault)
- [CAP Theorem](https://en.wikipedia.org/wiki/CAP_theorem)
