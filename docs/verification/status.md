# Evidence status — 2026-10-08

State: P9 — all frozen V01–V17 gates passed for the **bounded technical runtime
v0.1.0**, source `9a0439f5efbcea15d60fbfd17ffc8e4ec49b1b0b`. This establishes the
project's defined technical production-readiness gate, **not clinical deployment
readiness**, physical identity, or globally current truth.

[Machine-readable local evidence](local-evidence.json) binds the tested source
content and commands. Final publication checks rerun on the committed candidate.
The unpublished history was normalized to GitHub noreply author metadata before
publication; file contents and constitution-before-source ordering are preserved.

| Gate | Evidence |
|---|---|
| V01 | Constitution commit f0778a9 precedes production source; immutable hashes checked by verify.py |
| V02 | Primary-source landscape and ADR 0001; PCIM overlap disclosed; public portfolio/name rechecked |
| V03 | Protocol frozen at e25cd87; strict structural schema, parser bounds including exact 4096/4097-byte regression, committed wire vectors |
| V04 | RFC 8032, independent OpenSSL verification, every-field tampering, role/scope/time/revocation controls |
| V05 | Cedar/Birch delay, forks, missing ancestors, explicit unbinding, observation equivocation, supersession |
| V06 | 1,200 generated histories; pairwise oracle, three partitions, ordering and replay laws |
| V07 | Actual Go-to-Python CLI path; local numeric slot; context stripping/substitution; malformed original wire rejected before adaptation |
| V08 | Explicit local policy age, unknown issuer, revoked key, reconnection fork; global freshness always unknown |
| V09 | Wrong-tag software negative control; digitally unchanged evidence cannot establish physical identity |
| V10 | Corrected parser: fresh 60-second fuzz run, 76,484 executions, no crash; earlier run 240,089 executions did not find the later manual size defect |
| V11 | go vet, full race suite, govulncheck v1.8.0 (no vulnerabilities), public content/history scan; same-author adversarial review |
| V12 | Clean git clone, isolated caches, README commands verbatim, independent consumer module, clean worktree |
| V13 | Two byte-identical builds, packaged demo, SHA-256, SPDX 2.3 validated against upstream schema, source archive/local provenance |
| V14 | Original SVG mark, accurate architecture, captured real stdout and terminal replay; visual inspected and source/output hashes checked |
| V15 | OSS documents/templates; pinned CI, CodeQL and attestation/release workflows validated by actionlint v1.7.12 |
| V16 | Passed: [hosted Linux/macOS CI and CodeQL](https://github.com/Nima0101/carebind/actions/runs/37701720986); no open CodeQL alerts; anonymous clone quickstart/verification; actual public HTML and image verified |
| V17 | Passed: [v0.1.0 release](https://github.com/Nima0101/carebind/releases/tag/v0.1.0), [tagged hosted build](https://github.com/Nima0101/carebind/actions/runs/37702792138), all six attestations verified, anonymous downloads match, source archive and public Go module consumer pass |

Local platform: Go 1.27.1 darwin/arm64; Python 3.9.6 for main checks and 3.13 for
release/visual tooling; OpenSSL 3.6.4. No Windows support claim. Hosted platform
support requires the actual remote results. Local unsigned provenance is distinguished from the separately verified
hosted attestations. The [failure log](failure-log.md) records defects and retained
regressions; no gate was relaxed and no mandatory tests were skipped.

No clinical benefit, safety certification, adoption, physical trial, independent
human audit or real medical-device integration is established by this evidence.

Public repository: https://github.com/Nima0101/carebind. Private vulnerability
reporting is enabled. The initial public commit b05d55c passed all hosted jobs. The exact tagged source
also passed [Linux/macOS CI and CodeQL](https://github.com/Nima0101/carebind/actions/runs/37702163803)
before tagging. [Machine-readable release evidence](release-evidence.json) records
asset hashes, attestation policy and public-consumer checks. This status update is
a documentation-only follow-up; it does not move or overwrite the release tag.
