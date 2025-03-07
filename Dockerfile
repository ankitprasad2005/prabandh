FROM python:3.13-slim
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget unzip build-essential python3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app

RUN pip install poetry
RUN poetry config virtualenvs.create false

COPY pyproject.toml poetry.lock README.md /app/
COPY src /app/src
RUN poetry install

WORKDIR /app/src
EXPOSE 5555
# CMD ["gunicorn","-w", "4", "-k", "uvicorn.workers.UvicornWorker main:app", "--bind", "0.0.0.0:5555"]
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5555"]