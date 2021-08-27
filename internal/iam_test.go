// +build integration

package internal

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"testing"
)

// ARNMatchesUsername given an arn and an IAM username make sure the
// ARN looks right
func ARNMatchesUsername(iARN, iName string) error {
	a, err := arn.Parse(iARN)
	if err != nil {
		return err
	}
	if a.Service != "iam" {
		return errors.New("ARN is not for the IAM service")
	}
	if !strings.HasSuffix(a.Resource, iName) {
		msg := fmt.Sprintf("ARN doesn't match username: %s", iName)
		return errors.New(msg)
	}
	return nil
}
func TestGetAWSIdentity(t *testing.T) {
	myLogger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		logger *zerolog.Logger
	}
	tests := []struct {
		name        string
		args        args
		wantArn     string
		wantUserId  string
		wantAccount string
	}{
		{name: "valid", args: args{logger: &myLogger},
			wantAccount: "",
			wantArn:     "",
			wantUserId:  ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArn, _, _ := GetAWSIdentity(tt.args.logger)
			err := ARNMatchesUsername(gotArn, "test_easyaws")
			if err != nil {
				t.Error(err)
			}
		})
	}
}
