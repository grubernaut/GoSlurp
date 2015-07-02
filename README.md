GoSlurp
=======

Basic utility to slurp messages off of an SQS queue and place them into a database.

Goal: To be a fully configurable, lightweight, SQS polling service

## Config
Config is unmarshalled from a YAML file, currently called `config.yml`

Sample Config:
```yaml
---
region: "<aws_region>"
queue_url: "<sqs_queue url>"
message_attributes:
  - "Attribute"
  - "Another Attribute"
```

Currently config only supports String attributes

## TODO

* Create a daemon from CLI utility
* Expand SQS Configuration ability
* Debug why 10% of the time, no data is returned, yet no error thrown
* Expand Message Attribute Configuration, to match DataType of Attributes
* Only read Messages from top of queue (Or add configuration to read from bottom)
* Allow for deletion of messages once sent to storage
* Split single package into multiple packages for easier writing and expandability
* Add Database Output + Config for:
  * RDS
  * MySQL
  * DynamoDB
  * Postgres
* Better Documentation
* Test Suite
