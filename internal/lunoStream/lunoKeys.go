package lunoStream

import (
	"encoding/json"
	"fmt"
	"github.com/bhbosman/gocomms/impl"
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

func ProvideReadLunoKeys(data *lunoKeys) fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotated{Name: "LunoAPIKeyID", Target: impl.CreateStringContext(data.Key)}),
		fx.Provide(fx.Annotated{Name: "LunoAPIKeySecret", Target: impl.CreateStringContext(data.Secret)}),
	)
}
