package filter

import (
	"fmt"
	"testing"
)

func TestCompact(t *testing.T) {
	got := []string{"olá", "", "mundo"}
	got = Compact(got)
	want := []string{"olá", "mundo"}
	if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
		t.Errorf("got %q want %q\n", got, want)
	}
}

func TestSliptValues(t *testing.T) {
	t.Run("when there are some values", func(t *testing.T) {
		input := "um,dois,tres"
		got := SplitValues(&input)
		want := []string{"um", "dois", "tres"}
		if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
			t.Errorf("got %q want %q\n", got, want)
		}
	})
	t.Run("when there are commas without values", func(t *testing.T) {
		input := ",,,"
		got := SplitValues(&input)
		want := []string{}
		if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
			t.Errorf("got %q want %q\n", got, want)
		}
	})
	t.Run("when there is just an empty string", func(t *testing.T) {
		input := ""
		got := SplitValues(&input)
		want := []string{}
		if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
			t.Errorf("got %q want %q\n", got, want)
		}
	})
}
