FROM python:3.12.1

COPY ./celery_core /usr/src/celery_core
COPY ./config /usr/src/config
COPY ./service /usr/src/service
COPY ./images /usr/src/images


COPY ./requirements-celery.txt /usr/src/


RUN pip install -r /usr/src/requirements-celery.txt

WORKDIR /usr/src

# CMD ["/bin/bash"]
CMD celery -A celery_core.tasks worker -l INFO --pool=prefork --concurrency=10
 