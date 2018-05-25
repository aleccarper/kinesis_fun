# kinesis_fun

go build -o main producer/main.go  && zip main.zip main && sam local start-api --env-vars env.json

go run jabberjaw/main.go

AWS_ID=MYID AWS_TOKEN=MYTOKEN go run consumer/main.go
