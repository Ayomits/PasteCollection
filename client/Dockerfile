FROM node:24.4.1-alpine3.21

RUN apk add --no-cache runit

WORKDIR /app

ARG APP_ENV=dev

COPY ./client /app

COPY ./docker/runit/client/dev.run /tmp/runit-dev.run
COPY ./docker/runit/client/prod.run /tmp/runit-prod.run

RUN mkdir -p /etc/service/app /runit

RUN if [ "$APP_ENV" = "prod" ]; then \
  cp /tmp/runit-prod.run /etc/service/app/run; \
  else \
  cp /tmp/runit-dev.run /etc/service/app/run; \
  fi && \
  rm /tmp/runit-*.run


RUN chmod +x /etc/service/app/run

RUN npm install pnpm -g

CMD ["runsvdir", "/etc/service"]
