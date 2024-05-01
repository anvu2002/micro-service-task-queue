from fastapi import APIRouter, Request
import redis
from loguru import logger

from config.config import stages, REDIS_STORE_CONN_URI, STAGING_TIME
# Celery Tasks
from celery_core.tasks import similuate_buy_process,get_sim

redis_store = redis.Redis.from_url(REDIS_STORE_CONN_URI)
router = APIRouter(prefix="/api", tags=["api"])

@router.get("/buy/{name}")
async def buy(name: str):
    for i in range(0, 5):
        similuate_buy_process.apply_async((name, stages[i]), countdown=i*STAGING_TIME)
    return True


@router.get("/status/{name}")
async def status(name: str):
    return redis_store.get(name)


@router.post("/get_similarity")
async def get_similarity(request: Request):
    data = await request.json()
    if redis_store.get(data["prompt"]) is not None:
        return redis_store.get(data["prompt"])
    else:
        logger.info(f"Similarity Request Received for: {data}")

        try:
            get_sim.delay(data)
            return {
                        "Status":"Requested ML - Image Captioning Service"
                    }
        except:
            logger.error("ML DEAD :)")
