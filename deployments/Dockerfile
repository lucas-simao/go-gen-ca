FROM gcr.io/cloud-builders/go:1.18 AS build

ENV PORT 3000

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build -o api .

FROM node:18-alpine AS final

WORKDIR /app

COPY --from=build /app/api .
COPY --from=build /app/json-to-go.js .

EXPOSE $PORT
ENTRYPOINT ./api
