FROM pytorch/pytorch:latest
WORKDIR /app
RUN apt-get update
RUN apt-get install ffmpeg libsm6 libxext6 -y
COPY requirements.txt requirements.txt
RUN pip install -U pip
RUN pip install -r requirements.txt
COPY . .
RUN rm -r /app/api/uploads/*; exit 0
RUN rm -r /app/api/static/*; exit 0

# GOlang
RUN apt-get install -y curl
RUN curl -OL https://golang.org/dl/go1.22.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.22.2.linux-amd64.tar.gz
# RUN echo 'export PATH=$PATH:/usr/local/go/bin' | tee -a ~/.profile
# RUN source ~/.profile
ENV PATH=$PATH:/usr/local/go/bin
