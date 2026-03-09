from fastapi import FastAPI, UploadFile, File
from faster_whisper import WhisperModel
from pydantic import BaseModel


class TranscriptionWord(BaseModel):
    word: str
    start: float
    end: float


model = WhisperModel("/model", device='cpu', compute_type="int8")


app = FastAPI()


@app.post("/transcribe")
def transcribe(file: 'UploadFile' = File(...)) -> 'list[TranscriptionWord]':
    segments, _info = model.transcribe(file.file, word_timestamps=True, beam_size=2)

    result = []
    for segment in segments:
        if not segment.words:
            continue

        for word in segment.words:
            trans_word = TranscriptionWord(word=word.word, start=word.start, end=word.end)
            result.append(trans_word)
    return result
