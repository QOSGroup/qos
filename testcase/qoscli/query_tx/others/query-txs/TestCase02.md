# Description
```
指定错误的可选参数`--page`和`--limit`
```
# Input
- 指定错误的可选参数`--page`
```
$ qoscli query txs --tags "sender:address1hw43pwhtscealvu973r66vk83gus8myp40fy56" --page 0
```
- 指定错误的可选参数`--limit`
```
$ qoscli query txs --tags "sender:address1hw43pwhtscealvu973r66vk83gus8myp40fy56" --limit 0
```
# Output
- 指定错误的可选参数`--page`
```
$ qoscli query txs --tags "sender:address1hw43pwhtscealvu973r66vk83gus8myp40fy56" --page 0
ERROR: page must greater than 0
```
- 指定错误的可选参数`--limit`
```
$ qoscli query txs --tags "sender:address1hw43pwhtscealvu973r66vk83gus8myp40fy56" --limit 0
ERROR: limit must greater than 0
```