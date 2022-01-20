package lunoStream

import (
	"encoding/json"
	"fmt"
	"go.uber.org/fx"
	"os"
	"os/user"
)

func ProvideLunoKeys(fromFile bool, useThis *lunoKeys) fx.Option {
	if fromFile {
		return fx.Provide(
			fx.Annotated{
				Target: func() (*lunoKeys, error) {
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
				},
			})
	}
	if useThis == nil {
		return fx.Error(fmt.Errorf("provide LunoKeys == nil"))
	}
	return fx.Provide(
		fx.Annotated{
			Target: func() *lunoKeys {
				return useThis
			},
		})
}
