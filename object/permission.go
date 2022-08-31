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
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/util"
	"xorm.io/core"
)

type Permission struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`

	Users   []string `xorm:"mediumtext" json:"users"`
	Roles   []string `xorm:"mediumtext" json:"roles"`
	Domains []string `xorm:"mediumtext" json:"domains"`

	Model        string   `xorm:"varchar(100)" json:"model"`
	Adapter      string   `xorm:"varchar(100)" json:"adapter"`
	ResourceType string   `xorm:"varchar(100)" json:"resourceType"`
	Resources    []string `xorm:"mediumtext" json:"resources"`
	Actions      []string `xorm:"mediumtext" json:"actions"`
	Effect       string   `xorm:"varchar(100)" json:"effect"`
	IsEnabled    bool     `json:"isEnabled"`

	Submitter   string `xorm:"varchar(100)" json:"submitter"`
	Approver    string `xorm:"varchar(100)" json:"approver"`
	ApproveTime string `xorm:"varchar(100)" json:"approveTime"`
	State       string `xorm:"varchar(100)" json:"state"`
}

type PermissionRule struct {
	Ptype string `xorm:"varchar(100) index not null default ''" json:"ptype"`
	V0    string `xorm:"varchar(100) index not null default ''" json:"v0"`
	V1    string `xorm:"varchar(100) index not null default ''" json:"v1"`
	V2    string `xorm:"varchar(100) index not null default ''" json:"v2"`
	V3    string `xorm:"varchar(100) index not null default ''" json:"v3"`
	V4    string `xorm:"varchar(100) index not null default ''" json:"v4"`
	V5    string `xorm:"varchar(100) index not null default ''" json:"v5"`
	Id    string `xorm:"varchar(100) index not null default ''" json:"id"`
}

func GetPermissionCount(owner, field, value string) int {
	session := GetSession(owner, -1, -1, field, value, "", "")
	count, err := session.Count(&Permission{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetPermissions(owner string) []*Permission {
	permissions := []*Permission{}
	err := adapter.Engine.Desc("created_time").Find(&permissions, &Permission{Owner: owner})
	if err != nil {
		panic(err)
	}

	return permissions
}

func GetPaginationPermissions(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Permission {
	permissions := []*Permission{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&permissions)
	if err != nil {
		panic(err)
	}

	return permissions
}

func getPermission(owner string, name string) *Permission {
	if owner == "" || name == "" {
		return nil
	}

	permission := Permission{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&permission)
	if err != nil {
		panic(err)
	}

	if existed {
		return &permission
	} else {
		return nil
	}
}

func GetPermission(id string) *Permission {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getPermission(owner, name)
}

func UpdatePermission(id string, permission *Permission) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	oldPermission := getPermission(owner, name)
	if oldPermission == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(permission)
	if err != nil {
		panic(err)
	}

	if affected != 0 {
		removePolicies(oldPermission)
		if oldPermission.Adapter != "" && oldPermission.Adapter != permission.Adapter {
			isEmpty, _ := adapter.Engine.IsTableEmpty(oldPermission.Adapter)
			if isEmpty {
				err = adapter.Engine.DropTables(oldPermission.Adapter)
				if err != nil {
					panic(err)
				}
			}
		}
		addPolicies(permission)
	}

	return affected != 0
}

func AddPermission(permission *Permission) bool {
	affected, err := adapter.Engine.Insert(permission)
	if err != nil {
		panic(err)
	}

	if affected != 0 {
		addPolicies(permission)
	}

	return affected != 0
}

func DeletePermission(permission *Permission) bool {
	affected, err := adapter.Engine.ID(core.PK{permission.Owner, permission.Name}).Delete(&Permission{})
	if err != nil {
		panic(err)
	}

	if affected != 0 {
		removePolicies(permission)
		if permission.Adapter != "" && permission.Adapter != "permission_rule" {
			isEmpty, _ := adapter.Engine.IsTableEmpty(permission.Adapter)
			if isEmpty {
				err = adapter.Engine.DropTables(permission.Adapter)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return affected != 0
}

func (permission *Permission) GetId() string {
	return fmt.Sprintf("%s/%s", permission.Owner, permission.Name)
}

func GetPermissionsByUser(userId string) []*Permission {
	permissions := []*Permission{}
	err := adapter.Engine.Where("users like ?", "%"+userId+"%").Find(&permissions)
	if err != nil {
		panic(err)
	}

	return permissions
}

func GetPermissionsByRole(roleId string) []*Permission {
	permissions := []*Permission{}
	err := adapter.Engine.Where("roles like ?", "%"+roleId+"%").Find(&permissions)
	if err != nil {
		panic(err)
	}

	return permissions
}

func GetPermissionsBySubmitter(owner string, submitter string) []*Permission {
	permissions := []*Permission{}
	err := adapter.Engine.Desc("created_time").Find(&permissions, &Permission{Owner: owner, Submitter: submitter})
	if err != nil {
		panic(err)
	}

	return permissions
}

func MigratePermissionRule() {
	models := []*Model{}
	err := adapter.Engine.Find(&models, &Model{})
	if err != nil {
		panic(err)
	}

	isHit := false
	for _, model := range models {
		if strings.Contains(model.ModelText, "permission") {
			// update model table
			model.ModelText = strings.Replace(model.ModelText, "permission,", "", -1)
			UpdateModel(model.GetId(), model)
			isHit = true
		}
	}

	if isHit {
		// update permission_rule table
		sql := "UPDATE `permission_rule`SET V0 = V1, V1 = V2, V2 = V3, V3 = V4, V4 = V5 WHERE V0 IN (SELECT CONCAT(owner, '/', name) AS permission_id FROM `permission`)"
		_, err = adapter.Engine.Exec(sql)
		if err != nil {
			return
		}
	}
}
