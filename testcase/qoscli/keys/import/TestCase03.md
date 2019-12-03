# Description
```
导入密钥库中不存在的密钥
```
# Input
```
$ qoscli keys import test02
```
# Output
若私钥和密码均输入正确，导入成功：
```
$ qoscli keys import test02
> Enter ed25519 private key:
dzQ2ii+7KxVLzduw3PusyszjCtF/hgovYo+x4+ugfT7vRSenE/hvfbYElit/9lwwGT56J+c02RRzmQpIEnEzRQ==
> Enter a passphrase for your key:<输入密码>
> Repeat the passphrase:<再次输入密码>

```
若私钥输入出错，导入失败： 
``` 
qoscli keys import test02
> Enter ed25519 private key:
70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U=
> Enter a passphrase for your key:
> Repeat the passphrase:
panic: Expected PrivKeyEd25519 to include concatenated pubkey bytes

goroutine 1 [running]:
github.com/tendermint/tendermint/crypto/ed25519.PrivKeyEd25519.PubKey(0x7d6ff813a72745ef, 0x305cf67f2b9604b6, 0x14d934e7277a3e19, 0x45337112480a9973, 0x0, 0x0, 0x0, 0x0, 0x10, 0x10)
        C:/Users/wzj_s/go/pkg/mod/github.com/tendermint/tendermint@v0.32.0/crypto/ed25519/ed25519.go:75 +0x13a
github.com/QOSGroup/qbase/keys.Keybase.writeImportInfoKey(0x116b1c0, 0xc00000a148, 0xc0000bfc70, 0x115d4e0, 0xc0003960c0, 0xc000074088, 0x6, 0xc00002aad4, 0x8, 0x0, ...)
        C:/Users/wzj_s/go/pkg/mod/github.com/!q!o!s!group/qbase@v0.2.2-0.20190725073544-9c9f4bb8ffbf/keys/keybase.go:416 +0xcb
github.com/QOSGroup/qbase/keys.Keybase.CreateImportInfo(...)
        C:/Users/wzj_s/go/pkg/mod/github.com/!q!o!s!group/qbase@v0.2.2-0.20190725073544-9c9f4bb8ffbf/keys/keybase.go:97
github.com/QOSGroup/qbase/client/keys.importCommand.func1(0xc0009a4780, 0xc00099a6f0, 0x1, 0x1, 0x0, 0x0)
        C:/Users/wzj_s/go/pkg/mod/github.com/!q!o!s!group/qbase@v0.2.2-0.20190725073544-9c9f4bb8ffbf/client/keys/import.go:80 +0x4a0
github.com/spf13/cobra.(*Command).execute(0xc0009a4780, 0xc00099a6b0, 0x1, 0x1, 0xc0009a4780, 0xc00099a6b0)
        C:/Users/wzj_s/go/pkg/mod/github.com/spf13/cobra@v0.0.3/command.go:762 +0x46c
github.com/spf13/cobra.(*Command).ExecuteC(0x1878080, 0xb873a6, 0x1878080, 0xeaf6ac)
        C:/Users/wzj_s/go/pkg/mod/github.com/spf13/cobra@v0.0.3/command.go:852 +0x2f3
github.com/spf13/cobra.(*Command).Execute(...)
        C:/Users/wzj_s/go/pkg/mod/github.com/spf13/cobra@v0.0.3/command.go:800
github.com/tendermint/tendermint/libs/cli.Executor.Execute(0x1878080, 0x103a8a8, 0x3, 0xeb7142)
        C:/Users/wzj_s/go/pkg/mod/github.com/tendermint/tendermint@v0.32.0/libs/cli/setup.go:89 +0x43
main.main()
        C:/Users/wzj_s/go/src/qos/cmd/qoscli/main.go:81 +0x677
```
若密码输入错误，导出失败:
- 两次输入不匹配
```
$ qoscli keys import test03
> Enter ed25519 private key:
dzQ2ii+7KxVLzduw3PusyszjCtF/hgovYo+x4+ugfT7vRSenE/hvfbYElit/9lwwGT56J+c02RRzmQpIEnEzRQ==
> Enter a passphrase for your key:<输入密码>
> Repeat the passphrase:<输入密码>
ERROR: passphrases don't match
```
- 密码长度小于8位
```
qoscli keys import test03
> Enter ed25519 private key:
dzQ2ii+7KxVLzduw3PusyszjCtF/hgovYo+x4+ugfT7vRSenE/hvfbYElit/9lwwGT56J+c02RRzmQpIEnEzRQ==
> Enter a passphrase for your key:
ERROR: password must be at least 8 characters
```