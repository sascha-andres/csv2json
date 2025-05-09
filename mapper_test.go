package csv2json

import "testing"

// TestCreate tests the creation of Mapper instances using various configurations and validates expected errors or outcomes.
func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		options []OptionFunc
		want    *Mapper
		wantErr bool
	}{
		{
			name: "basic mapper",
			options: []OptionFunc{
				WithIn("input.csv"),
				WithOut("output.json"),
			},
			want: &Mapper{
				In:  "input.csv",
				Out: "output.json",
			},
			wantErr: false,
		},
		{
			name: "mapper with array and named",
			options: []OptionFunc{
				WithIn("input.csv"),
				WithOut("output.json"),
				WithArray(true),
				WithNamed(true),
			},
			want: &Mapper{
				In:    "input.csv",
				Out:   "output.json",
				Array: true,
				Named: true,
			},
			wantErr: false,
		},
		{
			name: "empty input error",
			options: []OptionFunc{
				WithIn(""),
				WithOut("output.json"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty output error",
			options: []OptionFunc{
				WithIn("input.csv"),
				WithOut(""),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMapper(tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareMappers(got, tt.want) {
				t.Errorf("NewMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

// compareMappers compares two Mapper instances for equality, considering their public fields: In, Out, Array, and Named.
// TODO switch to cmp
func compareMappers(a, b *Mapper) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.In == b.In &&
		a.Out == b.Out &&
		a.Array == b.Array &&
		a.Named == b.Named
}
