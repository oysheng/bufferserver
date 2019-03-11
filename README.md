# bufferserver

the bufferserver of bytom dapp.

## Requirements

[Go](https://golang.org/doc/install) version 1.8 or higher, with `$GOPATH` set to your preferred directory

Furthermore, the version of [Mysql 5.7](https://www.mysql.com/) and the latest stable version of [Redis](https://redis.io/) are needed.

The target of building source code is the `api` and `updater` binaries. The `api` provide the RPC request and response service for users, and the `updater` provide the synchronization services for `blockcenter` and `browser`. The `blockcenter` is the decentralized wallet server for `bytom`, and the `browser` is the `bytom` blockchain browser for searching transaction.

```bash
$ make all
```

then change directory to `target`, and you can find the binaries:

```bash
$ cd target
```

## setup

### create database

```bash
$ mysql -u root -p > database/dump.sql
```

Enter the correct of `root` password that will create the database successfully.

### modify config

```js
{
  "gin-gonic": {
    "listening_port": 3100,  // the port of API service
    "is_release_mode": false
  },
  "mysql": {
    "master": {
      "host": "127.0.0.1",  // the IP of database server
      "port": 3306,         // the port of database server, default is 3306
      "username": "root",   // the username of server 
      "password": "yang",   // the password of server
      "database": "deposit" // the name of database
    },
    "log_mode": true
  },
  "redis": {
    "endpoint": "52.82.31.236:6379",    // the IP and port of redis
    "pool_size": 10,    
    "database": 6,
    "password": "block12345",
    "cache_seconds": 600,
    "long_cache_hours": 24
  },
  "api": {
    "mysql_conns": {
      "max_open_conns": 20, // the max open connects
      "max_idle_conns": 10  // the max idle connects
    }
  },
  "updater": {
    "block_center": {
      "sync_seconds": 10,   // the synchronization interval
      "url": "http://127.0.0.1:3000",   // the url of blockcenter server
      "mysql_conns": {
        "max_open_conns": 5,
        "max_idle_conns": 5
      }
    },
    "browser": {
      "sync_seconds": 60,   // the synchronization interval
      "expiration_hours": 24,   // the expiration time of hours
      "url": "https://blockmeta.com/api/wisdom" // the url of browser server
    }
  }
}
```

### startup service

startup `api` and `updater` service.

```bash
$ ./target/api config_local.json
```

```bash
$ ./target/updater config_local.json
```