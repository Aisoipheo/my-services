FROM jenkins/jnlp-slave:4.13.2-1-alpine AS base

USER root

ENV WORKSPACE /home/jenkins

# =============== INSTALL GO ===============

FROM base AS install-go

ENV GOLANG_VERSION=1.18.4 \
	GOLANG_DOWNLOAD_SHA256=c9b099b68d93f5c5c8a8844a89f8db07eaa58270e3a1e01804f17f4cf8df02f5

ENV GOLANG_DOWNLOAD_URL=https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

# Install Go
RUN wget -q "$GOLANG_DOWNLOAD_URL" -O golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

# =============== INSTALL GOLANGCI-LINT ===============

FROM base AS install-golangci-lint

# Is unsafe better to ADD from local file
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.47.2

# =============== INSTALL PYTHON & PIP ===============

FROM base AS install-python

ENV PYTHON3_VERSION=3.9.5-r2 \
	PYTHON3_PIP_VERSION=20.3.4-r1

# Install Python3 and pip
RUN apk update \
	&& apk add --no-cache \
	python3=$PYTHON3_VERSION \
	py3-pip=$PYTHON3_PIP_VERSION

# =============== INSTALL YAMLLINT ===============

FROM install-python AS install-yamllint

ENV YAMLLINT_VERSION=1.27.1

RUN pip install --no-cache-dir \
	yamllint==$YAMLLINT_VERSION

# =============== INSTALL HASKELL & STACK ===============

FROM base AS install-haskell

ENV GHC_VERSION=8.10.7 \
	STACK_VERSION=2.7.5 \
	STACK_RESOLVER=lts-18.28 \
	GHCUP_VERSION=0.1.17.8 \
	GHCUP_DOWNLOAD_SHA256=d7c66d2c4762393b752634bd4a2d7deb4e22340a9d8af7bb95aa1ce6f6651235 \
	STACK_DOWNLOAD_SHA256=9bcd165358d4dcafd2b33320d4fe98ce72faaf62300cc9b0fb86a27eb670da50
ENV GHCUP_DOWNLOAD_URL=https://downloads.haskell.org/~ghcup/$GHCUP_VERSION/x86_64-linux-ghcup-$GHCUP_VERSION \
	STACK_DIRNAME=stack-$STACK_VERSION-linux-x86_64 \
	STACK_DOWNLOAD_URL=https://github.com/commercialhaskell/stack/releases/download/v$STACK_VERSION/stack-$STACK_VERSION-linux-x86_64.tar.gz

ENV GHCUP_INSTALL_BASE_PREFIX=/
ENV PATH=/.ghcup/bin:$PATH

# Install prerequisites
RUN apk update \
	&& apk add --no-cache \
	gcc \
	libc-dev \
	gmp-dev \
	xz \
	make \
	git \
	zlib \
	libffi

# Install ghcup
RUN wget -q "$GHCUP_DOWNLOAD_URL" -O /usr/local/bin/ghcup \
	&& echo "$GHCUP_DOWNLOAD_SHA256  /usr/local/bin/ghcup" | sha256sum -c - \
	&& chmod +x /usr/local/bin/ghcup \
	&& ghcup config set downloader Wget

# Install GHC
RUN ghcup install ghc $GHC_VERSION \
	&& ghcup set ghc $GHC_VERSION

# Install stack
RUN wget -q "$STACK_DOWNLOAD_URL" -O stack.tar.gz \
	&& echo "$STACK_DOWNLOAD_SHA256  stack.tar.gz" | sha256sum -c - \
	&& tar -xzf stack.tar.gz \
	&& cp -L $STACK_DIRNAME/stack /usr/local/bin/stack \
	&& rm -rf stack.tar.gz $STACK_DIRNAME

# Configure stack
RUN stack config set system-ghc true --global \
	&& stack config set install-ghc false --global

# =============== INSTALL HADOLINT ===============

FROM install-haskell AS install-hadolint

RUN git clone https://github.com/hadolint/hadolint \
	&& cd hadolint \
	&& stack install

RUN mv /root/.local/bin/hadolint /usr/local/bin/hadolint

# =============== FINAL IMAGE ===============

# otherwise doesnt work
FROM install-yamllint

COPY --from=install-go /usr/local/go /usr/local/go
COPY --from=install-golangci-lint /usr/local/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=install-hadolint /usr/local/bin/hadolint /usr/local/bin/hadolint
# COPY --from=install-yamllint /usr/bin/yamllint /usr/bin/yamllint

# Set up env
ENV GOPATH=$WORKSPACE/go\
	GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$GOPATH/bin/:$PATH

RUN mkdir -p $GOPATH \
	&& chown -R jenkins:jenkins $GOPATH

WORKDIR $WORKSPACE

# su jenkins
USER jenkins
