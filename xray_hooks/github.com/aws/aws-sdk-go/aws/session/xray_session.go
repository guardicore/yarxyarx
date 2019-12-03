package session

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-xray-sdk-go/xray"
)

func New(cfgs ...*aws.Config) *Session {
    return xray.AWSSession(session.New(cfgs...))
}

func NewSession(cfgs ...*aws.Config) (*Session, error) {
    session, err := session.NewSession(cfgs...)
    if err != nil {
        return nil, err
    }
    return xray.AWSSession(session), nil
}

func NewSessionWithOptions(opts Options) (*Session, error) {
    session, err := session.NewSessionWithOptions(opts)
    if err != nil {
        return nil, err
    }
    return xray.AWSSession(session), nil
}
