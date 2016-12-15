# Overview

IBM MQ Trigger Monitor for AWS Lambda (`mqlambdatm` for short) is a small component that glues together IBM MQ messaging infrastructure and Amazon AWS Lambda. It is intended to provide additional integration and deployment options for enterpise/corporate IT environments that gradually deploy new workloads on AWS. A typical requirement in such settings is pushing near-real time events generated inside the corporate network to consumers on AWS. Since IBM MQ is ubiquitous in corporate environments, a natural approach would be utilizing MQ messaging backbone to deliver the event messages. But MQ by itself will deliver the messages to a queue manager running on an EC2 instance, and it falls with the application developers and architects to provide "the last mile" solution - picking up messages from the queues and pushing them to the intended destinations.

The question then remains how those applications are to be deployed and managed. Of course there are traditional options - running them alongside the queue manager in the same EC2 instance or employing dedicated instances, with or without containers. `mqlambdatm` opens an additional option of *serverless deployment* of these applications.

A broad view of a possible deployment of MQ Trigger Monitor for AWS Lambda:

![MQ Lambda TM - deployment architecture](doc/mqlambdatm1.jpg)

This document currently does not provide explanations on IBM MQ triggering mechanism. You can find them in [IBM MQ documentation on triggering](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.0.0/com.ibm.mq.dev.doc/q026910_.htm)

# Obtaining MQ Trigger Monitor for AWS Lambda

You need Go development environment installed. See [here](https://golang.org/doc/install) for the details. (You don't need anything besides MQ Server installation to run `mqlambdatm`, but since it is not distributed in binary form you need to build it first.)

The provided instructions are for Linux. There is no reason it shouldn't work on Windows but it hasn't been tested.

After installing Go tools and setting the environment variables, run

``` 
go get github.com/yurist/mqlambdatm
```

You will find `mqlambdatm` executable in $GOPATH/bin.

# Configuring and running

TBD - sample

# Developing Lambda functions for use with MQ triggering

TBD - sample

# License

Apache License v.2.0
