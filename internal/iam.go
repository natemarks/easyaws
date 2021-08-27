package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog"
)

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
