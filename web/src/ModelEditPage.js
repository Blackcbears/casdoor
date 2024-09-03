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

import React from "react";
import {Button, Card, Col, Input, Row, Select} from "antd";
import * as ModelBackend from "./backend/ModelBackend";
import * as OrganizationBackend from "./backend/OrganizationBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

import {Controlled as CodeMirror} from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/addon/lint/lint.css";
import "codemirror/addon/lint/lint";
import {newModel} from "casbin";

require("codemirror/mode/properties/properties");

const {Option} = Select;

class ModelEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      organizationName: props.organizationName !== undefined ? props.organizationName : props.match.params.organizationName,
      modelName: props.match.params.modelName,
      model: null,
      organizations: [],
      users: [],
      mode: props.location.mode !== undefined ? props.location.mode : "edit",
    };
  }

  UNSAFE_componentWillMount() {
    this.getModel();
    this.getOrganizations();
  }

  getModel() {
    ModelBackend.getModel(this.state.organizationName, this.state.modelName)
      .then((res) => {
        if (res.data === null) {
          this.props.history.push("/404");
          return;
        }

        if (res.status === "error") {
          Setting.showMessage("error", res.msg);
          return;
        }

        this.setState({
          model: res.data,
        });
      });
  }

  getOrganizations() {
    OrganizationBackend.getOrganizations("admin")
      .then((res) => {
        this.setState({
          organizations: res.data || [],
        });
      });
  }

  parseModelField(key, value) {
    if ([""].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updateModelField(key, value) {
    value = this.parseModelField(key, value);

    const model = this.state.model;
    model[key] = value;
    this.setState({
      model: model,
    });
  }

  checkModelSyntax = (modelText) => {
    try {
      const model = newModel(modelText);
      if (!model.model.get("r") || !model.model.get("p") || !model.model.get("e")) {
        throw new Error("Model is missing one or more required sections (r, p, or e)");
      }
      return null;
    } catch (e) {
      return e.message;
    }
  };

  createLinter = (CodeMirror) => {
    CodeMirror.registerHelper("lint", "properties", (text) => {
      const error = this.checkModelSyntax(text);
      if (error) {
        const lineMatch = error.match(/line (\d+)/);
        if (lineMatch) {
          const lineNumber = parseInt(lineMatch[1], 10) - 1;
          return [{
            from: CodeMirror.Pos(lineNumber, 0),
            to: CodeMirror.Pos(lineNumber, text.split("\n")[lineNumber].length),
            message: error,
            severity: "error",
          }];
        } else {
          return [{
            from: CodeMirror.Pos(0, 0),
            to: CodeMirror.Pos(text.split("\n").length - 1),
            message: error,
            severity: "error",
          }];
        }
      }
      return [];
    });
  };

  renderModel() {
    return (
      <Card size="small" title={
        <div>
          {this.state.mode === "add" ? i18next.t("model:New Model") : i18next.t("model:Edit Model")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitModelEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: "20px"}} type="primary" onClick={() => this.submitModelEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
          {this.state.mode === "add" ? <Button style={{marginLeft: "20px"}} onClick={() => this.deleteModel()}>{i18next.t("general:Cancel")}</Button> : null}
        </div>
      } style={(Setting.isMobile()) ? {margin: "5px"} : {}} type="inner">
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Organization"), i18next.t("general:Organization - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Select virtual={false} style={{width: "100%"}} disabled={!Setting.isAdminUser(this.props.account) || Setting.builtInObject(this.state.model)} value={this.state.model.owner} onChange={(value => {this.updateModelField("owner", value);})}>
              {
                this.state.organizations.map((organization, index) => <Option key={index} value={organization.name}>{organization.name}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Name"), i18next.t("general:Name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input disabled={Setting.builtInObject(this.state.model)} value={this.state.model.name} onChange={e => {
              this.updateModelField("name", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Display name"), i18next.t("general:Display name - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.model.displayName} onChange={e => {
              this.updateModelField("displayName", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("general:Description"), i18next.t("general:Description - Tooltip"))} :
          </Col>
          <Col span={22} >
            <Input value={this.state.model.description} onChange={e => {
              this.updateModelField("description", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={(Setting.isMobile()) ? 22 : 2}>
            {Setting.getLabel(i18next.t("model:Model text"), i18next.t("model:Model text - Tooltip"))} :
          </Col>
          <Col span={22}>
            <div style={{width: "100%"}} >
              <CodeMirror
                value={this.state.model.modelText}
                options={{
                  mode: "properties",
                  theme: "default",
                  lineNumbers: true,
                  lint: true,
                }}
                onBeforeChange={(editor, data, value) => {
                  if (Setting.builtInObject(this.state.model)) {
                    return;
                  }
                  this.updateModelField("modelText", value);
                }}
                editorDidMount={(editor, value, cb) => {
                  this.createLinter(editor.constructor);
                }}
              />
            </div>
          </Col>
        </Row>
      </Card>
    );
  }

  submitModelEdit(exitAfterSave) {
    const model = Setting.deepCopy(this.state.model);
    ModelBackend.updateModel(this.state.organizationName, this.state.modelName, model)
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", i18next.t("general:Successfully saved"));
          this.setState({
            modelName: this.state.model.name,
          });

          if (exitAfterSave) {
            this.props.history.push("/models");
          } else {
            this.props.history.push(`/models/${this.state.model.owner}/${this.state.model.name}`);
          }
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to save")}: ${res.msg}`);
          this.updateModelField("name", this.state.modelName);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  deleteModel() {
    ModelBackend.deleteModel(this.state.model)
      .then((res) => {
        if (res.status === "ok") {
          this.props.history.push("/models");
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to delete")}: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  render() {
    return (
      <div>
        {
          this.state.model !== null ? this.renderModel() : null
        }
        <div style={{marginTop: "20px", marginLeft: "40px"}}>
          <Button size="large" onClick={() => this.submitModelEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: "20px"}} type="primary" size="large" onClick={() => this.submitModelEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
          {this.state.mode === "add" ? <Button style={{marginLeft: "20px"}} size="large" onClick={() => this.deleteModel()}>{i18next.t("general:Cancel")}</Button> : null}
        </div>
      </div>
    );
  }
}

export default ModelEditPage;
