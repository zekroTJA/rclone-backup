FROM python:3-alpine
WORKDIR /app
COPY backup/ backup/
COPY README.md .
COPY LICENSE .
RUN apk add rclone
ENTRYPOINT [ "python3", "backup/main.py" ]