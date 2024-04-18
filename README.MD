# sync block tx

## init db

```
psql   -U postgres -h localhost -p 5432 -c "create database sync_chain"
```

```
psql   -U postgres -h localhost -p 5432 -d sync_chain -f ./mpcCommon/mpcdao/migration/sync_chain.sql
```

## rpc

eth=https://sepolia.infura.io/v3/99f75c99059d42908ae78ca4b150795d,https://sepolia.infura.io/v3/32cf16fa990742ac889fbb7a34cbb496,https://sepolia.infura.io/v3/e00918a38d09462c8fcf7f08031e7f61,https://sepolia.infura.io/v3/98bf75163f454d5f85865d78a597170c

eth = https://eth.meowrpc.com
eht=https://robin.rangersprotocol.com/pubhub/api/jsonrpc
