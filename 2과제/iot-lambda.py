import base64
import json
from datetime import datetime

print('Loading function')


def lambda_handler(event, context):
    output = []

    for record in event['records']:
        print('record data: {}'.format(record['data'])) #
        
        print(record['recordId'])
        payload = base64.b64decode(record['data']).decode('utf-8')
        print('payload: {}'.format(payload))
        
        # Do custom processing on the payload here
        payload_dict = json.loads(payload)
        sensorid = payload_dict['sensor_id']
        temperature = payload_dict['current_temperature']
        status = payload_dict['status']
        eventtime = payload_dict['event_time']
        eventtime = eventtime.split(".")[0]
        eventtime = datetime.fromisoformat(eventtime)
        
        # Do custom processing on the payload here
        payload_dict['sensorid'] = sensorid
        payload_dict['temperature'] = int(temperature)
        payload_dict['year'] = eventtime.year
        payload_dict['month'] = eventtime.month
        payload_dict['day'] = eventtime.day
        payload_dict['hour'] = eventtime.hour
        payload_dict['minute'] = eventtime.minute
        payload_dict['second'] = eventtime.second
        
        # Do custom processing on the payload here
        partition_keys = {"sensorid": sensorid,
                          "year": eventtime.year,
                          "month": eventtime.month,
                          "day": eventtime.day,
                          "hour": eventtime.hour
                          }
        
        # Do custom processing on the payload here
        payload = json.dumps(payload_dict)
        
        output_record = {
            'recordId': record['recordId'],
            'result': 'Ok',
            'data': base64.b64encode(payload.encode('utf-8')).decode('utf-8'),
            'metadata': { 'partitionKeys': partition_keys }
        }
        output.append(output_record)

    print('Successfully processed {} records.'.format(len(event['records'])))

    return {'records': output}
