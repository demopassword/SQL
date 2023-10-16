package main

import (
        "encoding/json"
        "fmt"
        "log"
        "math/rand"
        "time"
        "strings"

        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/kinesis"
)

var (
        kdsName   = "taxi-stream"
        region    = "ap-northeast-2"
        clientKDS *kinesis.Kinesis
)

type TaxiData struct {
        ID               string  `json:"id"`
        VendorID         int     `json:"vendorId"`
        PickupDate       string  `json:"pickupDate"`
        DropoffDate      string  `json:"dropoffDate"`
        PassengerCount   int     `json:"passengerCount"`
        PickupLongitude  string  `json:"pickupLongitude"`
        PickupLatitude   string  `json:"pickupLatitude"`
        DropoffLongitude string  `json:"dropoffLongitude"`
        DropoffLatitude  string  `json:"dropoffLatitude"`
        StoreAndFwdFlag  int     `json:"storeAndFwdFlag"`
        GCDistance       int     `json:"gcDistance"`
        TripDuration     int     `json:"tripDuration"`
        GoogleDistance   int     `json:"googleDistance"`
        GoogleDuration   int     `json:"googleDuration"`
}

func init() {
        // AWS 세션 생성
        sess, err := session.NewSession(&aws.Config{
                Region: aws.String(region),
        })
        if err != nil {
                log.Fatal("Failed to create AWS session: ", err)
        }

        // Kinesis 클라이언트 생성
        clientKDS = kinesis.New(sess)
}

func getLatLon() string {
        latLonData := []string{
                "-73.98174286,40.71915817", "-73.98508453,40.74716568", "-73.97333527,40.76407242", "-73.99310303,40.75263214",
                "-73.98229218,40.75133133", "-73.96527863,40.80104065", "-73.97010803,40.75979996", "-73.99373627,40.74176025", "-73.98544312,40.73571014",
                "-73.97686005,40.68337631", "-73.9697876,40.75758362", "-73.99397278,40.74086761", "-74.00531769,40.72866058", "-73.99013519,40.74885178",
                "-73.9595108,40.76280975", "-73.99025726,40.73703384", "-73.99495697,40.745121", "-73.93579865,40.70730972", "-73.99046326,40.75100708",
                "-73.9536438,40.77526093", "-73.98226166,40.75159073", "-73.98831177,40.72318649", "-73.97222137,40.67683029", "-73.98626709,40.73276901",
                "-73.97852325,40.78910065", "-73.97612,40.74908066", "-73.98240662,40.73148727", "-73.98776245,40.75037384", "-73.97187042,40.75840378",
                "-73.87303925,40.77410507", "-73.9921875,40.73451996", "-73.98435974,40.74898529", "-73.98092651,40.74196243", "-74.00701904,40.72573853",
                "-74.00798798,40.74022675", "-73.99419403,40.74555969", "-73.97737885,40.75883865", "-73.97051239,40.79664993", "-73.97693634,40.7599144",
                "-73.99306488,40.73812866", "-74.00775146,40.74528885", "-73.98532867,40.74198914", "-73.99037933,40.76152802", "-73.98442078,40.74978638",
                "-73.99173737,40.75437927", "-73.96742249,40.78820801", "-73.97813416,40.72935867", "-73.97171021,40.75943375", "-74.00737,40.7431221",
                "-73.99498749,40.75517654", "-73.91600037,40.74634933", "-73.99924469,40.72764587", "-73.98488617,40.73621368", "-73.98627472,40.74737167",
        }

        rand.Seed(time.Now().UnixNano())
        randomNum := rand.Intn(len(latLonData))

        return latLonData[randomNum]
}

func getStore() int {
        taxiOptions := []int{0, 1}
        rand.Seed(time.Now().UnixNano())
        randomIndex := rand.Intn(len(taxiOptions))

        return taxiOptions[randomIndex]
}

func sendDataToKinesis() {
        for {
                taxiData := TaxiData{
                        ID:               fmt.Sprintf("id%d", rand.Intn(8888888-1665586+1)+1665586),
                        VendorID:         rand.Intn(2) + 1,
                        PickupDate:       time.Now().Format(time.RFC3339),
                        DropoffDate:      time.Now().Add(time.Minute * time.Duration(rand.Intn(100-30+1)+30)).Format(time.RFC3339),
                        PassengerCount:   rand.Intn(9) + 1,
                        PickupLongitude:  "",
                        PickupLatitude:   "",
                        DropoffLongitude: "",
                        DropoffLatitude:  "",
                        StoreAndFwdFlag:  getStore(),
                        GCDistance:       rand.Intn(7) + 1,
                        TripDuration:     rand.Intn(10000-8+1) + 8,
                        GoogleDistance:   0,
                        GoogleDuration:   0,
                }

                // 랜덤한 위도와 경도 가져오기
                pickupLocation := getLatLon()
                dropoffLocation := getLatLon()

                pickupLocationSplit := strings.Split(pickupLocation, ",")
                dropoffLocationSplit := strings.Split(dropoffLocation, ",")

                taxiData.PickupLongitude = pickupLocationSplit[0]
                taxiData.PickupLatitude = pickupLocationSplit[1]
                taxiData.DropoffLongitude = dropoffLocationSplit[0]
                taxiData.DropoffLatitude = dropoffLocationSplit[1]

                // Kinesis 스트림에 데이터 전송
                data, err := json.Marshal(taxiData)
                if err != nil {
                        log.Println("Failed to marshal JSON:", err)
                        continue
                }

                putRecordInput := &kinesis.PutRecordInput{
                        Data:         data,
                        StreamName:   aws.String(kdsName),
                        PartitionKey: aws.String(taxiData.ID),
                }

                _, err = clientKDS.PutRecord(putRecordInput)
                if err != nil {
                        log.Println("Failed to put record in Kinesis:", err)
                        continue
                }

                log.Printf("Taxi data sent to Kinesis: ID %s\n", taxiData.ID)

                // 1초마다 데이터 전송
                time.Sleep(1 * time.Second)
        }
}

func main() {
        go sendDataToKinesis()

        // 메인 스레드는 계속 실행됩니다.
        select {}
}
