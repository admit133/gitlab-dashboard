FROM node:lts-alpine AS gui

WORKDIR /var/lib/gui
ADD gui /var/lib/gui
RUN yarn
RUN yarn build

FROM golang:1.17 AS server

ADD server /go/src/server
WORKDIR /go/src/server/cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server


FROM alpine

COPY --from=server /go/bin/server /bin/server
COPY --from=gui /var/lib/gui/build /public

ENV GITLAB_BASE_URL=https://gitlab.com
ENV GITLAB_TOKEN=""
ENV GITLAB_PROJECT_IDS=""
ENV USER_LINK_TEMPLATE=https://gitlab.com/{username}
ENV ENVIRONMENT_UPDATE_DURATION=1m
ENV PUBLIC_DIR="/public"

EXPOSE 3001

CMD ["/bin/server"]
