# Concepts

## ParamSet

Other module parameter structures implement the following specifications:

```go
type ParamSet interface {
	// parameter key-value pairs
	KeyValuePairs() KeyValuePairs

	// validate single parameter
	ValidateKeyValue(key string, value string) (interface{}, btypes.Error)

	// module name
	GetParamSpace() string

	// validate
	Validate() btypes.Error
}
```

function `KeyValuePairs()` in `ParamSet` returns all the parameter key-value pairs in the module:
```go
type KeyValuePair struct {
	Key   []byte
	Value interface{}
}

type KeyValuePairs []KeyValuePair
```

QOS `stake`, `gov`, `distribution` modules have parameters.

## Parameter management

### Register mappers

Initialize parameter configuration when registering mappers and hooks:
```go
// app/app.go
app.mm.RegisterMapperAndHooks(app, params.ModuleName, &stake.Params{}, &distribution.Params{}, &gov.Params{})

// module/params/mapper/mapper.go
func (mapper Mapper) RegisterParamSet(ps ...qtypes.ParamSet) {
	for _, ps := range ps {
		if ps != nil {
			if _, ok := mapper.paramSets[ps.GetParamSpace()]; ok {
				panic(fmt.Sprintf("<%s> already registered ", ps.GetParamSpace()))
			}
			mapper.paramSets[ps.GetParamSpace()] = ps
		}
	}
}
```

### Initialization

Save parameter information when the module with parameters is initialized:
```go
func InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data types.GenesisState) []abci.ValidatorUpdate{
    ...
    mapper.SetParams(ctx, params)
    ...
}
```

### Update

Submit parameter change proposal to update parameters. 