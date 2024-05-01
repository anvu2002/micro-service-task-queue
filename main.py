# encoding:utf-8

import os
import time
import uvicorn
from loguru import logger
from fastapi import FastAPI
from starlette.middleware.cors import CORSMiddleware

from router import (service_endpoints_router)

app = FastAPI(title="ML_Serivces", docs_url="/", version="1.0.0")

os.environ["no_proxy"]="*"
os.environ['OBJC_DISABLE_INITIALIZE_FORK_SAFETY'] = 'YES'


log_dir = "logs"
log_path = os.path.join(log_dir, f'{time.strftime("%Y-%m-%d")}.log')
logger.add(log_path, rotation='0:00', enqueue=True, serialize=False,
           encoding="utf-8", retention="7 days", diagnose=False, backtrace=True)

origins = [
    "*",
]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,  # 允许访问的源
    allow_credentials=True,  # 支持 cookie
    allow_methods=["*"],  # 允许使用的请求方法
    allow_headers=["*"],  # 允许携带的 Headers
)

app.include_router(service_endpoints_router)

# Docker Port = 5000
if __name__ == "__main__":
    uvicorn.run("main:app", port=8850, log_level="info")
    
