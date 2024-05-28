"""
Celery Task Module
"""
from loguru import logger
from celery import Task
from celery.exceptions import MaxRetriesExceededError
from .app_worker import app
from .yolo import YoloModel
from .core import Similarity, KeyWordExtractor, TextToSpeech
from typing import List


# region Exception Class
class CeleryTasksException(Exception):
    """Catch Celery Tasks Execeptions"""

    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            setattr(self, k, v)

    def __str__(self):
        err_msg = ""
        for k, v in vars(self).items():
            err_msg += f"{k} : {v}, "
        return err_msg


# endregion

# region Celery Tasks Models
class PredictTask(Task):
    def __init__(self):
        super().__init__()
        self.model = None

    def __call__(self, *args, **kwargs):
        if not self.model:
            logger.info("Loading YOLO Model...")
            self.model = YoloModel()
            logger.info("YOLO Model loaded")
        return self.run(*args, **kwargs)


class SimilarityTask(Task):
    def __init__(self):
        super().__init__()
        self.model = None

    def __call__(self, *args, **kwargs):
        if not self.model:
            logger.info("Loading Similarity Model...")
            self.model = Similarity()
            logger.info("Similarity Model loaded")
        return self.run(*args, **kwargs)


class KeywordExtractionTask(Task):
    def __init__(self):
        super().__init__()
        self.model = None

    def __call__(self, *args, **kwargs):
        if not self.model:
            logger.info("Loading Keywords Model...")
            self.model = KeyWordExtractor()
            logger.info("Keywords Model loaded")
        return self.run(*args, **kwargs)


class TTSTask(Task):
    def __init__(self):
        super().__init__()
        self.model = None

    def __call__(self, *args, **kwargs):
        if not self.model:
            logger.info("Loading TextToSpeech Model...")
            self.model = TextToSpeech()
            logger.info("TextToSpeech Model loaded")
        return self.run(*args, **kwargs)


# endregion

# region register Celery Tasks
@app.task(ignore_result=False, bind=True, base=PredictTask)
def predict_image(self, data):
    try:
        data_pred = self.model.predict(data)
        # used to display log once this celery task is done
        return {"task_name": "predict_image", "status": "SUCCESS", "result": data_pred}
    except Exception as ex:
        try:
            self.retry(countdown=2)
        except MaxRetriesExceededError as ex:
            return {
                "task_name": "predict_image",
                "status": "FAIL",
                "result": "max retried achieved",
            }


@app.task(ignore_result=False, bind=True, base=SimilarityTask)
def get_sim(self, data):
    similairty_results: List[dict] = []
    try:
        image_urls, prompt_text = data.get("images"), data.get("prompt")
        logger.info(f"Processing similarity on keyword: {prompt_text}\n")
        # logger.debug(f"image_urls = {image_urls}")

        if not image_urls:
            raise CeleryTasksException(
                task_name=get_sim.__name__, err_str="Images URLs or Location are empty"
            )

        similairty_results = self.model.similarity(image_urls, prompt_text)
        if not similairty_results:
            raise CeleryTasksException(
                task_name=get_sim.__name__, err_str="EMPTY similarity result"
            )

        return {
            "task_name": "get_sim",
            "status": "SUCCESS",
            "result": similairty_results,
        }
    except (CeleryTasksException, Exception) as e:
        return {"task_name": "get_sim", "status": "FAIL", "error": str(e)}


@app.task(ignore_result=False, bind=True, base=KeywordExtractionTask)
def get_keywords(self, data):
    res: dict[list] = {}
    logger.info("Processing Keywords .. ")
    try:
        doc = data.get("doc")

        res["sentences"], res["keywords"] = self.model.filter_keywords(doc)
        return {
            "task_name": "get_keywords",
            "status": "SUCCESS",
            "result": res,
        }
    except (CeleryTasksException, Exception) as e:
        return {"task_name": "get_tts", "status": "FAIL", "error": str(e)}


@app.task(ignore_result=False, bind=True, base=TTSTask)
def get_tts(self, data):
    try:
        text, save_path = data.get("text"), data.get("save_path")

        logger.debug("text = ", text)
        logger.debug("save_path = ", save_path)

        self.model.text_to_speech(text, save_path)
        return {
            "task_name": "get_sim",
            "status": "SUCCESS",
            "result": save_path,
        }
    except (CeleryTasksException, Exception) as e:
        return {"task_name": "get_tts", "status": "FAIL", "error": str(e)}


# endregion
