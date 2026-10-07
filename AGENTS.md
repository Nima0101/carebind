# Working on CareBind
Read docs/engineering/constitution.md, docs/engineering/phase-model.md,
docs/verification/publication-gate.md, docs/verification/verification-matrix.md,
docs/architecture/protocol.md (once present), and SECURITY.md before changes.

Preserve origin bindings, generations, ambiguity, and physical-verification limits.
Never weaken frozen gates or tests for green. Add a regression before fixing a
semantic defect. Record protocol/security/privacy decisions in docs/decisions/.
Every claim needs an executable check and an explicit limit. Use synthetic data
only. No sub-agents. No hidden skips, fabricated evidence, or clinical decisions.
Run the documented verification command and report failures and unrun checks.
Changes to crypto, trust, parser, adapters or release logic need adversarial tests
and a threat-model update. Review as an adversary separately from implementation;
same-author review is not independent human assurance. No publication until the
prepublication gate passes. Production readiness also requires hosted evidence.
Constitution v1 is immutable; amendments may only strengthen its obligations.
