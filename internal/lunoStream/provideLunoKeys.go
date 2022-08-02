package lunoStream

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bhbosman/goLuno/internal/flags"
	"go.uber.org/fx"
	"io"
	"os"
	"os/user"
)

func ProvideLunoKeys() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (*lunoKeys, error) {
				if *flags.LunoSecret == "" && *flags.LunoKey == "" {
					var f io.ReadCloser
					if *flags.LunoKeyFile != "" {
						if _, err := os.Stat(*flags.LunoKeyFile); errors.Is(err, os.ErrNotExist) {
							return nil, err
						}
						var err error
						f, err = os.Open(*flags.LunoKeyFile)
						if err != nil {
							return nil, err
						}
						defer func() {
							_ = f.Close()
						}()
					} else {
						var err error
						var currentUser *user.User
						currentUser, err = user.Current()
						if err != nil {
							return nil, err
						}
						f, err = os.Open(fmt.Sprintf("%v/.luno/keys.json", currentUser.HomeDir))
						if err != nil {
							return nil, err
						}
						defer func() {
							_ = f.Close()
						}()
					}
					data := &lunoKeys{}
					decoder := json.NewDecoder(f)
					err := decoder.Decode(data)
					if err != nil {
						return nil, err
					}
					return data, nil
				} else if *flags.LunoSecret != "" && *flags.LunoKey != "" {
					data := &lunoKeys{
						Key:    *flags.LunoKey,
						Secret: *flags.LunoSecret,
					}
					return data, nil
				}
				return nil, fmt.Errorf("could not create Luno keys")
			},
		},
	)
}
