# _________ Component: Images Similarity _________
from PIL import Image
from transformers import BlipProcessor, BlipForConditionalGeneration
from sentence_transformers import SentenceTransformer
import requests
import numpy as np
from loguru import logger
from openai import OpenAI
from typing import List

# Init global models, processors, and clients on server boot
# logger.info("Loading Image Caption, Sentence Transformer....")
# i2t_processor = BlipProcessor.from_pretrained(
#     "Salesforce/blip-image-captioning-large"
# )
# i2t_model = BlipForConditionalGeneration.from_pretrained(
#     "Salesforce/blip-image-captioning-large"
# )
# t2v_model = SentenceTransformer("sentence-transformers/all-MiniLM-L6-v2")
# logger.info("Finished!")


# _________ Component: Text-to-Speech _________
# from config.config import OPENAI_KEY

# _________ Component: Keyword Extractor _________
from sklearn.feature_extraction.text import CountVectorizer
import nltk, string
from nltk.tokenize import sent_tokenize, word_tokenize
from nltk.corpus import stopwords

try:
    nltk.data.find("tokenizers/punkt")
except LookupError:
    nltk.download("punkt")

try:
    nltk.data.find("tokenizers/stopwords")
except LookupError:
    nltk.download("stopwords")

# from api.router.core_types import ImageScore
# from api.models import ImageScore


class Similarity:
    """
    ML Service: Image Captioning
    Usage: Caption images --> compare and determine similarity with the provided prompts / keywords
    """

    def __init__(self):
        # Init global models, processors, and clients on server boot
        self.i2t_processor = BlipProcessor.from_pretrained(
            "Salesforce/blip-image-captioning-large"
        )
        self.i2t_model = BlipForConditionalGeneration.from_pretrained(
            "Salesforce/blip-image-captioning-large"
        )
        self.t2v_model = SentenceTransformer("sentence-transformers/all-MiniLM-L6-v2")

    def image_to_text(self, image_url: str) -> str:
        raw_image: Image = None
        result: str = None

        try:
            # internet url
            if image_url.startswith("http"):
                raw_image = Image.open(
                    requests.get(image_url, stream=True).raw
                ).convert("RGB")
            # local image
            else:
                raw_image = Image.open(image_url).convert("RGB")

            # unconditional image captioning
            inputs = self.i2t_processor(raw_image, return_tensors="pt")

            out = self.i2t_model.generate(**inputs, max_new_tokens=400)
            result = self.i2t_processor.decode(out[0], skip_special_tokens=True)
            logger.info(f"Image {image_url}'s descrption is: {result}\n")

        except Exception as e:
            logger.error(
                f"Error '{e}' occured when processing captioning image {image_url}"
            )

        return result

    def text_to_vec(self, text: str) -> np.array:
        embeddings = np.array(self.t2v_model.encode(text))

        return embeddings

    def similarity(self, image_urls: List[str], prompt_text: str) -> List:
        sim_results: List[dict] = []
        prompt_embedding = self.text_to_vec(prompt_text)

        for image_url in image_urls:
            image_text = self.image_to_text(image_url)
            text_embedding = self.text_to_vec(image_text)
            score = float(np.linalg.norm(text_embedding - prompt_embedding))
            sim_results.append(
                {"url": image_url, "score": score, "description": image_text}
            )

        return sim_results if sim_results else None


class KeyWordExtractor:
    """
    ML Service: Phrases (tokens) and Keywords Generator
    Usage: Supply with a document
    """

    def __init__(self):
        # Possible modules: keyBERT, vlt5

        pass

    def filter_keywords(raw_text: str) -> list[str]:
        # fp = open(file, encoding='UTF-8')
        # raw_text = fp.read()

        # Tokenize --> Sentences
        sentences = sent_tokenize(raw_text)
        logger.info(f"N Sen = {len(sentences)}\n")
        stop_words = set(stopwords.words("english"))
        filtered_sentences = []

        # Extract Keywords per tokens / sentence
        for sentence in sentences:
            words = word_tokenize(sentence)
            filtered_sentence = [
                word for word in words if word.lower() not in stop_words
            ]
            filtered_sentences.append(filtered_sentence)
        trash = list(string.punctuation) + list(string.whitespace)
        preprocessed_text = [
            " ".join(word for word in sentence if word not in trash)
            for sentence in filtered_sentences
        ]

        return sentences, preprocessed_text


class TextToSpeech:
    """
    ML Service: Generate speech from text
    Usage: supply with text ( keywords from the doc)

    """

    def __init__(self):
        self.openai_client = OpenAI(api_key="")

    def text_to_speech(self, text: str, save_path: str) -> None:
        try:
            logger.info(f"Converting {text} to {save_path}")

            # response = self.openai_client.audio.speech.create(
            #     model="tts-1", voice="alloy", input=text
            # )

            # response.stream_to_file(save_path)

            #  Use GCP instead -- for now
            import time

            time.sleep(15)
            logger.info(f"{save_path} TTS created!")
        except Exception as e:
            logger.error(e)
