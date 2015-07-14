GoSlurp
=======

This is meant as a side project to teach myself Go. This definitely needs a lot of work, and I probably picked something a little too complex for my first foray into learning Go. Any help or contributions would be greatly appreciated :) :heart: 

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
export_as: "json"
export_path: "output.json"
```

You can currently export any messages found into a json file. Work will probably continue on this, so that json values aren't duplicated if the tool is ran as a daemon

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
