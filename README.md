# tfstate-lookup
Golang CLI for looking up what terraform resources are deployed to an AWS account. Useful for knowing what infrastructure we have deployed in terraform using terraform in AWS.

```
$ ./tfstate-lookup -h
CLI tool to query what terraform resources are deployed in an AWS account. Can be used to query any AWS account.

                Example:
                ./tfstate-lookup --account <account> --bucket_name <mytfbucket>

Usage:
  tfstate-lookup [flags]

Flags:
      --account string       AWS account to query terraform resources
      --bucket_name string   AWS bucket name where terraform state files are stored
  -h, --help                 help for tfstate-lookup
```

## Pre-requisities
This assumes you have valid AWS credentials to your account. Ensure that your ~/.aws/credentials is set up. Follow https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html to get started.

## Usage
Run the binary (example with *infrastructureci* account)
```
make build
cd bin
./tfstate-lookup --account <account> --bucket_name <mybucket>
```


