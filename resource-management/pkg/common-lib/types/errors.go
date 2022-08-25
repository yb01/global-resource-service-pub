/*
Copyright 2022 Authors of Global Resource Service.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import "errors"

const (
	ErrMsg_HostRequestExceedLimit     = "Requested host number exceeds limit"
	ErrMsg_HostRequestExceedCapacity  = "Requested hosts exceeds capacity"
	ErrMsg_HostRequestLessThanMiniaml = "Requested host number less than minimal request"

	ErrMsg_ClientIdExisted = "Client id exists"

	ErrMsg_FailedToProcessBookmarkEvent = "Failed to process bookmark events"

	ErrMsg_EndOfEventQueue = "Reach the end of event queue"

	ErrMsg_ObjectNotFound = "Object not found"
)

var Error_HostRequestExceedLimit = errors.New(ErrMsg_HostRequestExceedLimit)
var Error_HostRequestExceedCapacity = errors.New(ErrMsg_HostRequestExceedCapacity)
var Error_HostRequestLessThanMiniaml = errors.New(ErrMsg_HostRequestLessThanMiniaml)

var Error_ClientIdExisted = errors.New(ErrMsg_ClientIdExisted)

var Error_FailedToProcessBookmarkEvent = errors.New(ErrMsg_FailedToProcessBookmarkEvent)

var Error_EndOfEventQueue = errors.New(ErrMsg_EndOfEventQueue)

var Error_ObjectNotFound = errors.New(ErrMsg_ObjectNotFound)
