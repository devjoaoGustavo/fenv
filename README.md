# fenv

Fenv provides a simple way to fetch parameters from AWS Simple Systems Manager (SSM)

  ```shell
  $ fenv CoIsA

  // COISA = trem
  // OUTRA_COISA = 42
  ```

By default, the value of AWS_REGION is `us-east-2`, but it can be changed by prefixing the command like this:

  ```shell
  $ AWS_REGION=us-east-1 fenv my_var_name

  // MY_VAR_NAME = coisa
  // MY_VAR_NAME_OLD = trem
  // THIS_IS_MY_VAR_NAME = treco
  ```

Or exporting the environment variable

  ```shell
  $ export AWS_REGION=us-east-1
  $ fenv my_var_name

  // MY_VAR_NAME = coisa
  // MY_VAR_NAME_OLD = trem
  // THIS_IS_MY_VAR_NAME = treco
  ```

In order to access SSM, it is expected that AWS credentials are in place, or that AWS environment variables are set

```shell
$ cat ~/.aws/credentials

[default]
aws_access_key_id = <your access key id goes here>
aws_secret_access_key = <your secret access key goes here
```

In case you use `assume-role` in order to grant youself access to SSM and other AWS services, you might do something like this

```shell
unset AWS_ACCESS_KEY_ID
unset AWS_SECRET_ACCESS_KEY
unset AWS_SESSION_TOKEN

aws_credentials_json=$(aws sts assume-role --role-arn $ROLE_ARN --role-session-name $SESSION --region $AWS_REGION --duration-seconds 86400)

export AWS_ACCESS_KEY_ID=$(echo "$aws_credentials_json" | jq --exit-status --raw-output .Credentials.AccessKeyId)
export AWS_SECRET_ACCESS_KEY=$(echo "$aws_credentials_json" | jq --exit-status --raw-output .Credentials.SecretAccessKey)
export AWS_SESSION_TOKEN=$(echo "$aws_credentials_json" | jq --exit-status --raw-output .Credentials.SessionToken)
```

## Installation

The easiest way to install this program is using `go install` command.
```
go install github.com/devjoaoGustavo/fenv
```
Or you can download the binary directly from release page [here](https://github.com/devjoaoGustavo/fenv/releases/latest)
