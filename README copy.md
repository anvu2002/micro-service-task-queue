# ML Infra  - Task Queue

## Tech Stack
fastapi - celery - rabbitmq - redis -> Docker


## Start up the Container

```bash
docker-compose up -d --build
```

fastapi: 5000
celery flower: 5555


swagger docs - `http://localhost:5000/`

redoc - `http://localhost:5000/redoc`

celery flower - `http://localhost:5555`
