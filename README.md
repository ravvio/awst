
# AWST
A simple CLI to fetch and display data from AWS in a human readable format.

## Logs

### Examples
Search for log groups which name contains `lambda`, fetch all logs since April 12 2024,
until 1 week and 3 days ago, and start a live tail
```
awst logs search -e lambda --since 2024-04-12 --until 1w3d --all --tail
```
