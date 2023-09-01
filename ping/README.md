## Go Ping Package

The ping package provides a simple Go implementation for sending ICMP echo requests, commonly known as "ping" commands. This document covers the purpose and the inner workings of the provided code and sheds light on the underlying principles of ICMP and the ping mechanism.

## Table of Contents

1. [What is ICMP?](#what-is-icmp)
2. [How Does Ping Work?](#how-does-ping-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is ICMP?

ICMP, which stands for Internet Control Message Protocol, is an error-reporting protocol which can be used to generate error messages to the source IP address when network problems prevent delivery of IP packets. ICMP can also be used to relay query messages.

**ICMP Echo Request and Reply**: These are the message types that power the "ping" tool. The Echo Request is sent out to a host, and the Echo Reply is the response sent back to the origin.

**Message Structure**: Each ICMP message has a specific structure, including:

-   `Type`: Type of the message (e.g., Echo Request or Echo Reply).
-   `Code`: A sub-code for the type.
-   `Checksum`: A checksum to ensure the integrity of the message.
-   `Identifier`: A unique identifier.
-   `Sequence Number`: A sequence number for the message.

### How Does Ping Work?

The "ping" tool operates by sending ICMP Echo Request messages to a target host and waiting for an Echo Reply. The time between sending the request and receiving the reply is known as the round-trip time.

1. **Sending the Request**: The source machine sends an ICMP Echo Request message to the target host.
2. **Waiting for a Reply**: The source machine waits for a certain amount of time for a response from the target host. If the target host is reachable, it will respond with an ICMP Echo Reply message. If the target host is not reachable, the request will time out.
3. **Calculating the Round-Trip Time**: When (and if) the reply is received, the source machine calculates the time difference between sending the request and receiving the reply. This is the round-trip time. The round-trip time is a measure of the latency of the connection between the source and the target.
4. **Timeouts**: If the source machine does not receive a reply within a certain timeframe, the request is considered to have "timed out". This could be due to the target host being unreachable, or the request being lost in transit.

### Understanding the Code

#### Constants:

-   `icmpEchoReply and icmpEchoRequest`: ICMP types for echo reply and request, respectively.
-   `defaultTimeout`: The maximum duration to wait for a reply.

#### Data Structures:

`icmpMessage`: Represents the ICMP message structure, including type, code, checksum, and sequence numbers.

#### Functions:

-   `Ping(host string) (time.Duration, error)`: Sends an ICMP request to the provided host and waits for a reply.
    -   It first resolves the IP of the host.
    -   Establishes a connection to the host using net.DialIP.
    -   Constructs an ICMP request message.
    -   Computes the checksum and sends the message.
    -   A goroutine waits for the reply or a timeout.
-   `computeChecksum(data []byte) uint16`: Computes the ICMP checksum, essential for integrity checks of ICMP packets.

### Further Reading

-   [RFC 792 - ICMP](https://datatracker.ietf.org/doc/html/rfc792): This is the original specification for the ICMP protocol.
-   [Go net package](https://pkg.go.dev/net): The official Go documentation for the net package, which our ping package relies on.
