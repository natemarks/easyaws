// +build integration

package rds

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"reflect"
	"testing"
)

func TestGetLatestSnapshot(t *testing.T) {
	logger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		instance string
		log      *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid",args: args{
			instance: "some-instance-name",
			log:      &logger,
		},
		wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLatestSnapshotId(tt.args.instance, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRestorePGSnapshotId(t *testing.T) {
	logger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		input RestoreSnapshotIdInput
		log   *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid",
			args: args{
				input: RestoreSnapshotIdInput{
					DBInstanceIdentifier: "deleteme-some-instance-name",
					DBSnapshotIdentifier: "rds:some-instance-name-2021-11-23-02-15",
					DBSubnetGroupName: "some-subnet-group-name",
					VpcSecurityGroupIds: []string{
						"sg-111111111",
						"sg-2222222222"},
					Tags: []Tag{{
						Key:   "deleteme",
						Value: "true"},
						{
						Key: "deleteme_after",
						Value: "2021-11-23",
						},
					},
				},
				log:   &logger,
			},wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RestorePGSnapshotId(tt.args.input, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestorePGSnapshotId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetVPCSecurityGroups(t *testing.T) {
	logger := log.With().Str("test_key", "test_value").Logger()
	type args struct {
		instance string
		log      *zerolog.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "valid",args: args{
			instance: "some-instance-name",
			log:      &logger,
		},wantErr: false,
		want: []string{
			"sg-1111111111",
			"sg-2222222222"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVPCSecurityGroups(tt.args.instance, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVPCSecurityGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVPCSecurityGroups() got = %v, want %v", got, tt.want)
			}
		})
	}
}