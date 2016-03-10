# set author and base
FROM centos:centos7
MAINTAINER Luis Pabón <lpabon@redhat.com>

LABEL version="0.1"
LABEL description="Centos 7 docker image for ams"

# let's setup all the necessary environment variables
ENV BUILD_HOME=/build
ENV GOPATH=$BUILD_HOME/golang
ENV PATH=$GOPATH/bin:$PATH
ENV APLO_CONF_DIR=/etc/aplo

# install dependencies, build and cleanup
RUN mkdir $BUILD_HOME $GOPATH $APLO_CONF_DIR && \
    yum -q -y install golang git mercurial && \
    yum -q -y install make && \
    yum -q -y clean all && \
    cd $BUILD_HOME && \
    mkdir -p $GOPATH/src/github.com/gluster && \
    cd $GOPATH/src/github.com/gluster && \
    git clone https://github.com/gluster/aplo.git

RUN go get github.com/robfig/glock && \
    glock sync github.com/gluster/aplo/src/ams

RUN cd $GOPATH/src/github.com/aplo/src/ams && go build && \
    cp ams /usr/bin/ams && \
    cd && rm -rf $BUILD_HOME && \
    yum -q -y remove git golang make && \
    yum -q -y autoremove && \
    yum -q -y clean all

ADD ./ams.json /etc/aplo/ams.json

# post install config and volume setup
VOLUME /etc/aplo

# expose port, set user and set entrypoint with config option
ENTRYPOINT ["/usr/bin/ams"]
EXPOSE 8080

CMD ["-config=/etc/aplo/ams.json"]
