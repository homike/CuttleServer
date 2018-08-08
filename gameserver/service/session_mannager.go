package service

import "errors"

type SessionManager struct {
	Sessions map[uint32]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[uint32]*Session),
	}
}

func (self *SessionManager) AddSession(sess *Session) error {
	self.Sessions[sess.AccountID] = sess
	return nil
}

func (self *SessionManager) GetSession(accountID uint32) (*Session, error) {
	sess, ok := self.Sessions[accountID]
	if !ok {
		return nil, errors.New("cannot find session")
	}
	return sess, nil
}
