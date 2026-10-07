#!/usr/bin/env python3
"""Synthetic receiving adapter with numeric local slots. No clinical data.

Transport preserves the complete signed evidence. Local policy and slot expectation
are separately supplied by the consumer, never inferred from a current roster.
"""
import json
import subprocess
import sys
import tempfile
from pathlib import Path


def strict_json(data):
    def unique(pairs):
        result = {}
        for key, value in pairs:
            if key in result:
                raise ValueError("duplicate_field")
            result[key] = value
        return result
    return json.loads(data, object_pairs_hook=unique,
                      parse_constant=lambda value: (_ for _ in ()).throw(ValueError("invalid_number")))


def receive(binary, bundle, policy, slot):
    if type(slot) is not int or not 1 <= slot <= 999999:
        raise ValueError("invalid local slot")
    bundle = json.loads(json.dumps(bundle))
    bundle["query"]["target_scope"] = "clinic"
    bundle["query"]["target"] = "slot-" + str(slot)
    with tempfile.TemporaryDirectory() as directory:
        path = Path(directory) / "transport.json"
        path.write_text(json.dumps(bundle), encoding="utf-8")
        result = subprocess.run(
            [binary, "evaluate", str(path), str(policy)],
            check=True, capture_output=True, text=True,
        )
    report = json.loads(result.stdout)
    # Local storage never assigns a slot on an unresolved result. Historical
    # supersession is retained explicitly, never promoted to current truth.
    accepted = report["state"] in (
        "evidence-consistent-within-known-scope", "superseded-in-known-history"
    )
    return {"local_slot": slot if accepted else None,
            "origin_binding": report["binding"], "evidence": report}


if __name__ == "__main__":
    if len(sys.argv) != 5:
        raise SystemExit("usage: receive.py BINARY BUNDLE LOCAL_POLICY SLOT")
    source = Path(sys.argv[2])
    with source.open("rb") as stream:
        data = stream.read(1048577)
    if len(data) > 1048576:
        raise SystemExit("input_limit")
    print(json.dumps(receive(sys.argv[1], strict_json(data), sys.argv[3], int(sys.argv[4]))))
