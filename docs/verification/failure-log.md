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
- Final parser audit found the aggregate input cap did not enforce the frozen
  4,096-byte per-event wire cap when an event contained excess whitespace.
  TestPerEventWireLimit failed before the fix. Parse now checks raw event length
  before normalization. This is a production fix: fuzz, race, clean clone and
  release evidence are rerun; earlier green does not transfer to changed bytes.
- The end-to-end oversized-event regression then found the receiver could erase
  excessive wire whitespace during JSON normalization. The native parser rejected
  it but the receiver initially exited 0. The receiver now validates an immutable
  copy of the original bytes before decoding/adapting. Actual process regressions
  reject oversized, duplicate-field, mis-cased and overly nested wire input.
