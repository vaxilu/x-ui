#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

cur_dir=$(pwd)

# check root
[[ $EUID -ne 0 ]] && echo -e "${red}错误：${plain} 必须使用root用户运行此脚本！\n" && exit 1

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
    echo -e "${red}Versão do sistema não detectada, entre em contato com o autor do script!${plain}\n" && exit 1
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

echo "架构: ${arch}"

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
        echo -e "${red}Por favor, use o CentOS 7 ou superior!${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"ubuntu" ]]; then
    if [[ ${os_version} -lt 16 ]]; then
        echo -e "${red}Por favor, use o Ubuntu 16 ou superior!${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"debian" ]]; then
    if [[ ${os_version} -lt 8 ]]; then
        echo -e "${red}Por favor, use o Debian 8 ou superior!${plain}\n" && exit 1
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
        last_version=$(curl -Ls "https://api.github.com/repos/dutra01/x-ui/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [[ ! -n "$last_version" ]]; then
            echo -e "${red}Falha ao detectar a versão x-ui, pode ser que o limite da API do Github tenha sido excedido. Tente novamente mais tarde ou especifique manualmente a versão x-ui a ser instalada${plain}"
            exit 1
        fi
        echo -e "检测到 x-ui 最新版本：${last_version}，开始安装"
        wget -N --no-check-certificate -O /usr/local/x-ui-linux-${arch}.tar.gz https://github.com/dutra01/x-ui/releases/download/${last_version}/x-ui-linux-${arch}.tar.gz
        if [[ $? -ne 0 ]]; then
            echo -e "${red}Falha ao baixar x-ui, certifique-se de que seu servidor pode baixar arquivos do Github${plain}"
            exit 1
        fi
    else
        last_version=$1
        url="https://github.com/dutra01/x-ui/releases/download/${last_version}/x-ui-linux-${arch}.tar.gz"
        echo -e "开始安装 x-ui v$1"
        wget -N --no-check-certificate -O /usr/local/x-ui-linux-${arch}.tar.gz ${url}
        if [[ $? -ne 0 ]]; then
            echo -e "${red}Falha ao baixar x-ui v$1, verifique se esta versão existe${plain}"
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

    wget --no-check-certificate -O /usr/bin/x-ui https://raw.githubusercontent.com/DuTra01/x-ui/main/x-ui.sh
    chmod +x /usr/local/x-ui/x-ui.sh
    chmod +x /usr/bin/x-ui
    config_after_install
    
    #echo -e "如果是全新安装，默认网页端口为 ${green}54321${plain}，用户名和密码默认都是 ${green}admin${plain}"
    #echo -e "请自行确保此端口没有被其他程序占用，${yellow}并且确保 54321 端口已放行${plain}"
    #    echo -e "若想将 54321 修改为其它端口，输入 x-ui 命令进行修改，同样也要确保你修改的端口也是放行的"
    #echo -e ""
    #echo -e "如果是更新面板，则按你之前的方式访问面板"
    #echo -e ""
    systemctl daemon-reload
    systemctl enable x-ui
    systemctl start x-ui
    echo -e "${green}x-ui v${last_version}${plain} A instalação está concluída, o painel é lançado,"
    echo -e ""
    echo -e "x-ui Como usar o script de gerenciamento: "
    echo -e "----------------------------------------------"
    echo -e "x-ui              - Exibe o menu de gerenciamento (com mais funções)"
    echo -e "x-ui start        - Inicie o painel x-ui"
    echo -e "x-ui stop         - Parar painel x-ui"
    echo -e "x-ui restart      - Reinicie o painel x-ui"
    echo -e "x-ui status       - Ver o status do x-ui"
    echo -e "x-ui enable       - Defina o x-ui para iniciar automaticamente na inicialização"
    echo -e "x-ui disable      - Cancelar inicialização automática de inicialização x-ui"
    echo -e "x-ui log          - Ver registros x-ui"
    echo -e "x-ui v2-ui        - Migre os dados da conta v2-ui desta máquina para x-ui"
    echo -e "x-ui update       - Atualizar painel x-ui"
    echo -e "x-ui install      - Instale o painel x-ui"
    echo -e "x-ui uninstall    - desinstalar o painel x-ui"
    echo -e "----------------------------------------------"
}

echo -e "${green}iniciar a instalação${plain}"
install_base
install_x-ui $1
