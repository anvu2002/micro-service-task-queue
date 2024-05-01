import os
from os.path import join, dirname
from dotenv import load_dotenv

dotenv_path = join(dirname(__file__), '.env')
load_dotenv(dotenv_path)
OPENAI_KEY = os.environ.get("OPENAI_KEY")

# Reddis  - Celery - RabbitMQ

REDIS_HOST = os.environ.get('REDIS_HOST',"redis")
REDIS_PORT = os.environ.get('REDIS_PORT',6379)
REDIS_CELERY_DB_INDEX = os.environ.get('REDIS_CELERY_DB_INDEX',10)
REDIS_STORE_DB_INDEX = os.environ.get('REDIS_STORE_DB_INDEX',0)

RABBITMQ_HOST = os.environ.get('RABBITMQ_HOST',"redis")
RABBITMQ_USERNAME = os.environ.get('RABBITMQ_USERNAME',"guest")
RABBITMQ_PASSWORD = os.environ.get('RABBITMQ_PASSWORD',"guest")
RABBITMQ_PORT = os.environ.get('RABBITMQ_PORT',5672)

BROKER_CONN_URI = f"amqp://{RABBITMQ_USERNAME}:{RABBITMQ_PASSWORD}@{RABBITMQ_HOST}:{RABBITMQ_PORT}"
BACKEND_CONN_URI = f"redis://{REDIS_HOST}:{REDIS_PORT}/{REDIS_CELERY_DB_INDEX}"
REDIS_STORE_CONN_URI = f"redis://{REDIS_HOST}:{REDIS_PORT}/{REDIS_STORE_DB_INDEX}"

stages = ["confirmed", "shipped", "in transit", "arrived", "delivered"]
STAGING_TIME = 7 # seconds