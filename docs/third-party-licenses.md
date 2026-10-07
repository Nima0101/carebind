# Third-party licenses
Production Go module dependencies: none. Compiled binaries include the Go standard
library/runtime; its BSD 3-Clause text is in [GO-LICENSE](licenses/GO-LICENSE) and
binary archives. Go's bundled standard-library vendor components are distributed
under the Go project's BSD license. Record exact toolchain in each SBOM.

Python reference/adapter checks use the standard library. OpenSSL is an externally
installed test tool, not redistributed. Pillow 11.3.0 (HPND) is an optional external
recording renderer, not a runtime or redistributed library. GitHub Actions and
govulncheck are build/security tools, pinned in workflows, not product dependencies.
The original mark is code-native SVG authored for CareBind. Demo assets are actual
synthetic command outputs rendered by the recording script; no external artwork.
