FROM ubuntu:22.04

WORKDIR /

ARG USERNAME=vpnuser
ARG PASSWORD=qweasdQWEASD
ARG PORT=9007

RUN apt-get update && apt-get install wget curl tar -y

COPY . .
RUN chmod +x entrypoint.sh

RUN apt-get update && \
    apt-get install -yq tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Tehran /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata

WORKDIR /usr/local

RUN wget -N --no-check-certificate -O /usr/local/x-ui-linux-amd64.tar.gz https://github.com/vaxilu/x-ui/releases/download/0.3.2/x-ui-linux-amd64.tar.gz
RUN tar zxvf x-ui-linux-amd64.tar.gz
RUN rm x-ui-linux-amd64.tar.gz -f
RUN chmod +x /usr/local/x-ui/x-ui
RUN chmod +x /usr/local/x-ui/bin/xray-linux-amd64

RUN wget --no-check-certificate -O /usr/bin/x-ui https://raw.githubusercontent.com/vaxilu/x-ui/main/x-ui.sh
RUN chmod +x /usr/local/x-ui/x-ui.sh
RUN chmod +x /usr/bin/x-ui

RUN /usr/local/x-ui/x-ui setting -username $USERNAME -password $PASSWORD
RUN /usr/local/x-ui/x-ui setting -port $PORT

WORKDIR /

CMD ["./entrypoint.sh"]
