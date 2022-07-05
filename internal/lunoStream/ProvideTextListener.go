package lunoStream

import (
	"github.com/bhbosman/goCommsNetDialer"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/fx"
	"net/url"
)

func ProvideTextListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	settings *AppSettings,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl,
) fx.Option {
	var opts []fx.Option
	if settings.textListenerEnabled {
		localTextListenerUrl, err := url.Parse(settings.textListenerUrl)
		if err != nil {
			return fx.Error(err)
		}
		if err != nil {
			return fx.Error(err)
		}
		opts = append(
			opts,
			TextListener(
				serviceIdentifier,
				serviceDependentOn,
				ConsumerCounter,
				1024,
				false,
				nil,
				localTextListenerUrl,
				settings.pairs...))
	}
	return fx.Options(opts...)
}

func ProvideCompressedListener(
	serviceIdentifier model.ServiceIdentifier,
	serviceDependentOn model.ServiceIdentifier,
	settings *AppSettings,
	ConsumerCounter *goCommsNetDialer.CanDialDefaultImpl) fx.Option {
	var opts []fx.Option
	if settings.compressedListenerEnabled {
		localTextListenerUrl, err := url.Parse(settings.compressedListenerUrl)
		if err != nil {
			return fx.Error(err)
		}
		if err != nil {
			return fx.Error(err)
		}
		opts = append(opts,
			CompressedListener(
				serviceIdentifier,
				serviceDependentOn,
				ConsumerCounter,
				1024,
				false,
				nil,
				localTextListenerUrl,
				settings.pairs...))
	}
	return fx.Options(opts...)
}
