FROM python:3-alpine as builder

WORKDIR /app

COPY ./requirements.txt .

RUN pip install -r requirements.txt

FROM python:3-alpine as runner

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY ./app.py  .

ENTRYPOINT ["python3", "/app/app.py"]