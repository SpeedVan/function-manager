FROM alpine:3.10.2

RUN mkdir /app

COPY ./build/function-manager /app/function-manager

ENTRYPOINT [ "/app/function-manager" ] 