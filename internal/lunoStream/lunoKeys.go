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

func ProvideReadLunoKeys() fx.Option {
	data := &lunoKeys{}
	current, err := user.Current()
	if err != nil {
		return fx.Error(err)
	}
	f, err := os.Open(fmt.Sprintf("%v/.luno/keys.json", current.HomeDir))
	if err != nil {
		return fx.Error(err)
	}

	defer func() {
		_ = f.Close()
	}()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(data)
	if err != nil {
		return fx.Error(err)
	}
	return fx.Options(
		fx.Provide(fx.Annotated{Name: "LunoAPIKeyID", Target: impl.CreateStringContext(data.Key)}),
		fx.Provide(fx.Annotated{Name: "LunoAPIKeySecret", Target: impl.CreateStringContext(data.Secret)}),
	)
}
