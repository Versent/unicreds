# unicreds

unicreds is currently a pretty faithful port of [credstash](https://github.com/fugue/credstash) to [golang](https://golang.org/). This utility enables secure storage secrets in [DynamoDB](https://aws.amazon.com/dynamodb/) using [KMS](https://aws.amazon.com/kms/) to encrypt and sign these Credentials. Access to these keys is controlled using [IAM](https://aws.amazon.com/iam/).

# disclaimer

This is currently a work in progress, things are progressing but it is a side project.

# why

The number one reason for this port is platform support, getting credstash running on Windows and some older versions of Redhat Enterprise is a pain. Golang is fantastic at enabling simple deployment of core tools across a range of platforms with very little friction.

In addition to this we have some ideas about how this tool can be expanded to support some interesting use cases we have internally.

That said we have learnt a lot from how credstash worked and aim to remain compatible with it in the future where possible.

# license

This code is Copyright (c) 2015 Versent and released under the MIT license. All rights not explicitly granted in the MIT license are reserved. See the included LICENSE.md file for more details.
