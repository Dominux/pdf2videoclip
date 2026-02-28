FROM ollama/ollama:0.17.4-rocm

RUN apt update && apt install -y wget && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./run.sh ./
COPY ./Modelfile ./Modelfile

ENTRYPOINT [ "sh" ]
CMD [ "./run.sh" ]
