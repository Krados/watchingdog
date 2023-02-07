# watchingdog

Watchingdog is a repo that we can implement our own alerting system, it can watch anything and if something went wrong then alert the user.

## Learn from examples

### Dummy

```
$ go run examples/dummy/dummy_backend.go
```

### ETCD
With the help of etcd we can make sure that any given time, there will be only one leader exists and only leader will do watch and alert stuff.

create etcd docker container if you dont have

```
$ docker run -itd --name etcd -p 2379:2379 -e ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd
```

then run this command in different terminal

```
$ go run examples/etcd/etcd_backend.go
```