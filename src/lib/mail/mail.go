package mail

import (
	"github.com/fragmenta/view"
	"github.com/sendgrid/sendgrid-go"
)

// TODO - instead of package variables, use New() mailer and put the variables in the mailer?
// Try to be consistent about use of pkg variables

// The Mail service user (must be set before first sending)
var user string

// The Mail service secret key/password (must be set before first sending)
var secret string

var from string

// Setup sets the user and secret for use in sending mail (possibly later we should have a config etc)
func Setup(u string, s string, f string) {
	user = u
	secret = s
	from = f
}

// Send sends mail
func Send(recipients []string, subject string, template string, context map[string]interface{}) error {

	// At present we use sendgrid, we may later allow many mail services to be used
	//	sg := sendgrid.NewSendGridClient(user, secret)
	sg := sendgrid.NewSendGridClientWithApiKey(secret)

	message := sendgrid.NewMail()
	message.SetFrom(from)
	message.AddTos(recipients)
	message.SetSubject(subject)

	// TODO: consider best way to set arbitrary email headers
	// perhaps require mail.New() for object which we can then set headers on etc...
	if context["reply_to"] != nil {
		replyTo := context["reply_to"].(string)
		message.SetReplyTo(replyTo)
	}

	// Load the template, and substitute using context
	// We should possibly set layout from caller too?
	view := view.NewWithPath("", nil)
	view.Template(template)
	view.Context(context)

	html, err := view.RenderToString()
	if err != nil {
		return err
	}

	message.SetHTML(html)

	return sg.Send(message)
}

// SendOne sends email to ONE recipient only
func SendOne(recipient string, subject string, template string, context map[string]interface{}) error {
	return Send([]string{recipient}, subject, template, context)
}
