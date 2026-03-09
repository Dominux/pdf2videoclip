FROM python:3.11-slim

RUN apt-get update \
  && apt-get install -y ffmpeg git-lfs \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN git-lfs clone https://huggingface.co/dvislobokov/faster-whisper-large-v3-turbo-russian /model

COPY ./requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

COPY ./main.py ./

CMD ["fastapi", "run", "main.py"]
