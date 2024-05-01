FROM python:3.12.1

COPY ./celery_core /usr/src/celery_core
COPY ./config /usr/src/config

COPY ./requirements-celery.txt /usr/src/


RUN pip install -r /usr/src/requirements-celery.txt

RUN pip install flower==2.0.1

WORKDIR /usr/src

# CMD celery flower -A celery_core.tasks --broker=amqp://${RABBITMQ_USERNAME}:${RABBITMQ_PASSWORD}@${RABBITMQ_HOST}:${RABBITMQ_PORT}//

CMD celery -A celery_core.tasks --broker=redis://localhost:6379/0 flower


