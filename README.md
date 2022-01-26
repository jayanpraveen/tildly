<div align="center">
	<h1>
		<div>tildly</div>
		<img src="./.gitlab/assets/tildly-small.png" alt="tildly" width="480">
	</h1>
</div>

tildly is a simple yet efficient url-shortner which guarantees a unique short url for each request.

## Implementation

When an application instance connects to etcd a new counter range is assigned to it. For example when instance-1
and instance-2 connect to etcd, first a million keys is assigned to instance-1 and the next million keys is assigned
to instance-2.The hash (or the short-url) is generated with a 'counter value' prefixed to the orginal request which
increments after each use, by this way every hash is unique.

### Atomicity

Each counter increment is atomic and no two hashes can share the same 'counter value'. This is achieved using
[transactions in etcd](https://etcd.io/docs/v3.5/learning/api/#transaction) and by comparing the `ModRevision`
value before each increment.

## Running multiple instances

Run multiple instaces of the appliction with the `-port` flag

```makefile
go run main.go -port=8080
.
.
```

To run a multi-member etcd cluster and to test fault tolerance refer [this](https://etcd.io/docs/v3.5/dev-guide/local_cluster/).
or run `make etcd`
