package lunoStream

import (
	"github.com/bhbosman/goLuno/internal/listener"
	"go.uber.org/fx"
	"net/url"
)

func ProvideCompressedListener(
	settings *AppSettings,
) fx.Option {
	var opts []fx.Option
	localTextListenerUrl, err := url.Parse(settings.compressedListenerUrl)
	if err != nil {
		return fx.Error(err)
	}
	if err != nil {
		return fx.Error(err)
	}
	opts = append(opts,
		listener.CompressedListener(
			1024,
			false,
			nil,
			localTextListenerUrl,
		),
	)

	return fx.Options(opts...)
}
