This project is mostly about AWS integrations. The scripts are used to run tests that require setup teardown

To manualy test the iam assume role


```bash

export the credentials for an account that is created wiht no 
# source  an assume role convenience function (awsCreds)
source <(curl -sS "https://raw.githubusercontent.com/natemarks/pipeline-scripts/v0.0.9/scripts/utility.sh")
awsCreds 709310380790 iam_assume_role test_iam_assume_role

```