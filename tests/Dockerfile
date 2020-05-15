FROM fedora
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf install -y git golang make gcc gcc-c++ which
ARG REPOFOLDER
ADD $REPOFOLDER /peggy
RUN pushd /peggy/module/ && make && make install
RUN pushd /peggy/ && bash tests/setup-validators.sh
CMD peggyd start
