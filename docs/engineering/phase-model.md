# Phase and state model v1
P0 authorized -> P1 governance-frozen -> P2 researched -> P3 protocol-frozen ->
P4 implemented -> P5 adversarially-verified -> P6 clean-clone-verified ->
P7 publication-candidate -> P8 public-verified -> P9 production-ready (bounded
technical runtime only, never clinical suitability).

Transitions require evidence in docs/verification/status.md. A phase may be
blocked without stopping independent work. Regression reopens affected gates;
previous green does not transfer to changed bytes. Pending, failed, blocked,
passed and not-applicable are distinct. N/A needs a scope-based rationale made
before the check, never a workaround for failure. P7 requires all prepublication
gates; P8 requires actual remote checks/public clone; P9 requires both plus all
production gates. No release tag before hosted checks. A local candidate is not
a public release. Record source revision, tools, commands, outcomes and limits.
