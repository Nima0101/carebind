#!/usr/bin/env python3
import hashlib
import subprocess
import tarfile
import tempfile
from pathlib import Path
root = Path(__file__).resolve().parents[1]

def build():
    subprocess.run(["python3", "scripts/release.py"], cwd=root, check=True)
    return {p.name: hashlib.sha256(p.read_bytes()).hexdigest() for p in (root / "dist").iterdir()}

first, second = build(), build()
assert first == second, "release is not byte reproducible"
archives = list((root / "dist").glob("carebind-0.1.0-*.tar.gz"))
for path in archives:
    if "source" in path.name:
        continue
    with tempfile.TemporaryDirectory() as directory:
        with tarfile.open(path) as archive:
            archive.extractall(directory, filter="data")
        binary = Path(directory) / "carebind"
        subprocess.run([str(binary), "demo"], check=True, capture_output=True)
print("PASS: two identical release builds; extracted packaged binary runs all demo assertions")
