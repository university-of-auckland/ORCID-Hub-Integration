# Consumer simple message agnostic event consumer

a sample Lambda...
It can be triggered either by SQS or API Gateway directly.

## Event Message Flow
![ScreenShot](/flow.png?raw=true "Message Flow")


# Building

```
go build . && upx consumer && zip consumer.zip consumer
```
