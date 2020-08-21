package utils

import "testing"

func TestStringsContainString(t *testing.T) {
	type args struct {
		strings []string
		needle  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"found",
			args{
				strings: []string{"test", "test2", "test3"},
				needle:  "test3",
			},
			true,
		},
		{
			"notFound",
			args{
				strings: []string{"test", "test2", "test3"},
				needle:  "test4",
			},
			false,
		},
		{
			"empty",
			args{
				strings: []string{},
				needle:  "test4",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringsContainString(tt.args.strings, tt.args.needle); got != tt.want {
				t.Errorf("StringsContainString() = %v, want %v", got, tt.want)
			}
		})
	}
}
