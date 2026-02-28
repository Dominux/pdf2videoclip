# Downloading model files
if [ -f "/models/model.gguf" ]; then
  echo "Model file already exists"
else
  echo "Downloading model"

  wget -O /models/model.gguf https://huggingface.co/bartowski/Qwen2.5-7B-Instruct-GGUF/resolve/main/Qwen2.5-7B-Instruct-Q5_K_M.gguf

  echo "Successfully downloaded model"
fi

# creating models in background
sh -c "sleep 5 && ollama create "qwen-zoomer" -f ./Modelfile" &

# Running Ollama
/bin/ollama serve
