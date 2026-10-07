#!/usr/bin/env python3
"""Use only committed files, execute README commands verbatim, test external API."""
import os
import re
import shutil
import subprocess
import tempfile
from pathlib import Path

root = Path(__file__).resolve().parents[1]
if subprocess.check_output(["git", "status", "--porcelain"], cwd=root).strip():
    raise SystemExit("clean-clone gate requires committed source")
with tempfile.TemporaryDirectory(prefix="carebind-clean-") as directory:
    clone = Path(directory) / "checkout"
    subprocess.run(["git", "clone", "--no-local", str(root), str(clone)], check=True, capture_output=True)
    env = dict(os.environ)
    env["GOCACHE"] = str(Path(directory) / "build-cache")
    env["GOMODCACHE"] = str(Path(directory) / "module-cache")
    commands = re.findall(r"```sh\n(.*?)```", (clone / "README.md").read_text(), re.S)
    assert len(commands) == 2, "review changed README quickstart blocks"
    for block in commands:
        subprocess.run(["/bin/sh", "-ec", block], cwd=clone, env=env, check=True)
    consumer = Path(directory) / "consumer"
    consumer.mkdir()
    (consumer / "go.mod").write_text(
        "module consumer.example\n\ngo 1.27.1\n\nrequire github.com/Nima0101/carebind v0.0.0\n"
        "replace github.com/Nima0101/carebind => " + str(clone) + "\n")
    shutil.copyfile(clone / "examples/consumer_test.go", consumer / "consumer_test.go")
    subprocess.run(["go", "test", "./..."], cwd=consumer, env=env, check=True)
    assert not subprocess.check_output(["git", "status", "--porcelain"], cwd=clone).strip()
print("PASS: clean git clone; README commands verbatim; isolated caches; external module consumer; clean worktree")
