# ADR 0001 — bounded offline association evidence
Accepted before source. IHE PCIM already addresses history and conflicts. Build a
portable evaluator and conformance kit rather than an alleged replacement standard.
Use Go standard-library Ed25519/SHA-256/JSON, no runtime dependencies, no server or
database. This keeps transport, clinical content, key custody and physical checks
outside the core. A pure batch API makes partition/merge laws directly testable.
Two synthetic adapters (Go source channel and Python receiving slot) must preserve
signed envelopes. A separately expressed reference evaluator challenges semantics.
This diversity is not independent authorship or clinical assurance.
