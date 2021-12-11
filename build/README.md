# MySQLOnDocker5.7
MySQL5.7

# Use
```
# dumpfileをコンテナへコピー
docker cp {DumpFile}.sql {ContainerName:/tmp/{DumpFile}.sql


mysql -u {User} /tmp/{DumpFile}.sql < {DumpFile}.sql
```
