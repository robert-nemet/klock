package pkg_webhook

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isMatch(t *testing.T) {
	evaluator := NewEvaluator()
	type args struct {
		rule  string
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
		err  bool
	}{
		{
			name: "Simple match",
			args: args{
				rule:  "test",
				value: "test",
			},
			want: true,
		},
		{
			name: "Simple do not match",
			args: args{
				rule:  "test",
				value: "not_match",
			},
			want: false,
		},
		{
			name: "OR red|blue match",
			args: args{
				rule:  "red|blue",
				value: "blue",
			},
			want: true,
		},
		{
			name: "(red|green)&blue match",
			args: args{
				rule:  "(red|green)&blue",
				value: "green",
			},
			want: false,
		},
		{
			name: "red|green&!blue match",
			args: args{
				rule:  "red|green&^blue",
				value: "green",
			},
			want: true,
		},
		{
			name: "red|green&^|^|^blue error",
			args: args{
				rule:  "red|green&^|^|^blue",
				value: "green",
			},
			want: false,
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluator.IsMatch(tt.args.rule, tt.args.value)
			if tt.err {
				assert.Error(t, err)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_prepack(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			name: "Simple",
			args: "blue",
			want: []string{"blue"},
		},
		{
			name: "One or",
			args: "blue|green",
			want: []string{"blue", "|", "green"},
		},
		{
			name: "Or and And",
			args: "blue&green|violet",
			want: []string{"blue", "&", "green", "|", "violet"},
		},
		{
			name: "Some brackets",
			args: "green&(red|blue)",
			want: []string{"green", "&", "(", "red", "|", "blue", ")"},
		},
		{
			name: "Some brackets 2nd",
			args: "(red|blue)&green",
			want: []string{"(", "red", "|", "blue", ")", "&", "green"},
		},
		{
			name: "Some brackets 3rd",
			args: "(red|blue)&(green|black)",
			want: []string{"(", "red", "|", "blue", ")", "&", "(", "green", "|", "black", ")"},
		},
		{
			name: "Some brackets and not",
			args: "(red|blue)&^(green|black)",
			want: []string{"(", "red", "|", "blue", ")", "&", "^", "(", "green", "|", "black", ")"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepack(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepare(t *testing.T) {
	type args struct {
		input []string
		value string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple",
			args: args{
				input: []string{"blue"},
				value: "blue",
			},
			want: []string{"true"},
		},
		{
			name: "simple false",
			args: args{
				input: []string{"blue"},
				value: "red",
			},
			want: []string{"false"},
		},
		{
			name: "simple not",
			args: args{
				input: []string{"^", "blue"},
				value: "red",
			},
			want: []string{"^", "false"},
		},
		{
			name: "simple or",
			args: args{
				input: []string{"blue", "|", "red"},
				value: "blue",
			},
			want: []string{"true", "|", "false"},
		},
		{
			name: "simple and",
			args: args{
				input: []string{"blue", "&", "red"},
				value: "blue",
			},
			want: []string{"true", "&", "false"},
		},
		{
			name: "brakets: (a|b)&c",
			args: args{
				input: []string{"(", "blue", "|", "red", ")", "&", "green"},
				value: "blue",
			},
			want: []string{"(", "true", "|", "false", ")", "&", "false"},
		},
		{
			name: "brakets: (a|b)&^c",
			args: args{
				input: []string{"(", "blue", "|", "red", ")", "&", "^", "green"},
				value: "blue",
			},
			want: []string{"(", "true", "|", "false", ")", "&", "^", "false"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepare(tt.args.input, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepare() = %v, want %v", got, tt.want)
			}
		})
	}
}
