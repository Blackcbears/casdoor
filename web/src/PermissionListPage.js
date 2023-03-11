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
import {Link} from "react-router-dom";
import {Button, Switch, Table} from "antd";
import moment from "moment";
import * as Setting from "./Setting";
import * as PermissionBackend from "./backend/PermissionBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";
import PopconfirmModal from "./PopconfirmModal";

class PermissionListPage extends BaseListPage {
  newPermission() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.owner,
      name: `permission_${randomName}`,
      createdTime: moment().format(),
      displayName: `New Permission - ${randomName}`,
      users: [`${this.props.account.owner}/${this.props.account.name}`],
      roles: [],
      domains: [],
      resourceType: "Application",
      resources: ["app-built-in"],
      actions: ["Read"],
      effect: "Allow",
      isEnabled: true,
      submitter: this.props.account.name,
      approver: "",
      approveTime: "",
      state: "Pending",
    };
  }

  addPermission() {
    const newPermission = this.newPermission();
    PermissionBackend.addPermission(newPermission)
      .then((res) => {
        if (res.status === "ok") {
          this.props.history.push({pathname: `/permissions/${newPermission.owner}/${newPermission.name}`, mode: "add"});
          Setting.showMessage("success", i18next.t("general:Successfully added"));
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to add")}: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  deletePermission(i) {
    PermissionBackend.deletePermission(this.state.data[i])
      .then((res) => {
        if (res.status === "ok") {
          Setting.showMessage("success", i18next.t("general:Successfully deleted"));
          this.setState({
            data: Setting.deleteRow(this.state.data, i),
            pagination: {total: this.state.pagination.total - 1},
          });
        } else {
          Setting.showMessage("error", `${i18next.t("general:Failed to delete")}: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `${i18next.t("general:Failed to connect to server")}: ${error}`);
      });
  }

  renderTable(permissions) {
    const columns = [
      // https://github.com/ant-design/ant-design/issues/22184
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "150px",
        fixed: "left",
        sorter: true,
        ...this.getColumnSearchProps("name"),
        render: (text, record, index) => {
          return (
            <Link to={`/permissions/${record.owner}/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Organization"),
        dataIndex: "owner",
        key: "owner",
        width: "120px",
        sorter: true,
        ...this.getColumnSearchProps("owner"),
        render: (text, record, index) => {
          return (
            <Link to={`/organizations/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Created time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "160px",
        sorter: true,
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Display name"),
        dataIndex: "displayName",
        key: "displayName",
        width: "160px",
        sorter: true,
        ...this.getColumnSearchProps("displayName"),
      },
      {
        title: i18next.t("role:Sub users"),
        dataIndex: "users",
        key: "users",
        // width: '100px',
        sorter: true,
        ...this.getColumnSearchProps("users"),
        render: (text, record, index) => {
          return Setting.getTags(text);
        },
      },
      {
        title: i18next.t("role:Sub roles"),
        dataIndex: "roles",
        key: "roles",
        // width: '100px',
        sorter: true,
        ...this.getColumnSearchProps("roles"),
        render: (text, record, index) => {
          return Setting.getTags(text);
        },
      },
      {
        title: i18next.t("role:Sub domains"),
        dataIndex: "domains",
        key: "domains",
        sorter: true,
        ...this.getColumnSearchProps("domains"),
        render: (text, record, index) => {
          return Setting.getTags(text);
        },
      },
      {
        title: i18next.t("permission:Resource type"),
        dataIndex: "resourceType",
        key: "resourceType",
        filterMultiple: false,
        filters: [
          {text: "Application", value: "Application"},
        ],
        width: "170px",
        sorter: true,
      },
      {
        title: i18next.t("permission:Resources"),
        dataIndex: "resources",
        key: "resources",
        // width: '100px',
        sorter: true,
        ...this.getColumnSearchProps("resources"),
        render: (text, record, index) => {
          return Setting.getTags(text);
        },
      },
      {
        title: i18next.t("permission:Actions"),
        dataIndex: "actions",
        key: "actions",
        // width: '100px',
        sorter: true,
        ...this.getColumnSearchProps("actions"),
        render: (text, record, index) => {
          const tags = text.map((tag, i) => {
            switch (tag) {
            case "Read":
              return i18next.t("permission:Read");
            case "Write":
              return i18next.t("permission:Write");
            case "Admin":
              return i18next.t("permission:Admin");
            default:
              return null;
            }
          });
          return Setting.getTags(tags);
        },
      },
      {
        title: i18next.t("permission:Effect"),
        dataIndex: "effect",
        key: "effect",
        filterMultiple: false,
        filters: [
          {text: i18next.t("permission:Allow"), value: "Allow"},
          {text: i18next.t("permission:Deny"), value: "Deny"},
        ],
        width: "120px",
        sorter: true,
        render: (text, record, index) => {
          switch (text) {
          case "Allow":
            return Setting.getTag("success", i18next.t("permission:Allow"));
          case "Deny":
            return Setting.getTag("error", i18next.t("permission:Deny"));
          default:
            return null;
          }
        },
      },
      {
        title: i18next.t("general:Is enabled"),
        dataIndex: "isEnabled",
        key: "isEnabled",
        width: "120px",
        sorter: true,
        render: (text, record, index) => {
          return (
            <Switch disabled checkedChildren="ON" unCheckedChildren="OFF" checked={text} />
          );
        },
      },
      {
        title: i18next.t("permission:Submitter"),
        dataIndex: "submitter",
        key: "submitter",
        filterMultiple: false,
        width: "120px",
        sorter: true,
      },
      {
        title: i18next.t("permission:Approver"),
        dataIndex: "approver",
        key: "approver",
        filterMultiple: false,
        width: "120px",
        sorter: true,
      },
      {
        title: i18next.t("permission:Approve time"),
        dataIndex: "approveTime",
        key: "approveTime",
        filterMultiple: false,
        width: "120px",
        sorter: true,
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("permission:State"),
        dataIndex: "state",
        key: "state",
        filterMultiple: false,
        filters: [
          {text: i18next.t("permission:Approved"), value: "Approved"},
          {text: i18next.t("permission:Pending"), value: "Pending"},
        ],
        width: "120px",
        sorter: true,
        render: (text, record, index) => {
          switch (text) {
          case "Approved":
            return Setting.getTag("success", i18next.t("permission:Approved"));
          case "Pending":
            return Setting.getTag("error", i18next.t("permission:Pending"));
          default:
            return null;
          }
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        width: "170px",
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/permissions/${record.owner}/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <PopconfirmModal
                title={i18next.t("general:Sure to delete") + `: ${record.name} ?`}
                onConfirm={() => this.deletePermission(index)}
              >
              </PopconfirmModal>
            </div>
          );
        },
      },
    ];

    const paginationProps = {
      total: this.state.pagination.total,
      showQuickJumper: true,
      showSizeChanger: true,
      showTotal: () => i18next.t("general:{total} in total").replace("{total}", this.state.pagination.total),
    };

    return (
      <div>
        <Table scroll={{x: "max-content"}} columns={columns} dataSource={permissions} rowKey="name" size="middle" bordered pagination={paginationProps}
          title={() => (
            <div>
              {i18next.t("general:Permissions")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" onClick={this.addPermission.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }

  fetch = (params = {}) => {
    let field = params.searchedColumn, value = params.searchText;
    const sortField = params.sortField, sortOrder = params.sortOrder;
    if (params.type !== undefined && params.type !== null) {
      field = "type";
      value = params.type;
    }
    this.setState({loading: true});

    const getPermissions = Setting.isLocalAdminUser(this.props.account) ? PermissionBackend.getPermissions : PermissionBackend.getPermissionsBySubmitter;
    getPermissions(Setting.isAdminUser(this.props.account) ? "" : this.props.account.owner, params.pagination.current, params.pagination.pageSize, field, value, sortField, sortOrder)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            loading: false,
            data: res.data,
            pagination: {
              ...params.pagination,
              total: res.data2,
            },
            searchText: params.searchText,
            searchedColumn: params.searchedColumn,
          });
        } else {
          if (Setting.isResponseDenied(res)) {
            this.setState({
              loading: false,
              isAuthorized: false,
            });
          }
        }
      });
  };
}

export default PermissionListPage;
