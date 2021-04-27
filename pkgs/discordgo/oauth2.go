// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains functions related to Discord OAuth2 endpoints

package discordgo

// ------------------------------------------------------------------------------------------------
// Code specific to Discord OAuth2 Applications
// ------------------------------------------------------------------------------------------------

// The MembershipState represents whether the user is in the team or has been invited into it
type MembershipState int

// Constants for the different stages of the MembershipState
const (
	MembershipStateInvited MembershipState = iota + 1
	MembershipStateAccepted
)

// A TeamMember struct stores values for a single Team Member, extending the normal User data - note that the user field is partial
type TeamMember struct {
	User            *User           `json:"user"`
	TeamID          int64           `json:"team_id,string"`
	MembershipState MembershipState `json:"membership_state"`
	Permissions     []string        `json:"permissions"`
}

// A Team struct stores the members of a Discord Developer Team as well as some metadata about it
type Team struct {
	ID          int64         `json:"id,string"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	OwnerID     int64         `json:"owner_user_id,string"`
	Members     []*TeamMember `json:"members"`
}

// An Application struct stores values for a Discord OAuth2 Application
type Application struct {
	ID                  int64     `json:"id,string,omitempty"`
	Name                string    `json:"name"`
	Icon                string    `json:"icon,omitempty"`
	Description         string    `json:"description,omitempty"`
	RPCOrigins          *[]string `json:"rpc_origins,omitempty"`
	BotPublic           bool      `json:"bot_public"`
	BotRequireCodeGrant bool      `json:"bot_require_code_grant,omitempty"`
	Owner               *User     `json:"owner"`
	Summary             string    `json:"summary"`
	VerifyKey           string    `json:"verify_key"`
	Team                *Team     `json:"team"`
	GuildID             int64     `json:"guild_id,string"`
	PrimarySkuID        int64     `json:"primary_sku_id,string"`
	Slug                string    `json:"slug,omitempty"`
	Flags               int       `json:"flags,omitempty"`
}

// Application returns an Application structure of a specific Application
//   appID : The ID of an Application
func (s *Session) Application(appID int64) (st *Application, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointOAuth2Application(appID), nil, EndpointOAuth2Application(0))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// Applications returns all applications for the authenticated user
func (s *Session) Applications() (st []*Application, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointOAuth2Applications, nil, EndpointOAuth2Applications)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationCreate creates a new Application
//    name : Name of Application / Bot
//    uris : Redirect URIs (Not required)
func (s *Session) ApplicationCreate(ap *Application) (st *Application, err error) {

	data := struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RedirectURIs *[]string `json:"redirect_uris,omitempty"`
	}{ap.Name, ap.Description, ap.RPCOrigins}

	body, err := s.RequestWithBucketID("POST", EndpointOAuth2Applications, data, EndpointOAuth2Applications)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationUpdate updates an existing Application
//   var : desc
func (s *Session) ApplicationUpdate(appID int64, ap *Application) (st *Application, err error) {

	data := struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RedirectURIs *[]string `json:"redirect_uris,omitempty"`
	}{ap.Name, ap.Description, ap.RPCOrigins}

	body, err := s.RequestWithBucketID("PUT", EndpointOAuth2Application(appID), data, EndpointOAuth2Application(0))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

// ApplicationDelete deletes an existing Application
//   appID : The ID of an Application
func (s *Session) ApplicationDelete(appID int64) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointOAuth2Application(appID), nil, EndpointOAuth2Application(0))
	if err != nil {
		return
	}

	return
}

// Asset struct stores values for an asset of an application
type Asset struct {
	Type int    `json:"type"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ApplicationAssets returns an application's assets
func (s *Session) ApplicationAssets(appID int64) (ass []*Asset, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointOAuth2ApplicationAssets(appID), nil, EndpointOAuth2ApplicationAssets(0))
	if err != nil {
		return
	}

	err = unmarshal(body, &ass)
	return
}

// ------------------------------------------------------------------------------------------------
// Code specific to Discord OAuth2 Application Bots
// ------------------------------------------------------------------------------------------------

// ApplicationBotCreate creates an Application Bot Account
//
//   appID : The ID of an Application
//
// NOTE: func name may change, if I can think up something better.
func (s *Session) ApplicationBotCreate(appID int64) (st *User, err error) {

	body, err := s.RequestWithBucketID("POST", EndpointOAuth2ApplicationsBot(appID), nil, EndpointOAuth2ApplicationsBot(0))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}