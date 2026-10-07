package carebind

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"unicode/utf8"
)

// StrictDecode rejects duplicate, missing, unknown and mis-cased fields, nulls,
// invalid UTF-8 and excessive work before constructing protocol structures.
func StrictDecode(data []byte, out any) error {
	if len(data) > MaxInput {
		return fail("input_limit")
	}
	if !utf8.Valid(data) {
		return fail("invalid_utf8")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := scan(d, 0); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return fail("trailing_data")
	}
	var raw any
	d = json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if d.Decode(&raw) != nil {
		return fail("json_syntax")
	}
	t := reflect.TypeOf(out)
	if t == nil || t.Kind() != reflect.Pointer {
		return fail("decode_target")
	}
	if err := shape(raw, t.Elem()); err != nil {
		return err
	}
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if d.Decode(out) != nil {
		return fail("json_type")
	}
	return nil
}
func scan(d *json.Decoder, depth int) error {
	if depth > 8 {
		return fail("depth_limit")
	}
	t, err := d.Token()
	if err != nil {
		return fail("json_syntax")
	}
	if t == nil {
		return fail("null_forbidden")
	}
	if s, ok := t.(string); ok && len(s) > 1024 {
		return fail("string_limit")
	}
	if delim, ok := t.(json.Delim); ok {
		count := 0
		seen := map[string]bool{}
		for d.More() {
			count++
			if delim == '{' {
				if count > 32 {
					return fail("object_limit")
				}
				k, err := d.Token()
				if err != nil {
					return fail("json_syntax")
				}
				s, ok := k.(string)
				if !ok || len(s) > 64 {
					return fail("field_shape")
				}
				if seen[s] {
					return fail("duplicate_field")
				}
				seen[s] = true
			} else if delim == '[' {
				if count > 1024 {
					return fail("array_limit")
				}
			} else {
				return fail("json_syntax")
			}
			if err := scan(d, depth+1); err != nil {
				return err
			}
		}
		if _, err := d.Token(); err != nil {
			return fail("json_syntax")
		}
	}
	return nil
}
func shape(v any, t reflect.Type) error {
	switch t.Kind() {
	case reflect.Struct:
		m, ok := v.(map[string]any)
		if !ok || len(m) != t.NumField() {
			return fail("object_shape")
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			val, ok := m[f.Tag.Get("json")]
			if !ok {
				return fail("field_shape")
			}
			if err := shape(val, f.Type); err != nil {
				return err
			}
		}
	case reflect.Slice:
		a, ok := v.([]any)
		if !ok {
			return fail("array_shape")
		}
		for _, x := range a {
			if err := shape(x, t.Elem()); err != nil {
				return err
			}
		}
	}
	return nil
}
func Parse(data []byte) (Bundle, error) {
	var b Bundle
	if err := StrictDecode(data, &b); err != nil {
		return b, err
	}
	return b, validate(b)
}
func validate(b Bundle) error {
	if b.V != 1 {
		return fail("unsupported_version")
	}
	if b.Events == nil || len(b.Events) > MaxEvents {
		return fail("event_limit")
	}
	if err := b.Policy.Validate(); err != nil {
		return err
	}
	if !digestRE.MatchString(b.Query.Observation) || ((b.Query.TargetScope != "" || b.Query.Target != "") && (!token(b.Query.TargetScope) || !token(b.Query.Target))) {
		return fail("query_shape")
	}
	for _, e := range b.Events {
		if err := e.Event.Validate(); err != nil {
			return err
		}
		if !sigRE.MatchString(e.Signature) {
			return fail("signature_shape")
		}
	}
	return nil
}
