FROM golang:1.12.6-stretch

LABEL maintainer="zekro <contact@zekro.de>"

#### PREPARINGS ####

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash - &&\
    apt-get install -y \
        nodejs \
        git

RUN npm i -g @vue/cli

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/zekroTJA/slms

ADD . .

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/certs

#### BUILD BACK END ####

RUN dep ensure -v

RUN go build -v -o ./bin/slms -ldflags "\
		-X github.com/zekroTJA/slms/internal/static.AppVersion=$(git describe --tags) \
		-X github.com/zekroTJA/slms/internal/static.AppCommit=$(git rev-parse HEAD) \
        -X github.com/zekroTJA/slms/internal/static.Release=TRUE" \
        ./cmd/slms/*.go

#### BUILD FRONT END ####

RUN cd web &&\
    npm install &&\
    npm run build

RUN mkdir -p ./bin/web &&\
    mv ./web/dist ./bin/web/dist

#### EXPOSE AND RUN ####

EXPOSE 8080

WORKDIR ${GOPATH}/src/github.com/zekroTJA/slms/bin

CMD ./slms \
        -c /etc/config/config.yml \
        -addr :8080