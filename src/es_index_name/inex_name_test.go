package index_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	indexName "github.com/czerasz/go-lambda-sns-to-es/src/es_index_name"
)

func TestIndexName(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "simple",
			want:    fmt.Sprintf("sns-%d.%02d.%02d", now.Year(), now.Month(), now.Day()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := indexName.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexNameWithTemplate(t *testing.T) {
	os.Setenv(indexName.IdxVarName, "{{ .Env.PREFIX }}-{{ .Date.Year }}.{{ .Date.Month }}.{{ .Date.Day }}")
	os.Setenv("PREFIX", "test")
	defer func() {
		os.Unsetenv(indexName.IdxVarName)
		os.Unsetenv("PREFIX")
	}()

	now := time.Now()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "custom index template",
			want:    fmt.Sprintf("test-%d.%02d.%02d", now.Year(), now.Month(), now.Day()),
			wantErr: false,
		},
		{
			name:    "invalid index template",
			want:    fmt.Sprintf("test-%d.%02d.%02d", now.Year(), now.Month(), now.Day()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := indexName.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexNameWithInvalidTemplate(t *testing.T) {
	os.Setenv(indexName.IdxVarName, "{{ Env }}-{{ .Date.Year }}.{{ .Date.Month }}.{{ .Date.Day }}")
	defer func() {
		os.Unsetenv(indexName.IdxVarName)
	}()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "invalid index template",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := indexName.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexName() = %v, want %v", got, tt.want)
			}
		})
	}
}
