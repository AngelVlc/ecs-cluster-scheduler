package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type eventDetail struct {
	ClusterName string `json:"clusterName"`
	ServiceName string `json:"serviceName"`
	Action      string `json:"action"`
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	log.Printf("Event detail: %v\n", string(event.Detail))

	detail, err := getEventDetail(event)
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("error loading the default config: %v", err)
	}

	input := ecs.UpdateServiceInput{
		Cluster:      aws.String(detail.ClusterName),
		Service:      aws.String(detail.ServiceName),
		DesiredCount: getDesiredCount(detail),
	}

	ecsClient := ecs.NewFromConfig(cfg)

	_, err = ecsClient.UpdateService(ctx, &input)
	if err != nil {
		return fmt.Errorf("error updating the service: %v", err)
	}

	log.Printf("Updated desired count to '%v' in service '%v' of cluster '%v'\n", *input.DesiredCount, detail.ServiceName, detail.ClusterName)

	return nil
}

func main() {
	lambda.Start(handler)
}

func getEventDetail(event events.CloudWatchEvent) (*eventDetail, error) {
	detail := eventDetail{}

	err := json.NewDecoder(ioutil.NopCloser(bytes.NewReader(event.Detail))).Decode(&detail)
	if err != nil {
		return nil, fmt.Errorf("error decoding the event detail: %v", err)
	}

	return &detail, nil
}

func getDesiredCount(detail *eventDetail) *int32 {
	if detail.Action == "start" {
		return aws.Int32(1)
	}

	return aws.Int32(0)
}
