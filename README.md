[![Build Status](https://travis-ci.org/Versent/unicreds.svg?branch=master)](https://travis-ci.org/Versent/unicreds)

# unicreds

Unicreds is a command line tool to manage secrets within an AWS account, the aim is to keep securely stored
with your systems and data so you don't have to manage them externally. It uses [DynamoDB](https://aws.amazon.com/dynamodb/) and [KMS](https://aws.amazon.com/kms/) to store and
encrypt these secrets. Access to these keys is controlled using [IAM](https://aws.amazon.com/iam/).

Unicreds is written in [Go](https://golang.org/) and is based on [credstash](https://github.com/fugue/credstash).

# setup

1. Create a KMS key in IAM, using an aws profile you have configured in the aws CLI. You can ommit `--profile` if you use the Default profile.
```
aws --region ap-southeast-2 --profile [yourawsprofile] kms create-key --query 'KeyMetadata.KeyId'
```
**Note:** You will also need to assign permission to users other than the root account to access and use the key see [How to Help Protect Sensitive Data with AWS KMS](https://blogs.aws.amazon.com/security/post/Tx79IILINW04DC/How-to-Help-Protect-Sensitive-Data-with-AWS-KMS).
2. Assign the `credstash` alias to the key using the key id printed when you created the KMS key.
```
aws --region ap-southeast-2 --profile [yourawsprofile] kms create-alias --alias-name 'alias/credstash' --target-key-id "xxxx-xxxx-xxxx-xxxx-xxxx"
```
3. Run unicreds setup to create the dynamodb table in your region, ensure you have your credentials configured using the [awscli](https://aws.amazon.com/cli/).
```
unicreds setup --region ap-southeast-2 --profile [yourawsprofile]
```
**Note:** It is really important to tune DynamoDB to your read and write requirements if you're using unicreds with automation!

# demo

To illustrate how unicreds works I made a screen recording of list, put, get and delete.

![Image of screencast](docs/images/unicreds_recording.gif)

# usage

```
usage: unicreds [<flags>] <command> [<args> ...]

A credential/secret storage command line tool.

Flags:
      --help                     Show context-sensitive help (also try --help-long and
                                 --help-man).
  -c, --csv                      Enable csv output for table data.
  -d, --debug                    Enable debug mode.
  -j, --json                     Output results in JSON
  -r, --region=REGION            Configure the AWS region
  -p, --profile=PROFILE          Configure the AWS profile
  -R, --role=ROLE                Specify an AWS role ARN to assume
  -t, --table="credential-store"
                                 DynamoDB table.
  -k, --alias="alias/credstash"  KMS key alias.
  -E, --enc-context=ENC-CONTEXT ...
                                 Add a key value pair to the encryption context.
      --version                  Show application version.

Commands:
  help [<command>...]
    Show help.

  setup
    Setup the dynamodb table used to store credentials.

  get <credential> [<version>]
    Get a credential from the store.

  getall [<flags>]
    Get latest credentials from the store.

  list [<flags>]
    List latest credentials with names and version.

  put <credential> <value> [<version>]
    Put a credential into the store.

  put-file <credential> <value> [<version>]
    Put a credential from a file into the store.

  delete <credential>
    Delete a credential from the store.

  exec <command>...
    Execute a command with all secrets loaded as environment variables.
```

Unicreds supports the `AWS_*` environment variables, and configuration in `~/.aws/credentials` and `~/.aws/config`

# examples

* List secrets using default profile:
```
$ unicreds list
```

* List secrets using the default profile, in a different region:
```
$ unicreds -r us-east-2 list
$ AWS_REGION=us-east-2 unicreds list
```

* List secrets using profile MYPROFILE in `~/.aws/credentials` or `~/.aws/config`
```
$ unicreds -r us-west-2 -p MYPROFILE list
$ AWS_PROFILE=MYPROFILE unicreds list
```

* List secrets using a profile, but also assuming a role:
```
$ unicreds -r us-west-2 -p MYPROFILE -R arn:aws:iam::123456789012:role/MYROLE list
```

* Store a login for `test123` from unicreds using the encryption context feature.
```
$ unicreds -r us-west-2 put test123 -E 'stack:123' testingsup
   • stored                    name=test123 version=0000000000000000001
```

* Retrieve a login for `test123` from unicreds using the encryption context feature.
```
$ unicreds -r us-west-2 get test123 -E 'stack:123'
testingsup
```

* Example of a failed encryption context check.
```
$ unicreds -r us-west-2 get test123 -E 'stack:12'
   ⨯ failed                    error=InvalidCiphertextException:
	status code: 400, request id: 0fed8a0b-5ea1-11e6-b359-fd8168c3c784
```

* Execute `env` command, all secrets are loaded as environment variables.
```
$ unicreds -r us-west-2 exec -- env
```

# references

* [How to Protect the Integrity of Your Encrypted Data by Using AWS Key Management Service and EncryptionContext](https://blogs.aws.amazon.com/security/post/Tx2LZ6WBJJANTNW/How-to-Protect-the-Integrity-of-Your-Encrypted-Data-by-Using-AWS-Key-Management)

# install

If you're on OSX you can install unicreds using homebrew now!

```
brew tap versent/homebrew-taps
brew install unicreds
```

Otherwise grab an archive from the [github releases page](https://github.com/Versent/unicreds/releases).

# development

I use `scantest` to watch my code and run tests on save.

```
go get github.com/smartystreets/scantest
```

# testing
You can run unit tests which mock out the KMS and DynamoDB backend using `make test`.

There is an integration test you can run using `make integration`. You must set the `AWS_REGION` (default `us-west-2`), `UNICREDS_KEY_ALIAS` (default `alias/unicreds`), and `UNICREDS_TABLE_NAME` (default `credential-store`) environment variables to point to valid AWS resources.

# auto-versioning

If you've been using unicreds auto-versioning before September 2015, Unicreds had the [same](https://github.com/fugue/credstash/issues/51) [bug](https://github.com/Versent/unicreds/issues/34) as credstash when auto-versioning that results in a sorting error after ten versions. You should be able to run the [credstash-migrate-autoversion.py](https://github.com/fugue/credstash/blob/master/credstash-migrate-autoversion.py) script included in the root of the credstash repository to update your versions prior to using the latest version of unicreds.

# Docker ENTRYPOINT

It is possible to use `unicreds exec` as an [entrypoint](https://docs.docker.com/engine/reference/builder/#/entrypoint) for loading safely your secrets as environment variables inside your container in AWS ECS.

### Example
```
RUN curl -sL \
    https://github.com/Versent/unicreds/releases/download/v1.5.0/unicreds_1.5.0_linux_x86_64.tgz \
 | tar zx -C /usr/local/bin \
 && chmod +x /usr/local/bin/unicreds
ENTRYPOINT ["/usr/local/bin/unicreds", "exec", "--"]
```

With [IAM roles for Amazon ECS tasks](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html) you can give the necessary privileges to your container so that it can exploit `unicreds`.

# todo

* Add the ability to filter list / getall results using DynamoDB filters, at the moment I just use `| grep blah`.
* Work on the output layout.
* Make it easier to import files

# license

This code is Copyright (c) 2015 Versent and released under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.
