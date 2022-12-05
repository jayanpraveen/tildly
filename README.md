<div align="center">
	<h1>
		<div>tildly</div>
		<img src="./.gitlab/assets/tildly-small.png" alt="tildly" width="480">
	</h1>
</div>

tildly, a simple yet efficient url-shortner which gives a unique short url for each request.

## Implementation

When an application instance connects to etcd, a new counter range is assigned to it. Say two instance of the application, 'instance-1'
and instance-2' connects to etcd, first a million keys is assigned to instance-1 and the next million keys is assigned
to instance-2. The hash (or the 'short-url') is generated with a 'counter value' prefixed to the orginal request which
increments after each use, by this way each hash is _\*unique_\*.

### Atomicity

Each counter increment is atomic and no two hashes can share the same 'counter value'. This is achieved using
[transactions in etcd](https://etcd.io/docs/v3.5/learning/api/#transaction) and by comparing the [ModRevision](https://github.com/etcd-io/etcd/issues/6518)
before each increment.

## Sample Request

```sh
{
    "longUrl": "https://www.gnu.org",
    "exipreAt": 0
}
```

## Running locally

Run multiple instaces of the appliction with the `-port` flag, the default port is `:8080`. The replication factor can be modifed in the [cassandra.go](datastore/cassandra.go) file.

```makefile
go run main.go -port=8080
go run main.go -port=8081
go run main.go -port=8082
. ...
```

To run a multi-member etcd cluster and to test fault tolerance refer [this](https://etcd.io/docs/v3.5/dev-guide/local_cluster/).
or run `make etcd`. To test fault tolerance in Cassandra, simply bring down a node and perform a query in anoter node.

## Using Docker

using docker compose:

```docker
docker-compose up -d
```

This spins up the application(port `:8080`) with an etcd and cassandra cluster.
