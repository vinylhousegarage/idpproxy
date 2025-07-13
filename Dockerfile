FROM golang:1.24.5

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    make \
    && rm -rf /var/lib/apt/lists/*

CMD ["bash"]
