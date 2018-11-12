module.exports = {
    title: "QOS",
    description: "Documentation for the QOS.",
    dest: "./dist/docs",
    base: "/docs/",
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        lastUpdated: "Last Updated",
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
                    ["/start/installation", "Install"],
                    ["/start/networks", "Networks"],
                    ["/start/testnet", "Testnet"]
                ]
            },
            {
                title: "Client",
                collapsable: false,
                children: [
                    ["/client/txs", "TXs"],
                    ["/client/qcp", "QCP"]
                ]
            }
            ,
            {
                title: "SPEC",
                collapsable: false,
                children: [
                    ["/spec/account", "Account"],
                    ["/spec/genesis", "Genesis"],
                    ["/spec/txs/transfer", "Transfer"],
                    ["/spec/txs/approve", "Approve"],
                    ["/spec/txs/createqsc_issue_cli", "CreateQSCIssueQSCcli"],
                    ["/spec/txs/createqsc_issue_design", "CreateQSCIssueQSCdesign"]
                ]
            }
        ]
    }
}
