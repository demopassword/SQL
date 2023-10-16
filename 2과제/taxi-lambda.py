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
        id = payload_dict['id']
        vendorid = payload_dict['vendorId']
        pickupdate = payload_dict['pickupDate']
        dropoffdate = payload_dict['dropoffDate']
        
        # Do custom processing on the payload here
        # del payload_dict['pickupDate']
        # del payload_dict['dropoffDate']
        # del payload_dict['tripDuration']
        # payload_dict['id'] = id

        if vendorid == 1:
            payload_dict['vendor'] = "Hyundai"
        elif vendorid == 2:
            payload_dict['vendor'] = "KIA"
            
        parsed_datetime = datetime.fromisoformat(pickupdate)
        formatted_datetime = parsed_datetime.strftime("%Y-%m-%d %H:%M:%S")
        payload_dict['pickupDate'] = formatted_datetime
        parsed_datetime = datetime.fromisoformat(dropoffdate)
        formatted_datetime = parsed_datetime.strftime("%Y-%m-%d %H:%M:%S")
        payload_dict['dropoffDate'] = formatted_datetime
        time_difference = datetime.fromisoformat(dropoffdate) - datetime.fromisoformat(pickupdate)
        payload_dict['tripDuration'] = time_difference.total_seconds()

        # Do custom processing on the payload here
        partition_keys = {"taxi_id": id,
                          "year": parsed_datetime.year,
                          "month": parsed_datetime.month,
                          "day": parsed_datetime.day,
                          "hour": parsed_datetime.hour
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
