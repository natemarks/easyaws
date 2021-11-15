package secrets

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog"
)

type GetRemoteCredentialsInput struct {
	// The AWS SecretManager secret identifier. ex: /path/to/y/secret
	AWSSMSecretID string
	// The JSON key to look up
	JSONKey string
}

// GetRemoteCredentialsOutput The output includes the credentiols and the sha256
// some of the username, toke and the AWS SM secret document. the sha256sum makes
// it easier to test and troubleshoot bad values without compromising the secure
// data
type GetRemoteCredentialsOutput struct {
	//The sha256sum of the AWS SM secret json document
	SecretSha256sum string
	// Value retrieved using the JSON key
	JSONValue string
	// Sha256sum of the  retrieved using the JSON key
	JSONValueSha256sum string
}

func Sha256sum(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

// GetRemoteCredentials returns the remote credentials and sha256sums
func GetRemoteCredentials(i GetRemoteCredentialsInput, log *zerolog.Logger) (GetRemoteCredentialsOutput, error) {

	// Setup the client
	log.Info().Msg("setting up the AWS Secret Manager client")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err)
	}

	SecretClient := *secretsmanager.NewFromConfig(cfg)

	SecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(i.AWSSMSecretID),
		VersionId:    nil,
		VersionStage: nil,
	}
	// Get the secret doc from AWS

	log.Info().Msg("getting the secret doc from AWS SM")
	secretDoc, err := SecretClient.GetSecretValue(context.TODO(), SecretInput)
	if err != nil {
		log.Fatal().Err(err)
	}

	// unmarshal the JSON secret doc into a map. If the structure isn't a map this will fail
	log.Info().Msg("Unmarshalling credentials from AWSSM secret doc")
	var objmap map[string]string
	err = json.Unmarshal([]byte(*secretDoc.SecretString), &objmap)
	if err != nil {
		log.Fatal().Err(err)
	}
	// Use the provided username and token key names to get the credential values
	JSONValue := objmap[i.JSONKey]
	result := GetRemoteCredentialsOutput{
		SecretSha256sum:   Sha256sum(*secretDoc.SecretString),
		JSONValue:          JSONValue,
		JSONValueSha256sum: Sha256sum(JSONValue),
	}
	log.Debug().Msgf("SecretJSON Document(sha256): %s", Sha256sum(result.SecretSha256sum))
	log.Debug().Msgf("JSONValue(sha256): %s", Sha256sum(result.JSONValue))
	return result, err
}
