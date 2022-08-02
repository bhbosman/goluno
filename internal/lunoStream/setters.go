package lunoStream

//type addCurrencyPair struct {
//	pair *common.PairInformation
//}

//func AddCurrencyPair(pair *common.PairInformation) *addCurrencyPair {
//	return &addCurrencyPair{pair: pair}
//}

//func (self *addCurrencyPair) apply(settings *AppSettings) error {
//	settings.pairs = append(settings.pairs, self.pair)
//	return nil
//}

//type textListenerUrl struct {
//	url                string
//	ServiceIdentifier  model.ServiceIdentifier
//	ServiceDependentOn model.ServiceIdentifier
//}

//func TextListenerUrl(url string,
//	ServiceIdentifier model.ServiceIdentifier,
//	serviceDependentOn model.ServiceIdentifier) *textListenerUrl {
//	return &textListenerUrl{url: url,
//		ServiceIdentifier:  ServiceIdentifier,
//		ServiceDependentOn: serviceDependentOn,
//	}
//}

//func (self *textListenerUrl) apply(settings *AppSettings) error {
//	settings.textListenerUrl = self.url
//	settings.textListenerEnabled = true
//	return nil
//}

type compressedListenerUrl struct {
	url string
}

func CompressedListenerUrl(
	url string,
) *compressedListenerUrl {
	return &compressedListenerUrl{
		url: url,
	}
}

func (self *compressedListenerUrl) apply(settings *AppSettings) error {
	settings.compressedListenerUrl = self.url
	settings.compressedListenerEnabled = true
	return nil
}

type Errors struct {
	errors []error
}

func NewErrors(errors ...error) *Errors {
	return &Errors{
		errors: errors,
	}
}

func (self *Errors) apply(settings *AppSettings) error {
	settings.errors = append(settings.errors, self.errors...)
	return nil
}
