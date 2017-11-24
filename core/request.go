package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ghettovoice/gosip/log"
)

// Request RFC 3261 - 7.1.
type Request interface {
	Message
	Method() RequestMethod
	SetMethod(method RequestMethod)
	Recipient() Uri
	SetRecipient(recipient Uri)
	/* Common Helpers */
	IsInvite() bool
	IsAck() bool
}

type request struct {
	message
	method    RequestMethod
	recipient Uri
}

func NewRequest(
	method RequestMethod,
	recipient Uri,
	sipVersion string,
	hdrs []Header,
	body string,
) Request {
	req := new(request)
	req.SetSipVersion(sipVersion)
	req.headers = newHeaders(hdrs)
	req.SetMethod(method)
	req.SetRecipient(recipient)
	req.SetLog(log.StandardLogger())

	if strings.TrimSpace(body) != "" {
		req.SetBody(body)
	}

	return req
}

func (req *request) Method() RequestMethod {
	return req.method
}
func (req *request) SetMethod(method RequestMethod) {
	req.method = method
}

func (req *request) Recipient() Uri {
	return req.recipient
}
func (req *request) SetRecipient(recipient Uri) {
	req.recipient = recipient
}

// StartLine returns Request Line - RFC 2361 7.1.
func (req *request) StartLine() string {
	var buffer bytes.Buffer

	// Every SIP request starts with a Request Line - RFC 2361 7.1.
	buffer.WriteString(
		fmt.Sprintf(
			"%s %s %s",
			string(req.method),
			req.Recipient(),
			req.SipVersion(),
		),
	)

	return buffer.String()
}

func (req *request) Short() string {
	var buffer bytes.Buffer

	buffer.WriteString(req.StartLine())

	if cseq, ok := req.CSeq(); ok {
		buffer.WriteString(fmt.Sprintf(" (%s)", cseq))
	}
	if callId, ok := req.CallId(); ok {
		buffer.WriteString(fmt.Sprintf(" (%s)", callId))
	}
	if from, ok := req.From(); ok {
		buffer.WriteString(fmt.Sprintf(" (%s)", from))
	}
	if to, ok := req.To(); ok {
		buffer.WriteString(fmt.Sprintf(" (%s)", to))
	}

	return buffer.String()
}

func (req *request) String() string {
	var buffer bytes.Buffer

	// write message start line
	buffer.WriteString(req.StartLine() + "\r\n")
	// Write the headers.
	buffer.WriteString(req.headers.String())
	// If the request has a message body, add it.
	buffer.WriteString("\r\n" + req.Body())

	return buffer.String()
}

func (req *request) Clone() Message {
	clone := NewRequest(
		req.Method(),
		req.Recipient().Clone(),
		req.SipVersion(),
		req.headers.CloneHeaders(),
		req.Body(),
	)
	clone.SetLog(req.Log())
	return clone
}

func (req *request) IsInvite() bool {
	return req.Method() == INVITE
}

func (req *request) IsAck() bool {
	return req.Method() == ACK
}
