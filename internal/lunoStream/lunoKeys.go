package lunoStream

import (
	"encoding/json"
	"fmt"
	"go.uber.org/fx"
	"os"
	"os/user"
)

type lunoKeys = struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

/*
Example of the keys.json file
{
  "key": "(some key)",
  "secret": "(some secret)"
}
*/

func ReadLunoKeys() (*lunoKeys, error) {
	data := &lunoKeys{}
	current, err := user.Current()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fmt.Sprintf("%v/.luno/keys.json", current.HomeDir))
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ProvideReadLunoKeys() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func() (*lunoKeys, error) {
					data, err := ReadLunoKeys()
					if err != nil {
						return nil, err
					}
					return data, nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Name: "LunoAPIKeyID",
				Target: func(data *lunoKeys) string {
					return data.Key
				},
			}),
		fx.Provide(
			fx.Annotated{
				Name: "LunoAPIKeySecret",
				Target: func(data *lunoKeys) string {
					return data.Secret
				},
			}),
	)
}
