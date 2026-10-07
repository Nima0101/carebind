#!/usr/bin/env python3
"""Deterministic local release candidate; never publishes or creates a tag."""
import gzip
import hashlib
import io
import json
import os
import platform
import subprocess
import tarfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
os.chdir(ROOT)
if subprocess.check_output(["git", "status", "--porcelain"]).strip():
    raise SystemExit("release requires a clean committed source tree")
revision = subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip()
go_version = subprocess.check_output(["go", "version"], text=True).strip()
goos = subprocess.check_output(["go", "env", "GOOS"], text=True).strip()
goarch = subprocess.check_output(["go", "env", "GOARCH"], text=True).strip()
out = ROOT / "dist"
out.mkdir(exist_ok=True)
binary = out / ("carebind.exe" if goos == "windows" else "carebind")
subprocess.run(["go", "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w -buildid=",
                "-o", str(binary), "./cmd/carebind"], check=True)

def digest(data):
    return hashlib.sha256(data).hexdigest()


def archive(name, members):
    raw = io.BytesIO()
    with tarfile.open(fileobj=raw, mode="w", format=tarfile.USTAR_FORMAT) as tar:
        for path, data, mode in sorted(members):
            info = tarfile.TarInfo(path)
            info.size, info.mode, info.mtime = len(data), mode, 0
            tar.addfile(info, io.BytesIO(data))
    (out / name).write_bytes(gzip.compress(raw.getvalue(), mtime=0))

archive("carebind-0.1.0-" + goos + "-" + goarch + ".tar.gz", [
    (binary.name, binary.read_bytes(), 0o755),
    ("LICENSE", (ROOT / "LICENSE").read_bytes(), 0o644),
    ("NOTICE", (ROOT / "NOTICE").read_bytes(), 0o644),
    ("GO-LICENSE", (ROOT / "docs/licenses/GO-LICENSE").read_bytes(), 0o644),
    ("README.md", (ROOT / "README.md").read_bytes(), 0o644),
])
source = subprocess.check_output(["git", "archive", "--format=tar", revision])
(out / "carebind-0.1.0-source.tar.gz").write_bytes(gzip.compress(source, mtime=0))
creation = {"creators": ["Tool: carebind-release-v1"], "created": "2026-10-08T00:00:00Z"}
sbom = {"spdxVersion": "SPDX-2.3", "dataLicense": "CC0-1.0", "SPDXID": "SPDXRef-DOCUMENT",
        "name": "CareBind-0.1.0-" + goos + "-" + goarch,
        "documentNamespace": "https://github.com/Nima0101/carebind/sbom/" + revision + "/" + goos + "/" + goarch,
        "creationInfo": creation,
        "packages": [
            {"name": "carebind", "SPDXID": "SPDXRef-CareBind", "versionInfo": "0.1.0",
             "downloadLocation": "NOASSERTION", "filesAnalyzed": False, "licenseConcluded": "MIT",
             "licenseDeclared": "MIT", "copyrightText": "Copyright 2026 Nima Khaki"},
            {"name": "Go-standard-library-and-runtime", "SPDXID": "SPDXRef-Go",
             "versionInfo": go_version.split()[2], "downloadLocation": "https://go.dev/dl/",
             "filesAnalyzed": False, "licenseConcluded": "BSD-3-Clause", "licenseDeclared": "BSD-3-Clause",
             "copyrightText": "Copyright The Go Authors"}],
        "relationships": [{"spdxElementId": "SPDXRef-DOCUMENT", "relationshipType": "DESCRIBES", "relatedSpdxElement": "SPDXRef-CareBind"},
                          {"spdxElementId": "SPDXRef-CareBind", "relationshipType": "DEPENDS_ON", "relatedSpdxElement": "SPDXRef-Go"}]}
(out / "sbom.spdx.json").write_text(json.dumps(sbom, sort_keys=True, indent=2) + "\n")
artifacts = sorted(p for p in out.iterdir() if p.name not in ("SHA256SUMS", "provenance.json"))
provenance = {"source_revision": revision, "toolchain": go_version,
              "command": "go build -trimpath -buildvcs=false -ldflags='-s -w -buildid=' ./cmd/carebind",
              "claim": "local unsigned build record; hosted attestation is a separate gate",
              "artifacts": {p.name: digest(p.read_bytes()) for p in artifacts}}
(out / "provenance.json").write_text(json.dumps(provenance, sort_keys=True, indent=2) + "\n")
(out / "SHA256SUMS").write_text("".join(digest(p.read_bytes()) + "  " + p.name + "\n"
                                       for p in sorted(out.iterdir()) if p.name != "SHA256SUMS"))
if subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip() != revision or subprocess.check_output(["git", "status", "--porcelain"]).strip():
    raise SystemExit("source changed during release; discard candidate artifacts")
print("Built deterministic candidate at source " + revision)
