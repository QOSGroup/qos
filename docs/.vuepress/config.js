module.exports = {
    title: "QOS",
    description: "Documentation for the QOS.",
    dest: "./dist/docs",
    base: "/qos/",
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        nav: [{text: "Back to QOS", link: "https://www.github.com/QOSGroup/qos"}],
        sidebar: [
            {
                title: "Introduction",
                collapsable: false,
                children: [
                    ["/introduction/qos", "QOS"]
                ]
            },
            {
                title: "Getting Started",
                collapsable: false,
                children: [
                    ["/install/installation", "Install"],
                    ["/install/networks", "Networks"],
                    ["/install/testnet", "Testnet"]
                ]
            },
            {
                title: "Client",
                collapsable: false,
                children: [
                    ["/client/qoscli", "qoscli"],
                    ["/client/qosd", "qosd"],
                ]
            }
            ,
            {
                title: "Spec",
                collapsable: false,
                children: [
                    ["/spec/txs/staking", "Stake"],
                    ["/spec/account", "Account"],
                    ["/spec/genesis", "Genesis"],
                    ["/spec/txs/transfer", "Transfer"],
                    ["/spec/txs/approve", "Approve"],
                    ["/spec/txs/qsc", "QSC"],
                    ["/spec/txs/qcp", "QCP"],
                    ["/spec/ca", "CA"]
                ]
            }
        ]
    }
}
