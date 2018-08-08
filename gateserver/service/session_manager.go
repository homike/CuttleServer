package service

import "errors"

type SessionManager struct {
	currentSessionID uint32
	Sessions         map[uint32]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions:         make(map[uint32]*Session),
		currentSessionID: 1,
	}
}

func (self *SessionManager) AssociateSession(sess *Session) error {
	if sess == nil {
		return errors.New("session is nil")
	}
	self.Sessions[self.currentSessionID] = sess
	sess.SetSessionID(self.currentSessionID)
	self.currentSessionID++
	return nil
}

func (self *SessionManager) AddSession(sess *Session) error {
	self.Sessions[sess.AccountID] = sess
	return nil
}

func (self *SessionManager) GetSession(accountID uint32) (*Session, error) {
	return self.Sessions[accountID], nil
}
