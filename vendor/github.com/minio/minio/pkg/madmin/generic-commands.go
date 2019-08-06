/*
 * Minio Cloud Storage, (C) 2016, 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package madmin

import (
	"encoding/json"
	"net/http"
)

// SetCredsReq - xml to send to the server to set new credentials
type SetCredsReq struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// SetAdminCredentials - Call Set Credentials API to set new access and
// secret keys in the specified Minio server
func (adm *AdminClient) SetAdminCredentials(access, secret string) error {
	// Setup request's body
	body, err := json.Marshal(SetCredsReq{access, secret})
	if err != nil {
		return err
	}

	ebody, err := EncryptData(adm.secretAccessKey, body)
	if err != nil {
		return err
	}

	// Setup new request
	reqData := requestData{
		relPath: "/v1/config/credential",
		content: ebody,
	}

	// Execute GET on bucket to list objects.
	resp, err := adm.executeMethod("PUT", reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	// Return error to the caller if http response code is
	// different from 200
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}
