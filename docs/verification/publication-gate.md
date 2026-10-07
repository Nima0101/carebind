# Publication and production gate v1
Before creating/pushing any public repository: V01–V15 all pass on the candidate
revision, current public-name/portfolio check passes, exact source has no private
content, and status links each gate to evidence. Failures block publication.
Owner authorization permits public creation/push only after these checks.

After public push: V16 must pass before a release tag; V17 must pass before calling
the bounded runtime production-ready. This sequencing avoids claiming hosted CI
before a remote exists. Remote failure keeps state public-unverified; fix it with
full local checks. No clinical production readiness, safety certification or
patient verification claim is authorized at any state.

The complete gate includes strict protocol/parser bounds and compatibility,
crypto/key lifecycle, trust scoping, replay/conflicts, reconnection, privacy,
adversarial/property/fuzz/reference tests, dependency audit, clean consumers,
reproducible release/SBOM/provenance, actual CI, accurate docs and no critical-path
unfinished work. External blockers must name the required action and exact failed
check; continue every independent task. Governance cannot be relaxed to publish.
