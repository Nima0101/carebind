# Research landscape — 2026-10-08
CareBind is an association-evidence runtime and synthetic conformance kit, not a
new healthcare identity standard. Prior art substantially overlaps the problem.

| Primary source inspected | Existing capability | CareBind boundary / limitation |
|---|---|---|
| [IHE PCIM revision 2.4, May 2026](https://profiles.ihe.net/DEV/PCIM/index.html) | Device/patient association history, corrections, conflicts, human validation; retrospective queries are outside its current scope; security policy is site-specific | Strongest alternative. CareBind supplies a small offline replayable signed-history evaluator and adversarial adapter corpus. It does not implement PCIM messages or claim that PCIM implementations cannot do this. |
| [FHIR R5 DeviceAssociation](https://hl7.org/fhir/R5/deviceassociation.html) | Association resource with subject, device, status and period | A representation to map deliberately; it alone is not CareBind's history evaluation algorithm. No FHIR conformance claimed. |
| [FHIR R6 ballot](https://build.fhir.org/deviceassociation.html) | Evolving association model | Preview, not used as a stable dependency. |
| [FHIR Provenance](https://fhir.hl7.org/fhir/provenance.html) | Resource lineage and signatures | Complementary carrier; no invented gap in provenance. |
| [FHIR Specimen](https://fhir.hl7.org/fhir/specimen.html) | Specimen subject and collection context | Synthetic item association tests only; no specimen/custody standard implementation. |
| [IHE SDPi](https://profiles.ihe.net/DEV/SDPi/index.html) | SDC context validation and gateway mapping requirements | Source context/version semantics already exist. CareBind does not replace SDC or claim live-device compatibility. |
| [sdc11073](https://github.com/Draegerwerk/sdc11073), [PyPI](https://pypi.org/project/sdc11073/) | Concrete SDC implementation for testing and demonstration | Prefer it for actual SDC testing; CareBind tests a smaller portable evidence contract without device networking. |
| [W3C VC 2.0](https://www.w3.org/TR/vc-data-model-2.0/) | Issuer claims, validity and status | Authentic statements need verifier policy and do not establish physical truth. CareBind uses explicit local policy, no DID/global identity system. |
| [W3C PROV](https://www.w3.org/TR/prov-overview/) | General provenance vocabulary | Does not prescribe this evaluator; useful conceptual foundation. |
| [IHE PMIR](https://profiles.ihe.net/ITI/PMIR/index.html) | Master identity lifecycle | Deliberately excluded: CareBind never merges a patient index. |
| [RFC 8032](https://www.rfc-editor.org/rfc/rfc8032) | Ed25519 algorithm and test vectors | Use Go standard implementation, not custom cryptography. |
| [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785) | General JSON canonicalization | CareBind uses a narrower fixed-field ASCII wire encoding, not JCS compatibility. |

Search log: web and GitHub queries included “PCIM patient device”, “device patient
association”, “offline patient association”, “sdc11073”, and registry searches
“patient device association” on PyPI/npm. GitHub's first three returned no matching
repositories; sdc11073 returned the upstream plus an intrusion-detection project.
Registry search found sdc11073 and neighboring device-normalization packages.
Absence of search matches is not absence of competing products. Standards and
proprietary implementations can compose the same capabilities. No global novelty,
patent clearance, interoperability certification or superiority is claimed.

Decision: proceed with a materially useful executable contribution: source-bound
history carried intact through two different synthetic adapters, deterministic
partition reconciliation, a separate reference interpretation, and failure vectors
including physical misbinding. Reuse the established ideas openly. If real adapters
cannot capture origin at source, this contract cannot repair their guesses.

Public portfolio inspected via GitHub API: contingram, hidweave, persistscope,
profile and website. These address recovery planning, HID contracts and crash
consistency respectively; CareBind's association semantics are distinct. Name
carebind returned 404 under the authorized account at inspection; recheck before
publication. No private repository source was inspected.
