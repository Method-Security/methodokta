/*
Okta Admin Management

Allows customers to easily access the Okta Management APIs

Copyright 2018 - Present Okta, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

API version: 2024.06.1
Contact: devex-public@okta.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.
package okta

import (
	"encoding/json"
)

// JwkUse struct for JwkUse
type JwkUse struct {
	Use *string `json:"use,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _JwkUse JwkUse

// NewJwkUse instantiates a new JwkUse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJwkUse() *JwkUse {
	this := JwkUse{}
	return &this
}

// NewJwkUseWithDefaults instantiates a new JwkUse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJwkUseWithDefaults() *JwkUse {
	this := JwkUse{}
	return &this
}

// GetUse returns the Use field value if set, zero value otherwise.
func (o *JwkUse) GetUse() string {
	if o == nil || o.Use == nil {
		var ret string
		return ret
	}
	return *o.Use
}

// GetUseOk returns a tuple with the Use field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JwkUse) GetUseOk() (*string, bool) {
	if o == nil || o.Use == nil {
		return nil, false
	}
	return o.Use, true
}

// HasUse returns a boolean if a field has been set.
func (o *JwkUse) HasUse() bool {
	if o != nil && o.Use != nil {
		return true
	}

	return false
}

// SetUse gets a reference to the given string and assigns it to the Use field.
func (o *JwkUse) SetUse(v string) {
	o.Use = &v
}

func (o JwkUse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Use != nil {
		toSerialize["use"] = o.Use
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *JwkUse) UnmarshalJSON(bytes []byte) (err error) {
	varJwkUse := _JwkUse{}

	err = json.Unmarshal(bytes, &varJwkUse)
	if err == nil {
		*o = JwkUse(varJwkUse)
	} else {
		return err
	}

	additionalProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &additionalProperties)
	if err == nil {
		delete(additionalProperties, "use")
		o.AdditionalProperties = additionalProperties
	} else {
		return err
	}

	return err
}

type NullableJwkUse struct {
	value *JwkUse
	isSet bool
}

func (v NullableJwkUse) Get() *JwkUse {
	return v.value
}

func (v *NullableJwkUse) Set(val *JwkUse) {
	v.value = val
	v.isSet = true
}

func (v NullableJwkUse) IsSet() bool {
	return v.isSet
}

func (v *NullableJwkUse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJwkUse(val *JwkUse) *NullableJwkUse {
	return &NullableJwkUse{value: val, isSet: true}
}

func (v NullableJwkUse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJwkUse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

