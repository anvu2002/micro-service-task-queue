import redis
import os
import sys
import json

from .worker import celery
from config.config import REDIS_STORE_CONN_URI
from loguru import logger
# from service.core import similarity

redis_store = redis.Redis.from_url(REDIS_STORE_CONN_URI)



@celery.task
def similuate_buy_process(name, stage):
    redis_store.set(name, stage)
    return stage

@celery.task
def get_sim(data: dict) -> list:
    logger.debug(f"Loading Image Captioning Engine ...")
    sys.path.append(os.getcwd())
    from service.core import similarity

    prompt  = data["prompt"]
    logger.debug("Processing Similarity on Keyword: %s\n",prompt)

    
    image_urls, prompt_text = data["images"], data["prompt"]
    results = similarity(image_urls, prompt_text)
    logger.debug(f"[*] DONE - Keyword: {prompt}\n")
    logger.debug(f"results = {results}")

    redis_store.set(prompt,json.dumps(results))

    return results
 