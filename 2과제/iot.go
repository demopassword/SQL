package main

import (
        "encoding/json"
        "fmt"
        "math/rand"
        "time"

        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/kinesis"
)

const (
        streamName = "iot-stream"
        awsRegion  = "ap-northeast-2" // 본인의 AWS 지역에 맞게 설정
)

type SensorData struct {
        SensorID           int       `json:"sensor_id"`
        CurrentTemperature float64   `json:"current_temperature"`
        Status             string    `json:"status"`
        EventTime          time.Time `json:"event_time"`
}

func getRandomData() SensorData {
        currentTemperature := 10 + rand.Float64()*170
        var status string

        if currentTemperature > 160 {
                status = "ERROR"
        } else if currentTemperature > 140 || rand.Intn(100) > 80 {
                status = []string{"WARNING", "ERROR"}[rand.Intn(2)]
        } else {
                status = "OK"
        }

        return SensorData{
                SensorID:           rand.Intn(100) + 1,
                CurrentTemperature: currentTemperature,
                Status:             status,
                EventTime:          time.Now(),
        }
}

func sendDataToKinesis(streamName string, kinesisClient *kinesis.Kinesis) {
        for {
                data := getRandomData()
                partitionKey := fmt.Sprintf("%d", data.SensorID)

                dataBytes, err := json.Marshal(data)
                if err != nil {
                        fmt.Println("Failed to marshal data:", err)
                        continue
                }

                putRecordInput := &kinesis.PutRecordInput{
                        StreamName:   aws.String(streamName),
                        Data:         dataBytes,
                        PartitionKey: aws.String(partitionKey),
                }

                _, err = kinesisClient.PutRecord(putRecordInput)
                if err != nil {
                        fmt.Println("Failed to put record in Kinesis:", err)
                        continue
                }

                fmt.Printf("Sent data to Kinesis: SensorID %d\n", data.SensorID)

                time.Sleep(time.Second)
        }
}

func main() {
        // AWS 지역(region) 설정 추가
        session, err := session.NewSession(&aws.Config{
                Region: aws.String(awsRegion),
        })
        if err != nil {
                fmt.Println("Failed to create AWS session:", err)
                return
        }

        kinesisClient := kinesis.New(session)
        sendDataToKinesis(streamName, kinesisClient)
}
