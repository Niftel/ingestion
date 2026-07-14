# ingestion

Praetor's **ingestion** service — the HTTP intake for everything a running job
reports back to the control plane.

It exposes the run-scoped endpoints the executor and host-runner call: job events,
log chunks (bytes to the JetStream object store, an index notification onto the
[`eventbus`](https://github.com/praetordev/eventbus)), and the credential-resolve
endpoint that hands a run its decrypted injectors. A published event's 2xx is what
tells the host-runner's syncer it may advance its cursor, so the endpoint blocks on
durable persistence — that's the contract that makes an outage lossless.

It is a leaf deployable: nothing imports it in production. It depends only on the
shared `praetordev/*` libraries (`eventbus`, `events`, `models`, `store`, `db`,
`credentials`, `objectstore`, `render`, `runtoken`, `metrics`, `env`, `plog`).

## Layout

```
main.go            entrypoint + auth middleware
core/              IngestionService: events, log chunks, credential resolve
handler/           chi HTTP handlers + metrics
inventoryrender/   renders inventory for ansible-inventory
```

## Build the image

```
docker build -t praetor-ingestion:latest .
```

Stable image name (`praetor-ingestion`) so the Helm chart and k3d/kind load step
are unaffected by the repo split. Serves on `:8081`.

## Tests

Unit tests (auth middleware) run standalone. The DB/NATS-backed log-chunk
integration test is gated on `TEST_DATABASE_URL` / `TEST_NATS_URL` and skips
without them; it exercises the full path through to the consumer's indexing.

```
go test ./...
```
