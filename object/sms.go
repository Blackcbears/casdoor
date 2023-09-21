// Copyright 2021 The Casdoor Authors. All Rights Reserved.
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

package object

import (
	"strings"

	sender "github.com/casdoor/go-sms-sender"
)

func getSmsClient(provider *Provider) (sender.SmsClient, error) {
	var client sender.SmsClient
	var err error

	if provider.Type == sender.HuaweiCloud || provider.Type == sender.AzureACS {
		client, err = sender.NewSmsClient(provider.Type, provider.ClientId, provider.ClientSecret, provider.SignName, provider.TemplateCode, provider.ProviderUrl, provider.AppId)
	} else {
		client, err = sender.NewSmsClient(provider.Type, provider.ClientId, provider.ClientSecret, provider.SignName, provider.TemplateCode, provider.AppId)
	}
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SendSms(provider *Provider, params map[string]string, phoneNumbers ...string) error {
	client, err := getSmsClient(provider)
	if err != nil {
		return err
	}

	if provider.Type == sender.Twilio {
		if provider.AppId != "" {
			phoneNumbers = append([]string{provider.AppId}, phoneNumbers...)
		}
	} else if provider.Type == sender.Aliyun {
		for i, number := range phoneNumbers {
			phoneNumbers[i] = strings.TrimPrefix(number, "+86")
		}
	}

	err = client.SendMessage(params, phoneNumbers...)
	return err
}

func CodeToSmsParams(provider *Provider, content string) map[string]string {
	params := map[string]string{}
	if provider.Type == sender.TencentCloud {
		params["0"] = content
	} else {
		params["code"] = content
	}
	return params
}
