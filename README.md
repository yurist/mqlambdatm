# Overview

Trigger Monitor for IBM MQ and AWS Lambda (`mqlambdatm` for short) is a component that glues together IBM MQ messaging infrastructure and Amazon AWS Lambda. It is intended to provide additional integration and deployment options for enterpise/corporate IT environments that gradually deploy new workloads on AWS. A typical requirement in such settings is pushing near-real time events generated inside the corporate network to consumers on AWS. Since IBM MQ is ubiquitous in corporate environments, a natural approach would be utilizing MQ messaging backbone to deliver the event messages. But MQ by itself will deliver the messages to a queue manager running on an EC2 instance, and it falls with the application developers and architects to provide "the last mile" solution - picking up messages from the queues and pushing them to the intended destinations.

The question then remains how those applications are to be deployed and managed. Of course there are traditional options - running them alongside the queue manager in the same EC2 instance or employing dedicated instances, with or without containers. `mqlambdatm` opens an additional option of *serverless deployment* of these applications.

A broad view of a possible deployment of MQ Trigger Monitor for AWS Lambda:

![MQ Lambda TM - deployment architecture](doc/mqlambdatm1.jpg)

This document currently does not provide explanations on IBM MQ triggering mechanism. You can find them in [IBM MQ documentation on triggering](https://www.ibm.com/support/knowledgecenter/SSFKSJ_9.0.0/com.ibm.mq.dev.doc/q026910_.htm)

All the development and testing has been done on Linux. There are no specific distribution or version dependencies besides the gclib version. Check the latest release regarding which version of gclib it was built on, and whether it is compatible with the Linux kernel and gclib you need it to run on. You can build the executable yourself from the sources on the necessary Linux level, see below. Check out [this container](https://gitbub.com/yurist/mqgobuild_container) as an basis or example of doing it in a Docker container.

# Obtaining Trigger Monitor for IBM MQ and AWS Lambda

Download the latest release [here](https://github.com/yurist/mqlambdatm/releases/latest). Download `mqlambdatm.tar.gz` and unzip it using `tar -zxvf`. Copy the executable `mqlambdatm` to any directory according to your conventions.

# Running 

    Usage of mqlambdatm:
    -log-level string
            log level (DEBUG, INFO, WARN, ERROR, FATAL, PANIC) (default "info")
    -m string
            queue manager to connect, default queue manager if not given
    -q string
            initiation queue to serve

The recommended way of running trigger monitors is by using MQ services, see [example configuration](https://github.com/yurist/mqlambdatm-demo/blob/master/demo/cloud-init/docker/config.mqsc#L27-L36) in MQ Lambda Trigger Monitor sample repository.

# Building from sources

The machine for building `mqlambdatm` must have IBM MQ SDK installed in the default location.

You need Go development environment. See [here](https://golang.org/doc/install) for the details. You also need `gcc` compiler.

After installing Go tools and setting the environment variables, run

``` 
go get github.com/yurist/mqlambdatm
```

You will find `mqlambdatm` executable in $GOPATH/bin.

# Developing Lambda functions for use with MQ triggering

See [MQ Lambda Trigger Monitor Sample project](https://github.com/yurist/mqlambdatm-demo) for a complete sample, including a CloudFormation template with Lambda configuration and related resources.

*Note:* Only Lambda functions in Java have been tested. Other languages require IBM MQ native clients, which are way more difficult to run in the (mostly undocumented) Lambda container environment. Also, getting IBM support for such a configuration will be probably hard if not impossible. The Java clients on the other hand do not presume virtually anything about the environment and can be run anywhere JVM is supported. (This should not be construed as a formal support statement which I'm not in a position to provide for IBM anyway. Clarify with IBM support what their position is regarding containers, Lambda and MQ Java clients if you need such a formal statement.)

# Use of IBM software

This project uses open source [IBM library](https://github.com/ibm-messaging/mq-golang) that provides Go MQI support. Other than that, no IBM licensed software is incorporated in this repository. To use this project, as well as the accompanying [sample project](https://github.com/yurist/mqlambdatm-demo), you need to acquire necessary licenses for IBM MQ server and clients.

# Legalese

This is my personal open source contribution, unrelated to any other party, including IBM and Amazon. It is provided in good faith, in the hope of being useful, but with no warranty of any kind whatsoever.

The product names and other terms in this document and in the source code may be trademarks of IBM, Amazon or other parties. All such usage is acknowledged.

# Support

I will be glad to provide clarifications and fixes on the best effort basis. Please [open an issue](https://github.com/yurist/mqlambdatm/issues) if you need help.

# License

[Apache License, Version 2.0](http://apache.org/licenses/LICENSE-2.0.html)
