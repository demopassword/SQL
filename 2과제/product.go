package main

import (
        "fmt"
        "log"
        "math/rand"
        "net/http"
        "os"
        "time"

        "github.com/gin-gonic/gin"
)

var products = []string{"Product A", "Product B", "Product C", "Product D", "Product E"}

func GetProductHandler(c *gin.Context) {
        rand.Seed(time.Now().UnixNano())
        randomProduct := products[rand.Intn(len(products))]

        c.JSON(http.StatusOK, gin.H{"product": randomProduct})
}

func Logger() gin.HandlerFunc {
        return func(c *gin.Context) {
                // 현재 시간과 클라이언트 IP 주소를 포함하여 접속 로그 생성
                logEntry := fmt.Sprintf("%s - [%s] \"%s %s %s\" %d %s",
                        c.ClientIP(),
                        time.Now().Format("02/Jan/2006:15:04:05 -0700"),
                        c.Request.Method,
                        c.Request.RequestURI,
                        c.Request.Proto,
                        c.Writer.Status(),
                        c.Request.UserAgent(),
                )

                // 로그를 "access.log" 파일에 추가
                file, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                if err != nil {
                        log.Fatal(err)
                }
                defer file.Close()
                log.SetOutput(file)

                // 접속 로그 기록
                log.Println(logEntry)

                // 요청 처리 계속
                c.Next()
        }
}

func HealthCheckHandler(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
        r := gin.Default()

        // 접속 로그 미들웨어 추가
        r.Use(Logger())

        r.GET("/v1/product", GetProductHandler)
        r.GET("/healthcheck", HealthCheckHandler) // /healthcheck 엔드포인트 추가

        fmt.Println("Server is running on :8080")
        r.Run(":8080")
}
