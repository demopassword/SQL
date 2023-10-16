import base64
import json

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
        datetime = payload_dict['datetime']
        method = payload_dict['method']
        path = payload_dict['path']
        statuscode = payload_dict['code']
        
        # Do custom processing on the payload here
        payload_dict['year'] = datetime.split('/')[0]
        payload_dict['month'] = datetime.split('/')[1]
        payload_dict['day'] = datetime.split('/')[2].split(' ')[0]
        payload_dict['hour'] = datetime.split(' ')[1].split(':')[0]
        payload_dict['minute'] = datetime.split(' ')[1].split(':')[1]
        payload_dict['second'] = datetime.split(' ')[1].split(':')[2]
        payload_dict['statuscode'] = statuscode
        
        # Do custom processing on the payload here
        payload = json.dumps(payload_dict)
        
        output_record = {
            'recordId': record['recordId'],
            'result': 'Ok',
            'data': base64.b64encode(payload.encode('utf-8')).decode('utf-8')
        }
        output.append(output_record)

    print('Successfully processed {} records.'.format(len(event['records'])))

    return {'records': output}
