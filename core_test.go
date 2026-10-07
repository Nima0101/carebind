package carebind_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cb "github.com/Nima0101/carebind"
	"github.com/Nima0101/carebind/internal/demo"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func result(t *testing.T, b cb.Bundle) cb.Result {
	t.Helper()
	r, err := cb.Evaluate(b)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func state(t *testing.T, b cb.Bundle, want string) {
	t.Helper()
	r := result(t, b)
	if r.State != want {
		t.Fatalf("want %s got %+v", want, r)
	}
}
func TestFlagshipCases(t *testing.T) {
	for _, c := range demo.Cases() {
		t.Run(c.Name, func(t *testing.T) {
			r := result(t, c.Bundle)
			if r.State != c.Want {
				t.Fatal(r)
			}
			if r.PhysicalIdentity != "unverified" || r.GlobalFreshness != "unknown" {
				t.Fatal("unsafe confidence")
			}
			if c.Name == "delayed-after-reassignment" && r.Origin != "Cedar" {
				t.Fatal("rebound history")
			}
		})
	}
}
func TestRFC8032(t *testing.T) {
	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	pub := decode("d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a")
	sig := decode("e5564300c360ac729086e2cc806e828a84877f1eb8e5d974d873e065224901555fb8821590a33bacc61e39701cf9b46bd25bf5f0595bbe24655141438e7a100b")
	if !ed25519.Verify(pub, nil, sig) {
		t.Fatal("RFC 8032 test 1")
	}
	sig[0] ^= 1
	if ed25519.Verify(pub, nil, sig) {
		t.Fatal("tamper")
	}
}
func TestTrustControls(t *testing.T) {
	for _, name := range []string{"wrong-key", "role", "scope", "expired", "future", "revoked", "rotation"} {
		t.Run(name, func(t *testing.T) {
			b := demo.Base()
			switch name {
			case "wrong-key":
				b.Policy.Keys[0].Public = b.Policy.Keys[1].Public
			case "role":
				b.Policy.Keys[0].Roles = []string{"link"}
			case "scope":
				b.Policy.Keys[0].Scope = "elsewhere"
			case "expired":
				b.Policy.Keys[0].NotAfter = 999
			case "future":
				b.Policy.AsOf = 100
				b.Policy.KnownAt = 100
			case "revoked":
				b.Policy.Keys[0].Revoked = true
			case "rotation":
				b.Policy.Keys[0].Key = "new-key"
			}
			state(t, b, "evidence-invalid")
		})
	}
	b := demo.Base()
	b.Policy.Keys = append(b.Policy.Keys, b.Policy.Keys[0])
	if _, e := cb.Evaluate(b); e == nil {
		t.Fatal("ambiguous policy accepted")
	}
}
func TestTamperEveryField(t *testing.T) {
	e := demo.Base().Events[1]
	v := reflect.ValueOf(&e.Event).Elem()
	for i := 0; i < v.NumField(); i++ {
		t.Run(v.Type().Field(i).Name, func(t *testing.T) {
			b := demo.Base()
			x := reflect.ValueOf(&b.Events[1].Event).Elem().Field(i)
			if x.Kind() == reflect.String {
				x.SetString(x.String() + "a")
			} else {
				x.SetInt(x.Int() + 1)
			}
			b.Query.Observation = b.Events[1].Event.ID()
			r, err := cb.Evaluate(b)
			if err == nil && r.State != "evidence-invalid" {
				t.Fatal(r)
			}
		})
	}
}
func TestHistories(t *testing.T) {
	for _, name := range []string{"gap", "parent-channel", "parent-generation", "backwards-time", "equivocation", "late-source", "clear", "mismatch", "missing", "query-kind"} {
		t.Run(name, func(t *testing.T) {
			b := demo.Base()
			want := "unknown-stale"
			e := demo.Binding("next", "Birch", 2, b.Events[0].Event.ID()).Event
			switch name {
			case "gap":
				e.Previous = strings.Repeat("0", 64)
			case "parent-channel":
				p := b.Events[0].Event
				p.Channel = "other"
				b.Events[0] = demo.Sign(p)
			case "parent-generation":
				e.Generation = 3
			case "backwards-time":
				e.Issued = 1
				want = "conflicting"
			case "equivocation":
				e.Nonce = b.Events[0].Event.Nonce
				want = "conflicting"
			case "late-source":
				o := b.Events[1].Event
				o.Issued = e.Issued
				b.Events[1] = demo.Sign(o)
				b.Query.Observation = o.ID()
				want = "conflicting"
			case "clear":
				e.Subject = ""
				n := demo.Sign(e)
				clearObservation := demo.Observation(n).Event
				clearObservation.Nonce = "observation-clear"
				b.Events = append(b.Events, n, demo.Sign(clearObservation))
				b.Query.Observation = b.Events[len(b.Events)-1].Event.ID()
				want = "unbound"
			case "mismatch":
				o := b.Events[1].Event
				o.Channel = "other"
				b.Events[1] = demo.Sign(o)
				b.Query.Observation = o.ID()
				want = "evidence-invalid"
			case "missing":
				b.Events = b.Events[1:]
				want = "unbound"
			case "query-kind":
				b.Query.Observation = b.Events[0].Event.ID()
				want = "evidence-invalid"
			}
			b.Events = append(b.Events, demo.Sign(e))
			if name == "parent-channel" {
				want = "unbound"
			}
			state(t, b, want)
		})
	}
}
func TestLinks(t *testing.T) {
	for _, name := range []string{"match", "missing", "fork", "wrong-origin", "equivocation", "early", "no-transitivity"} {
		t.Run(name, func(t *testing.T) {
			b := demo.Base()
			b.Query.TargetScope = "clinic"
			b.Query.Target = "slot-42"
			link := demo.Link(b.Events[0], "slot-42")
			want := "evidence-consistent-within-known-scope"
			switch name {
			case "missing":
				state(t, b, "unbound")
				return
			case "fork":
				b.Events = append(b.Events, demo.Link(b.Events[0], "slot-43"))
				want = "conflicting"
			case "wrong-origin":
				link.Event.Subject = "Birch"
				link = demo.Sign(link.Event)
				want = "evidence-invalid"
			case "equivocation":
				e := link.Event
				e.Target = "slot-43"
				b.Events = append(b.Events, demo.Sign(e))
				want = "conflicting"
			case "early":
				link.Event.Issued = 1
				link = demo.Sign(link.Event)
				want = "evidence-invalid"
			case "no-transitivity":
				b.Query.TargetScope = "third"
				want = "unbound"
			}
			b.Events = append(b.Events, link)
			state(t, b, want)
		})
	}
}

// A deliberately simple independent history oracle: pairwise generation checks
// and predecessor searches, without the implementation's generation index.
func oracle(events []cb.Envelope, o cb.Event) string {
	var selected *cb.Event
	for i := range events {
		if events[i].Event.ID() == o.Binding {
			e := events[i].Event
			selected = &e
		}
	}
	if selected == nil {
		return "unbound"
	}
	s := *selected
	gap, old := false, false
	for _, x := range events {
		a := x.Event
		if a.Kind != "binding" || a.Scope != s.Scope || a.Channel != s.Channel {
			continue
		}
		for _, y := range events {
			b := y.Event
			if b.Kind == "binding" && b.Scope == a.Scope && b.Channel == a.Channel && b.Generation == a.Generation && a.ID() != b.ID() {
				return "conflicting"
			}
		}
		if a.Generation > 1 {
			found := false
			for _, y := range events {
				b := y.Event
				if b.ID() == a.Previous && b.Kind == "binding" && b.Channel == a.Channel && b.Scope == a.Scope && b.Generation+1 == a.Generation && b.Issued <= a.Issued {
					found = true
				}
			}
			if !found {
				gap = true
			}
		}
		if a.Generation > s.Generation {
			old = true
		}
	}
	if gap {
		return "unknown-stale"
	}
	if old {
		return "superseded-in-known-history"
	}
	return "evidence-consistent-within-known-scope"
}
func TestGeneratedMergeLawsAndReference(t *testing.T) {
	rng := rand.New(rand.NewSource(8032))
	for trial := 0; trial < 1200; trial++ {
		b := demo.Base()
		prev := b.Events[0].Event.ID()
		n := 1 + rng.Intn(7)
		for j := 2; j <= n; j++ {
			e := demo.Binding(fmt.Sprintf("b-%d", j), fmt.Sprintf("subject-%d", j), int64(j), prev)
			b.Events = append(b.Events, e)
			prev = e.Event.ID()
		}
		switch trial % 4 {
		case 0:
			if n > 1 {
				b.Events = append(b.Events, demo.Binding("fork", "Maple", 2, b.Events[0].Event.ID()))
			}
		case 1:
			if n > 2 {
				b.Events = append(b.Events[:2], b.Events[3:]...)
			}
		case 2:
			if n > 1 {
				b.Events = b.Events[1:]
			}
		}
		var obs cb.Event
		for _, e := range b.Events {
			if e.Event.Kind == "observation" {
				obs = e.Event
			}
		}
		expected := oracle(b.Events, obs)
		first := result(t, b)
		if first.State != expected {
			t.Fatalf("trial %d reference=%s actual=%+v", trial, expected, first)
		}
		// Union three partitions in different groupings/orders. Dedup is set semantics.
		p, q, r := []cb.Envelope{}, []cb.Envelope{}, []cb.Envelope{}
		for i, e := range b.Events {
			switch i % 3 {
			case 0:
				p = append(p, e)
			case 1:
				q = append(q, e)
			case 2:
				r = append(r, e)
			}
		}
		b.Events = append(append(append([]cb.Envelope{}, r...), p...), q...)
		rng.Shuffle(len(b.Events), func(i, j int) { b.Events[i], b.Events[j] = b.Events[j], b.Events[i] })
		second := result(t, b)
		if !reflect.DeepEqual(first, second) {
			t.Fatalf("order/partition law trial %d", trial)
		}
		b.Events = append(b.Events, b.Events...)
		third := result(t, b)
		third.Duplicates = 0
		if !reflect.DeepEqual(first, third) {
			t.Fatalf("idempotence trial %d", trial)
		}
	}
}
func TestParserStrictness(t *testing.T) {
	b := demo.Base()
	raw, _ := json.Marshal(b)
	cases := [][]byte{append(append([]byte{}, raw...), []byte("{}")...), []byte("null"), []byte(`{"v":1,"v":1}`), bytes.Replace(raw, []byte(`"v":1`), []byte(`"V":1`), 1), bytes.Replace(raw, []byte(`"v":1`), []byte(`"v":1.0`), 1), bytes.Replace(raw, []byte(`"v":1`), []byte(`"v":2`), 1), bytes.Replace(raw, []byte(`"v":1`), []byte(`"v":1,"extra":0`), 1), bytes.Replace(raw, []byte(`"events":[`), []byte(`"events":null,"x":[`), 1), []byte(strings.Repeat("[", 10) + strings.Repeat("]", 10)), []byte(strings.Repeat(" ", cb.MaxInput+1)), append([]byte{255}, raw...)}
	for i, v := range cases {
		if _, err := cb.Parse(v); err == nil {
			t.Fatalf("accepted malformed %d", i)
		}
	}
	if _, err := cb.Parse(raw); err != nil {
		t.Fatal(err)
	}
	for _, n := range []int{513, 1025} {
		b.Events = make([]cb.Envelope, n)
		for i := range b.Events {
			b.Events[i] = demo.Base().Events[0]
		}
		v, _ := json.Marshal(b)
		if _, err := cb.Parse(v); err == nil {
			t.Fatal("limit")
		}
	}
}
func TestTypedLimitsAndSign(t *testing.T) {
	b := demo.Base()
	b.Events = make([]cb.Envelope, 513)
	if _, err := cb.Evaluate(b); err == nil {
		t.Fatal("typed limit")
	}
	if _, err := cb.Sign(demo.Base().Events[0].Event, nil); err == nil {
		t.Fatal("key length")
	}
}
func TestRejectedCannotPoisonAccepted(t *testing.T) {
	b := demo.Base()
	e := b.Events[0]
	e.Signature = strings.Repeat("0", 128)
	b.Events = append(b.Events, e)
	state(t, b, "evidence-consistent-within-known-scope")
}
func FuzzParse(f *testing.F) {
	raw, _ := json.Marshal(demo.Base())
	f.Add(raw)
	f.Add([]byte(`{"v":1,"v":1}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, data []byte) {
		b, err := cb.Parse(data)
		if err == nil {
			r, err := cb.Evaluate(b)
			if err != nil {
				t.Fatal(err)
			}
			if r.PhysicalIdentity != "unverified" {
				t.Fatal("physical claim")
			}
		}
	})
}

func TestObservationEquivocation(t *testing.T) {
	b := demo.Base()
	other := b.Events[1].Event
	other.Payload = strings.Repeat("0", 64)
	b.Events = append(b.Events, demo.Sign(other))
	state(t, b, "conflicting")
}

func TestLimitsAtBoundary(t *testing.T) {
	b := demo.Base()
	b.Events = make([]cb.Envelope, 512)
	for i := range b.Events {
		b.Events[i] = demo.Base().Events[0]
	}
	b.Events[511] = demo.Base().Events[1]
	raw, _ := json.Marshal(b)
	parsed, err := cb.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	state(t, parsed, "evidence-consistent-within-known-scope")
	b.Policy.KnownAt = b.Policy.AsOf - b.Policy.MaxAge
	state(t, b, "evidence-consistent-within-known-scope")
	b.Policy.KnownAt--
	state(t, b, "unknown-stale")
	huge := []byte(`{"` + strings.Repeat("a", 65) + `":0}`)
	if _, err := cb.Parse(huge); err == nil {
		t.Fatal("field bound")
	}
	for _, value := range []string{strings.Repeat("x", 1025), "\u0000", "\ufffd"} {
		base := demo.Base()
		base.Events[0].Event.Subject = value
		raw, _ := json.Marshal(base)
		if _, err := cb.Parse(raw); err == nil {
			t.Fatal("invalid subject")
		}
	}
	base := demo.Base()
	base.Policy.Keys = make([]cb.Key, 65)
	for i := range base.Policy.Keys {
		base.Policy.Keys[i] = demo.Key(fmt.Sprintf("issuer-%d", i), "field", "binding")
	}
	if _, err := cb.Evaluate(base); err == nil {
		t.Fatal("key limit")
	}
}
func TestMaliciousAcceptedIssuerAndIsolation(t *testing.T) {
	b := demo.Base()
	b.Policy.Keys = append(b.Policy.Keys, demo.Key("malicious", "field", "binding"))
	lie := b.Events[0].Event
	lie.Issuer = "malicious"
	lie.Subject = "Birch"
	b.Events = append(b.Events, demo.Sign(lie))
	state(t, b, "conflicting")
	// A separate source-instance must not poison the selected channel.
	b.Events[len(b.Events)-1].Event.Channel = "other-source:1"
	b.Events[len(b.Events)-1] = demo.Sign(b.Events[len(b.Events)-1].Event)
	state(t, b, "evidence-consistent-within-known-scope")
}
func TestWrongTagIsUndetectable(t *testing.T) {
	b := demo.Base()
	cedarPhysicalProp := result(t, b)
	// Moving a copy of the carrier changes no protocol bytes. There is no physical
	// observation input from which the runtime could establish the mismatch.
	birchPhysicalProp := result(t, b)
	if !reflect.DeepEqual(cedarPhysicalProp, birchPhysicalProp) || birchPhysicalProp.PhysicalIdentity != "unverified" {
		t.Fatal("invented physical assurance")
	}
}
