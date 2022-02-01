CREATE KEYSPACE IF NOT EXISTS tildly
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE tildly.url (
  hash text,
  LongUrl text,
  createdat bigint,
  exipireat bigint,
  PRIMARY KEY ((hash), createdat)
)
WITH CLUSTERING ORDER BY (createdat DESC);
