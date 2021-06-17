{{define "inboundInfoStream"}}
<p>传输: <a-tag color="green">[[ inbound.network ]]</a-tag></p>

<template v-if="inbound.isTcp || inbound.isWs || inbound.isH2">
    <p v-if="inbound.host">host: <a-tag color="green">[[ inbound.host ]]</a-tag></p>
    <p v-else>host: <a-tag color="orange">无</a-tag></p>

    <p v-if="inbound.path">path: <a-tag color="green">[[ inbound.path ]]</a-tag></p>
    <p v-else>path: <a-tag color="orange">无</a-tag></p>
</template>

<template v-if="inbound.isQuic">
    <p>quic 加密: <a-tag color="green">[[ inbound.quicSecurity ]]</a-tag></p>
    <p>quic 密码: <a-tag color="green">[[ inbound.quicKey ]]</a-tag></p>
    <p>quic 伪装: <a-tag color="green">[[ inbound.quicType ]]</a-tag></p>
</template>

<template v-if="inbound.isKcp">
    <p>kcp 加密: <a-tag color="green">[[ inbound.kcpType ]]</a-tag></p>
    <p>kcp 密码: <a-tag color="green">[[ inbound.kcpSeed ]]</a-tag></p>
</template>

<template v-if="inbound.isGrpc">
    <p>grpc serviceName: <a-tag color="green">[[ inbound.serviceName ]]</a-tag></p>
</template>

<template v-if="inbound.tls || inbound.xtls">
    <p v-if="inbound.tls">tls: <a-tag color="green">开启</a-tag></p>
    <p v-if="inbound.xtls">xtls: <a-tag color="green">开启</a-tag></p>
</template>
<template v-else>
    <p>tls: <a-tag color="red">关闭</a-tag></p>
</template>
<p v-if="inbound.tls">
    tls域名: <a-tag :color="inbound.serverName ? 'green' : 'orange'">[[ inbound.serverName ? inbound.serverName : "无" ]]</a-tag>
</p>
<p v-if="inbound.xtls">
    xtls域名: <a-tag :color="inbound.serverName ? 'green' : 'orange'">[[ inbound.serverName ? inbound.serverName : "无" ]]</a-tag>
</p>
{{end}}


{{define "component/inboundInfoComponent"}}
<div>
    <p>协议: <a-tag color="green">[[ dbInbound.protocol ]]</a-tag></p>
    <p>地址: <a-tag color="blue">[[ dbInbound.address ]]</a-tag></p>
    <p>端口: <a-tag color="green">[[ dbInbound.port ]]</a-tag></p>

    <template v-if="dbInbound.isVMess">
        <p>uuid: <a-tag color="green">[[ inbound.uuid ]]</a-tag></p>
        <p>alterId: <a-tag color="green">[[ inbound.alterId ]]</a-tag></p>
    </template>

    <template v-if="dbInbound.isVLess">
        <p>uuid: <a-tag color="green">[[ inbound.uuid ]]</a-tag></p>
        <p v-if="inbound.isXTls">flow: <a-tag color="green">[[ inbound.flow ]]</a-tag></p>
    </template>

    <template v-if="dbInbound.isTrojan">
        <p>密码: <a-tag color="green">[[ inbound.password ]]</a-tag></p>
    </template>

    <template v-if="dbInbound.isSS">
        <p>加密: <a-tag color="green">[[ inbound.method ]]</a-tag></p>
        <p>密码: <a-tag color="green">[[ inbound.password ]]</a-tag></p>
    </template>

    <template v-if="dbInbound.isSocks">
        <p>用户名: <a-tag color="green">[[ inbound.username ]]</a-tag></p>
        <p>密码: <a-tag color="green">[[ inbound.password ]]</a-tag></p>
    </template>

    <template v-if="dbInbound.isHTTP">
        <p>用户名: <a-tag color="green">[[ inbound.username ]]</a-tag></p>
        <p>密码: <a-tag color="green">[[ inbound.password ]]</a-tag></p>
    </template>

    <template v-if="dbInbound.isVMess || dbInbound.isVLess || dbInbound.isTrojan || dbInbound.isSS">
        {{template "inboundInfoStream"}}
    </template>
</div>
{{end}}

{{define "component/inboundInfo"}}
<script>
    Vue.component('inbound-info', {
        delimiters: ['[[', ']]'],
        props: ["dbInbound", "inbound"],
        template: `{{template "component/inboundInfoComponent"}}`,
    });
</script>
{{end}}