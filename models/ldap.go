package models

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/astaxie/beego"

	"gopkg.in/ldap.v2"
)

var lcfg = beego.BConfig.LDAPConfig

type LDAP struct {
	Conn *ldap.Conn
}

type LDAP_RESULT struct {
	DN         string              `json:"dn"`
	Attributes map[string][]string `json:"attributes"`
}

func (lc *LDAP) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

func (lc *LDAP) Connect() (err error) {
	if lcfg.TLS {
		lc.Conn, err = ldap.DialTLS("tcp", lcfg.Addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		lc.Conn, err = ldap.Dial("tcp", lcfg.Addr)
	}

	if err != nil {
		return
	}

	if !lcfg.TLS && lcfg.Starttls {
		err = lc.Conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			lc.Conn.Close()
			return err
		}
	}

	err = lc.Conn.Bind(lcfg.Binddn, lcfg.Bindpass)
	if err != nil {
		lc.Conn.Close()
		return err
	}

	return nil
}

func (lc *LDAP) Auth(name string, password string) (err error) {
	searchRequset := ldap.NewSearchRequest(
		lcfg.Basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", name),
		[]string{"cn", "uid"},
		nil,
	)

	sr, err := lc.Conn.Search(searchRequset)
	if err != nil {
		return
	}

	if len(sr.Entries) == 0 {
		err = errors.New("cannot find such user")
		return
	}

	err = lc.Conn.Bind(sr.Entries[0].DN, password)
	if err != nil {
		return
	}
	return
}

func LDAP_Auth(name string, password string) (err error) {
	lc := &LDAP{}
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		return
	}

	err = lc.Auth(name, password)
	return
}
