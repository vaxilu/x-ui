FROM centos:7.7.1908

LABEL maintainer="FaintGhost <zhang.yaowei@live.com>"

ENV XUI_VERSION=0.3.2

COPY systemctl.py /usr/bin/systemctl

RUN cd /root \
    &&  yum install wget curl tar -y \
    &&  wget --no-check-certificate -c https://github.com/vaxilu/x-ui/releases/download/${XUI_VERSION}/x-ui-linux-amd64.tar.gz \
    &&  tar zxvf x-ui-linux-amd64.tar.gz  \
    &&  chmod +x x-ui/x-ui x-ui/bin/xray-linux-* x-ui/x-ui.sh \
    &&  cp x-ui/x-ui.sh /usr/bin/x-ui \
    &&  cp -f x-ui/x-ui.service /etc/systemd/system/ \
    &&  mv x-ui/ /usr/local/ \
    &&  systemctl daemon-reload 

RUN systemctl enable x-ui

EXPOSE 54321

ENTRYPOINT /usr/bin/systemctl