"""
Web Application -- TransVLN -- 4567
"""

from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates

from starlette.middleware.cors import CORSMiddleware
import os

app = FastAPI(title="TransVLN", docs_url="/docs")
app.mount("/static", StaticFiles(directory="templates/static"), name="static")
templates = Jinja2Templates(directory="templates")

os.environ["no_proxy"] = "*"
os.environ["OBJC_DISABLE_INITIALIZE_FORK_SAFETY"] = "YES"

origins = [
    "*",
]

app.add_middleware(
    CORSMiddleware,
    # sources allow to access
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/", response_class=HTMLResponse)
async def index(request: Request):
    return templates.TemplateResponse("index.html", context={"request": request})
