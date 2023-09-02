// Copyright 2023 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package notification

import "github.com/nikoksr/notify"

func GetNotificationProvider(typ string, appId string, receiver string, method string, title string) (notify.Notifier, error) {
	if typ == "Telegram" {
		return NewTelegramProvider(appId, receiver)
	} else if typ == "Custom HTTP" {
		return NewCustomHttpProvider(receiver, method, title)
	} else if typ == "DingTalk" {
		return NewDingTalkProvider(appId, receiver)
	} else if typ == "Lark" {
		return NewLarkProvider(receiver)
	} else if typ == "Microsoft Teams" {
		return NewMicrosoftTeamsProvider(receiver)
	} else if typ == "Bark" {
		return NewBarkProvider(receiver)
	} else if typ == "Pushover" {
		return NewPushoverProvider(appId, receiver)
	}

	return nil, nil
}
