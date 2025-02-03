
# AWST
A simple CLI to fetch and display data from AWS in a human readable format.

## Logs
Available subcommands are:
- *list* - list available log groups
- *get* - retrieve logs of a log group given its name
- *search* - retrieve logs of a list of logs groups from a prefix or pattern search

### Examples
Retrieve up to 100 logs of a `/ecs/example` log group since 11 hours ago
```
awst logs get /ecs/example --since 11h --limit 100
```

Search for all log groups which name contains `lambda`, retrieve all logs since
April 12 2024, until 1 week and 3 days ago, and start a live tail:
```
awst logs search -e lambda --since 2024-04-12 --until 1w3d --all --tail
```

## S3

## DynamoDB
Available subcommands are:
- *list* - retrieve list of DynamoDB tables

### Examples
Retrieve up to 20 tables:
```
awst ddb list --limit 20
```
