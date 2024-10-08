// This file was auto-generated by Fern from our API Definition.

package methodokta

import (
	json "encoding/json"
	fmt "fmt"
	core "github.com/method-security/methodokta/generated/go/core"
	time "time"
)

type Application struct {
	Uid        string       `json:"uid" url:"uid"`
	Name       string       `json:"name" url:"name"`
	Label      string       `json:"label" url:"label"`
	Url        *string      `json:"url,omitempty" url:"url,omitempty"`
	Status     StatusType   `json:"status" url:"status"`
	Created    time.Time    `json:"created" url:"created"`
	AuthMethod AuthType     `json:"authMethod" url:"authMethod"`
	Groups     []*GroupInfo `json:"groups,omitempty" url:"groups,omitempty"`
	Users      []*UserInfo  `json:"users,omitempty" url:"users,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (a *Application) GetExtraProperties() map[string]interface{} {
	return a.extraProperties
}

func (a *Application) UnmarshalJSON(data []byte) error {
	type embed Application
	var unmarshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed: embed(*a),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*a = Application(unmarshaler.embed)
	a.Created = unmarshaler.Created.Time()

	extraProperties, err := core.ExtractExtraProperties(data, *a)
	if err != nil {
		return err
	}
	a.extraProperties = extraProperties

	a._rawJSON = json.RawMessage(data)
	return nil
}

func (a *Application) MarshalJSON() ([]byte, error) {
	type embed Application
	var marshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed:   embed(*a),
		Created: core.NewDateTime(a.Created),
	}
	return json.Marshal(marshaler)
}

func (a *Application) String() string {
	if len(a._rawJSON) > 0 {
		if value, err := core.StringifyJSON(a._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(a); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", a)
}

type ApplicationReport struct {
	Org          string         `json:"org" url:"org"`
	Applications []*Application `json:"applications,omitempty" url:"applications,omitempty"`
	Errors       []string       `json:"errors,omitempty" url:"errors,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (a *ApplicationReport) GetExtraProperties() map[string]interface{} {
	return a.extraProperties
}

func (a *ApplicationReport) UnmarshalJSON(data []byte) error {
	type unmarshaler ApplicationReport
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*a = ApplicationReport(value)

	extraProperties, err := core.ExtractExtraProperties(data, *a)
	if err != nil {
		return err
	}
	a.extraProperties = extraProperties

	a._rawJSON = json.RawMessage(data)
	return nil
}

func (a *ApplicationReport) String() string {
	if len(a._rawJSON) > 0 {
		if value, err := core.StringifyJSON(a._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(a); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", a)
}

type AuthType string

const (
	AuthTypeBrowserPlugin AuthType = "BROWSER_PLUGIN"
	AuthTypeOpenidConnect AuthType = "OPENID_CONNECT"
	AuthTypeSaml20        AuthType = "SAML_2_0"
	AuthTypeUnknown       AuthType = "UNKNOWN"
)

func NewAuthTypeFromString(s string) (AuthType, error) {
	switch s {
	case "BROWSER_PLUGIN":
		return AuthTypeBrowserPlugin, nil
	case "OPENID_CONNECT":
		return AuthTypeOpenidConnect, nil
	case "SAML_2_0":
		return AuthTypeSaml20, nil
	case "UNKNOWN":
		return AuthTypeUnknown, nil
	}
	var t AuthType
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

func (a AuthType) Ptr() *AuthType {
	return &a
}

type ApplicationInfo struct {
	Uid  string `json:"uid" url:"uid"`
	Name string `json:"name" url:"name"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (a *ApplicationInfo) GetExtraProperties() map[string]interface{} {
	return a.extraProperties
}

func (a *ApplicationInfo) UnmarshalJSON(data []byte) error {
	type unmarshaler ApplicationInfo
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*a = ApplicationInfo(value)

	extraProperties, err := core.ExtractExtraProperties(data, *a)
	if err != nil {
		return err
	}
	a.extraProperties = extraProperties

	a._rawJSON = json.RawMessage(data)
	return nil
}

func (a *ApplicationInfo) String() string {
	if len(a._rawJSON) > 0 {
		if value, err := core.StringifyJSON(a._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(a); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", a)
}

type GroupInfo struct {
	Uid  string `json:"uid" url:"uid"`
	Name string `json:"name" url:"name"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (g *GroupInfo) GetExtraProperties() map[string]interface{} {
	return g.extraProperties
}

func (g *GroupInfo) UnmarshalJSON(data []byte) error {
	type unmarshaler GroupInfo
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*g = GroupInfo(value)

	extraProperties, err := core.ExtractExtraProperties(data, *g)
	if err != nil {
		return err
	}
	g.extraProperties = extraProperties

	g._rawJSON = json.RawMessage(data)
	return nil
}

func (g *GroupInfo) String() string {
	if len(g._rawJSON) > 0 {
		if value, err := core.StringifyJSON(g._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(g); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", g)
}

type RoleInfo struct {
	Uid         string     `json:"uid" url:"uid"`
	Name        *string    `json:"name,omitempty" url:"name,omitempty"`
	Type        *RoleType  `json:"type,omitempty" url:"type,omitempty"`
	Description *string    `json:"description,omitempty" url:"description,omitempty"`
	Created     *time.Time `json:"created,omitempty" url:"created,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (r *RoleInfo) GetExtraProperties() map[string]interface{} {
	return r.extraProperties
}

func (r *RoleInfo) UnmarshalJSON(data []byte) error {
	type embed RoleInfo
	var unmarshaler = struct {
		embed
		Created *core.DateTime `json:"created,omitempty"`
	}{
		embed: embed(*r),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*r = RoleInfo(unmarshaler.embed)
	r.Created = unmarshaler.Created.TimePtr()

	extraProperties, err := core.ExtractExtraProperties(data, *r)
	if err != nil {
		return err
	}
	r.extraProperties = extraProperties

	r._rawJSON = json.RawMessage(data)
	return nil
}

func (r *RoleInfo) MarshalJSON() ([]byte, error) {
	type embed RoleInfo
	var marshaler = struct {
		embed
		Created *core.DateTime `json:"created,omitempty"`
	}{
		embed:   embed(*r),
		Created: core.NewOptionalDateTime(r.Created),
	}
	return json.Marshal(marshaler)
}

func (r *RoleInfo) String() string {
	if len(r._rawJSON) > 0 {
		if value, err := core.StringifyJSON(r._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(r); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", r)
}

type RoleType string

const (
	RoleTypeSuperAdmin RoleType = "SUPER_ADMIN"
	RoleTypeAdmin      RoleType = "ADMIN"
	RoleTypeUser       RoleType = "USER"
	RoleTypeReadOnly   RoleType = "READ_ONLY"
	RoleTypeUnknown    RoleType = "UNKNOWN"
)

func NewRoleTypeFromString(s string) (RoleType, error) {
	switch s {
	case "SUPER_ADMIN":
		return RoleTypeSuperAdmin, nil
	case "ADMIN":
		return RoleTypeAdmin, nil
	case "USER":
		return RoleTypeUser, nil
	case "READ_ONLY":
		return RoleTypeReadOnly, nil
	case "UNKNOWN":
		return RoleTypeUnknown, nil
	}
	var t RoleType
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

func (r RoleType) Ptr() *RoleType {
	return &r
}

type StatusType string

const (
	StatusTypeActive   StatusType = "ACTIVE"
	StatusTypeInactive StatusType = "INACTIVE"
	StatusTypeUnknown  StatusType = "UNKNOWN"
)

func NewStatusTypeFromString(s string) (StatusType, error) {
	switch s {
	case "ACTIVE":
		return StatusTypeActive, nil
	case "INACTIVE":
		return StatusTypeInactive, nil
	case "UNKNOWN":
		return StatusTypeUnknown, nil
	}
	var t StatusType
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

func (s StatusType) Ptr() *StatusType {
	return &s
}

type UserInfo struct {
	Uid   string `json:"uid" url:"uid"`
	Email string `json:"email" url:"email"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (u *UserInfo) GetExtraProperties() map[string]interface{} {
	return u.extraProperties
}

func (u *UserInfo) UnmarshalJSON(data []byte) error {
	type unmarshaler UserInfo
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*u = UserInfo(value)

	extraProperties, err := core.ExtractExtraProperties(data, *u)
	if err != nil {
		return err
	}
	u.extraProperties = extraProperties

	u._rawJSON = json.RawMessage(data)
	return nil
}

func (u *UserInfo) String() string {
	if len(u._rawJSON) > 0 {
		if value, err := core.StringifyJSON(u._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(u); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", u)
}

type Device struct {
	Uid          string      `json:"uid" url:"uid"`
	Name         string      `json:"name" url:"name"`
	Platform     string      `json:"platform" url:"platform"`
	Manufacturer *string     `json:"manufacturer,omitempty" url:"manufacturer,omitempty"`
	Model        *string     `json:"model,omitempty" url:"model,omitempty"`
	OsVersion    *string     `json:"osVersion,omitempty" url:"osVersion,omitempty"`
	Status       StatusType  `json:"status" url:"status"`
	Created      time.Time   `json:"created" url:"created"`
	Users        []*UserInfo `json:"users,omitempty" url:"users,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (d *Device) GetExtraProperties() map[string]interface{} {
	return d.extraProperties
}

func (d *Device) UnmarshalJSON(data []byte) error {
	type embed Device
	var unmarshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed: embed(*d),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*d = Device(unmarshaler.embed)
	d.Created = unmarshaler.Created.Time()

	extraProperties, err := core.ExtractExtraProperties(data, *d)
	if err != nil {
		return err
	}
	d.extraProperties = extraProperties

	d._rawJSON = json.RawMessage(data)
	return nil
}

func (d *Device) MarshalJSON() ([]byte, error) {
	type embed Device
	var marshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed:   embed(*d),
		Created: core.NewDateTime(d.Created),
	}
	return json.Marshal(marshaler)
}

func (d *Device) String() string {
	if len(d._rawJSON) > 0 {
		if value, err := core.StringifyJSON(d._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(d); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", d)
}

type DeviceReport struct {
	Org     string    `json:"org" url:"org"`
	Devices []*Device `json:"devices,omitempty" url:"devices,omitempty"`
	Errors  []string  `json:"errors,omitempty" url:"errors,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (d *DeviceReport) GetExtraProperties() map[string]interface{} {
	return d.extraProperties
}

func (d *DeviceReport) UnmarshalJSON(data []byte) error {
	type unmarshaler DeviceReport
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*d = DeviceReport(value)

	extraProperties, err := core.ExtractExtraProperties(data, *d)
	if err != nil {
		return err
	}
	d.extraProperties = extraProperties

	d._rawJSON = json.RawMessage(data)
	return nil
}

func (d *DeviceReport) String() string {
	if len(d._rawJSON) > 0 {
		if value, err := core.StringifyJSON(d._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(d); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", d)
}

type Group struct {
	Uid          string             `json:"uid" url:"uid"`
	Name         string             `json:"name" url:"name"`
	Type         GroupType          `json:"type" url:"type"`
	Description  *string            `json:"description,omitempty" url:"description,omitempty"`
	Created      time.Time          `json:"created" url:"created"`
	Applications []*ApplicationInfo `json:"applications,omitempty" url:"applications,omitempty"`
	Roles        []*RoleInfo        `json:"roles,omitempty" url:"roles,omitempty"`
	Users        []*UserInfo        `json:"users,omitempty" url:"users,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (g *Group) GetExtraProperties() map[string]interface{} {
	return g.extraProperties
}

func (g *Group) UnmarshalJSON(data []byte) error {
	type embed Group
	var unmarshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed: embed(*g),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*g = Group(unmarshaler.embed)
	g.Created = unmarshaler.Created.Time()

	extraProperties, err := core.ExtractExtraProperties(data, *g)
	if err != nil {
		return err
	}
	g.extraProperties = extraProperties

	g._rawJSON = json.RawMessage(data)
	return nil
}

func (g *Group) MarshalJSON() ([]byte, error) {
	type embed Group
	var marshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed:   embed(*g),
		Created: core.NewDateTime(g.Created),
	}
	return json.Marshal(marshaler)
}

func (g *Group) String() string {
	if len(g._rawJSON) > 0 {
		if value, err := core.StringifyJSON(g._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(g); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", g)
}

type GroupReport struct {
	Org    string   `json:"org" url:"org"`
	Groups []*Group `json:"groups,omitempty" url:"groups,omitempty"`
	Errors []string `json:"errors,omitempty" url:"errors,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (g *GroupReport) GetExtraProperties() map[string]interface{} {
	return g.extraProperties
}

func (g *GroupReport) UnmarshalJSON(data []byte) error {
	type unmarshaler GroupReport
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*g = GroupReport(value)

	extraProperties, err := core.ExtractExtraProperties(data, *g)
	if err != nil {
		return err
	}
	g.extraProperties = extraProperties

	g._rawJSON = json.RawMessage(data)
	return nil
}

func (g *GroupReport) String() string {
	if len(g._rawJSON) > 0 {
		if value, err := core.StringifyJSON(g._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(g); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", g)
}

type GroupType string

const (
	GroupTypeOktaGroup GroupType = "OKTA_GROUP"
	GroupTypeAppGroup  GroupType = "APP_GROUP"
	GroupTypeBuiltIn   GroupType = "BUILT_IN"
	GroupTypeUnknown   GroupType = "UNKNOWN"
)

func NewGroupTypeFromString(s string) (GroupType, error) {
	switch s {
	case "OKTA_GROUP":
		return GroupTypeOktaGroup, nil
	case "APP_GROUP":
		return GroupTypeAppGroup, nil
	case "BUILT_IN":
		return GroupTypeBuiltIn, nil
	case "UNKNOWN":
		return GroupTypeUnknown, nil
	}
	var t GroupType
	return "", fmt.Errorf("%s is not a valid %T", s, t)
}

func (g GroupType) Ptr() *GroupType {
	return &g
}

type Login struct {
	User        *UserInfo        `json:"user,omitempty" url:"user,omitempty"`
	Application *ApplicationInfo `json:"application,omitempty" url:"application,omitempty"`
	Count       int              `json:"count" url:"count"`
	TimeFrame   int              `json:"timeFrame" url:"timeFrame"`
	Last        time.Time        `json:"last" url:"last"`
	ScanDate    time.Time        `json:"scanDate" url:"scanDate"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (l *Login) GetExtraProperties() map[string]interface{} {
	return l.extraProperties
}

func (l *Login) UnmarshalJSON(data []byte) error {
	type embed Login
	var unmarshaler = struct {
		embed
		Last     *core.DateTime `json:"last"`
		ScanDate *core.DateTime `json:"scanDate"`
	}{
		embed: embed(*l),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*l = Login(unmarshaler.embed)
	l.Last = unmarshaler.Last.Time()
	l.ScanDate = unmarshaler.ScanDate.Time()

	extraProperties, err := core.ExtractExtraProperties(data, *l)
	if err != nil {
		return err
	}
	l.extraProperties = extraProperties

	l._rawJSON = json.RawMessage(data)
	return nil
}

func (l *Login) MarshalJSON() ([]byte, error) {
	type embed Login
	var marshaler = struct {
		embed
		Last     *core.DateTime `json:"last"`
		ScanDate *core.DateTime `json:"scanDate"`
	}{
		embed:    embed(*l),
		Last:     core.NewDateTime(l.Last),
		ScanDate: core.NewDateTime(l.ScanDate),
	}
	return json.Marshal(marshaler)
}

func (l *Login) String() string {
	if len(l._rawJSON) > 0 {
		if value, err := core.StringifyJSON(l._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(l); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", l)
}

type LoginReport struct {
	Org    string   `json:"org" url:"org"`
	Logins []*Login `json:"logins,omitempty" url:"logins,omitempty"`
	Errors []string `json:"errors,omitempty" url:"errors,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (l *LoginReport) GetExtraProperties() map[string]interface{} {
	return l.extraProperties
}

func (l *LoginReport) UnmarshalJSON(data []byte) error {
	type unmarshaler LoginReport
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*l = LoginReport(value)

	extraProperties, err := core.ExtractExtraProperties(data, *l)
	if err != nil {
		return err
	}
	l.extraProperties = extraProperties

	l._rawJSON = json.RawMessage(data)
	return nil
}

func (l *LoginReport) String() string {
	if len(l._rawJSON) > 0 {
		if value, err := core.StringifyJSON(l._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(l); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", l)
}

type OrgInfo struct {
	Uid         string     `json:"uid" url:"uid"`
	CompanyName *string    `json:"companyName,omitempty" url:"companyName,omitempty"`
	Url         string     `json:"url" url:"url"`
	Status      StatusType `json:"status" url:"status"`
	Created     time.Time  `json:"created" url:"created"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (o *OrgInfo) GetExtraProperties() map[string]interface{} {
	return o.extraProperties
}

func (o *OrgInfo) UnmarshalJSON(data []byte) error {
	type embed OrgInfo
	var unmarshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed: embed(*o),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*o = OrgInfo(unmarshaler.embed)
	o.Created = unmarshaler.Created.Time()

	extraProperties, err := core.ExtractExtraProperties(data, *o)
	if err != nil {
		return err
	}
	o.extraProperties = extraProperties

	o._rawJSON = json.RawMessage(data)
	return nil
}

func (o *OrgInfo) MarshalJSON() ([]byte, error) {
	type embed OrgInfo
	var marshaler = struct {
		embed
		Created *core.DateTime `json:"created"`
	}{
		embed:   embed(*o),
		Created: core.NewDateTime(o.Created),
	}
	return json.Marshal(marshaler)
}

func (o *OrgInfo) String() string {
	if len(o._rawJSON) > 0 {
		if value, err := core.StringifyJSON(o._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(o); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", o)
}

type OrgReport struct {
	Org          string   `json:"org" url:"org"`
	Organization *OrgInfo `json:"organization,omitempty" url:"organization,omitempty"`
	Errors       []string `json:"errors,omitempty" url:"errors,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (o *OrgReport) GetExtraProperties() map[string]interface{} {
	return o.extraProperties
}

func (o *OrgReport) UnmarshalJSON(data []byte) error {
	type unmarshaler OrgReport
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*o = OrgReport(value)

	extraProperties, err := core.ExtractExtraProperties(data, *o)
	if err != nil {
		return err
	}
	o.extraProperties = extraProperties

	o._rawJSON = json.RawMessage(data)
	return nil
}

func (o *OrgReport) String() string {
	if len(o._rawJSON) > 0 {
		if value, err := core.StringifyJSON(o._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(o); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", o)
}

type User struct {
	Uid             string             `json:"uid" url:"uid"`
	Firstname       *string            `json:"firstname,omitempty" url:"firstname,omitempty"`
	Lastname        *string            `json:"lastname,omitempty" url:"lastname,omitempty"`
	Email           string             `json:"email" url:"email"`
	Status          StatusType         `json:"status" url:"status"`
	Created         time.Time          `json:"created" url:"created"`
	PasswordChanged *time.Time         `json:"passwordChanged,omitempty" url:"passwordChanged,omitempty"`
	Applications    []*ApplicationInfo `json:"applications,omitempty" url:"applications,omitempty"`
	Groups          []*GroupInfo       `json:"groups,omitempty" url:"groups,omitempty"`
	Roles           []*RoleInfo        `json:"roles,omitempty" url:"roles,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (u *User) GetExtraProperties() map[string]interface{} {
	return u.extraProperties
}

func (u *User) UnmarshalJSON(data []byte) error {
	type embed User
	var unmarshaler = struct {
		embed
		Created         *core.DateTime `json:"created"`
		PasswordChanged *core.DateTime `json:"passwordChanged,omitempty"`
	}{
		embed: embed(*u),
	}
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*u = User(unmarshaler.embed)
	u.Created = unmarshaler.Created.Time()
	u.PasswordChanged = unmarshaler.PasswordChanged.TimePtr()

	extraProperties, err := core.ExtractExtraProperties(data, *u)
	if err != nil {
		return err
	}
	u.extraProperties = extraProperties

	u._rawJSON = json.RawMessage(data)
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	type embed User
	var marshaler = struct {
		embed
		Created         *core.DateTime `json:"created"`
		PasswordChanged *core.DateTime `json:"passwordChanged,omitempty"`
	}{
		embed:           embed(*u),
		Created:         core.NewDateTime(u.Created),
		PasswordChanged: core.NewOptionalDateTime(u.PasswordChanged),
	}
	return json.Marshal(marshaler)
}

func (u *User) String() string {
	if len(u._rawJSON) > 0 {
		if value, err := core.StringifyJSON(u._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(u); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", u)
}

type UserReport struct {
	Org    string   `json:"org" url:"org"`
	Users  []*User  `json:"users,omitempty" url:"users,omitempty"`
	Errors []string `json:"errors,omitempty" url:"errors,omitempty"`

	extraProperties map[string]interface{}
	_rawJSON        json.RawMessage
}

func (u *UserReport) GetExtraProperties() map[string]interface{} {
	return u.extraProperties
}

func (u *UserReport) UnmarshalJSON(data []byte) error {
	type unmarshaler UserReport
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*u = UserReport(value)

	extraProperties, err := core.ExtractExtraProperties(data, *u)
	if err != nil {
		return err
	}
	u.extraProperties = extraProperties

	u._rawJSON = json.RawMessage(data)
	return nil
}

func (u *UserReport) String() string {
	if len(u._rawJSON) > 0 {
		if value, err := core.StringifyJSON(u._rawJSON); err == nil {
			return value
		}
	}
	if value, err := core.StringifyJSON(u); err == nil {
		return value
	}
	return fmt.Sprintf("%#v", u)
}
