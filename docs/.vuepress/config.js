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
                title: "Command",
                collapsable: false,
                children: [
                    ["/command/qoscli", "qoscli"],
                    ["/command/qosd", "qosd"],
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
                    ["/spec/account", "Account"],
                    ["/spec/genesis", "Genesis"],
                    ["/spec/transfer", "Transfer"],
                    ["/spec/approve", "Approve"],
                    ["/spec/governance", "Governance"],
                    ["/spec/guardian", "Guardian"],
                    ["/spec/indexing", "Index"],
                    ["/spec/qsc", "QSC"],
                    ["/spec/qcp", "QCP"],
                    ["/spec/ca", "CA"]
                ]
            }
        ]
    }
}
