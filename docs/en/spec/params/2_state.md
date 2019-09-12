# State

## Param

```go
type Mapper struct {
	*mapper.BaseMapper
	
	paramSets map[string]qtypes.ParamSet
}
```
`MapperName` is `params`. Other module parameter sets registered to params module when the QOS App is created, and saved to the database during initialization.


param: `moduleName keyName -> keyValue`
