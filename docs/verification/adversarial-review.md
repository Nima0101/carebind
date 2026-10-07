# Adversarial review — same author, 2026-10-08
No independent human audit claimed. Reviewed the frozen specification against
parser, policy, history evaluation, CLI and consumer code after implementation.

- Tested issuer/key rotation IDs, nonce equivocation across events, generation
  forks, future/backwards timestamps, policy expiry and revoked historical keys.
- Verified unknown issuers and malformed/tampered evidence cannot silently supply
  accepted associations; rejected variants do not poison an otherwise valid event.
- Confirmed an accepted malicious issuer can force visible conflict within its
  authorized scope. It cannot be made truthful by signatures. Other source channels
  remain isolated; no current-roster fallback exists.
- Tested the consumer's numeric-slot mapping against missing evidence and a
  substituted destination. Embedded transport policy cannot override local policy.
- Audited parser bounds: finite byte/depth/key/array/token limits, strict struct
  fields, null/duplicate/unknown rejection, and typed API revalidation. No network,
  implicit clock, recursion through history or mutable global state.
- Merge tests compare a pairwise reference oracle over 1,200 generated histories,
  reorder three partitions and repeat envelopes. Python reference checks the
  accepted event set separately; OpenSSL challenges signature serialization.
- Wrong-tag test deliberately returns the same digital result for different
  fictional physical props. Physical identity remains unverified in every result.

Remaining boundaries: hidden evidence cannot be detected; real source adapters,
key provisioning, transport encryption/access control and physical human checks
are external. No latency, medical safety, standards conformance or user-study
claim. O(E²) worst-case evaluation is bounded by 512 input envelopes. Clinical
benefit and real-world usability remain unestablished and outside this release.
