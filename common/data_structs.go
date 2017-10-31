// Copyright 2017 Google Inc.
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
package common

type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Kind  string `json:"kind,omitempty"`
}

type ErrorResponse struct {
	ResponseCode    string `json:"response_code,omitempty"`
	ResponseMessage string `json:"response_message,omitempty"`
	StatusCode      int    `json:"-"`
	Kind            string `json:"kind,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.ResponseMessage
}

type ApiProduct struct {
	Id            string `db:"id"`
	Name          string `db:"name"`
	DisplayName   string `db:"display_name"`
	Description   string `db:"description"`
	ApiResources  string `db:"api_resources"`
	ApprovalType  string `db:"approval_type"`
	Scopes        string `db:"scopes"`
	Proxies       string `db:"proxies"`
	Environments  string `db:"environments"`
	Quota         string `db:"quota"`
	QuotaTimeUnit string `db:"quota_time_unit"`
	QuotaInterval int64  `db:"quota_interval"`
	CreatedAt     string `db:"created_at"`
	CreatedBy     string `db:"created_by"`
	UpdatedAt     string `db:"updated_at"`
	UpdatedBy     string `db:"updated_by"`
	TenantId      string `db:"tenant_id"`
}

type App struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	DisplayName string `db:"display_name"`
	AccessType  string `db:"access_type"`
	CallbackUrl string `db:"callback_url"`
	Status      string `db:"status"`
	AppFamily   string `db:"app_family"`
	CompanyId   string `db:"company_id"`
	DeveloperId string `db:"developer_id"`
	ParentId    string `db:"parent_id"`
	Type        string `db:"type"`
	CreatedAt   string `db:"created_at"`
	CreatedBy   string `db:"created_by"`
	UpdatedAt   string `db:"updated_at"`
	UpdatedBy   string `db:"updated_by"`
}
