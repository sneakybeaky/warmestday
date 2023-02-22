package main_test

import (
	"testing"
	"warmestday/cmd/web"
)

func TestNumDecimalPlaces(t *testing.T) {

	t.Parallel()

	cases := map[string]struct {
		value float64
		want  int
	}{
		"No floating part": {
			value: 1,
			want:  0,
		},
		"One decimal place": {
			value: 0.1,
			want:  1,
		},
		"Two decimal places": {
			value: 0.01,
			want:  2,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := main.NumDecPlaces(tc.value)

			if got != tc.want {
				t.Fatalf("wanted %d got %d", tc.want, got)
			}

		})
	}

}
