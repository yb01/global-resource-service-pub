package types

import "errors"

const (
	ErrMsg_HostRequestExceedLimit     = "Requested host number exceeds limit"
	ErrMsg_HostRequestExceedCapacity  = "Requested hosts exceeds capacity"
	ErrMsg_HostRequestLessThanMiniaml = "Requested host number less than minimal request"

	ErrMsg_ClientIdExisted = "Client id exists"

	ErrMsg_FailedToProcessBookmarkEvent = "Failed to process bookmark events"
)

var Error_HostRequestExceedLimit = errors.New(ErrMsg_HostRequestExceedLimit)
var Error_HostRequestExceedCapacity = errors.New(ErrMsg_HostRequestExceedCapacity)
var Error_HostRequestLessThanMiniaml = errors.New(ErrMsg_HostRequestLessThanMiniaml)

var Error_ClientIdExisted = errors.New(ErrMsg_ClientIdExisted)

var Error_FailedToProcessBookmarkEvent = errors.New(ErrMsg_FailedToProcessBookmarkEvent)
