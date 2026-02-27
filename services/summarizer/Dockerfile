FROM ollama/ollama:0.17.4-rocm

RUN apt update && apt install -y curl git-lfs

COPY ./run.sh /
COPY ./modelfiles /modelfiles

ENTRYPOINT [ "sh" ]
CMD [ "/run.sh" ]
