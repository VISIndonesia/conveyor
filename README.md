# Conveyor

Conveyor is an economic way of using BigQuery as warehouse of near real time events. It reads and writes events via Google Pubsub and then uploads the event files to BigQuery in a batch as a JSON file upload.


## Usage

### Publish

```go
	entity := "transaction"
	data := struct {
		ID   string `json:"ID"`
		From string `json:"from"`
	}{
		ID:   "20",
		From: "Mantesh",
	}
	topicID := "debicred"
	projectID := "edgenetworks-150209"
	credFilePath := "./debicred.json"
	err := conveyor.Publish(topicID, entity, data, projectID, credFilePath)

```

### Subscribe and Upload

#### configs
- Place the config file in the folder from where the app is triggered, as config.yaml    
OR
- Place the config file this way - $HOME/.conveyor/config.yaml

Run consumer
```
go run cmd/subscriber/periodic.go
```
