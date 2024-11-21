
# AWST
A simple CLI to fetch and display data from AWS in a human readable format.

## Logs

### Examples
Search for log groups which name contains `lambda`, fetch all logs from 12 hourse and 40 minutes
ago, and start a live tail
```
awst logs search -e lambda -f 12h50s --all --tail
```
