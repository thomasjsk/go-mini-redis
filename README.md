# go-mini-redis

Go implementation of a simple Redis server
### Supports:
1. Basic configuration like `listening-port`
2. Concurrent client connections
3. Replication from a master instance
- performing full master-slave handshake
- sending empty rdb replication file
- propagating `SET` commands from master to replica
4. Implementing TCP [redis serialization protocol](https://redis.io/docs/latest/develop/reference/protocol-spec/)
- reading `SimpleString`, `BulkString`, `Array` and `Null` types
5. Handling commands:
- [PING](https://redis.io/docs/latest/commands/ping/)
- [ECHO](https://redis.io/docs/latest/commands/echo/)
- [SET](https://redis.io/docs/latest/commands/set/) with expiry (i.e. `SET foo bar px 1000`)
- [GET](https://redis.io/docs/latest/commands/get/)
- [INFO](https://redis.io/docs/latest/commands/info/)
