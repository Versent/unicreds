[![Build Status](https://travis-ci.org/Versent/unicreds.svg?branch=master)](https://travis-ci.org/Versent/unicreds)

# unicreds

Unicreds is currently a pretty faithful port of [credstash](https://github.com/fugue/credstash) to [Go](https://golang.org/).

# overview

This command line utility automates the storage of encrypted secrets in [DynamoDB](https://aws.amazon.com/dynamodb/) using [KMS](https://aws.amazon.com/kms/) to encrypt and sign these Credentials. Access to these keys is controlled using [IAM](https://aws.amazon.com/iam/).

# setup

1. Add and configure a KMS key in IAM with the alias `credstash`, ensure this is created in the correct region as the user interface for this is quite confusing.
2. Run `unicreds setup` to create the dynamodb table in your region, ensure you have your credentials configured using the [awscli](https://aws.amazon.com/cli/).

# usage

```
usage: unicreds [<flags>] <command> [<args> ...]

A credential/secret storage command line tool.

Flags:
      --help                     Show context-sensitive help (also try --help-long and --help-man).
  -c, --csv                      Enable csv output for table data.
  -d, --debug                    Enable debug mode.
  -r, --region=REGION            Configure the AWS region
      --alias="alias/credstash"  KMS key alias.
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
