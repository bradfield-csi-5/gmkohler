# DNS Records and Messages

## Resource Records
four-tuple containing the following fields:
```
(Name, Value, Type, TTL)
```

* If `Type=A`, `Name` is a hostname, `Value` is IP for the hostname.  e.g.:
  ```
  (relay1.foo.bar.com, 145.37.93.126, A, TTL)
  ```
  
* If `Type=NS`, `Name` is a domain (e.g. `foo.com`), `Value` is hostname of an
authoritative DNS server that knows how to obtain the IP address for hosts in 
the domain.  This record is used to route further along the chain.  e.g.:
  ```
  (foo.com, dns.foo.com, NS, TTL)
  ```

* if `Type=CNAME`, `Value` is a canonical hostname for the alias hostname
`Name`.  Provides querying hosts canonical name for a hostname.  e.g.:
  ```
  (foo.com, relay1.bar.foo.com, CNAME, TTL)
  ```

* If `Type=MX`, then `Value` is a canonical name of a mail server that has an
alias hostname `Name`.  e.g.: 
  ```
  (foo.com, mail.bar.foo.com, MX, TTL)
  ```
  .  The DNS client must query for an `MX` record to get canonical name of mail 
server, whereas it must query for a `CNAME` record for a web server.

If DNS is authoritative for a hostname, it has the `A` record for the hostname.
If it's not authoritative it still may have the `A` record cached.  If not 
authoritative for a hostname, it'll have the `NS` record and a type `A` record
to provide the IP address for the authoritative DNS server.

## DNS Messages
There are two kinds of messages, *query* and *reply*.  Both messages have the 
same format:

```
|----------------|----------------| \
| Identification |      Flags     |   |
|----------------|----------------|   |  "Header section"
|  # Questions   |   # Ans RRs    |   |-     12 bytes
|----------------|----------------|   |
|  # Auth RRs    |  # Add'l RRs   |   |
|----------------|----------------| /
|  Questions (variable number)    |
|---------------------------------|
|   Answers (variable # RRs)      |
|---------------------------------|
|   Authority (variable # RRs)    |
|---------------------------------|
|   Add'l info (variable # RRs)   |
|---------------------------------|
```

### Header

* Identification is 16-bit number identifying query.  Copied into reply message
  to allow client to match query/reply.
* Flags
  * 1-bit query(0)/reply(1)
  * 1-bit authoritative flag, set when DNS server is authoritative server for a
    queried name
  * 1-bit recursion-desired flag is set when client (host or DNS server) desires
    that the DNS server perform recursion when it doesn't have the record
  * 1-bit recursion-available flag (set in reply if DNS server supports
    recursion).
* Number-of fields indicate number of occurrences of the four types of data
  that follow the header

### Question Section
1. Name field that contains the name being queried
2. Type field indicates typeof question being asked about the name
   e.g. host address associated with the name (Type A), or mail server for a
   name (Type MX)

### Answer section

Resource records for the name that was originally queried.

### Authority section

records of other authoritative servers

### Additional section

other helpful records.  e.g. the answer field in a reply to an MX query contains
the resource record providing canonical hostname.  Additional section may
contain teh Type A record with the IP address for the canonical hostname.