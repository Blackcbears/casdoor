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

package controllers

import (
	"encoding/json"

	"github.com/beego/beego/utils/pagination"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

// GetTokens
// @Title GetTokens
// @Tag Token API
// @Description get tokens
// @Param   owner     query    string  true        "The owner of tokens"
// @Param   pageSize     query    string  true        "The size of each page"
// @Param   p     query    string  true        "The number of the page"
// @Success 200 {array} object.Token The Response object
// @router /get-tokens [get]
func (c *ApiController) GetTokens() {
	owner := c.Input().Get("owner")
	limit := c.Input().Get("pageSize")
	page := c.Input().Get("p")
	field := c.Input().Get("field")
	value := c.Input().Get("value")
	sortField := c.Input().Get("sortField")
	sortOrder := c.Input().Get("sortOrder")
	if limit == "" || page == "" {
		tokens := object.GetTokens(owner)
		c.ResponseOk(tokens)
	} else {
		limit := util.ParseInt(limit)
		paginator := pagination.SetPaginator(c.Ctx, limit, int64(object.GetTokenCount(owner, field, value)))
		tokens := object.GetPaginationTokens(owner, paginator.Offset(), limit, field, value, sortField, sortOrder)
		c.ResponseOk(tokens, paginator.Nums())
	}
}

// GetToken
// @Title GetToken
// @Tag Token API
// @Description get token
// @Param   id     query    string  true        "The id ( owner/name ) of token"
// @Success 200 {object} object.Token The Response object
// @router /get-token [get]
func (c *ApiController) GetToken() {
	id := c.Input().Get("id")

	token := object.GetToken(id)
	c.ResponseOk(token)
}

// UpdateToken
// @Title UpdateToken
// @Tag Token API
// @Description update token
// @Param   id     query    string  true        "The id ( owner/name ) of token"
// @Param   body    body   object.Token  true        "Details of the token"
// @Success 200 {object} controllers.Response The Response object
// @router /update-token [post]
func (c *ApiController) UpdateToken() {
	id := c.Input().Get("id")

	var token object.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	response := wrapActionResponse(object.UpdateToken(id, &token))
	c.ResponseOk(response)
}

// AddToken
// @Title AddToken
// @Tag Token API
// @Description add token
// @Param   body    body   object.Token  true        "Details of the token"
// @Success 200 {object} controllers.Response The Response object
// @router /add-token [post]
func (c *ApiController) AddToken() {
	var token object.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	response := wrapActionResponse(object.AddToken(&token))
	c.ResponseOk(response)
}

// DeleteToken
// @Tag Token API
// @Title DeleteToken
// @Description delete token
// @Param   body    body   object.Token  true        "Details of the token"
// @Success 200 {object} controllers.Response The Response object
// @router /delete-token [post]
func (c *ApiController) DeleteToken() {
	var token object.Token
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &token)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	response := wrapActionResponse(object.DeleteToken(&token))
	c.ResponseOk(response)
}

// GetOAuthCode
// @Title GetOAuthCode
// @Tag Token API
// @Description get OAuth code
// @Param   id     query    string  true        "The id ( owner/name ) of user"
// @Param   client_id     query    string  true        "OAuth client id"
// @Param   response_type     query    string  true        "OAuth response type"
// @Param   redirect_uri     query    string  true        "OAuth redirect URI"
// @Param   scope     query    string  true        "OAuth scope"
// @Param   state     query    string  true        "OAuth state"
// @Success 200 {object} object.TokenWrapper The Response object
// @router /login/oauth/code [post]
func (c *ApiController) GetOAuthCode() {
	userId := c.Input().Get("user_id")
	clientId := c.Input().Get("client_id")
	responseType := c.Input().Get("response_type")
	redirectUri := c.Input().Get("redirect_uri")
	scope := c.Input().Get("scope")
	state := c.Input().Get("state")
	nonce := c.Input().Get("nonce")

	challengeMethod := c.Input().Get("code_challenge_method")
	codeChallenge := c.Input().Get("code_challenge")

	if challengeMethod != "S256" && challengeMethod != "null" && challengeMethod != "" {
		c.ResponseError(c.T("auth:Challenge method should be S256"))
		return
	}
	host := c.Ctx.Request.Host

	oauthCode := object.GetOAuthCode(userId, clientId, responseType, redirectUri, scope, state, nonce, codeChallenge, host, c.GetAcceptLanguage())
	c.ResponseOk(oauthCode)
}

// GetOAuthToken
// @Title GetOAuthToken
// @Tag Token API
// @Description get OAuth access token
// @Param   grant_type     query    string  true        "OAuth grant type"
// @Param   client_id     query    string  true        "OAuth client id"
// @Param   client_secret     query    string  true        "OAuth client secret"
// @Param   code     query    string  true        "OAuth code"
// @Success 200 {object} object.TokenWrapper The Response object
// @Success 400 {object} object.TokenError The Response object
// @Success 401 {object} object.TokenError The Response object
// @router /login/oauth/access_token [post]
func (c *ApiController) GetOAuthToken() {
	grantType := c.Input().Get("grant_type")
	refreshToken := c.Input().Get("refresh_token")
	clientId := c.Input().Get("client_id")
	clientSecret := c.Input().Get("client_secret")
	code := c.Input().Get("code")
	verifier := c.Input().Get("code_verifier")
	scope := c.Input().Get("scope")
	username := c.Input().Get("username")
	password := c.Input().Get("password")
	tag := c.Input().Get("tag")
	avatar := c.Input().Get("avatar")

	if clientId == "" && clientSecret == "" {
		clientId, clientSecret, _ = c.Ctx.Request.BasicAuth()
	}
	if clientId == "" {
		// If clientID is empty, try to read data from RequestBody
		var tokenRequest TokenRequest
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tokenRequest); err == nil {
			clientId = tokenRequest.ClientId
			clientSecret = tokenRequest.ClientSecret
			grantType = tokenRequest.GrantType
			refreshToken = tokenRequest.RefreshToken
			code = tokenRequest.Code
			verifier = tokenRequest.Verifier
			scope = tokenRequest.Scope
			username = tokenRequest.Username
			password = tokenRequest.Password
			tag = tokenRequest.Tag
			avatar = tokenRequest.Avatar
		}
	}
	host := c.Ctx.Request.Host

	oauthToken := object.GetOAuthToken(grantType, clientId, clientSecret, code, verifier, scope, username, password, host, refreshToken, tag, avatar, c.GetAcceptLanguage())
	c.ResponseOk(oauthToken)
	c.SetTokenErrorHttpStatus()
}

// RefreshToken
// @Title RefreshToken
// @Tag Token API
// @Description refresh OAuth access token
// @Param   grant_type     query    string  true        "OAuth grant type"
// @Param	refresh_token	query	string	true		"OAuth refresh token"
// @Param   scope     query    string  true        "OAuth scope"
// @Param   client_id     query    string  true        "OAuth client id"
// @Param   client_secret     query    string  false        "OAuth client secret"
// @Success 200 {object} object.TokenWrapper The Response object
// @Success 400 {object} object.TokenError The Response object
// @Success 401 {object} object.TokenError The Response object
// @router /login/oauth/refresh_token [post]
func (c *ApiController) RefreshToken() {
	grantType := c.Input().Get("grant_type")
	refreshToken := c.Input().Get("refresh_token")
	scope := c.Input().Get("scope")
	clientId := c.Input().Get("client_id")
	clientSecret := c.Input().Get("client_secret")
	host := c.Ctx.Request.Host

	if clientId == "" {
		// If clientID is empty, try to read data from RequestBody
		var tokenRequest TokenRequest
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tokenRequest); err == nil {
			clientId = tokenRequest.ClientId
			clientSecret = tokenRequest.ClientSecret
			grantType = tokenRequest.GrantType
			scope = tokenRequest.Scope
			refreshToken = tokenRequest.RefreshToken
		}
	}

	token := object.RefreshToken(grantType, refreshToken, scope, clientId, clientSecret, host)
	c.ResponseOk(token)
	c.SetTokenErrorHttpStatus()
}

// IntrospectToken
// @Title IntrospectToken
// @Description The introspection endpoint is an OAuth 2.0 endpoint that takes a
// parameter representing an OAuth 2.0 token and returns a JSON document
// representing the meta information surrounding the
// token, including whether this token is currently active.
// This endpoint only support Basic Authorization.
//
// @Param token formData string true "access_token's value or refresh_token's value"
// @Param token_type_hint formData string true "the token type access_token or refresh_token"
// @Success 200 {object} object.IntrospectionResponse The Response object
// @Success 400 {object} object.TokenError The Response object
// @Success 401 {object} object.TokenError The Response object
// @router /login/oauth/introspect [post]
func (c *ApiController) IntrospectToken() {
	tokenValue := c.Input().Get("token")
	clientId, clientSecret, ok := c.Ctx.Request.BasicAuth()
	if !ok {
		clientId = c.Input().Get("client_id")
		clientSecret = c.Input().Get("client_secret")
		if clientId == "" || clientSecret == "" {
			c.ResponseError(c.T("token:Empty clientId or clientSecret"))
			tokenError := &object.TokenError{
				Error: object.InvalidRequest,
			}
			c.ResponseOk(tokenError)
			c.SetTokenErrorHttpStatus()
			return
		}
	}
	application := object.GetApplicationByClientId(clientId)
	if application == nil || application.ClientSecret != clientSecret {
		c.ResponseError(c.T("token:Invalid application or wrong clientSecret"))
		tokenError := &object.TokenError{
			Error: object.InvalidClient,
		}
		c.ResponseOk(tokenError)
		c.SetTokenErrorHttpStatus()
		return
	}
	token := object.GetTokenByTokenAndApplication(tokenValue, application.Name)
	if token == nil {
		response := &object.IntrospectionResponse{Active: false}
		c.ResponseOk(response)
		return
	}
	jwtToken, err := object.ParseJwtTokenByApplication(tokenValue, application)
	if err != nil || jwtToken.Valid() != nil {
		// and token revoked case. but we not implement
		// TODO: 2022-03-03 add token revoked check, when we implemented the Token Revocation(rfc7009) Specs.
		// refs: https://tools.ietf.org/html/rfc7009
		response := &object.IntrospectionResponse{Active: false}
		c.ResponseOk(response)
		return
	}

	response := &object.IntrospectionResponse{
		Active:    true,
		Scope:     jwtToken.Scope,
		ClientId:  clientId,
		Username:  token.User,
		TokenType: token.TokenType,
		Exp:       jwtToken.ExpiresAt.Unix(),
		Iat:       jwtToken.IssuedAt.Unix(),
		Nbf:       jwtToken.NotBefore.Unix(),
		Sub:       jwtToken.Subject,
		Aud:       jwtToken.Audience,
		Iss:       jwtToken.Issuer,
		Jti:       jwtToken.Id,
	}
	c.ResponseOk(response)
}
