package carebind

import "sort"

// Evaluate rechecks typed inputs, authenticates the set, and evaluates historical
// association. The caller must provide its own trusted local policy.
func Evaluate(b Bundle) (Result, error) {
	if err := validate(b); err != nil {
		return Result{}, err
	}
	r := Result{V: 1, State: "unbound", Reason: "missing_observation", KnownAt: b.Policy.KnownAt, GlobalFreshness: "unknown", PhysicalIdentity: "unverified", Accepted: []string{}, Rejected: []Rejection{}}
	finish := func(state, reason string) (Result, error) { r.State = state; r.Reason = reason; return r, nil }
	accepted := map[string]Event{}
	seen := map[string]bool{}
	rejectSeen := map[string]bool{}
	nonces := map[string]map[string]bool{}
	for _, env := range b.Events {
		id := env.Event.ID()
		wire := id + env.Signature
		if seen[wire] {
			r.Duplicates++
			continue
		}
		seen[wire] = true
		if reason := verify(env, b.Policy); reason != "" {
			key := id + reason
			if !rejectSeen[key] {
				r.Rejected = append(r.Rejected, Rejection{id, reason})
				rejectSeen[key] = true
			}
			continue
		}
		accepted[id] = env.Event
		n := env.Event.Issuer + "/" + env.Event.Nonce
		if nonces[n] == nil {
			nonces[n] = map[string]bool{}
		}
		nonces[n][id] = true
	}
	for id := range accepted {
		r.Accepted = append(r.Accepted, id)
	}
	sort.Strings(r.Accepted)
	sort.Slice(r.Rejected, func(i, j int) bool {
		if r.Rejected[i].ID == r.Rejected[j].ID {
			return r.Rejected[i].Reason < r.Rejected[j].Reason
		}
		return r.Rejected[i].ID < r.Rejected[j].ID
	})
	equiv := func(e Event) bool { return len(nonces[e.Issuer+"/"+e.Nonce]) > 1 }
	o, ok := accepted[b.Query.Observation]
	if !ok {
		for _, rej := range r.Rejected {
			if rej.ID == b.Query.Observation {
				return finish("evidence-invalid", rej.Reason)
			}
		}
		return r, nil
	}
	if o.Kind != "observation" {
		return finish("evidence-invalid", "query_not_observation")
	}
	binding, ok := accepted[o.Binding]
	if !ok {
		return finish("unbound", "missing_binding")
	}
	if binding.Kind != "binding" || binding.Scope != o.Scope || binding.Channel != o.Channel || binding.Generation != o.Generation || o.Issued < binding.Issued {
		return finish("evidence-invalid", "origin_mismatch")
	}
	r.OriginScope = binding.Scope
	r.Origin = binding.Subject
	r.Binding = o.Binding
	conflict := equiv(o)
	gap := false
	superseded := false
	generations := map[int64]string{}
	for _, id := range r.Accepted {
		e := accepted[id]
		if e.Kind != "binding" || e.Scope != binding.Scope || e.Channel != binding.Channel {
			continue
		}
		if prior, ok := generations[e.Generation]; ok && prior != id {
			conflict = true
		}
		generations[e.Generation] = id
		if equiv(e) {
			conflict = true
		}
		if e.Generation > 1 {
			parent, ok := accepted[e.Previous]
			if !ok || parent.Kind != "binding" || parent.Scope != e.Scope || parent.Channel != e.Channel || parent.Generation != e.Generation-1 || parent.Issued > e.Issued {
				gap = true
			}
		}
		if e.Generation > binding.Generation {
			superseded = true
		}
		if e.Previous == o.Binding && o.Issued >= e.Issued {
			conflict = true
		}
	}
	if conflict {
		return finish("conflicting", "history_conflict")
	}
	if gap {
		return finish("unknown-stale", "incomplete_history")
	}
	if binding.Subject == "" {
		return finish("unbound", "explicit_unbound")
	}
	if b.Query.TargetScope != "" {
		if b.Query.TargetScope == binding.Scope {
			if b.Query.Target != binding.Subject {
				return finish("conflicting", "consumer_context_mismatch")
			}
		} else {
			targets := map[string]bool{}
			badLink := false
			linkEquiv := false
			for _, id := range r.Accepted {
				e := accepted[id]
				if e.Kind != "link" || e.Binding != o.Binding || e.TargetScope != b.Query.TargetScope {
					continue
				}
				if e.Scope != binding.Scope || e.Subject != binding.Subject || e.Issued < binding.Issued {
					badLink = true
				}
				if equiv(e) {
					linkEquiv = true
				}
				targets[e.Target] = true
			}
			if linkEquiv || len(targets) > 1 {
				return finish("conflicting", "link_conflict")
			}
			if badLink {
				return finish("evidence-invalid", "link_origin_mismatch")
			}
			if len(targets) == 0 {
				return finish("unbound", "missing_link")
			}
			if !targets[b.Query.Target] {
				return finish("conflicting", "consumer_context_mismatch")
			}
		}
		r.TargetScope = b.Query.TargetScope
		r.Target = b.Query.Target
	}
	if b.Policy.AsOf-b.Policy.KnownAt > b.Policy.MaxAge {
		return finish("unknown-stale", "stale_policy")
	}
	if superseded {
		return finish("superseded-in-known-history", "origin_preserved")
	}
	return finish("evidence-consistent-within-known-scope", "origin_preserved")
}
