https://mozillascience.github.io/working-open-workshop/contributing/



## Testing

### Unit tests

wil be run early and often using 
```bash
make test
```

### Integration tests

These are more cumbersome and require substantial fixture setup/teardown

```bash
make i_test
```
This target kicks off a flow like

-> Setup the project fixtures (ex. IAM account used by all tests)

----> Test 1: Setup the Test 1 fixtures
------> Run  Test 1
----> Test 1: Teardown the Test 1 fixtures

----> Test 1: Setup the Test 1 fixtures
------> Run  Test 1
----> Test 1: Teardown the Test 1 fixtures

...

-> Teardown the project fixtures (IAM account used by all tests)


#### Project Fixture Setup

Run the terraform module:  deployments/test_user
This creates a test user with a name based on the locals.default_tags_project (easyaws), so the IAM user is names "test_easyaws". It also creates credentials for that user and stores them in a secret based on the name of the IAM user (test_easyaws_credentials)

The username is used by each integration test. Test that use the project account add put this user into a group that's permitted to assume a role that grants permissions required for the test. 

To get the credentials, run:
```bash
aws secretsmanager get-secret-value --secret-id test_easyaws_credentials \
--query SecretString --output text
```

You can also source a convenience function from the pipeline-scripts project to automatically export the creds form the secret:
```bash
source <(curl -s "https://raw.githubusercontent.com/natemarks/pipeline-scripts/v0.0.10/scripts/utility.sh")
❯ credsFromSecretManager test_easyaws_credentials
❯ aws sts get-caller-identity
{
    "UserId": "AIDA2KJRRPL3DLKBX75QD",
    "Account": "709310380790",
    "Arn": "arn:aws:iam::709310380790:user/deleteme/users/test_easyaws"
}
```

#### Test Fixture Setup
The fixture setup, test execution and  teardown are all done in a single script. One example is scripts/test_iam_assume_role.sh.
NOTE: the name format of the script is important , becuase it's used to find the relevant terraform module in deployments/

The important capabilities of the script are:
 - run the terraform install/setup
 - run the relevant terraform modules
 - run the tests
 - run the terraform destroy. NOTE: the script template traps exits to make sure the terraform destroy always runs
