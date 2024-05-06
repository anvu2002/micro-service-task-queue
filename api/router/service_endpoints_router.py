import sys
import os

sys.path.insert(0, os.path.realpath(os.path.pardir))
from fastapi import Request, APIRouter, File, UploadFile
from fastapi.responses import JSONResponse

from loguru import logger
from typing import List
import uuid

from celery_tasks.tasks import predict_image, get_sim
from celery.result import AsyncResult
from models import Task, Prediction


UPLOAD_FOLDER = "uploads"
isdir = os.path.isdir(UPLOAD_FOLDER)
if not isdir:
    os.makedirs(UPLOAD_FOLDER)

router = APIRouter(prefix="/api", tags=["api"])


@router.post("/get_similarity")
async def get_similarity(request: Request):
    """
    Usage: Image selection through similarity against prompt keywords
    request format:
    {
        "images": ["./images/1.jpg", "https://test.com/1.jpg"],
        "prompt": "google search prompt"
    }
    """
    tasks = []
    try:
        d = {}
        data = await request.json()
        # Need Data Validation (data models)
        task_id = get_sim.delay(data)
        d["task_id"] = str(task_id)
        d["task_name"] = get_similarity.__name__
        d["status"] = "PROCESSING"
        d["url_result"] = f"/api/result/{task_id}"
        d["requested_data"] = data
        tasks.append(d)
        return JSONResponse(status_code=202, content=tasks)
    except Exception as err:
        logger.error(err)
        return JSONResponse(
            status_code=400, content={"TASK": "get_similarity", "status": "FAILED"}
        )


@router.post("/process")
async def process(files: List[UploadFile] = File(...)):
    tasks = []
    try:
        for file in files:
            d = {}
            try:
                # Save Uploaded File
                name = str(uuid.uuid4()).split("-")[0]
                ext = file.filename.split(".")[-1]
                file_name = f"{UPLOAD_FOLDER}/{name}.{ext}"
                with open(file_name, "wb+") as f:
                    f.write(file.file.read())
                f.close()

                # start task prediction
                task_id = predict_image.delay(os.path.join("api", file_name))
                d["task_id"] = str(task_id)
                d["task_name"] = process.__name__
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
            content={
                "task_id": str(task_id),
                "status": task.status,
                "result": "[PENDING...]",
            },
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
        content={"task_id": str(task_id), "status": task.status},
    )
