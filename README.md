# Conveyor

Conveyor is an economic way of using BigQuery as warehouse of near real time events. It reads and writes events via Google Pubsub and then uploads the event files to BigQuery in a batch as a JSON file upload.


## Usage

### Publish

```go
	entity := "transaction" // should match the name of the table in BigQuery
	data := struct {
		ID   string `json:"ID"` // json tag should match the field name in the above BigQuery table 
		From string `json:"from"`
	}{
		ID:   "20",
		From: "Mantesh",
	}
	topicID := "debicred"
	projectID := "example-project-150209"
	credFilePath := "./gcp-creds.json"
	err := conveyor.Publish(topicID, entity, data, projectID, credFilePath)

```

### Subscribe and Upload

#### configs
- Place the config file in the folder from where the app is triggered, as config.yaml    
OR
- Place the config file this way - $HOME/.conveyor/config.yaml

An example config can be viewed [here](https://github.com/VISIndonesia/conveyor/blob/master/config/sample-config.yaml)    

Run consumer
```
go run cmd/subscriber/periodic.go
```
