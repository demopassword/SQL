create kinesis analytics table
```sql
CREATE TABLE log_table (
  `timestamp` STRING,
  `status` INT,
  `response_time` STRING,
  `ip` STRING,
  `method` STRING,
  `path` STRING
)
WITH (
  'connector' = 'kinesis',
  'stream' = 'wscc-kinesis-input-stream',
  'aws.region' = 'ap-northeast-2',
  'scan.stream.initpos' = 'LATEST',
  'format' = 'json',
  'json.timestamp-format.standard' = 'ISO-8601'
);
```

example query
```
%flink.ssql(type=update)
SELECT `timestamp` as request_time, status, ip, `method` as client_method, path FROM log_table;
```
![image](https://github.com/demopassword/SQL/assets/145639874/4c800f1f-50d4-49d3-ad45-8979ae10dce6)


---

### CONCAT 합치는 함수
%flink.ssql(type=update)

SELECT
    CONCAT(SPLIT_INDEX(log, ' ', 1), ' ', SPLIT_INDEX(log, ' ', 2), ' ', SPLIT_INDEX(log, ' ', 3)) as time_request
FROM log_table
GROUP BY
    CONCAT(SPLIT_INDEX(log, ' ', 1), ' ', SPLIT_INDEX(log, ' ', 2), ' ', SPLIT_INDEX(log, ' ', 3))

![image](https://github.com/demopassword/SQL/assets/145639874/9058f0ed-430e-4375-9b83-db71bebe8ff8)


### TRIM 공백제거 함수
TRIM(SPLIT_INDEX(log, '|', 3)) as ip

### SPLIT_INDEX 특정 문자열 기준으로 나누는 함수
레코드가 이런식으로 존재한다고 가정
```json
{"log":"[GIN] 2023/10/03 - 12:43:17 | 200 |      41.797\u00b5s |     10.20.14.80 | GET      \"/health\""}
```
```sql
SPLIT_INDEX(log, ' ', 1)
```

출력
```
2023/10/03
```
