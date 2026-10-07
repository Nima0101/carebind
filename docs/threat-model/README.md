# Threat model
Assets: original source context, bounded uncertainty, local trust policy and private
association metadata. Attackers can reorder, replay, strip, modify or fabricate
transport; possess unknown keys; compromise accepted keys; supply wrong physical
tags; and hide remote history. Trusted boundary: local policy, source capture,
Ed25519/SHA-256 implementation, host OS and library caller. No key proves truth.

| Threat | Control | Remaining limit |
|---|---|---|
| Reassignment/delay | Binding digest and generation at source | Adapter can sign a lie or capture too late |
| Mutation/substitution | Domain-separated signature, payload digest, receiver expectation | Digest does not prove payload meaning |
| Unknown/unauthorized issuer | Explicit key/role/scope policy | Operator may authorize a bad issuer |
| Key compromise | Local revocation, finite validity, rotation IDs | Offline peers may not know compromise |
| Replay/equivocation | Set deduplication, nonce collisions and generation conflict | No global delivery or completeness oracle |
| Missing history | Explicit unknown/unbound, stale policy | Selectively hidden unseen branches cannot be detected |
| Parser/resource abuse | Hard bounds, strict shape, fuzzing | Host resource starvation outside process remains possible |
| Cross-site tracking | Scoped links, no names/demographics or transitive merge | Pseudonyms/digests are linkable, not anonymous |
| Wrong physical tag | Mandatory unverified output and negative control | Cannot detect with digital evidence alone |
| Untrusted policy injection | CLI requires separate local policy | Caller must not accept transport trust configuration |

Keys should be generated with OS randomness, held outside exchanged bundles,
provisioned out of band, rotated with distinct key IDs and bounded validity.
Never log private keys; library accepts them in memory and demo uses conspicuously
public deterministic seeds only. Revocation is a local policy update; reevaluate
all history. No signature timestamp assurance, PKI discovery, HSM integration or
instant offline revocation is claimed. No real-world deployment authorized.
