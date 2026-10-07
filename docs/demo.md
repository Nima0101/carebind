# Reproduce the demo
Run `go build -trimpath -o bin/carebind ./cmd/carebind` then `bin/carebind demo`.
Each printed result follows a real evaluator call and expected-state assertion.
Fixtures are emitted by `bin/carebind fixture NAME`; names appear in the demo.
Committed wire vectors live in testdata/. No patient/clinical data is included.

For PNG/GIF and asciicast capture, install Pillow 11.3.0 in an isolated Python
virtual environment, then run `python3 scripts/record_demo.py`. Optionally set
CAREBIND_RECORD_FONT to a local monospaced font. Rendering is code-generated
terminal replay of captured stdout; timing is slowed to 3.5 seconds per case.
Actual process timing is retained in demo.cast. recording.json binds source,
binary and stdout digests. Verify the text by rerunning the command.

Wrong-tag case models moving an unchanged signed Cedar carrier to a fictional
Birch prop. The unchanged software input necessarily produces the same digital
result. This is an explicit software blind spot, not a performed physical trial.
A printable carrier is available in docs/assets/carrier.html; the full signed
bundle must accompany it. A printed digest alone is not accepted evidence.

Presentation reference inspected 2026-10-08: Flutter's current public repository.
Adopted concise product definition, direct documentation/install paths and visible
contributor/security links. No Flutter branding, claims, wording or artwork copied.
