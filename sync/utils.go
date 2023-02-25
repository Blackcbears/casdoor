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

package sync

import (
	"bytes"
	"log"

	"github.com/xorm-io/xorm"
)

func GetUpdateSql(schemaName string, tableName string, columnNames []string, newColumnVal []string, oldColumnVal []string) string {
	var bt bytes.Buffer
	bt.WriteString("update " + schemaName + "." + tableName + " set ")
	for i, columnName := range columnNames {
		if i+1 < len(columnNames) {
			bt.WriteString(columnName + "=" + newColumnVal[i] + ", ")
		} else {
			bt.WriteString(columnName + "=" + newColumnVal[i] + " where ")
		}
	}

	for i, columnName := range columnNames {
		if i+1 < len(columnNames) {
			bt.WriteString(columnName + "=" + oldColumnVal[i] + " and ")
		} else {
			bt.WriteString(columnName + "=" + oldColumnVal[i] + ";")
		}
	}
	return bt.String()
}

func GetInsertSql(schemaName string, tableName string, columnNames []string, columnValue []string) string {
	var bt bytes.Buffer
	bt.WriteString("insert into " + schemaName + "." + tableName + " (")
	for i, columnName := range columnNames {
		if i+1 < len(columnNames) {
			bt.WriteString(columnName + ", ")
		} else {
			bt.WriteString(columnName + ") values (")
		}
	}
	for i, val := range columnValue {
		if i+1 < len(columnNames) {
			bt.WriteString(val + ", ")
		} else {
			bt.WriteString(val + ");")
		}
	}
	return bt.String()
}

func GetdeleteSql(schemaName string, tableName string, columnNames []string, columnValue []string) string {
	var bt bytes.Buffer
	bt.WriteString("delete from " + schemaName + "." + tableName + " where ")
	for i, columnName := range columnNames {
		if i+1 < len(columnNames) {
			bt.WriteString(columnName + " = " + columnValue[i] + " and ")
		} else {
			bt.WriteString(columnName + " = " + columnValue[i] + ";")
		}
	}
	return bt.String()
}

func CreateEngine(dataSourceName string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", dataSourceName)

	if err != nil {
		log.Fatal("connection mysql fail……", err)
	}

	// ping mysql
	err = engine.Ping()
	if err != nil {
		panic(err)
		return nil, err
	}

	log.Println("mysql connection success……")
	return engine, nil
}
