//
// Description:
// Author: Rodrigo Freitas
// Crested at: Mon Apr 29 16:23:52 -03 2019
//
package rtsp

const (
	// Success
	StatusLowStorageSpace = 250

	// Error
	StatusParameterNotUnderstood        = 451
	StatusConferenceNotFound            = 452
	StatusNotEnoughBandwidth            = 453
	StatusSessionNotFound               = 454
	StatusMethodNotValidInThisState     = 455
	StatusHeaderFieldNotValid           = 456
	StatusInvalidRange                  = 457
	StatusParameterReadOnly             = 458
	StatusAggregateOperationNotAllowed  = 459
	StatusOnlyAggregateOperationAllowed = 460
	StatusUnsupportedTransport          = 461
	StatusDestinationUnreachable        = 462
	StatusOptionNotSupported            = 551
)

var statusText = map[int]string{
	StatusLowStorageSpace:               "Low on Storage Space",
	StatusParameterNotUnderstood:        "Parameter Not Understood",
	StatusConferenceNotFound:            "Conference Not Found",
	StatusNotEnoughBandwidth:            "Not Enough Bandwidth",
	StatusSessionNotFound:               "Session Not Found",
	StatusMethodNotValidInThisState:     "Method Not Valid in This State",
	StatusHeaderFieldNotValid:           "Header Field Not Valid for Resource",
	StatusInvalidRange:                  "Invalid Range",
	StatusParameterReadOnly:             "Parameter Is Read-Only",
	StatusAggregateOperationNotAllowed:  "Aggregate Operation Not Allowed",
	StatusOnlyAggregateOperationAllowed: "Only Aggregate Operation Allowed",
	StatusUnsupportedTransport:          "Unsupported Transport",
	StatusDestinationUnreachable:        "Destination Unreachable",
	StatusOptionNotSupported:            "Option not supported",
}

func StatusText(code int) string {
	return statusText[code]
}
