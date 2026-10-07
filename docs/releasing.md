# Release process
1. Complete V01–V15 on a clean candidate; record exact commands, revision and failures.
2. Run scripts/reproduce.py: two byte-identical archives, checksums, SPDX SBOM,
   source archive and local provenance. Inspect license inventory and public leak scan.
3. Recheck public portfolio/name. Only the authorized owner may create/push after
   all prepublication gates pass. Local success is not hosted success.
4. Hosted CI and CodeQL must pass on the exact main revision. Verify an anonymous
   clone, README quickstart, demo and scans. This completes V16.
5. Create version tag v0.1.0 at that verified revision and dispatch release workflow.
   It re-verifies source and refuses a tag/revision mismatch, builds packages,
   attaches hosted provenance and uploads immutable workflow artifacts. Download
   and verify them before creating a GitHub release with those exact assets.
6. V17 requires actual attestation validation and downloaded assets; only then may
   status describe bounded technical production readiness. No clinical suitability.

Versioning: wire v1 rejects unknown versions/fields. Breaking wire changes require
new version and vectors; release versions follow SemVer. Keep v1 compatibility
fixtures. Do not overwrite published tags/assets; ship a new patch for fixes.
Security fixes reopen affected gates. No external package registry publish is needed.
