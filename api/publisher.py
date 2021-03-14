import time

from kafka import KafkaProducer

f = open('/home/tbeppe/Documents/tcc/pre-processing/values.csv', 'r')

prd = KafkaProducer(bootstrap_servers='localhost:9092',compression_type='gzip')
for line in f.readlines():
    prd.send('b3-simulate',bytes(line, 'utf-8'))
    time.sleep(1)
f.close()