package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog"
)

// GetAWSIdentity Get the STS identity for the current session
func GetAWSIdentity(logger *zerolog.Logger) (Arn, UserId, Account string) {

	// get the aws sdk client config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sts.NewFromConfig(cfg)

	input := &sts.GetCallerIdentityInput{}

	p, err := client.GetCallerIdentity(context.TODO(), input)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return *p.Arn, *p.UserId, *p.Account
}

func GetAssumeRoleCreds(assumeRoleARN string, logger *zerolog.Logger) (*stscreds.AssumeRoleProvider, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatal().Err(err)
		return nil, err
	}
	stsSvc := sts.NewFromConfig(cfg)
	if err != nil {
		logger.Fatal().Err(err)
		return nil, err
	}
	creds := stscreds.NewAssumeRoleProvider(stsSvc, assumeRoleARN)
	return creds, err
}

// GetAWSIdentityWithAssumeRole Get the STS identity for the current session
func GetAWSIdentityWithAssumeRole(assumeRole string, logger *zerolog.Logger) (string, error) {

	// get the aws sdk client config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatal().Err(err)
		return "", err
	}
	creds, err := GetAssumeRoleCreds(assumeRole, logger)
	cfg.Credentials = aws.NewCredentialsCache(creds)

	// Create service client value configured for credentials
	// from assumed role.
	svc := sts.NewFromConfig(cfg)

	input := &sts.GetCallerIdentityInput{}

	p, err := svc.GetCallerIdentity(context.TODO(), input)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return *p.Arn, err
}
