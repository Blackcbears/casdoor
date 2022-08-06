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

import {createButton} from "react-social-login-buttons";
import {StaticBaseUrl} from "../Setting";

function Icon() {
  return <img src={`${StaticBaseUrl}/buttons/gitee.svg`} alt="Sign in with Gitee" />;
}

const config = {
  text: "Sign in with Gitee",
  icon: Icon,
  iconFormat: name => `fa fa-${name}`,
  style: {background: "rgb(199,29,35)"},
  activeStyle: {background: "rgb(147,22,26)"},
};

const GiteeLoginButton = createButton(config);

export default GiteeLoginButton;
