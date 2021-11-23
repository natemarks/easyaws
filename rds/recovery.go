package rds

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog"
)

type Tag struct {
	Key string
	Value string
}

type RestoreSnapshotIdInput struct {
	// DBInstanceIdentifier: name of the instance created form the snapshot
	DBInstanceIdentifier string
	// DBSnapshotIdentifier: Snapshot ID to restore
	DBSnapshotIdentifier string
	// DBSubnetGroupName: Subnet group for the restored database
	DBSubnetGroupName string
	VpcSecurityGroupIds []string
	Tags []Tag
}

type RestoreSnapshotIdOutput struct {
	DBInstanceIdentifier string
}

func tagsToAWSTags( mytags []Tag) ([]types.Tag) {
	var res []types.Tag
	for _, v := range mytags {
		res = append(res,
			types.Tag{
			Key:   aws.String(v.Key),
			Value: aws.String(v.Value),
		})

	}
	return res
}

// filterLatestSnapshot Given a slice of DBSnapshots, return the onne wiht the latest InstanceCreateTime
func filterLatestSnapshot(ss []types.DBSnapshot) types.DBSnapshot {
	var latest types.DBSnapshot
	for i, v := range ss {
		if i == 0 {
			latest = v
			continue
		}
		if v.OriginalSnapshotCreateTime.After(*latest.OriginalSnapshotCreateTime){
			latest = v
		}
	}
	return latest
}

// GetLatestSnapshotId restore a given snapshot to a given instnace namd
func GetLatestSnapshotId(instance string, log *zerolog.Logger) (string, error) {

	// Setup the client
	log.Info().Msgf("looking up latest snaphot for RDS instance: %s", instance)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err)
	}

	RDSCLient := *rds.NewFromConfig(cfg)

	input := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(instance),
	}
	// Get the secret doc from AWS

	log.Info().Msg("getting the secret doc from AWS SM")
	ssOutput, err := RDSCLient.DescribeDBSnapshots(context.TODO(), input)
	if err != nil {
		log.Fatal().Err(err)
	}
	latest := filterLatestSnapshot(ssOutput.DBSnapshots)
	return *latest.DBSnapshotIdentifier, err

}

// GetSubnetGroup Given and instance ID, return the subnet group id
func GetSubnetGroup(instance string, log *zerolog.Logger) (string, error) {

	// Setup the client
	log.Info().Msgf("looking up latest snaphot for RDS instance: %s", instance)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err)
	}

	RDSCLient := *rds.NewFromConfig(cfg)

	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}
	// Get the secret doc from AWS

	log.Info().Msg("getting the secret doc from AWS SM")
	ssOutput, err := RDSCLient.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		log.Fatal().Err(err)
	}

	return *ssOutput.DBInstances[0].DBSubnetGroup.DBSubnetGroupName, err

}

// GetVPCSecurityGroups Given and instance ID, return the security groups
func GetVPCSecurityGroups(instance string, log *zerolog.Logger) ([]string, error) {
	var res []string
	// Setup the client
	log.Info().Msgf("looking up latest snaphot for RDS instance: %s", instance)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err)
	}

	RDSCLient := *rds.NewFromConfig(cfg)

	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}
	// Get the secret doc from AWS

	log.Info().Msg("getting the secret doc from AWS SM")
	ssOutput, err := RDSCLient.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		log.Fatal().Err(err)
	}
	for _, v := range ssOutput.DBInstances[0].VpcSecurityGroups{
		res = append(res, *v.VpcSecurityGroupId)
	}
	return res, err

}

// RestorePGSnapshotId Restore a snapshot to a enw RDS postgres instance
func RestorePGSnapshotId(input RestoreSnapshotIdInput, log *zerolog.Logger) (RestoreSnapshotIdOutput, error){
	// Setup the client
	log.Info().Msgf("Restoring instance: %s from snapshot ID: %s",
		input.DBInstanceIdentifier,
		input.DBSnapshotIdentifier)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err)
	}

	RDSCLient := *rds.NewFromConfig(cfg)

	rInput := rds.RestoreDBInstanceFromDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(input.DBSnapshotIdentifier),
		DBInstanceIdentifier: aws.String(input.DBInstanceIdentifier),
		Tags: tagsToAWSTags(input.Tags),
		DBSubnetGroupName: aws.String(input.DBSubnetGroupName),
		VpcSecurityGroupIds: input.VpcSecurityGroupIds,
	}
	rOutput, err := RDSCLient.RestoreDBInstanceFromDBSnapshot(context.TODO(), &rInput)
	if err != nil {
		log.Fatal().Err(err)
	}

	return RestoreSnapshotIdOutput{
		*rOutput.DBInstance.DBInstanceIdentifier,
	}, err
}