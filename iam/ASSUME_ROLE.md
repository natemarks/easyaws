I want to be able to use assume role in golang  to minimize passing assumed role creds through shells/shellscripts in pipelines

Testing:
SETUP:
 - resources should be named/tagged 'DELETEME'
 - create an assumable role  that grants parameter store access
 - create user, group and role that allows the user to assume the assumable role



Try to create  a log with a test user. this should fail
use AssumeRole
Try agian. this should work

TEARDOWN:
 - delete the IAM and parameter store resource

## TODO

throw-away go modules that write and read json data to parameter store with unrestricted rights
throw-away module that assumes role. use it fromt the parameter store test module

use terraform to create the test fixtures
create intrgration test scripts
## terraform create
 - IAM role  than can create and write to a parameter store
 - IAM role that can read from the created log group
 - IAM user, group, role that  
## terraform destroy