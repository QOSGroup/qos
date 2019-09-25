# Description
```
正常查询QSC联盟币
```
# Input
- 查询的QSC不存在
```
$ qoscli query qsc star --indent
```
- 查询的QSC存在
```
$ qoscli query qsc STAR --indent
```
# Output
- 查询的QSC不存在
```
$ qoscli query qsc star --indent
ERROR: star not exists.
```
- 查询的QSC存在
```
$ qoscli query qsc STAR --indent
{
  "name": "STAR",
  "chain_id": "capricorn-3000",
  "extrate": "1",
  "description": "",
  "banker": "address1ap6myv248ell0gkwnhzapjnyvt38584dk3e0te"
}
```