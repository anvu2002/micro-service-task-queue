"""
VLNML API Endpoints -- 8000
"""
import os
import sys

sys.path.insert(0, os.path.realpath(os.path.pardir))
import time
import uvicorn
from loguru import logger
from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from starlette.middleware.cors import CORSMiddleware

from router import service_endpoints_router

UPLOAD_FOLDER = "uploads"
STATIC_FOLDER = "static/results"

isdir = os.path.isdir(UPLOAD_FOLDER)
if not isdir:
    os.makedirs(UPLOAD_FOLDER)

isdir = os.path.isdir(STATIC_FOLDER)
if not isdir:
    os.makedirs(STATIC_FOLDER)


app = FastAPI(title="ML_Services", docs_url="/", version="1.0.0")
app.mount("/static", StaticFiles(directory=STATIC_FOLDER), name="static")
# app.mount("/img", StaticFiles(directory="images/"), name="img")


os.environ["no_proxy"] = "*"
os.environ["OBJC_DISABLE_INITIALIZE_FORK_SAFETY"] = "YES"

origins = [
    "*",
]

app.add_middleware(
    CORSMiddleware,
    # sources are allow to access
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

log_dir = "logs"
log_path = os.path.join(log_dir, f'{time.strftime("%Y-%m-%d")}.log')
logger.add(
    log_path,
    rotation="0:00",
    enqueue=True,
    serialize=False,
    encoding="utf-8",
    retention="7 days",
    diagnose=False,
    backtrace=True,
)

app.include_router(service_endpoints_router)
