## 분단위로 잘라서 카운트 출력
```
SELECT
    DATE_TRUNC('minute', timestamp) AS my_time,
    count(*)
FROM
    "output"
GROUP BY
    DATE_TRUNC('minute', timestamp)
```
![image](https://github.com/demopassword/SQL/assets/145639874/e6ee9058-d481-4c35-99eb-fd5396da2d19)
