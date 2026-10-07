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
