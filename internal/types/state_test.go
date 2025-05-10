package types

import "testing"

func TestString(t *testing.T) {
	t.Run("returns open for StateOpen", func(t *testing.T) {
		state := StateOpen
		want := "open"

		got := state.String()
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns closed for StateClosed", func(t *testing.T) {
		state := StateClosed
		want := "closed"

		got := state.String()
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("returns JSON for StateOpen", func(t *testing.T) {
		state := StateOpen
		want := `"open"`

		got, err := state.MarshalJSON()
		if err != nil {
			t.Fatal("expected nil, got error")
		}

		if string(got) != want {
			t.Errorf("got %q, want %q", string(got), want)
		}
	})

	t.Run("returns JSON for StateClosed", func(t *testing.T) {
		state := StateClosed
		want := `"closed"`

		got, err := state.MarshalJSON()
		if err != nil {
			t.Fatal("expected nil, got error")
		}

		if string(got) != want {
			t.Errorf("got %q, want %q", string(got), want)
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("returns error for invalid JSON", func(t *testing.T) {
		data := []byte(`example`)
		var state State

		if err := state.UnmarshalJSON(data); err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns StateOpen for open", func(t *testing.T) {
		data := []byte(`"open"`)
		var state State
		want := StateOpen

		if err := state.UnmarshalJSON(data); err != nil {
			t.Fatal("expected nil, got error")
		}

		if state != want {
			t.Errorf("got %v, want %v", state, want)
		}
	})

	t.Run("returns StateClosed for closed", func(t *testing.T) {
		data := []byte(`"closed"`)
		var state State
		want := StateClosed

		if err := state.UnmarshalJSON(data); err != nil {
			t.Fatal("expected nil, got error")
		}

		if state != want {
			t.Errorf("got %v, want %v", state, want)
		}
	})
}
