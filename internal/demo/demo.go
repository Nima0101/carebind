// Package demo contains public, synthetic keys and nonmedical fixtures only.
package demo

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	cb "github.com/Nima0101/carebind"
)

func Private(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("PUBLIC SYNTHETIC FIXTURE ONLY: " + label))
	return ed25519.NewKeyFromSeed(seed[:])
}
func Key(issuer, scope string, roles ...string) cb.Key {
	p := Private(issuer).Public().(ed25519.PublicKey)
	return cb.Key{Issuer: issuer, Key: "demo-v1", Public: hex.EncodeToString(p), Scope: scope, Roles: roles, NotBefore: 0, NotAfter: 2000}
}
func Sign(e cb.Event) cb.Envelope {
	env, err := cb.Sign(e, Private(e.Issuer))
	if err != nil {
		panic(err)
	}
	return env
}
func Binding(nonce, subject string, generation int64, previous string) cb.Envelope {
	return Sign(cb.Event{V: 1, Kind: "binding", Issuer: "cedar-source", Key: "demo-v1", Nonce: nonce, Scope: "field", Channel: "emulator-1:channel-1", Generation: generation, Subject: subject, Previous: previous, Issued: 100 + generation*10})
}
func Observation(binding cb.Envelope) cb.Envelope {
	e := binding.Event
	payload := sha256.Sum256([]byte("SYNTHETIC opaque token; no clinical value"))
	return Sign(cb.Event{V: 1, Kind: "observation", Issuer: "cedar-source", Key: "demo-v1", Nonce: "observation-1", Scope: e.Scope, Channel: e.Channel, Generation: e.Generation, Binding: e.ID(), Payload: hex.EncodeToString(payload[:]), Issued: e.Issued + 1})
}
func Base() cb.Bundle {
	b := Binding("binding-cedar", "Cedar", 1, "")
	o := Observation(b)
	return cb.Bundle{V: 1, Events: []cb.Envelope{b, o}, Policy: cb.Policy{AsOf: 1000, KnownAt: 990, MaxAge: 60, Keys: []cb.Key{Key("cedar-source", "field", "binding", "observation"), Key("receiving-site", "clinic", "link")}}, Query: cb.Query{Observation: o.Event.ID()}}
}
func Link(b cb.Envelope, target string) cb.Envelope {
	return Sign(cb.Event{V: 1, Kind: "link", Issuer: "receiving-site", Key: "demo-v1", Nonce: "link-" + target, Scope: b.Event.Scope, Subject: b.Event.Subject, Binding: b.Event.ID(), TargetScope: "clinic", Target: target, Issued: 150})
}

type Case struct {
	Name   string
	Bundle cb.Bundle
	Want   string
	Note   string
}

func Cases() []Case {
	base := Base()
	cedar := base.Events[0]
	birch := Binding("binding-birch", "Birch", 2, cedar.Event.ID())
	delayed := Base()
	delayed.Events = append(delayed.Events, birch)
	linked := Base()
	linked.Events = append(linked.Events, birch, Link(cedar, "slot-42"))
	linked.Query.TargetScope = "clinic"
	linked.Query.Target = "slot-42"
	conflict := Base()
	conflict.Events = append(conflict.Events, birch, Binding("offline-fork", "Maple", 2, cedar.Event.ID()))
	unknown := Base()
	unknown.Policy.Keys = unknown.Policy.Keys[1:]
	tampered := Base()
	tampered.Events[1].Event.Payload = fmt.Sprintf("%064d", 0)
	tampered.Query.Observation = tampered.Events[1].Event.ID()
	replay := Base()
	replay.Events = append(replay.Events, replay.Events...)
	stale := Base()
	stale.Policy.KnownAt = 0
	strip := Base()
	strip.Events = strip.Events[1:]
	substitute := Base()
	substitute.Query.TargetScope = "field"
	substitute.Query.Target = "Birch"
	revoked := Base()
	revoked.Policy.Keys[0].Revoked = true
	item := Base()
	ib := item.Events[0].Event
	ib.Channel = "synthetic-item:sample-1"
	item.Events[0] = Sign(ib)
	item.Events[1] = Observation(item.Events[0])
	item.Query.Observation = item.Events[1].Event.ID()
	return []Case{
		{"source-capture", base, "evidence-consistent-within-known-scope", "Cedar source context captured before delivery"},
		{"delayed-after-reassignment", delayed, "superseded-in-known-history", "Naive current roster says Birch; origin remains Cedar"},
		{"cross-site-link", linked, "superseded-in-known-history", "Explicit receiving handle slot-42; no global merge"},
		{"offline-reconnection", conflict, "conflicting", "Partition histories union exposes Birch / Maple fork"},
		{"unknown-issuer", unknown, "evidence-invalid", "Unaccepted issuer gives no accepted association"},
		{"tampered-envelope", tampered, "evidence-invalid", "Payload digest modified after signing"},
		{"replay-duplicate", replay, "evidence-consistent-within-known-scope", "Duplicate count 2; accepted evidence count stays 2"},
		{"stale-policy", stale, "unknown-stale", "Offline policy age is explicit; global freshness unknown"},
		{"stripped-context", strip, "unbound", "Missing binding cannot be replaced by a current roster"},
		{"substituted-destination", substitute, "conflicting", "Consumer context differs from signed origin"},
		{"revoked-key", revoked, "evidence-invalid", "Local revocation invalidates historical uses too"},
		{"synthetic-item", item, "evidence-consistent-within-known-scope", "Opaque synthetic item association; no specimen interpretation"},
		{"wrong-physical-tag", base, "evidence-consistent-within-known-scope", "MODEL: Cedar tag attached to Birch prop; software CANNOT detect it"},
	}
}
