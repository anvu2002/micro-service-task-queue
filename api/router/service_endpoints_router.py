import sys
import os

sys.path.insert(0, os.path.realpath(os.path.pardir))
from fastapi import Request, APIRouter, File, UploadFile
from fastapi.responses import JSONResponse

from loguru import logger
from typing import List
import uuid

from celery_tasks.tasks import predict_image
from celery.result import AsyncResult
from models import Task, Prediction


UPLOAD_FOLDER = "uploads"
isdir = os.path.isdir(UPLOAD_FOLDER)
if not isdir:
    os.makedirs(UPLOAD_FOLDER)

router = APIRouter(prefix="/api", tags=["api"])


@router.post("/process")
async def process(files: List[UploadFile] = File(...)):
    tasks = []
    try:
        for file in files:
            d = {}
            try:
                name = str(uuid.uuid4()).split("-")[0]
                ext = file.filename.split(".")[-1]
                file_name = f"{UPLOAD_FOLDER}/{name}.{ext}"
                with open(file_name, "wb+") as f:
                    f.write(file.file.read())
                f.close()

                # start task prediction
                task_id = predict_image.delay(os.path.join("api", file_name))
                d["task_id"] = str(task_id)
                d["status"] = "PROCESSING"
                d["url_result"] = f"/api/result/{task_id}"
            except Exception as ex:
                logger.info(ex)
                d["task_id"] = str(task_id)
                d["status"] = "ERROR"
                d["url_result"] = ""
            tasks.append(d)
        return JSONResponse(status_code=202, content=tasks)
    except Exception as ex:
        logger.info(ex)
        return JSONResponse(status_code=400, content=[])


@router.get("/result/{task_id}", response_model=Prediction)
async def result(task_id: str):
    task = AsyncResult(task_id)

    # Task Not Ready
    if not task.ready():
        return JSONResponse(
            status_code=202,
            content={"task_id": str(task_id), "status": task.status, "result": ""},
        )

    # Task done: return the value
    task_result = task.get()
    result = task_result.get("result")
    return JSONResponse(
        status_code=200,
        content={
            "task_id": str(task_id),
            "status": task_result.get("status"),
            "result": result,
        },
    )


@router.get("/status/{task_id}", response_model=Prediction)
async def status(task_id: str):
    task = AsyncResult(task_id)
    return JSONResponse(
        status_code=200,
        content={"task_id": str(task_id), "status": task.status, "result": ""},
    )
