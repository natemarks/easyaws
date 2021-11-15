package secrets

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

// IsValid checks the GetRemoteCredentialsOutput and returns true if it looks ok
func IsValid(output GetRemoteCredentialsOutput) bool {
	if len(output.SecretSha256sum) != 64 {
		return false
	}
	if len(output.JSONValueSha256sum) != 64 {
		return false
	}
	if output.JSONValue == "" {
		return false
	}
	return true
}

func IsValidJSON(input string) bool {
	var objmap map[string]string
	err := json.Unmarshal([]byte(input), &objmap)
	if err != nil {
		return false
	}
	return true
}

// TestGetRemoteCredentials This test gets credentials based on  a secret manager config proided through env vars
// It executes a loose check on the returned output
func TestGetRemoteCredentials(t *testing.T) {
	logger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		i   GetRemoteCredentialsInput
		log *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid", args: args{
			i: GetRemoteCredentialsInput{
				AWSSMSecretID: os.Getenv("AWSSMSECRETID"),
				JSONKey:       os.Getenv("JSONKEY"),
			},
			log: &logger,
		},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteCredentials(tt.args.i, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !IsValid(got) {
				t.Errorf("GetRemoteCredentials() output doesn't look right")
			}
		})
	}
}

func TestGetSecretJSON(t *testing.T) {
	logger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		i   GetSecretJSONInput
		log *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "valid", args: args{
			i: GetSecretJSONInput{
				AWSSMSecretID: os.Getenv("AWSSMSECRETID"),
			},
			log: &logger,
		},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSecretJSON(tt.args.i, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecretJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !IsValidJSON(got) {
				t.Errorf("GetSecretJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupJSONKey(t *testing.T) {
	type args struct {
		key     string
		JSONDoc string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "valid", args: args{
			key:     "ggg",
			JSONDoc: "{\"ggg\": \"ooo\"}",
		},wantErr: false,want: "ooo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LookupJSONKey(tt.args.key, tt.args.JSONDoc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupJSONKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LookupJSONKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}