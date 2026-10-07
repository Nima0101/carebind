# Evidence status — 2026-10-08

State: P5 adversarially verified; clean-clone and release gates are in progress.
Not production-ready. No public repository or release has been published yet.

| Gates | Evidence / current state |
|---|---|
| V01 | Passed: governance commit e4c77e5 precedes production source; frozen hashes match |
| V02 | Passed: primary-source landscape and ADR 0001; substantial PCIM overlap disclosed |
| V03–V05 | Passed locally: strict parser, protocol vectors, Ed25519/OpenSSL, trust controls, historical and conflict tests |
| V06 | Passed locally: 1,200 generated histories with pairwise reference and merge laws |
| V07–V09 | Passed locally: real Go/Python process adapters; offline states; explicit wrong-tag software blind spot |
| V10 | Passed: 60-second native fuzz, 240,089 executions, no crash; finite exploration only |
| V11 | Passed locally: go vet, race, govulncheck v1.8.0 (no vulnerabilities), public-content scan; no independent audit |
| V12 | Pending clean committed clone and separate consumer module |
| V13 | Pending deterministic packaged release verification |
| V14 | Actual CLI stdout/cast/PNG/GIF and original SVG mark; recording provenance retained; visual inspected |
| V15 | OSS documents, templates, pinned CI/CodeQL/release workflows present; hosted execution pending |
| V16–V17 | Pending publication, hosted checks, anonymous clone, tagged attested release |

Local tools: Go 1.27.1 darwin/arm64, Python 3.9.6 (main checks), Python 3.13
(release/visual tools), OpenSSL 3. Fuzz/race passed before test-only additions;
new adapter regression and subsequent gates are run explicitly. Core statement
coverage observed 91.1% before additional boundary tests; coverage is descriptive,
not a frozen acceptance threshold. Process integration coverage is not included
in that Go unit-test metric. See failure-log.md for discovered/corrected issues.

All examples are synthetic; no clinical deployment, patient identity, safety
benefit, adoption, standards certification or real medical-device integration is
claimed. Local evidence does not establish hosted platform support.
