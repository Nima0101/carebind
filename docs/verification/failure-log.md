# Failure and correction log

- First build: default Go cache was outside writable workspace. Use isolated
  temporary build/module caches in this session; no product change.
- First semantic suite: backwards-time fixture expected unknown-stale but also
  placed the observation after a declared successor. Frozen conflict precedence
  requires conflicting. Corrected assertion to the stronger state and retained
  backwards-time plus missing-parent controls. No production change.
- Explicit-unbound fixture reused the prior observation nonce with different
  bytes, correctly triggering issuer equivocation. Assign a fresh nonce to isolate
  unbinding; retain a separate observation-equivocation regression. No production
  change or gate relaxation.
- Adapter review found Python's default JSON loader could collapse duplicate keys
  before the strict Go parser saw the wire input. Added duplicate-key rejection at
  the receiver boundary and a regression; preserved the core parser's strict gate.
- Release-workflow review added explicit Actions read permission for its hosted-CI
  lookup. Without it, an otherwise authorized release job could not verify the gate.
- Release review bound the source archive to the captured revision and added a
  final source-stability check. A concurrent source change must invalidate a build,
  not yield a provenance record for mixed revisions. Reproduction is rerun only
  after source changes are committed.
