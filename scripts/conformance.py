#!/usr/bin/env python3
"""Wire/semantic checks against committed vectors and OpenSSL Ed25519."""
import json
import os
import subprocess
import tempfile
from pathlib import Path
from reference import classify, canonical, event_id

ROOT = Path(__file__).resolve().parents[1]
BINARY = str(ROOT / "bin/carebind")
with tempfile.TemporaryDirectory() as directory:
    tmp = Path(directory)
    for path in sorted((ROOT / "testdata").glob("*.json")):
        b = json.loads(path.read_text())
        policy = tmp / "policy.json"
        policy.write_text(json.dumps(b["policy"]))
        actual = json.loads(subprocess.check_output([BINARY, "evaluate", str(path), str(policy)]))
        expected = classify(b, actual["accepted"])
        assert expected is None or actual["state"] == expected, (path.name, expected, actual)
        for env in b["events"]:
            event = env["event"]
            key = next((k for k in b["policy"]["keys"] if
                        (k["issuer"], k["key"]) == (event["issuer"], event["key"])), None)
            if key is None:
                continue
            # RFC 8410 SubjectPublicKeyInfo prefix for Ed25519.
            (tmp / "key.der").write_bytes(bytes.fromhex("302a300506032b6570032100" + key["public"]))
            (tmp / "sig.bin").write_bytes(bytes.fromhex(env["signature"]))
            (tmp / "message").write_bytes(b"CareBind/v1\n" + canonical(event))
            cmd = [os.environ.get("OPENSSL", "openssl"), "pkeyutl", "-verify", "-pubin", "-keyform", "DER",
                   "-inkey", str(tmp / "key.der"), "-rawin", "-in", str(tmp / "message"), "-sigfile", str(tmp / "sig.bin")]
            signature_valid = subprocess.run(cmd, capture_output=True).returncode == 0
            bad_signature = any(r["id"] == event_id(event) and r["reason"] == "bad_signature"
                                for r in actual["rejected"])
            assert signature_valid != bad_signature, (path.name, event["kind"])
            (tmp / "message").write_bytes(b"wrong-domain\n")
            assert subprocess.run(cmd, capture_output=True).returncode != 0
print("PASS: committed protocol vectors, separate semantic oracle, independent OpenSSL signature verification and negative control")
