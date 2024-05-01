from PIL import Image
from transformers import BlipProcessor, BlipForConditionalGeneration
from sentence_transformers import SentenceTransformer
import requests
import numpy as np
from loguru import logger
from openai import OpenAI
from typing import List
import os

from config.config import OPENAI_KEY
# from router.core_types import ImageScore

#For keyword extractor
from sklearn.feature_extraction.text import CountVectorizer
import pandas as pd
import nltk
from nltk.tokenize import sent_tokenize, word_tokenize
from nltk.corpus import stopwords


# try:
#     nltk.data.find('tokenizers/punkt')
# except LookupError:
#     nltk.download('punkt')

# try:
#     nltk.data.find('tokenizers/stopwords')
# except LookupError:
#     nltk.download('stopwords')

# Init global models, processors, and clients on server boot
i2t_processor = BlipProcessor.from_pretrained("Salesforce/blip-image-captioning-large")
i2t_model = BlipForConditionalGeneration.from_pretrained("Salesforce/blip-image-captioning-large")
t2v_model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')
openai_client = OpenAI(api_key=OPENAI_KEY)

def image_to_text(image_url: str) -> str:
    raw_image: Image = None
    result: str = None

    try:
        # internet url
        if image_url.startswith("http"): 
            raw_image = Image.open(requests.get(image_url, stream=True).raw).convert('RGB')
        # local image
        else:
            raw_image = Image.open(image_url).convert('RGB')

        # unconditional image captioning
        inputs = i2t_processor(raw_image, return_tensors="pt")

        out = i2t_model.generate(**inputs, max_new_tokens=400)
        result = i2t_processor.decode(out[0], skip_special_tokens=True)
        logger.info(f"Image {image_url}'s descrption is: {result}\n")

    except Exception as e:
        logger.error(f"{e} occured when processing captioning image {image_url}")

    return result

def text_to_vec(text: str) -> np.array:
    embeddings = np.array(t2v_model.encode(text))

    return embeddings

def similarity(image_urls: List[str], prompt_text: str):
    sim_results=[]
    prompt_embedding = text_to_vec(prompt_text)

    for image_url in image_urls:
        image_text = image_to_text(image_url)
        if image_text is None:
            logger.error(f"Cannot Find {image_url}")
            logger.debug(f"WORKING DIR: {os.getcwd()}")
            logger.debug(f"Check Postman: ""./images/[image_name]")
            continue
        text_embedding = text_to_vec(image_text)
        score = float(np.linalg.norm(text_embedding-prompt_embedding))
        sim_results.append({
            "score":score, "url":image_url, "description":image_text
            })

    return sim_results if len(sim_results) > 0 else None

def text_to_speech(text: str, save_path: str) -> str:
    logger.info(f"Converting {text} to {save_path}")
    
    response = openai_client.audio.speech.create(
        model="tts-1",
        voice="alloy",
        input=text
    )

    response.stream_to_file(save_path)
    logger.info(f"{save_path} created!")


#REMOVE STOPWORDS AND RETURN A PREPROCESSED STRING
def filter_keywords(raw_text: str) -> list[str]:

    sentences = sent_tokenize(raw_text)

    stop_words = set(stopwords.words('english'))

    filtered_sentences = []

    for sentence in sentences:
        words = word_tokenize(sentence)
        filtered_sentence = [word for word in words if word.lower() not in stop_words]
        filtered_sentences.append(filtered_sentence)

    #up until now filtered_sentences is a list of keywords
    print(filtered_sentences)

    #join the keywords to return a single string object
    preprocessed_text = [' '.join(sentence) for sentence in filtered_sentences]

    # print(preprocessed_text) #list
    # for sentence in preprocessed_text:
    #     print(sentence)

    return preprocessed_text



