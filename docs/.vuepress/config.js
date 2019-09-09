module.exports = {
    dest: "./dist/docs",
    base: "/qos/",
    markdown: {
        lineNumbers: true
    },
    locales: {
        '/' : {
            lang: '简体中文',
            selectText: '选择语言',
            label: '简体中文',
            title: 'QOS官方文档',
            description: 'QOS官方文档'
        },
        '/en/' : {
            lang: 'English',
            selectText: 'Languages',
            label: 'English',
            title: 'QOS Documentation',
            description: 'Documentation for the QOS.'
        },
    },
    themeConfig: {
        sidebarDepth: 3,
        locales: {
            '/' : {
                selectText: '选择语言',
                nav: [{text: "QOS官网", link: "https://www.qoschain.io/"},{text: "白皮书", link: "https://github.com/QOSGroup/whitepaper"}],
                sidebar: [
                    {
                        title: "简介",
                        collapsable: false,
                        children: [
                            ["/introduction/qos", "QOS"]
                        ]
                    },
                    {
                        title: "开始",
                        collapsable: false,
                        children: [
                            ["/install/installation", "安装"],
                            ["/install/networks", "本地运行"],
                            ["/install/testnet", "测试网络"]
                        ]
                    },
                    {
                        title: "功能",
                        collapsable: false,
                        children: [
                            ["/command/qoscli", "qoscli 命令集"],
                            ["/command/qosd", "qosd 命令集"],
                        ]
                    }
                    ,
                    {
                        title: "设计",
                        collapsable: false,
                        children: [
                            ["/spec/eco_module.v1", "经济模型"],
                            ["/spec/validators/all_about_validators", "验证节点详解"],
                            ["/spec/staking", "验证节点"],
                            ["/spec/bank/", "账户"],
                            ["/spec/genesis", "创世"],
                            ["/spec/approve/", "预授权"],
                            ["/spec/governance", "治理"],
                            ["/spec/guardian", "系统用户"],
                            ["/spec/indexing", "索引"],
                            ["/spec/qsc", "代币"],
                            ["/spec/qcp", "联盟链"],
                            ["/spec/ca", "证书"]
                        ]
                    }
                ]
            },
            '/en/' : {
                selectText: 'Languages',
                nav: [{text: "Back to QOS.", link: "https://www.qoschain.io/"},{text: "White Paper", link: "https://github.com/QOSGroup/whitepaper"}],
                sidebar: [
                    {
                        title: "Introduction",
                        collapsable: false,
                        children: [
                            ["/en/introduction/qos", "QOS"]
                        ]
                    },
                    {
                        title: "Getting Started",
                        collapsable: false,
                        children: [
                            ["/en/install/installation", "Install"],
                            ["/en/install/networks", "Networks"],
                            ["/en/install/testnet", "Testnet"]
                        ]
                    },
                    {
                        title: "Command",
                        collapsable: false,
                        children: [
                            ["/en/command/qoscli", "qoscli"],
                            ["/en/command/qosd", "qosd"],
                        ]
                    }
                    ,
                    {
                        title: "Spec",
                        collapsable: false,
                        children: [
                            ["/spec/eco_module.v1", "ECO"],
                            ["/spec/validators/all_about_validators", "Validators"],
                            ["/spec/staking", "Stake"],
                            ["/spec/bank/", "Bank"],
                            ["/spec/genesis", "Genesis"],
                            ["/en/spec/approve/", "Approve"],
                            ["/spec/governance", "Governance"],
                            ["/spec/guardian", "Guardian"],
                            ["/spec/indexing", "Index"],
                            ["/spec/qsc", "QSC"],
                            ["/spec/qcp", "QCP"],
                            ["/spec/ca", "CA"]
                        ]
                    }
                ]
            },
        }
    }
}


