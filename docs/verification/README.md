# Verification
Run `python3 scripts/verify.py` for frozen-hash checks, public-content/link scan,
unit/property tests, vet, formatting, build, two adapters, reference/OpenSSL checks
and demo assertions. Go and OpenSSL must be on PATH; OPENSSL can select OpenSSL 3.

Additional mandatory checks:

```sh
go test -race ./...
go test -run '^$' -fuzz FuzzParse -fuzztime=60s -parallel=2
go install golang.org/x/vuln/cmd/govulncheck@v1.8.0
govulncheck ./...
python3 scripts/reproduce.py
```

Release reproduction requires Python 3.12+ and a clean committed tree. Main tests
and adapter scripts require only Python 3.9+. No tests are intentionally skipped.
Fuzzing's duration is fixed before the run; its execution count is observed, never
an invented coverage claim. The generated-history oracle executes 1,200 trials.

| Claim | Executable evidence | Limit |
|---|---|---|
| Original context survives reassignment | TestFlagshipCases, integration.py | Honest source capture required |
| Merge laws | TestGeneratedMergeLawsAndReference | Finite generated histories; trust policy fixed |
| Scoped trust and crypto | TestTrustControls, TestRFC8032, conformance.py | Local policy, no physical truth |
| Strict parser | TestParserStrictness, FuzzParse | Bounded exploration |
| Conflict/gaps/links | TestHistories, TestLinks, TestObservationEquivocation | Supplied history only |
| Context-preserving adapters | integration.py, consumer test | Synthetic implementations, not EHR/device certification |
| Wrong tag blind spot | wrong-physical-tag demo and TestFlagshipCases | Explicitly cannot detect physical mismatch |
| Reproducible package | reproduce.py | Same toolchain/platform; unsigned local provenance |

[Matrix](verification-matrix.md) · [Gate](publication-gate.md) ·
[Status](status.md) · [Failure log](failure-log.md) · [Claim policy](claim-evidence-policy.md)
