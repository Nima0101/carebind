#!/usr/bin/env python3
"""Real process/API consumer tests: no mocked evaluator."""
import importlib.util
import json
import subprocess
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
BINARY = str(ROOT / "bin/carebind")
spec = importlib.util.spec_from_file_location("receiver", ROOT / "examples/receiver/receive.py")
receiver = importlib.util.module_from_spec(spec)
spec.loader.exec_module(receiver)


def call(*args):
    return json.loads(subprocess.check_output([BINARY, *args], text=True))


with tempfile.TemporaryDirectory() as directory:
    policy = Path(directory) / "policy.json"
    policy.write_text(json.dumps(call("demo-policy")))
    source = call("fixture", "cross-site-link")
    original = json.dumps(source, sort_keys=True)
    report = receiver.receive(BINARY, source, policy, 42)
    assert report["local_slot"] == 42
    assert report["evidence"]["origin"] == "Cedar"
    assert report["evidence"]["state"] == "superseded-in-known-history"
    assert json.dumps(source, sort_keys=True) == original
    assert receiver.receive(BINARY, source, policy, 43)["local_slot"] is None
    stripped = json.loads(original)
    stripped["events"] = [e for e in stripped["events"] if e["event"]["kind"] != "binding"]
    assert receiver.receive(BINARY, stripped, policy, 42)["evidence"]["state"] == "unbound"
    # Exercise the actual receiver process, including its original-wire boundary.
    original_wire = subprocess.check_output([BINARY, "fixture", "cross-site-link"])
    wire_path = Path(directory) / "wire.json"
    command = ["python3", str(ROOT / "examples/receiver/receive.py"), BINARY,
               str(wire_path), str(policy), "42"]
    wire_path.write_bytes(original_wire)
    received = json.loads(subprocess.check_output(command))
    assert received["local_slot"] == 42
    for malformed in (
        original_wire.replace(b'"event":{', b'"event":{' + b" " * 4096, 1),
        original_wire.replace(b'"v":1', b'"v":1,"v":1', 1),
        original_wire.replace(b'"v":1', b'"V":1', 1),
        b"[" * 100 + b"]" * 100,
    ):
        wire_path.write_bytes(malformed)
        rejected = subprocess.run(command, capture_output=True)
        assert rejected.returncode == 2, "receiver normalized invalid wire evidence"
    # Transport cannot grant itself authority by changing embedded policy.
    untrusted = json.loads(original)
    untrusted["policy"]["keys"][0]["revoked"] = False
    local = call("demo-policy")
    local["keys"][0]["revoked"] = True
    policy.write_text(json.dumps(local))
    assert receiver.receive(BINARY, untrusted, policy, 42)["local_slot"] is None
try:
    receiver.strict_json('{"query":1,"query":2}')
except ValueError:
    pass
else:
    raise AssertionError("receiver collapsed duplicate JSON fields")
print("PASS: Go source -> signed transport -> Python numeric-slot consumer; stripping, substitution, local-policy override")
