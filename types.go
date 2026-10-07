// Package carebind evaluates bounded, source-bound association evidence.
// It never verifies physical identity or makes clinical decisions.
package carebind

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"regexp"
)

const Version = "0.1.0"
const MaxInput = 1 << 20
const MaxEvents = 512
const domain = "CareBind/v1\n"

type Error struct{ Code string }

func (e *Error) Error() string { return e.Code }
func fail(code string) error   { return &Error{code} }

type Event struct {
	V           int    `json:"v"`
	Kind        string `json:"kind"`
	Issuer      string `json:"issuer"`
	Key         string `json:"key"`
	Nonce       string `json:"nonce"`
	Scope       string `json:"scope"`
	Channel     string `json:"channel"`
	Generation  int64  `json:"generation"`
	Subject     string `json:"subject"`
	Previous    string `json:"previous"`
	Binding     string `json:"binding"`
	Payload     string `json:"payload"`
	TargetScope string `json:"target_scope"`
	Target      string `json:"target"`
	Issued      int64  `json:"issued"`
}
type Envelope struct {
	Event     Event  `json:"event"`
	Signature string `json:"signature"`
}
type Key struct {
	Issuer    string   `json:"issuer"`
	Key       string   `json:"key"`
	Public    string   `json:"public"`
	Scope     string   `json:"scope"`
	Roles     []string `json:"roles"`
	NotBefore int64    `json:"not_before"`
	NotAfter  int64    `json:"not_after"`
	Revoked   bool     `json:"revoked"`
}
type Policy struct {
	AsOf    int64 `json:"as_of"`
	KnownAt int64 `json:"known_at"`
	MaxAge  int64 `json:"max_age"`
	Keys    []Key `json:"keys"`
}
type Query struct {
	Observation string `json:"observation"`
	TargetScope string `json:"target_scope"`
	Target      string `json:"target"`
}
type Bundle struct {
	V      int        `json:"v"`
	Events []Envelope `json:"events"`
	Policy Policy     `json:"policy"`
	Query  Query      `json:"query"`
}
type Rejection struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}
type Result struct {
	V                int         `json:"v"`
	State            string      `json:"state"`
	Reason           string      `json:"reason"`
	OriginScope      string      `json:"origin_scope"`
	Origin           string      `json:"origin"`
	Binding          string      `json:"binding"`
	TargetScope      string      `json:"target_scope"`
	Target           string      `json:"target"`
	KnownAt          int64       `json:"known_at"`
	GlobalFreshness  string      `json:"global_freshness"`
	PhysicalIdentity string      `json:"physical_identity"`
	Accepted         []string    `json:"accepted"`
	Rejected         []Rejection `json:"rejected"`
	Duplicates       int         `json:"duplicates"`
}

var tokenRE = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,64}$`)
var digestRE = regexp.MustCompile(`^[0-9a-f]{64}$`)
var sigRE = regexp.MustCompile(`^[0-9a-f]{128}$`)

func token(s string) bool    { return tokenRE.MatchString(s) }
func timestamp(n int64) bool { return n >= 0 && n <= 4102444800 }
func kind(s string) bool     { return s == "binding" || s == "observation" || s == "link" }
func (e Event) Validate() error {
	if e.V != 1 || !kind(e.Kind) || !token(e.Issuer) || !token(e.Key) || !token(e.Nonce) || !token(e.Scope) || !timestamp(e.Issued) {
		return fail("event_shape")
	}
	switch e.Kind {
	case "binding":
		if !token(e.Channel) || e.Generation < 1 || e.Generation > 2147483647 || (e.Subject != "" && !token(e.Subject)) || (e.Generation == 1 && e.Previous != "") || (e.Generation > 1 && !digestRE.MatchString(e.Previous)) || e.Binding != "" || e.Payload != "" || e.TargetScope != "" || e.Target != "" {
			return fail("binding_shape")
		}
	case "observation":
		if !token(e.Channel) || e.Generation < 1 || e.Generation > 2147483647 || !digestRE.MatchString(e.Binding) || !digestRE.MatchString(e.Payload) || e.Subject != "" || e.Previous != "" || e.TargetScope != "" || e.Target != "" {
			return fail("observation_shape")
		}
	case "link":
		if !token(e.Subject) || !token(e.TargetScope) || !token(e.Target) || e.TargetScope == e.Scope || !digestRE.MatchString(e.Binding) || e.Generation != 0 || e.Channel != "" || e.Previous != "" || e.Payload != "" {
			return fail("link_shape")
		}
	}
	return nil
}

// Canonical returns the protocol's fixed-field ASCII encoding, not general JCS.
func (e Event) Canonical() ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(e)
}
func (e Event) ID() string {
	b, _ := json.Marshal(e)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func Sign(e Event, private ed25519.PrivateKey) (Envelope, error) {
	b, err := e.Canonical()
	if err != nil {
		return Envelope{}, err
	}
	if len(private) != ed25519.PrivateKeySize {
		return Envelope{}, fail("private_key_length")
	}
	return Envelope{e, hex.EncodeToString(ed25519.Sign(private, append([]byte(domain), b...)))}, nil
}
func (p Policy) Validate() error {
	if !timestamp(p.AsOf) || !timestamp(p.KnownAt) || p.KnownAt > p.AsOf || p.MaxAge < 0 || p.MaxAge > 86400 || len(p.Keys) > 64 || p.Keys == nil {
		return fail("policy_shape")
	}
	seen := map[string]bool{}
	for _, k := range p.Keys {
		name := k.Issuer + "/" + k.Key
		if !token(k.Issuer) || !token(k.Key) || !token(k.Scope) || !digestRE.MatchString(k.Public) || !timestamp(k.NotBefore) || !timestamp(k.NotAfter) || k.NotBefore > k.NotAfter || len(k.Roles) == 0 || len(k.Roles) > 3 || seen[name] {
			return fail("key_shape")
		}
		seen[name] = true
		roles := map[string]bool{}
		for _, r := range k.Roles {
			if !kind(r) || roles[r] {
				return fail("key_roles")
			}
			roles[r] = true
		}
	}
	return nil
}
func verify(e Envelope, p Policy) string {
	for _, k := range p.Keys {
		if k.Issuer != e.Event.Issuer || k.Key != e.Event.Key {
			continue
		}
		pub, _ := hex.DecodeString(k.Public)
		sig, _ := hex.DecodeString(e.Signature)
		b, _ := e.Event.Canonical()
		if !ed25519.Verify(pub, append([]byte(domain), b...), sig) {
			return "bad_signature"
		}
		if k.Revoked {
			return "revoked_key"
		}
		if p.AsOf < k.NotBefore || p.AsOf > k.NotAfter || e.Event.Issued < k.NotBefore || e.Event.Issued > k.NotAfter || e.Event.Issued > p.AsOf {
			return "key_time"
		}
		scope := e.Event.Scope
		if e.Event.Kind == "link" {
			scope = e.Event.TargetScope
		}
		if k.Scope != scope {
			return "unauthorized_scope"
		}
		for _, r := range k.Roles {
			if r == e.Event.Kind {
				return ""
			}
		}
		return "unauthorized_role"
	}
	return "unknown_issuer"
}
