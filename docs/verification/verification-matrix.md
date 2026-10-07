# Frozen verification matrix v1
All rows mandatory for prepublication unless explicitly labeled post-publication.

| Gate | Required experiment and acceptance |
|---|---|
| V01 governance | Hashes match pre-source freeze; no weakened obligations |
| V02 research | Primary standards + strongest prior art + package/GitHub synonyms; material useful boundary documented, no novelty assertion |
| V03 wire | Versioned schema, strict unknown/duplicate/type/depth/size rejection; compatibility vectors |
| V04 crypto | Ed25519 known-answer verification, tampering, wrong key, revoked/expired/unknown issuer, role/scope violations; key lifecycle docs |
| V05 history | Delayed Cedar observation after Birch reassignment retains Cedar; same-generation conflicts/forks, gaps and supersession explicit |
| V06 merge | Generated order/duplicate/permutation/partition laws against separately expressed reference semantics, at least 1000 generated histories |
| V07 adapters | Two distinct local representations preserve full context to consumer; stripping/substitution fail; scoped link and synthetic item transfer |
| V08 offline | Explicit policy time/age, missing/revoked evidence, reconnection conflict, no global-current assertion |
| V09 blind spot | Valid evidence with wrong fictional physical tag stays digitally consistent and explicitly physically unverified |
| V10 parser fuzz | At least 60 seconds native fuzz plus seeded malformed corpus; zero unresolved crashes or uncontrolled resource failures |
| V11 security | Threat model, static analysis, race check, dependency/vulnerability audit, secret and public-content leak scans; no unresolved high/critical issues |
| V12 clean consumer | Fresh git clone, README commands verbatim, public API consumer and both adapters work; no reliance on workspace files |
| V13 release | Deterministic release rebuilt twice byte-identically; SHA-256, SBOM, source provenance, license inventory, packaging/install check |
| V14 presentation | Original mark, accurate architecture, real command recording and screenshot/animation, all required hard cases, links valid |
| V15 repository | License, notices, governance, contributing, conduct, security, changelog, issue/PR templates, CI/security/release workflows |
| V16 remote | Post-publication: hosted CI on every claimed platform, anonymous clone quickstart, actual remote README/demo/checks verified |
| V17 provenance | Post-publication: tagged source, hosted artifact attestation and downloadable verified release artifacts |

No throughput target is claimed. Start with one local platform; broader support
needs actual CI evidence. No network service/concurrency persistence claim: batch
immutable evidence is the boundary; race checks still mandatory for library use.
