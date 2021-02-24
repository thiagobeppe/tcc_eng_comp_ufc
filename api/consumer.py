from kafka import KafkaConsumer

consumer = KafkaConsumer('b3-simulate',group_id='b3-simulate-python')
for msg in consumer:
    print (msg)