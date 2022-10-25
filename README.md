<p align="center">
 <h1 align="center">X-UI [PT-BR]</h1>
 <p align="center">Painel de Xray com suporte a multiusuário e multiprotocolo</p>
</p>
  <br/>
  <p align="center">
    <a href="https://github.com/TelksBr/x-ui_br/releases">
      <img alt="releases" src="https://img.shields.io/github/downloads/telksbr/x-ui_br/total.svg" />
    </a>
    <a href="https://github.com/TelksBr/x-ui_br/network/members">
      <img alt="forks" src="https://img.shields.io/github/forks/telksbr/x-ui_br.svg" />
    </a>
    <a href="https://github.com/anuraghazra/github-readme-stats/issues">
      <img alt="stars" src="https://img.shields.io/github/stars/telksbr/x-ui_br.svg" />
    </a>
    <a href="https://github.com/TelksBr/x-ui_br/watchers">
      <img alt="watchers" src="https://img.shields.io/github/watchers/telksbr/x-ui_br.svg" />
    </a>
    <a href="https://github.com/TelksBr/x-ui_br/watchers">
      <img alt="pr-closed" src="https://img.shields.io/github/issues-pr-closed/telksbr/x-ui_br.svg" />
    </a>
    <a href="https://github.com/TelksBr/">
      <img alt="mainteined" src="https://img.shields.io/badge/Maintained%3F-yes-green.svg" />
    </a>
  </p>

# Características e Funções

- Monitoramento do estado do sistema
- Suporte multi-protocolo multiusuário, modo de visualização em página da web
- Protocolos suportados: vmess, vless, trojan, shadowsocks, dokodemo-door, socks, http
- Suporte para configurar mais configurações de transmissão
- Estatísticas de tráfego, limite de tráfego, limite de tempo de expiração
- Modelos de configuração de x-ray personalizáveis
- Suporte ao painel de acesso https (traga seu próprio nome de domínio + certificado ssl)
- Suporte a aplicação de certificado SSL com um clique e renovação automática
- Para itens de configuração mais avançados, consulte o painel para obter detalhes



## Sistemas Operacionais Suportados

- CentOS 7+
- Ubuntu 16+
- Debian 8+



> O script funciona no Oracle Linux para Arm64, porém não é confirmado como estável.

# Instalar e Atualizar (recomendado)

````
bash <(curl -Ls https://raw.githubusercontent.com/TelksBr/x-ui_br/main/install.sh)
````

## Instalação e atualização manuais

1. Primeiro baixe o pacote compactado mais recente de [releases](https://github.com/TelksBr/x-ui_br/releases), geralmente escolha a arquitetura `amd64`
2. Em seguida, carregue o pacote compactado para o diretório `/root/` do servidor e use o usuário `root` para efetuar login no servidor

> Se a arquitetura da CPU do seu servidor não for `amd64`, substitua `amd64` no comando por outra arquitetura

````
cd /root/
rm x-ui/ /usr/local/x-ui/ /usr/bin/x-ui -rf
tar zxvf x-ui-linux-amd64.tar.gz
chmod +x x-ui/x-ui x-ui/bin/xray-linux-* x-ui/x-ui.sh
cp x-ui/x-ui.sh /usr/bin/x-ui
cp -f x-ui/x-ui.service /etc/systemd/system/
mv x-ui/ /usr/local/
systemctl daemon-reload
systemctl habilitar x-ui
systemctl reiniciar x-ui
````

## Instale usando o docker (atualmente cria o painel em chinês)

> Este tutorial do docker e a imagem do docker são fornecidos por [Chasing66](https://github.com/Chasing66)

1. Instale o Docker

```shell
curl -fsSL https://get.docker.com | sh
````

2. Instale o x-ui

```shell
mkdir x-ui && cd x-ui
docker run -itd --network=host \
    -v $PWD/db/:/etc/x-ui/ \
    -v $PWD/cert/:/root/cert/ \
    --name x-ui --restart=unless-stopped \
    enwaiax/x-ui:latest
````

> Construa sua própria imagem

```shell
docker build -t x-ui .
````

## Aplicação de certificado SSL (quebrado)

> Esta função e tutorial são fornecidos por [FranzKafkaYu](https://github.com/FranzKafkaYu)

O script tem uma função de aplicativo de certificado SSL integrada. Para usar este script para solicitar um certificado, as seguintes condições devem ser atendidas:

- Conheça o endereço de e-mail registrado na Cloudflare.
- Conheça a chave de API global da Cloudflare.
- O nome de domínio foi resolvido para o servidor atual por meio do cloudflare.


Ao usar, basta digitar `domain name`, `email`, `API KEY`.

Precauções:

- O script usa a API DNS para solicitação de certificado
- Use Let'sEncrypt como a parte CA por padrão
- O diretório de instalação do certificado é o diretório /root/cert
- Os certificados solicitados por este script são todos os certificados de nome de domínio genérico

## Telegram Bot (em desenvolvimento, temporariamente indisponível)

> Esta função e tutorial são fornecidos por [FranzKafkaYu](https://github.com/FranzKafkaYu)

X-UI suporta notificação diária de tráfego, lembrete de login do painel e outras funções através do robô Tg. Para usar o robô Tg, você precisa se inscrever por conta própria.

Para tutoriais de aplicativos específicos, consulte [link do blog](https://coderfan.net/how-to-use-telegram-bot-to-alarm-you-when-someone-login-into-your-vps.html )

Instruções de uso: Defina os parâmetros relacionados ao robô no fundo do painel, incluindo:

- Token do Bot no Telegram
- ID do chat do Bot no Telegram
- Tempo de execução do ciclo do Bot Telegram, na sintaxe do crontab

Sintaxe de referência:
- `30 * * * * *` => Notificar aos 30s de cada ponto
- `@hourly` => Notificações de hora em hora
- `@daily` => Notificação diária (00:00 AM)
- `@every 8h` => Notificação a cada 8 horas

Conteúdo da notificação no Telegram:
- Uso de tráfego de conexões
- Lembrete de login do painel
- Lembrete de expiração da conexão
- Alertas de tráfego

Mais recursos estão planejados...

## Agradecimentos

[Niduka Akalanka](https://github.com/NidukaAkalanka)

[othmx](https://github.com/othmx) (pela compilação)


