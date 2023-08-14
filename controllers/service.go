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

// Casdoor will expose its providers as services to SDK
// We are going to implement those services as APIs here

package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

type EmailForm struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Sender    string   `json:"sender"`
	Receivers []string `json:"receivers"`
	Provider  string   `json:"provider"`
}

type SmsForm struct {
	Content   string   `json:"content"`
	Receivers []string `json:"receivers"`
	OrgId     string   `json:"organizationId"` // e.g. "admin/built-in"
}

// SendEmail
// @Title SendEmail
// @Tag Service API
// @Description This API is not for Casdoor frontend to call, it is for Casdoor SDKs.
// @Param   clientId    query    string  true        "The clientId of the application"
// @Param   clientSecret    query    string  true    "The clientSecret of the application"
// @Param   from    body   controllers.EmailForm    true         "Details of the email request"
// @Success 200 {object}  Response object
// @router /api/send-email [post]
func (c *ApiController) SendEmail() {
	var emailForm EmailForm

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &emailForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	var provider *object.Provider
	if emailForm.Provider != "" {
		// called by frontend's TestEmailWidget, provider name is set by frontend
		provider, err = object.GetProvider(util.GetId("admin", emailForm.Provider))
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

	} else {
		// called by Casdoor SDK via Client ID & Client Secret, so the used Email provider will be the application' Email provider or the default Email provider
		provider, err = c.GetProviderFromContext("Email")
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
	}

	// when receiver is the reserved keyword: "TestSmtpServer", it means to test the SMTP server instead of sending a real Email
	if len(emailForm.Receivers) == 1 && emailForm.Receivers[0] == "TestSmtpServer" {
		err := object.DailSmtpServer(provider)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk()
	}

	if util.IsStringsEmpty(emailForm.Title, emailForm.Content, emailForm.Sender) {
		c.ResponseError(fmt.Sprintf(c.T("service:Empty parameters for emailForm: %v"), emailForm))
		return
	}

	invalidReceivers := []string{}
	for _, receiver := range emailForm.Receivers {
		if !util.IsEmailValid(receiver) {
			invalidReceivers = append(invalidReceivers, receiver)
		}
	}

	if len(invalidReceivers) != 0 {
		c.ResponseError(fmt.Sprintf(c.T("service:Invalid Email receivers: %s"), invalidReceivers))
		return
	}

	code := "123456"
	// "You have requested a verification code at Casdoor. Here is your code: %s, please enter in 5 minutes."
	content := fmt.Sprintf(emailForm.Content, code)
	for _, receiver := range emailForm.Receivers {
		err = object.SendEmail(provider, emailForm.Title, content, receiver, emailForm.Sender)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
	}

	c.ResponseOk()
}

// SendSms
// @Title SendSms
// @Tag Service API
// @Description This API is not for Casdoor frontend to call, it is for Casdoor SDKs.
// @Param   clientId    query    string  true        "The clientId of the application"
// @Param   clientSecret    query    string  true    "The clientSecret of the application"
// @Param   from    body   controllers.SmsForm    true           "Details of the sms request"
// @Success 200 {object}  Response object
// @router /api/send-sms [post]
func (c *ApiController) SendSms() {
	provider, err := c.GetProviderFromContext("SMS")
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	var smsForm SmsForm
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &smsForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if provider.Type != "Custom HTTP SMS" {
		invalidReceivers := getInvalidSmsReceivers(smsForm)
		if len(invalidReceivers) != 0 {
			c.ResponseError(fmt.Sprintf(c.T("service:Invalid phone receivers: %s"), strings.Join(invalidReceivers, ", ")))
			return
		}
	}

	err = object.SendSms(provider, smsForm.Content, smsForm.Receivers...)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk()
}
