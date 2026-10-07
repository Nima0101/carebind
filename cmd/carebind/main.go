package main

import (
	"encoding/json"
	"fmt"
	cb "github.com/Nima0101/carebind"
	"github.com/Nima0101/carebind/internal/demo"
	"io"
	"os"
)

func read(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file_open")
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, cb.MaxInput+1))
	if err != nil {
		return nil, fmt.Errorf("file_read")
	}
	if len(b) > cb.MaxInput {
		return nil, fmt.Errorf("input_limit")
	}
	return b, nil
}
func emit(v any) error { return json.NewEncoder(os.Stdout).Encode(v) }
func run(args []string) error {
	if len(args) == 1 && args[0] == "version" {
		fmt.Println("carebind " + cb.Version + " wire=1")
		return nil
	}
	if len(args) == 1 && args[0] == "demo" {
		fmt.Println("CareBind | The channel changes. The origin does not.")
		fmt.Println("SYNTHETIC ONLY | global freshness: unknown | physical identity: unverified")
		for _, c := range demo.Cases() {
			r, err := cb.Evaluate(c.Bundle)
			if err != nil {
				return err
			}
			if r.State != c.Want {
				return fmt.Errorf("demo_assertion: %s", c.Name)
			}
			fmt.Printf("\n%s\n  %s | origin=%s | accepted=%d duplicates=%d\n  %s\n", c.Name, r.State, r.Origin, len(r.Accepted), r.Duplicates, c.Note)
		}
		return nil
	}
	if len(args) == 2 && args[0] == "fixture" {
		for _, c := range demo.Cases() {
			if c.Name == args[1] {
				return emit(c.Bundle)
			}
		}
		return fmt.Errorf("unknown_fixture")
	}
	if len(args) == 1 && args[0] == "demo-policy" {
		return emit(demo.Base().Policy)
	}
	if len(args) == 3 && args[0] == "evaluate" {
		data, err := read(args[1])
		if err != nil {
			return err
		}
		b, err := cb.Parse(data)
		if err != nil {
			return err
		}
		data, err = read(args[2])
		if err != nil {
			return err
		}
		var p cb.Policy
		if err = cb.StrictDecode(data, &p); err != nil {
			return err
		}
		b.Policy = p
		r, err := cb.Evaluate(b)
		if err != nil {
			return err
		}
		return emit(r)
	}
	return fmt.Errorf("usage: carebind version | demo | fixture NAME | demo-policy | evaluate BUNDLE LOCAL_POLICY")
}
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
