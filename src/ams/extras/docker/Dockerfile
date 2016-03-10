# set author and base
FROM centos:centos7
MAINTAINER Luis Pabón <lpabon@redhat.com>

LABEL version="0.1"
LABEL description="Centos 7 docker image for ams"

# let's setup all the necessary environment variables
ENV BUILD_HOME=/build
ENV GOPATH=$BUILD_HOME/golang
ENV PATH=$GOPATH/bin:$PATH
ENV HEKETIC_CONF_DIR=/etc/ams

# install dependencies, build and cleanup
RUN mkdir $BUILD_HOME $GOPATH $HEKETI_CONF_DIR && \
    yum -q -y install golang git && \
    yum -q -y install make && \
    yum -q -y clean all && \
    cd $BUILD_HOME && \
    go get github.com/ams/ams && \
    go get github.com/robfig/glock && \
    glock sync github.com/ams/ams && \
    cd $GOPATH/src/github.com/ams/ams && make && \
    cp ams /usr/bin/ams && \
    cd && rm -rf $BUILD_HOME && \
    yum -q -y remove git golang make && \
    yum -q -y autoremove && \
    yum -q -y clean all

# post install config and volume setup
VOLUME /etc/ams

# expose port, set user and set entrypoint with config option
ENTRYPOINT ["/usr/bin/ams"]
EXPOSE 8080

CMD ["-config=/etc/aplo/ams.json"]
