Example of how to use Go's `log/slog` in AWS Lambda so the logs are structured JSON and automatically include request/invocation context.

How Lambda logging works:

- Anything your fuction writes to stdout/stderr ends up in CloudWatch Logs.
- Lambda also emits its own platform/system logs; you can configure those to be JSON or text (separate from your app logs).

Test locally:

```
sam build
mkdir events
sam local generate-event apigateway aws-proxy > ./events/apigw.json
sam local invoke
```

Deploy:

```
sam build
sam deploy # uses samconfig.yaml (created with --guided)
```

Delete:

```
sam delete
```
