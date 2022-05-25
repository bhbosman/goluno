package lunoStream

import (
	"github.com/bhbosman/gocomms/netDial"
	"go.uber.org/fx"
)

func ProvideTextListener(settings *AppSettings, ConsumerCounter *netDial.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.textListenerEnabled {
		opts = append(opts, TextListener(
			ConsumerCounter,
			1024,
			settings.textListenerUrl,
			settings.pairs...))
	}
	return fx.Options(opts...)
}

func ProvideCompressedListener(settings *AppSettings, ConsumerCounter *netDial.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.compressedListenerEnabled {
		opts = append(opts,
			CompressedListener(
				ConsumerCounter,
				1024,
				settings.compressedListenerUrl,
				settings.pairs...))
	}
	return fx.Options(opts...)
}
