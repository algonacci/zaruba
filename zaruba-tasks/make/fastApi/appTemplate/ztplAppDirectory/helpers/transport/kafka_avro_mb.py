from typing import Any, Callable, List, Mapping, TypedDict
from helpers.transport.interface import MessageBus
from helpers.transport.kafka_helper import create_kafka_topic
from helpers.transport.kafka_avro_config import KafkaAvroEventMap
from confluent_kafka.avro import AvroProducer, AvroConsumer

import threading
import traceback


def create_kafka_avro_connection_parameters(bootstrap_servers: str, schema_registry: str = '', sasl_mechanism: str = '', sasl_plain_username: str = '', sasl_plain_password: str = '', security_protocol='', **kwargs) -> Mapping[str, Any]:
    if sasl_mechanism == '':
        sasl_mechanism = 'PLAIN'
    if security_protocol == '':
        security_protocol = 'SASL_PLAINTEXT'
    # https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
    return {
        'bootstrap.servers': bootstrap_servers,
        'schema.registry.url': schema_registry,
        'sasl.mechanism': sasl_mechanism,
        'sasl.username': sasl_plain_username,
        'sasl.password': sasl_plain_password,
        # 'topic.metadata.propagation.max.ms': '100',
        # 'security.protocol': security_protocol,
        **kwargs
    }


class KafkaAvroMessageBus(MessageBus):

    def __init__(self, kafka_avro_connection_parameters: Mapping[str, Any], kafka_avro_event_map: KafkaAvroEventMap):
        self.kafka_avro_connection_parameters = kafka_avro_connection_parameters
        self.kafka_avro_event_map = kafka_avro_event_map
        self.consumers: Mapping[str, AvroConsumer] = {}
        self.event_map = kafka_avro_event_map
        self.is_shutdown = False
        self.error_count = 0

    def get_error_count(self) -> int:
        return self.error_count

    def shutdown(self):
        if self.is_shutdown:
            return
        for event_name, consumer in self.consumers.items():
            print('stop listening to {event_name}'.format(event_name=event_name))
            consumer.close()
        self.is_shutdown = True

    def handle(self, event_name: str) -> Callable[..., Any]:
        def register_event_handler(event_handler: Callable[[Any], Any]):
            topic = self.event_map.get_topic(event_name)
            # create topic if not exists
            self._create_kafka_topic(topic)
            # create consumer
            group_id = self.event_map.get_group_id(event_name)
            consumer_args = {**self.kafka_avro_connection_parameters}
            consumer_args['group.id'] = group_id
            consumer = AvroConsumer(consumer_args)
            # register consumer
            self.consumers[event_name] = consumer
            # start consume
            thread = threading.Thread(target=self._handle, args=[consumer, consumer_args, event_name, topic, group_id, event_handler], daemon = True)
            thread.start()
        return register_event_handler

    def _handle(self, consumer: AvroConsumer, consumer_args, event_name: str, topic: str, group_id: str, event_handler: Callable[[Any], Any]):
        for _ in range(3):
            try:
                print({'action': 'subscribe_kafka_avro_topic', 'topic': topic})
                consumer.subscribe([topic])
                break
            except:
                print(traceback.format_exc())
                consumer = AvroConsumer(consumer_args)
                self.consumers[event_name] = consumer
        while True:
            try:
                serialized_message = consumer.poll(1)
                if serialized_message is None:
                    continue
                if serialized_message.error():
                    print("AvroConsumer error: {}".format(serialized_message.error()))
                    continue
                if serialized_message.value():
                    print({'action': 'handle_kafka_avro_event', 'event_name': event_name, 'value': serialized_message.value(), 'key': serialized_message.key(), 'topic': serialized_message.topic(), 'partition': serialized_message.partition(), 'offset': serialized_message.offset(), 'group_id': group_id})
                    message = self.event_map.get_decoder(event_name)(serialized_message.value())
                    event_handler(message)
            except:
                self.error_count += 1
                print(traceback.format_exc())
                consumer = AvroConsumer(consumer_args)
                self.consumers[event_name] = consumer
                print({'action': 're_subscribe_kafka_avro_topic', 'topic': topic})
                consumer.subscribe([topic])

    def publish(self, event_name: str, message: Any) -> Any:
        try:
            topic = self.event_map.get_topic(event_name)
            # create topic if not exist
            # self._create_kafka_topic(topic)
            # time.sleep(3)
            key_schema = self.event_map.get_key_schema(event_name)
            value_schema = self.event_map.get_value_schema(event_name)
            producer = AvroProducer(self.kafka_avro_connection_parameters, default_key_schema=key_schema, default_value_schema=value_schema)
            key_maker = self.event_map.get_key_maker(event_name)
            key = key_maker(message)
            serialized_message = self.event_map.get_encoder(event_name)(message)
            print({'action': 'publish_kafka_avro_event', 'event_name': event_name, 'key': key, 'message': message, 'topic': topic, 'serialized': serialized_message})
            producer.produce(topic=topic, key=key, value=serialized_message, callback=_produce_callback)
            producer.flush()
        except Exception as e:
            self.error_count += 1
            raise e

    def _create_kafka_topic(self, topic: str):
        create_kafka_topic(topic, {
            'bootstrap.servers': self.kafka_avro_connection_parameters['bootstrap.servers'],
            'sasl.mechanism': self.kafka_avro_connection_parameters['sasl.mechanism'],
            'sasl.username': self.kafka_avro_connection_parameters['sasl.username'],
            'sasl.password': self.kafka_avro_connection_parameters['sasl.password']
        })


def _produce_callback(err, msg):
    if err is not None:
        print("Failed to deliver message: %s: %s" % (str(msg), str(err)))
    else:
        print("Message produced: %s" % (str(msg)))