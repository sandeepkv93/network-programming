## Paxos

Paxos is a family of protocols for solving consensus in a network of unreliable processors. It guarantees safety (never returning incorrect result) but not liveness (may not terminate).

## What is Paxos?

Paxos was introduced by Leslie Lamport and is one of the most important consensus algorithms in distributed systems.

**Roles**:
- **Proposer**: Proposes values
- **Acceptor**: Accepts or rejects proposals
- **Learner**: Learns the chosen value

**Two Phases**:

**Phase 1: Prepare**
1. Proposer selects proposal number n
2. Sends Prepare(n) to acceptors
3. Acceptors promise not to accept proposals < n
4. Acceptors respond with highest-numbered proposal they've accepted

**Phase 2: Accept**
1. If quorum responds, proposer sends Accept(n, v) to acceptors
2. v = highest-numbered accepted value from Phase 1, or proposer's value
3. Acceptors accept if they haven't promised to ignore this number
4. When quorum accepts, consensus is reached

**Properties**:
- **Safety**: Only one value is chosen
- **Liveness**: Eventually a value is chosen (with assumptions)
- **Fault Tolerance**: Works with < n/2 failures

## Applications

- Google Chubby (lock service)
- Apache ZooKeeper (coordination)
- Distributed databases
- Cloud storage systems

## Further Reading

- [Paxos Made Simple - Leslie Lamport](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf)
- [The Part-Time Parliament](https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf)
- [Paxos - Wikipedia](https://en.wikipedia.org/wiki/Paxos_(computer_science))
