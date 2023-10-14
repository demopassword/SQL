#windows

%flink.ssql(type=update)

SELECT 
    CAST(TO_TIMESTAMP(`timestamp`, 'yyyy/MM/dd HH:mm:ss') AS TIMESTAMP) AS window_start,
    path,
    COUNT(path) AS path_cnt
FROM 
    kinesis_stream
GROUP BY 
    CAST(TO_TIMESTAMP(`timestamp`, 'yyyy/MM/dd HH:mm:ss') AS TIMESTAMP), 
    path;
