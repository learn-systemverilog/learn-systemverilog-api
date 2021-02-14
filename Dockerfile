FROM verilator/verilator:4.108

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive \
    && apt-get install --no-install-recommends -y wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ARG EMSCRIPTEN_VERSION=2.0.13

RUN git clone --depth 1 --branch ${EMSCRIPTEN_VERSION} https://github.com/emscripten-core/emsdk.git
RUN cd emsdk \
    && ./emsdk install ${EMSCRIPTEN_VERSION} \
    && ./emsdk activate ${EMSCRIPTEN_VERSION} \
    && echo 'source "/work/emsdk/emsdk_env.sh"' >> $HOME/.bashrc

ARG GOLANG_VERSION=1.15.8

RUN wget https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
    && rm go${GOLANG_VERSION}.linux-amd64.tar.gz \
    && echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.bashrc

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR learn-systemverilog-api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENTRYPOINT ["bash"]
