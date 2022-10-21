#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

cur_dir=$(pwd)

# check root
[[ $EUID -ne 0 ]] && echo -e "${red}erro: ${plain} Este script deve ser executado como usuário root!\n" && exit 1

# check os
if [[ -f /etc/redhat-release ]]; then
    release="centos"
elif cat /etc/issue | grep -Eqi "debian"; then
    release="debian"
elif cat /etc/issue | grep -Eqi "ubuntu"; then
    release="ubuntu"
elif cat /etc/issue | grep -Eqi "centos|red hat|redhat"; then
    release="centos"
elif cat /proc/version | grep -Eqi "debian"; then
    release="debian"
elif cat /proc/version | grep -Eqi "ubuntu"; then
    release="ubuntu"
elif cat /proc/version | grep -Eqi "centos|red hat|redhat"; then
    release="centos"
else
    echo -e "${red}image.png！${plain}\n" && exit 1
fi

arch=$(arch)

if [[ $arch == "x86_64" || $arch == "x64" || $arch == "amd64" ]]; then
    arch="amd64"
elif [[ $arch == "aarch64" || $arch == "arm64" ]]; then
    arch="arm64"
elif [[ $arch == "s390x" ]]; then
    arch="s390x"
else
    arch="amd64"
    echo -e "${red}Falha ao detectar o esquema, use o esquema padrão: ${arch}${plain}"
fi

echo "Arquitetura: ${arch}"

if [ $(getconf WORD_BIT) != '32' ] && [ $(getconf LONG_BIT) != '64' ]; then
    echo "Este software não suporta sistema de 32 bits (x86), use o sistema de 64 bits (x86_64), se a detecção estiver errada, entre em contato com o autor"
    exit -1
fi

os_version=""

# os version
if [[ -f /etc/os-release ]]; then
    os_version=$(awk -F'[= ."]' '/VERSION_ID/{print $3}' /etc/os-release)
fi
if [[ -z "$os_version" && -f /etc/lsb-release ]]; then
    os_version=$(awk -F'[= ."]+' '/DISTRIB_RELEASE/{print $2}' /etc/lsb-release)
fi

if [[ x"${release}" == x"centos" ]]; then
    if [[ ${os_version} -le 6 ]]; then
        echo -e "${red}Por favor, use o CentOS 7 ou superior! ${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"ubuntu" ]]; then
    if [[ ${os_version} -lt 16 ]]; then
        echo -e "${red}Por favor, use o Ubuntu 16 ou superior! ${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"debian" ]]; then
    if [[ ${os_version} -lt 8 ]]; then
        echo -e "${red}Por favor, use o Debian 8 ou superior! ${plain}\n" && exit 1
    fi
fi

install_base() {
    if [[ x"${release}" == x"centos" ]]; then
        yum install wget curl tar -y
    else
        apt install wget curl tar -y
    fi
}

#This function will be called when user installed x-ui out of sercurity
config_after_install() {
    echo -e "${yellow}Por motivos de segurança, é necessário modificar à força a senha da porta e da conta após a conclusão da instalação/atualização.${plain}"
    read -p "Confirme se deseja continuar?[y/n]": config_confirm
    if [[ x"${config_confirm}" == x"y" || x"${config_confirm}" == x"Y" ]]; then
        read -p "Por favor, defina o nome da sua conta:" config_account
        echo -e "${yellow}O nome da sua conta será definido como:${config_account}${plain}"
        read -p "Por favor, defina a senha da sua conta:" config_password
        echo -e "${yellow}A senha da sua conta será definida como:${config_password}${plain}"
        read -p "Por favor, defina a porta de acesso ao painel:" config_port
        echo -e "${yellow}Sua porta de acesso ao painel será configurada para:${config_port}${plain}"
        echo -e "${yellow}Confirmar configuração, configuração${plain}"
        /usr/local/x-ui/x-ui setting -username ${config_account} -password ${config_password}
        echo -e "${yellow}Configuração de senha da conta concluída${plain}"
        /usr/local/x-ui/x-ui setting -port ${config_port}
        echo -e "${yellow}Configuração da porta do painel concluída${plain}"
    else
        echo -e "${red}Cancelado, todos os itens de configuração são configurações padrão, modifique a tempo${plain}"
    fi
}

install_x-ui() {
    systemctl stop x-ui
    cd /usr/local/

    if [ $# == 0 ]; then
        last_version=$(curl -Ls "https://api.github.com/repos/TelksBr/x-ui_br/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [[ ! -n "$last_version" ]]; then
            echo -e "${red}Falha ao detectar a versão x-ui, pode ser que o limite da API do Github tenha sido excedido. Tente novamente mais tarde ou especifique manualmente a versão x-ui a ser instalada${plain}"
            exit 1
        fi
        echo -e "x-ui versão mais recente detectada：${last_version}，iniciar a instalação"
        wget -N --no-check-certificate -O /usr/local/x-ui-linux-${arch}.tar.gz https://github.com/TelksBr/x-ui_br/releases/download/${last_version}/x-ui-linux-${arch}.tar.gz
        if [[ $? -ne 0 ]]; then
            echo -e "${red}Falha ao baixar x-ui, certifique-se de que seu servidor pode baixar arquivos do Github${plain}"
            exit 1
        fi
    else
        last_version=$1
        url="https://github.com/TelksBr/x-ui_br/releases/download/${last_version}/x-ui-linux-${arch}.tar.gz"
        echo -e "iniciar a instalação x-ui v$1"
        wget -N --no-check-certificate -O /usr/local/x-ui-linux-${arch}.tar.gz ${url}
        if [[ $? -ne 0 ]]; then
            echo -e "${red}baixar x-ui v$1 falhou, verifique se esta versão existe${plain}"
            exit 1
        fi
    fi

    if [[ -e /usr/local/x-ui/ ]]; then
        rm /usr/local/x-ui/ -rf
    fi

    tar zxvf x-ui-linux-${arch}.tar.gz
    rm x-ui-linux-${arch}.tar.gz -f
    cd x-ui
    chmod +x x-ui bin/xray-linux-${arch}
    cp -f x-ui.service /etc/systemd/system/
    wget --no-check-certificate -O /usr/bin/x-ui https://raw.githubusercontent.com/TelksBr/x-ui_br/main/x-ui.sh
    chmod +x /usr/local/x-ui/x-ui.sh
    chmod +x /usr/bin/x-ui
    config_after_install
    #echo -e "Se for uma instalação nova, a porta web padrão é ${green}54321${plain}, e o nome de usuário e senha padrão são ${green}admin${plain}"
     #echo -e "Por favor, certifique-se de que esta porta não esteja ocupada por outros programas, ${yellow} e certifique-se de que a porta 54321 foi liberada ${plain}"
     # echo -e "Se você quiser modificar 54321 para outra porta, digite o comando x-ui para modificá-lo e também certifique-se de que a porta que você modificou também seja permitida"
     #echo -e ""
     #echo -e "Se atualizar o painel, acesse o painel como você fez antes"
     #echo -e ""
    systemctl daemon-reload
    systemctl enable x-ui
    systemctl start x-ui
    echo -e "${green}x-ui v${last_version}${plain}A instalação está concluída, o painel é lançado,"
    echo -e ""
    echo -e "x-ui Como usar o script de gerenciamento: "
    echo -e "----------------------------------------------"
    echo -e "x-ui              - Mostrar menu de gerenciamento (mais funções)"
    echo -e "x-ui start        - Inicie o painel x-ui"
    echo -e "x-ui stop         - parar painel x-ui"
    echo -e "x-ui restart      - reinicie o painel x-ui"
    echo -e "x-ui status       - Ver o status do x-ui"
    echo -e "x-ui enable       - Defina o x-ui para iniciar automaticamente na inicialização"
    echo -e "x-ui disable      - Cancelar inicialização automática de inicialização x-ui"
    echo -e "x-ui log          - Ver registros x-ui"
    echo -e "x-ui v2-ui        - Migre os dados da conta v2-ui desta máquina para x-ui"
    echo -e "x-ui update       - Atualize o painel x-ui"
    echo -e "x-ui install      - Instale o painel x-ui"
    echo -e "x-ui uninstall    - Desinstale o painel x-ui"
    echo -e "----------------------------------------------"
}

echo -e "${green}iniciar a instalação${plain}"
install_base
install_x-ui $1
