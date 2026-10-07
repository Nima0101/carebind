# CareBind engineering constitution v1
Frozen before production source, 2026-10-08. Scope: a payload-opaque association
evidence runtime and conformance kit for synthetic integration engineering.

## Invariants
I1. Observations name immutable source binding IDs and source generations. Current
channel assignment never rewrites their origin. Missing context is unbound.
I2. Merge is commutative, associative and idempotent over accepted evidence;
conflicting histories stay explicit, including forks and issuer equivocation.
I3. Signatures bind bytes and keys, never physical truth. Every result carries
physical identity unverified and scope/freshness limitations.
I4. Acceptance requires explicit local issuer, key, role and scope policy. Unknown,
revoked, expired or unauthorized keys cannot contribute accepted associations.
I5. Cross-site links are scoped assertions, never global/transitive patient merges.
I6. Clinical content stays outside the runtime. Payload digests/references are
opaque. No global patient index, EHR UI, diagnosis, dosing, prioritization or care
access decisions. Missing evidence never instructs denial of human care.
I7. Inputs, evidence count and computation have explicit tested finite bounds.
I8. Every decision is reproducible from supplied evidence and explicit policy/time;
no clock/order dependent hidden tie breaking, last-write-wins or silent guessing.

## No-cheating contract
Never delete/weaken failing tests, move acceptance thresholds after results,
ignore flaky/skipped checks, or substitute mocks for claimed end-to-end evidence.
Never turn a compile/CI pass into production assurance. Never invent badges,
adoption, benchmarks, screenshots, recordings, integrations, security audits or
clinical benefit. Never hide negative controls, failures or unresolved findings.
Never label an association patient verified, globally current, or physically true.
Never use real patient data, credentials, hospital records or live medical devices.
Never publish machine paths, private mission notes or secrets. Record exact failed
commands and evidence scope. No sub-agents; no claim of independent authorship.
Tests may be corrected only for a documented specification error with retained
counterexample and a stronger replacement; the frozen obligations cannot relax.

## Evidence and change authority
The human owns mission and publication authorization. The implementation owner
may choose engineering details and execute authorized checks. No deployment to
clinical use is authorized. Freeze these governance files with SHA-256 before
source. Preserve their original bytes. Record strengthening amendments separately.
Protocol freeze follows primary-source research and an ADR. Any later protocol
change increments its version or documents wire-compatible stricter validation.
Failures advance engineering work, never marketing claims. A true external
blocker leaves the relevant gate blocked while independent work continues.
