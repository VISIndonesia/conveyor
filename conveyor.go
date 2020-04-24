package conveyor

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
)

func Consume(subID string, timeout time.Duration, projectID, credFile string) (map[string][][]byte, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Println("Conveyor: Error getting client: ", err)
		return nil, err
	}
	defer client.Close()

	sub := client.Subscription(subID)

	sub.ReceiveSettings.Synchronous = true

	// Receive messages for 5 seconds.
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cm := make(chan *pubsub.Message, 1)
	eventsMap := make(map[string][][]byte)
	go func() {
		for {
			select {
			case msg := <-cm:
				entity, ok := msg.Attributes["entity"]
				if !ok {
					log.Println("Warning: entity name missing. Data:" + string(msg.Data))
					continue
				}
				if _, ok := eventsMap[entity]; !ok {
					eventsMap[entity] = make([][]byte, 0)
				}
				eventsMap[entity] = append(eventsMap[entity], msg.Data)
				msg.Ack()
			case <-ctx.Done():
				return
			}
		}
	}()

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil && status.Code(err) != codes.Canceled {
		log.Println("Conveyor: Error receiving message:", err)
		return nil, err
	}
	close(cm)
	return eventsMap, nil
}

func UploadEvents(eventsMap map[string][][]byte, dataset, projectID, credFile string) []error {
	ctx := context.Background()
	errList := make([]error, 0)
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credFile))
	if err != nil {
		errList = append(errList, err)
		return errList
	}

	for eventType, eventList := range eventsMap {
		var events string
		for _, event := range eventList {
			events = events + string(event) + "\n"
		}
		b := bytes.NewBufferString(events)
		rs := bigquery.NewReaderSource(b)
		rs.SourceFormat = bigquery.JSON
		rs.IgnoreUnknownValues = true
		ds := client.Dataset(dataset)
		loader := ds.Table(eventType).LoaderFrom(rs)
		job, err := loader.Run(ctx)
		if err != nil {
			errList = append(errList, err)
			continue
		}
		status, err := job.Wait(ctx)
		if err != nil {
			errList = append(errList, err)
			continue
		}
		if status.Err() != nil {
			errList = append(errList, status.Err())
			continue
		}
	}
	return errList
}

func Publish(topicID, entity string, data interface{}, projectID, credFile string) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Println("Conveyor: Error marshaling data:", err)
		return err
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Println("Conveyor: Error getting client: ", err)
		return err
	}
	defer client.Close()
	topic := client.Topic(topicID)
	defer topic.Stop()
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(dataJson),
		Attributes: map[string]string{
			"entity": entity,
		},
	})
	_, err = result.Get(ctx)
	if err != nil {
		log.Println("Conveyor: Error publishing event", err)
		return err
	}
	return nil
}
