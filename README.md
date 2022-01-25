<div align="center">
	<h1>
		<div>tildly</div>
		<img src="./.gitlab/assets/tildly-small.png" alt="tildly" width="480">
	</h1>
</div>

tildly is a simple yet efficient url-shortner which guarantees a unique short url for each request.

## Implementation

The hash (or the short-url) is generated with a 'counter value' prefixed to the orginal request which increments after each use, by this way every hash is unique.

### Atomicity

Each counter increment is atomic and no two instances can share the same 'counter value' for generating the hash. This is achieved using
[transactions in etcd](https://etcd.io/docs/v3.5/learning/api/#transaction) and by comparing the `ModRevision` value before each increment.

> #### From the etcd docs
>
> A transaction is an atomic If/Then/Else construct over the key-value store... Transactions can be used for protecting keys from unintended concurrent updates, building compare-and-swap operations, and developing higher-level concurrency control.

## Known issues

In case of the `ModRevision`'s do not match, the system fails silenty without retrying.
