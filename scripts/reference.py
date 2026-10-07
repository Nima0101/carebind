#!/usr/bin/env python3
"""Independent semantic interpretation of *preauthenticated* events.

This is a reference test oracle, not a cryptographic verifier or production API.
Its test driver supplies accepted IDs from the Go verifier, making the shared
crypto boundary explicit. Independent OpenSSL verification covers wire signatures.
"""
import hashlib
import json


FIELDS = "v kind issuer key nonce scope channel generation subject previous binding payload target_scope target issued".split()


def canonical(e):
    return json.dumps({key: e[key] for key in FIELDS}, separators=(",", ":")).encode()


def event_id(e):
    return hashlib.sha256(canonical(e)).hexdigest()


def classify(bundle, accepted):
    evidence = {event_id(x["event"]): x["event"] for x in bundle["events"]
                if event_id(x["event"]) in accepted}
    q, p = bundle["query"], bundle["policy"]
    observation = evidence.get(q["observation"])
    if observation is None:
        return None  # Trust rejection is outside this oracle.
    if observation["kind"] != "observation":
        return "evidence-invalid"
    origin = evidence.get(observation["binding"])
    if origin is None:
        return "unbound"
    if origin["kind"] != "binding" or any(origin[k] != observation[k] for k in
        ("scope", "channel", "generation")) or observation["issued"] < origin["issued"]:
        return "evidence-invalid"
    values = list(evidence.values())
    def equiv(e):
        return sum(x["issuer"] == e["issuer"] and x["nonce"] == e["nonce"] for x in values) > 1
    history = [e for e in values if e["kind"] == "binding" and
               (e["scope"], e["channel"]) == (origin["scope"], origin["channel"])]
    if equiv(observation) or any(equiv(e) for e in history):
        return "conflicting"
    if len({e["generation"] for e in history}) != len(history):
        return "conflicting"
    if any(e["previous"] == observation["binding"] and
           observation["issued"] >= e["issued"] for e in history):
        return "conflicting"
    for e in history:
        if e["generation"] == 1:
            continue
        parents = [x for x in history if event_id(x) == e["previous"] and
                   x["generation"] + 1 == e["generation"] and x["issued"] <= e["issued"]]
        if not parents:
            return "unknown-stale"
    if not origin["subject"]:
        return "unbound"
    if q["target_scope"]:
        if q["target_scope"] == origin["scope"]:
            if q["target"] != origin["subject"]:
                return "conflicting"
        else:
            links = [e for e in values if e["kind"] == "link" and
                     e["binding"] == observation["binding"] and e["target_scope"] == q["target_scope"]]
            targets = {e["target"] for e in links}
            if len(targets) > 1 or any(equiv(e) for e in links):
                return "conflicting"
            if any(e["scope"] != origin["scope"] or e["subject"] != origin["subject"] or
                   e["issued"] < origin["issued"] for e in links):
                return "evidence-invalid"
            if not links:
                return "unbound"
            if q["target"] not in targets:
                return "conflicting"
    if p["as_of"] - p["known_at"] > p["max_age"]:
        return "unknown-stale"
    if any(e["generation"] > origin["generation"] for e in history):
        return "superseded-in-known-history"
    return "evidence-consistent-within-known-scope"
