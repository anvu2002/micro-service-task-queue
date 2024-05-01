from celery import Celery
from config.config import BACKEND_CONN_URI, BROKER_CONN_URI

celery = Celery('hello', broker=BROKER_CONN_URI, backend=BACKEND_CONN_URI)