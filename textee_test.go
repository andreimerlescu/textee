package textee

import (
	"reflect"
	"testing"

	"github.com/andreimerlescu/gematria"
)

func TestTextee_ParseString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want map[string]int32
	}{
		{
			name: "test 1",
			args: args{
				input: "All right let's move from this point on 16 March 84, let's move in time to our second location which is a specific building near where you are now. Are you ready? Just a minute. All. right, I will wait. All right, move now from this area to the front ground level of the building known as the Menara Building, to the front of, on the ground, the Menara Building.",
			},
			want: map[string]int32{
				// Populate this map with the expected substrings and their counts.
				"all right lets": 1,
				// ... more substrings and counts ...
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tt, err := NewTextee(tc.args.input)
			if err != nil {
				t.Errorf("ParseString() error = %v", err)
			}
			tt, err = tt.ParseString(tc.args.input)
			if err != nil {
				t.Errorf("ParseString() error = %v", err)
			}

			for substring, expectedCount := range tc.want {
				if count, ok := tt.Substrings[substring]; !ok || count.Load() != expectedCount {
					t.Errorf("Substring count for %q = %v, want %v", substring, count.Load(), expectedCount)
				}
			}
		})
	}
}

func TestTextee_CalculateGematria(t *testing.T) {
	t.Run("testing calculate gematria", func(t *testing.T) {
		tt, err := NewTextee("manifesting three six nine")
		if err != nil {
			t.Errorf("CalculateGematria() error = %v", err)
		}
		got, err2 := tt.CalculateGematria()
		if err2 != nil {
			t.Errorf("CalculateGematria() error = %v", err2)
		}
		want := gematria.Gematria{
			Jewish:  337,
			English: 702,
			Simple:  117,
		}
		if !reflect.DeepEqual(got.Gematrias["manifesting"], want) {
			t.Errorf("CalculateGematria() = %v, want %v", got, want)
		}
	})
}
