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
2. Assign the `credstash` alias to the key using the key id printed when you created the KMS key.
```
aws --region ap-southeast-2 --profile [yourawsprofile] kms create-alias --alias-name 'alias/credstash' --target-key-id "arn:aws:kms:ap-southeast-2:xxxx:key/xxxx-xxxx-xxxx-xxxx-xxxx"
```
3. Run unicreds setup to create the dynamodb table in your region, ensure you have your credentials configured using the [awscli](https://aws.amazon.com/cli/).
```
unicreds setup --region ap-southeast-2 --profile [yourawsprofile]
```
NOTE: It is really important to tune DynamoDB to your read and write requirements if your using unicreds with automation!

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

  get <credential>
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

```

# install

If your on OSX you can install unicreds using homebrew now!

```
brew tap versent/homebrew-taps
brew install unicreds
```

Otherwise grab an archive from the [github releases page](https://github.com/Versent/unicreds/releases).

# why

The number one reason for this port is platform support, getting credstash running on Windows and some older versions of Redhat Enterprise is a pain. Go enables deployment of tools across a range of platforms with very little friction.

In addition to this we have some ideas about how this tool can be expanded to support some interesting use cases we have internally.

That said we have learnt a lot from how credstash worked and aim to remain compatible with it in the future where possible.

# development

I use `scantest` to watch my code and run tests on save.

```
go get github.com/smartystreets/scantest
```

# todo

* Add the ability to filter list / getall results using DynamoDB filters, at the moment I just use `| grep blah`.
* Work on the output layout.
* Make it easier to import files

# license

This code is Copyright (c) 2015 Versent and released under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.
