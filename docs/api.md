# Library and CLI
Import `github.com/Nima0101/carebind` with `go get github.com/Nima0101/carebind@v0.1.0`,
a pinned source revision, or a local module replacement. The public tagged module
consumer path has been tested with checksum-database verification enabled. `Parse([]byte)` validates bounded JSON into Bundle; `Evaluate(Bundle)`
revalidates typed input and returns a Result or typed `*Error`. Supply your own
trusted Policy before evaluation. `Sign(Event, ed25519.PrivateKey)` creates an
envelope. Generate operational keys externally with OS randomness; no CLI key
persistence or enrollment service is implemented.

A result is a snapshot, not permission to act. Inspect state, reason, origin,
binding, target, rejected evidence and scope/freshness/physical limitations. The
CLI `evaluate BUNDLE LOCAL_POLICY` emits JSON and exits 0 for every well-formed
technical outcome, including conflicting/unbound/evidence-invalid. Exit 2 means
invalid arguments, input/policy or I/O. Never interpret exit 0 as identity verified.

`version`, `demo`, `fixture NAME`, and `demo-policy` provide synthetic demonstrations.
Shell redirection overwrites its destination before execution; use temporary files.
No result contains a private key, clinical payload, telemetry or automatic action.

Threading: Evaluate owns no mutable global state. Callers must not mutate input
slices while evaluation is in progress. There is no persistence or distributed
transaction claim. Merge by evidence union then re-evaluate under local policy;
input bounds include duplicates. No last-write-wins or hidden network/clock calls.
