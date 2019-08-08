# Description
```
获取列表(JSON格式)
```
# Input
```
$ qoscli keys list -o json
```
# Output
注意，这里对返回结果做了手动格式化， 原始结果是未格式化的json.
```
$ qoscli keys list -o json
[
    {
        "name": "test", 
        "type": "local", 
        "address": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56", 
        "pub_key": "heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="
    }, 
    {
        "name": "test01", 
        "type": "local", 
        "address": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m", 
        "pub_key": "70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="
    }, 
    {
        "name": "test02", 
        "type": "import", 
        "address": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m", 
        "pub_key": "70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="
    }, 
    {
        "name": "test03", 
        "type": "import", 
        "address": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m", 
        "pub_key": "70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="
    }
]
```

- 原始结果：
```
$ qoscli keys list -o json
[{"name":"test","type":"local","address":"address1hw43pwhtscealvu973r66vk83gus8myp40fy56","pub_key":"heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="},{"name":"test01","type":"local","address":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m","pub_key":"70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="},{"name":"test02","type":"import","address":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m","pub_key":"70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="},{"name":"test03","type":"import","address":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m","pub_key":"70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="}]
```