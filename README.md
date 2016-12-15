# Overview

IBM MQ Trigger Monitor for AWS Lambda is a small component that glues together IBM MQ messaging infrastructure and Amazon AWS Lambda. It is intended to provide additional integration and deployment options for enterpise/corporate IT environments that gradually deploy new workloads on AWS. A typical requirement in such settings is pushing near-real time events generated inside the corporate network to consumers on AWS. Since IBM MQ is ubiquitous in corporate environments, a natural approach would be utilizing MQ messaging backbone to deliver the event messages. But MQ by itself will deliver the messages to an queue manager running on an EC2 instance, and it falls with the application developers and architects to provide "the last mile" solution - picking up messages from the queues and pushing them to the intended destinations.

The question then remains how those applications are to be deployed and managed. Of course there are traditional options - running them alongside the queue manager in the same EC2 instance or employing dedicated instances, with or without containers. MQ Trigger Monitor for AWS Lambda opens an additional option of *serverless deployment* of these applications.

A broad view of possible deployment of MQ Trigger Monitor for AWS Lambda:

![MQ Lambda TM - deployment architecture](doc/mqlambdatm1.jpg)



