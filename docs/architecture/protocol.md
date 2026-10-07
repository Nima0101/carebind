# CareBind wire protocol 1
Status: frozen v1 before implementation; see protocol-freeze.json. Library release
version and wire version are separate. Unsupported versions fail closed. Additive
fields require a new wire version because v1 rejects unknown fields. Existing v1
vectors remain regression fixtures. No standards conformance claim is made.

## Signed event
Canonical event is compact UTF-8 JSON, exactly these fields in this order, with
no omitted fields: v, kind, issuer, key, nonce, scope, channel, generation, subject,
previous, binding, payload, target_scope, target, issued. v=1. kind is binding,
observation or link. All names/handles use ASCII [A-Za-z0-9_.:-], length 1–64.
Empty values are allowed only for fields irrelevant to the kind. Digests use
64 lowercase hex characters. Integers are decimal without exponent, 0–2147483647
for generation, 0–4102444800 for issued (Unix seconds). generation is at least 1
for binding/observation and 0 for link. nonce is issuer-unique across key rotation.
A binding subject is a local encounter handle or empty for an explicit unbound
channel generation. Channels include a source-instance namespace; resets require
a new channel identifier, never reusing generation 1 with new meaning.

binding: scope/channel/generation/subject; previous is empty at generation 1,
otherwise the prior binding digest. binding/payload/target_scope/target empty.
observation: scope/channel/generation plus binding digest and opaque SHA-256 payload
digest; subject/previous/target_scope/target empty. No clinical payload in protocol.
link: binding digest, source scope/subject, target_scope/target; channel/previous/
payload empty. It maps only this binding into the stated target scope. Not
transitive, not a permanent person identity, and cannot rewrite the origin.

Event ID = SHA-256 of canonical event bytes. Signature = Ed25519 over the ASCII
domain separator `CareBind/v1` followed by one newline then canonical event bytes.
Envelope fields: event (object), signature (128 lowercase hex characters). Verify
canonical form after strict decoding. No algorithms or keys are taken from remote
URLs. Synthetic demo keys are public fixtures and must never be operational keys.

## Input bundle and local policy
Bundle fields: v=1, events (array), policy, query. Policy is a separately trusted
local file/configuration at the CLI boundary; do not use a sender's trust policy.
CLI accepts bundle AND an explicit separate policy file and replaces bundle policy.
Library caller supplies trusted policy. Embedded policy is for reproducible test
vectors only. Policy fields: as_of, known_at, max_age, keys. Integers 0–4102444800;
known_at <= as_of; max_age <= 86400. No implicit wall clock.
Key fields: issuer, key, public (64 hex), scope, roles (unique kinds array),
not_before, not_after, revoked. One issuer/key pair per policy. Link keys need
scope equal target_scope; other events need scope equal event scope. Roles are
explicit; no wildcard. Keys must be active at as_of and at issued. Revocation
invalidates all uses in this evaluation, including historical ones; no trusted
signing timestamp exists. Future-issued events fail. Policy older than max_age
makes otherwise accepted results unknown-stale. Global freshness always unknown.
Query fields: observation (digest), target_scope and target (both empty or both
valid tokens). These are consumer expectations, not sender-provided identity.

## Evaluation
Reject malformed bundles atomically with typed error. Parse all event shapes
strictly even if their issuer is unknown. Well-formed but untrusted or bad-signature
events are quarantined with ID/reason, never counted in the accepted set. Exact
duplicates add no evidence; expose duplicate count separately. Decisions depend
only on accepted unique event IDs plus explicit policy/query. Invalid evidence
cannot poison unrelated accepted history. Querying a rejected ID returns
evidence-invalid. Missing observation or binding returns unbound.

Issuer+nonce collisions with different accepted IDs are equivocation: any selected
observation, chain binding or selected link using that nonce is conflicting.
Find the observation's exact referenced binding; scope/channel/generation must
match. Preserve that binding's subject forever; never consult a current roster.
For all accepted bindings of that channel/scope: multiple IDs at a generation
are conflicting, even if their subjects match. Missing parent, mismatched parent
channel/scope/generation or decreasing issued time makes history unknown-stale.
Walk all relevant chains; no arbitrary tie-break. Higher valid generation makes
an older binding superseded-in-known-history; its origin remains unchanged.
A referenced binding with empty subject returns unbound. Observation issued must
be >= binding issued; observations at/after a known successor's issued time are
conflicting (half-open association interval). Delayed delivery is not new issuance.

A target in the origin scope must equal origin subject. Cross-scope target needs
one distinct matching link target for this exact binding and requested scope;
missing is unbound, multiple targets are conflicting. Links with wrong source
scope/subject are evidence-invalid for that query. Link issued before binding is
invalid. No reverse/transitive closure. Repeated equivalent links add no confidence.

Precedence after observation/binding validation: conflict, incomplete history,
unbound subject/link, stale policy, superseded, evidence-consistent-within-known-
scope. Cross-scope invalid links can make result evidence-invalid. Result includes
origin, binding, local target, reason, policy knowledge time, global freshness
unknown, physical identity unverified, accepted IDs and quarantine. Rejected or
unbound results never direct denial of care. Equivocation is local evidence only.

## Hard limits
Input <= 1 MiB, depth <= 8, every object <= 32 keys, every array <= 1024 elements,
events <= 512, keys <= 64, string tokens <= 1024 bytes, event <= 4096 bytes. Reject
duplicate/unknown/mis-cased fields, null, invalid UTF-8, trailing values, negative/
overflow/floating integers and invalid shape. Bounded O(E² + E*K) evaluation is
acceptable at E<=512. No recursion over graph depth. No hidden file/network access.
