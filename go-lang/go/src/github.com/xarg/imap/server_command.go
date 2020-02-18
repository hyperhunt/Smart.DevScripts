// Copyright 2014 The imapsrv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the imapsrv.LICENSE file.

package imap

import (
	"fmt"
)

// An IMAP command
type command interface {
	// Execute the command and return an imap response
	execute(s *session) *serverResponse
}

//------------------------------------------------------------------------------

type noop struct {
	tag string
}

// Execute a noop
func (c *noop) execute(s *session) *serverResponse {
	return ok(c.tag, "NOOP Completed")
}

//------------------------------------------------------------------------------

// A CAPABILITY command
type capability struct {
	tag string
}

// Execute a capability
func (c *capability) execute(s *session) *serverResponse {
	// The IMAP server is assumed to be running over SSL and so
	// STARTTLS is not supported and LOGIN is not disabled
	return ok(c.tag, "CAPABILITY completed").
		extra("CAPABILITY IMAP4rev1")
}

//------------------------------------------------------------------------------

// A LOGIN command
type login struct {
	tag      string
	userId   string
	password string
}

// Login command
func (c *login) execute(sess *session) *serverResponse {

	// Has the user already logged in?
	if sess.st != notAuthenticated {
		message := "LOGIN already logged in"
		sess.log(message)
		return bad(c.tag, message)
	}

	// TODO: implement login
	if c.userId == "test" {
		sess.st = authenticated
		return ok(c.tag, "LOGIN completed")
	}

	// Fail by default
	return no(c.tag, "LOGIN failure")
}

//------------------------------------------------------------------------------

// A LOGOUT command
type logout struct {
	tag string
}

// Logout command
func (c *logout) execute(sess *session) *serverResponse {

	sess.st = notAuthenticated
	return ok(c.tag, "LOGOUT completed").
		extra("BYE IMAP4rev1 Server logging out").
		shouldClose()
}

//------------------------------------------------------------------------------

// A SELECT command
type selectMailbox struct {
	tag     string
	mailbox string
}

// Select command
func (c *selectMailbox) execute(sess *session) *serverResponse {

	// Is the user authenticated?
	if sess.st != authenticated {
		message := "SELECT not authenticated"
		sess.log(message)
		return bad(c.tag, message)
	}

	// Select the mailbox
	exists, err := sess.selectMailbox(c.mailbox)

	if err != nil {
		return internalError(sess, c.tag, "SELECT", err)
	}

	if !exists {
		return no(c.tag, "SELECT No such mailbox")
	}

	// Build a response that includes mailbox information
	res := ok(c.tag, "SELECT completed")

	err = sess.addMailboxInfo(res)

	if err != nil {
		return internalError(sess, c.tag, "SELECT", err)
	}

	return res
}

//------------------------------------------------------------------------------

// An unknown/unsupported command
type unknownCommand struct {
	tag string
	cmd string
}

// Report an error for an unknown command
func (c *unknownCommand) execute(s *session) *serverResponse {
	message := fmt.Sprintf("%s unknown command", c.cmd)
	s.log(message)
	return bad(c.tag, message)
}

//------ Helper functions ------------------------------------------------------

// Log an error and return an response
func internalError(sess *session, tag string, commandName string, err error) *serverResponse {
	message := commandName + " " + err.Error()
	sess.log(message)
	return no(tag, message).shouldClose()
}
