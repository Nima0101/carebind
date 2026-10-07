#!/usr/bin/env python3
"""Fast, fail-closed local gate. Full release checks are separately explicit."""
import hashlib
import json
import re
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def run(*cmd):
    print("+ " + " ".join(cmd), flush=True)
    subprocess.run(cmd, cwd=ROOT, check=True)


for name in ("governance-freeze.json", "protocol-freeze.json"):
    manifest = json.loads((ROOT / "docs/verification" / name).read_text())
    for path, digest in manifest["sha256"].items():
        assert hashlib.sha256((ROOT / path).read_bytes()).hexdigest() == digest, path
print("PASS: pre-source governance and protocol freezes", flush=True)
files = subprocess.check_output(["git", "ls-files", "--cached", "--others", "--exclude-standard"], cwd=ROOT, text=True).splitlines()
for name in files:
    path = ROOT / name
    if not path.is_file() or path.suffix in (".gif", ".png"):
        continue
    data = path.read_text(errors="strict")
    for pattern in (r"/Users/", r"/home/[A-Za-z]", r"gh[pousr]_" + r"[A-Za-z0-9]{20,}",
                    r"-----BEGIN " + r"(?:RSA |OPENSSH |EC )?PRIVATE KEY-----",
                    r"AKIA" + r"[A-Z0-9]{16}"):
        # The scanner contains patterns, not leaked paths or secrets.
        if name != "scripts/verify.py":
            assert not re.search(pattern, data), (name, "public-content scan")
    assert not name.startswith(".mission/"), name
    if path.suffix == ".go":
        assert not re.search(r"\bt\.Skip|\bf\.Skip|TODO|FIXME", data), (name, "hidden skip/unfinished code")
    if path.suffix == ".md":
        for link in re.findall(r"\]\(([^)]+)\)", data):
            if "://" in link or link.startswith("#") or link.startswith("mailto:"):
                continue
            target = (path.parent / link.split("#")[0]).resolve()
            assert target.exists(), (name, link)
run("go", "test", "./...")
run("go", "vet", "./...")
go_files = [x for x in files if x.endswith(".go")]
assert not subprocess.check_output(["gofmt", "-l", *go_files], cwd=ROOT).strip(), "gofmt"
run("go", "build", "-trimpath", "-o", "bin/carebind", "./cmd/carebind")
run("python3", "scripts/integration.py")
run("python3", "scripts/conformance.py")
run("bin/carebind", "demo")
print("PASS: fast verification; race, 60s fuzz, vulnerability audit and clean release remain separate gates")
