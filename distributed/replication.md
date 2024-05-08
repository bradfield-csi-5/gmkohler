# Replicating our KV Store

The objective of this project step is to improve the fault tolerance, throughput, and latency of your KV store by way
of replication.  With just a single data store, we are able to avoid many of the challenges of distributed systems, but
also fail to take advantage of their benefits.  By replicating to a second location (or third, or many) we will "buy"
fault tolerance, higher throughput, and potentially lower latency, but it will "cost" the additional complexity of
maintaining a certain degree of consistency.

## 1. Manual recovery from a backup

Currently, the server writes to a single copy of th data store (most likely, a file).  If this node were to fail, we may
be in trouble!  We will certainly incur downtime, and may even lose data in the event of a disk failure or file system
corruption.

As a first step, implement a system to _back up_ data to a secondary location, which we could manually switch over to if
needed.

Some considerations:

* You will need to decide between synchronous and asynchronous replication.  Consider the tradeoffs!

* You may need to design encoding formats and a communication protocol for replication, or you may be able to
appropriate the protocol you use between the client and server.

* For now, assume that the backup will only be used if an administrator decides to manually switch to it.  That doesn't
mean you should entirely ignore replication lag, just that it's an issue that only will arise when the primary fails.

## 2. Scaling reads

While maintaining a "hot standby" replica will improve fault tolerance, it does nothing for throughput or latency.  As a
next step, modify your system to permit reads from multiple locations (writes will still go to a single leader).

Some considerations:

* What consistency guarantees would you like to make?  What suits the system?  Note that you may not be able to
implement your preferred consistency scheme until after our class on consistency and consensus.

* How would your design change if you modelled PUT requests as insertions into an append-only log, rather than a
mutable store?

* What is it about your specific use-case that makes it challenging or straightforward to scale reads horizontally with
replication?  Can you think of use cases that would be more or less challenging?

* Do you need to do load balancing?  If so, how does this change your architecture?

## 3. Automatic fail over

While it may not always be a good idea to fail over automatically, there are some contexts and systems where it may be
appropriate.  For instance, Dynamo assumes nodes will periodically fail and effectively "fails over" to replicas
automatically byt he way of sloppy quorum and hinted handoff.

As an optional stretch goal, implement a system where reads still succeed even if the primary replica is down.  As a
further stretch goal, implement a system to automatically fail over to another node for _writes_, if the primary is
down.  Note that this may introduce consistency issues; for now it may be sufficient to be aware of the issues rather
than designing robust solutions.  We will revisit the objective in our class on consistency and consensus.