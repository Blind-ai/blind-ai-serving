## launch APIs

`go run cmd/AI/lungh/main.go`<br/>
`go run cmd/AI/skin/main.go`

Routes:

`/api/lungh/evaluate/image`<br/>
input, "file":file<br/>
output:
{"Type":"","Probability": float,"ProcessingTime":""}<br/><br/>
`/api/skin/evaluate/image`<br/>
input, "file":file<br/>
{"Type":"","Probability": float,"ProcessingTime":""}

## launch AIs

See readme.md in `./AIs/lung` and `./AIs/skin`
