package examples_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	cb "github.com/Nima0101/carebind"
	"strings"
	"testing"
)

// This consumer uses only exported API, independently of demo/internal helpers.
func TestPublicConsumer(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	binding := cb.Event{V: 1, Kind: "binding", Issuer: "example", Key: "key-1", Nonce: "bind-1", Scope: "site", Channel: "source:1", Generation: 1, Subject: "Cedar", Issued: 10}
	observation := cb.Event{V: 1, Kind: "observation", Issuer: "example", Key: "key-1", Nonce: "obs-1", Scope: "site", Channel: "source:1", Generation: 1, Binding: binding.ID(), Payload: strings.Repeat("0", 64), Issued: 11}
	first, err := cb.Sign(binding, priv)
	if err != nil {
		t.Fatal(err)
	}
	second, err := cb.Sign(observation, priv)
	if err != nil {
		t.Fatal(err)
	}
	bundle := cb.Bundle{V: 1, Events: []cb.Envelope{first, second}, Policy: cb.Policy{AsOf: 20, KnownAt: 20, MaxAge: 0, Keys: []cb.Key{{Issuer: "example", Key: "key-1", Public: hex.EncodeToString(pub), Scope: "site", Roles: []string{"binding", "observation"}, NotBefore: 0, NotAfter: 30}}}, Query: cb.Query{Observation: observation.ID(), TargetScope: "site", Target: "Cedar"}}
	result, err := cb.Evaluate(bundle)
	if err != nil {
		t.Fatal(err)
	}
	if result.State != "evidence-consistent-within-known-scope" || result.Origin != "Cedar" {
		t.Fatal(result)
	}
}
