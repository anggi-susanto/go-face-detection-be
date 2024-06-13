FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /face_detection_be

# Install Python dependencies
RUN apk add --no-cache python3-dev gcc musl-dev
RUN pip3 install opencv-python-headless

EXPOSE 8080

CMD [ "/face_detection_be" ]
