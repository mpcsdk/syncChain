# sync block tx

## init db

```
psql   -U postgres -h localhost -p 5432 -c "create database syncChain"
```

```
psql   -U postgres -h localhost -p 5432 -d syncChain -f ./manifest/migration/syncChain.sql
```
