# bufferserver

the bufferserver of bytom dapp.

## Requirements

Building requires [Go](https://golang.org/doc/install) version 1.8 or higher, with `$GOPATH` set to your preferred directory. Furthermore, the version of [Mysql 5.7](https://www.mysql.com/) and the latest stable version of [Redis](https://redis.io/) are needed.

The target of building source code are the `api` and `updater` binaries. The `api` provide the RPC request and response service for users, and the `updater` provide the synchronization services for `blockcenter` and `browser`. The `blockcenter` is the decentralized wallet server for `bytom`, and the `browser` is the `bytom` blockchain browser for searching transaction.

```bash
$ make all
```

then change directory to `target`, and you can find the compiled results.

## Getting Started

### Create database and tables

`dump.sql` contain the SQL of creating database and tables.

```bash
$ mysql -u root -p < database/dump.sql
```

Enter the correct of `root` password that will create database and tables successfully.

### Modify config

the config file is `config_local.json`, these parameters can be changed according to developer needs. 

```js
{
  "gin-gonic": {
    "listening_port": 3100,   // the port of API service
    "is_release_mode": false  // the release mode of gin-gonic, it's debug mode with false
  },
  "mysql": {
    "master": {
      "host": "127.0.0.1",  // the IP of database server
      "port": 3306,         // the port of database server, default is 3306
      "username": "root",   // the username of server 
      "password": "root",   // the password of server
      "database": "deposit" // the name of database
    },
    "log_mode": true        // the log mode, true print detailed logs, false only print error logs 
  },
  "redis": {
    "endpoint": "127.0.0.1:6379", // the IP and port of redis
    "pool_size": 10,              // the pool size
    "database": 6,                // the category of database
    "password": "password",       // the passwprd
    "cache_seconds": 600,         // the expiration of cache
    "long_cache_hours": 24        // the long expiration of cache
  },
  "api": {
    "mysql_conns": {
      "max_open_conns": 20,   // the max open connects
      "max_idle_conns": 10    // the max idle connects
    }
  },
  "updater": {
    "block_center": {
      "sync_seconds": 10,               // the synchronization interval
      "url": "http://127.0.0.1:3000",   // the url of blockcenter server
      "mysql_conns": {
        "max_open_conns": 5,    // the max open connects
        "max_idle_conns": 5     // the max idle connects
      }
    }
  }
}
```

### Startup service

startup `api` and `updater` service with the target binaries.

```bash
$ ./target/api config_local.json

$ ./target/updater config_local.json
```

The RPC description refer to [JSON API](https://github.com/oysheng/bufferserver/wiki).