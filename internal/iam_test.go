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
	if !strings.HasPrefix(a.Resource, "user") {
		return errors.New("ARN doesn't seem to refer to a user object")
	}
	return nil
}

// ARNMatchesAssumedRole Check the expected assume role name
// Make sure the arn looks right based on the role name. the are names are prefixed, so
// partially random
// example:
// arn:aws:sts::709310380790:assumed-role/easyaws_iam_assume_role/aws-go-sdk-1630070768490628000
func ARNMatchesAssumedRole(arARN, assumedRoleName string) error {
	a, err := arn.Parse(arARN)
	if err != nil {
		return err
	}
	if a.Service != "sts" {
		return errors.New("ARN is not for the STS service")
	}
	if !strings.Contains(a.Resource, assumedRoleName) {
		msg := fmt.Sprintf("ARN doesn't contain assumed role name: %s", assumedRoleName)
		return errors.New(msg)
	}
	if !strings.HasPrefix(a.Resource, "assumed-role") {
		return errors.New("ARN doesn't seem to refer to an assumed role object")
	}
	return nil
}

// TestGetAWSIdentity Run STS get-caller-identity and compare it to the user
// account created by the project fixture setup
// (scripts/setup_project_fixtures.sh) The script that runs go test against this
// must change to the expected credentials before running go test in order for
// this test to succeed
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

func TestGetAssumeRoleCreds(t *testing.T) {
	myLogger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		assumeRoleARN string
		logger        *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid", args: args{
			assumeRoleARN: "arn:aws:iam::709310380790:role/easyaws_iam_assume_role",
			logger:        &myLogger,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAssumeRoleCreds(tt.args.assumeRoleARN, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAssumeRoleCreds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAWSIdentityWithAssumeRole(t *testing.T) {
	myLogger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		assumeRole string
		logger     *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid",
			args: args{
				assumeRole: "arn:aws:iam::709310380790:role/easyaws_iam_assume_role",
				logger:     &myLogger,
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAWSIdentityWithAssumeRole(tt.args.assumeRole, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAWSIdentityWithAssumeRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = ARNMatchesAssumedRole(got, "easyaws_iam_assume_role")
			if err != nil {
				t.Errorf("Returned ARN (%s) doesn't seem to match role name: %s", got, "easyaws_iam_assume_role")
			}
		})
	}
}
