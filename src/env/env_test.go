package env_test

import (
	"os"
	"testing"

	env "github.com/czerasz/go-lambda-sns-to-es/src/env"
)

func TestAllVars(t *testing.T) {
	envVars := os.Environ()

	if len(envVars) != len(env.AllVars()) {
		t.Errorf("AllVars() should have the same ammount of items as os.Environ()")
	}
}

func TestGetEnv(t *testing.T) {
	extVarName := "EXISTING_VARIABLE"
	extVarVal := "existing variable"
	os.Setenv(extVarName, extVarVal)
	defer os.Unsetenv(extVarName)

	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "environment variable exists",
			args: args{
				key:      extVarName,
				fallback: "fallback",
			},
			want: extVarVal,
		},
		{
			name: "environment variable does NOT exists",
			args: args{
				key:      "NOT_EXISTING_VARIABLE",
				fallback: "fallback",
			},
			want: "fallback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := env.GetEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
