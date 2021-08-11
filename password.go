package ldaputil

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func UserChangePass(cfg Config, username, oldPassword, password string) error {
	conn, err := ldap.Dial("tcp", cfg.Server)
	if err != nil {
		log.Printf("fail to connect ldap server: %s.", err)
		return err
	}
	defer conn.Close()
	err = conn.Bind(fmt.Sprintf("cn=%s,%s", username, cfg.BaseDN), oldPassword)
	if err != nil {
		log.Printf("fail to bind credentials: %s.", err)
		return err
	}

	_, err = conn.PasswordModify(ldap.NewPasswordModifyRequest("", oldPassword, password))
	if err != nil {
		log.Printf("fail to modify password: %s.", err)
		return err
	}

	return nil
}

func UserSetPass(cfg Config, username, password string) error {
	log.Printf("set pass for user \"%s\" with password \"%s\".", username, password)
	if password == "" {
		b := make([]byte, 128)
		rand.Read(b)
		hash := md5.New()
		hash.Write(b)
		password = hex.EncodeToString(hash.Sum(nil))
		log.Printf("empty password, replaced with \"%s\".", password)
	}

	conn, err := ldap.Dial("tcp", cfg.Server)
	if err != nil {
		log.Printf("fail to connect ldap server: %s.", err)
		return err
	}
	defer conn.Close()
	err = conn.Bind(cfg.Bind.DN, cfg.Bind.Secret)
	if err != nil {
		log.Printf("fail to bind credentials: %s.", err)
		return err
	}

	_, err = conn.PasswordModify(ldap.NewPasswordModifyRequest(fmt.Sprintf("cn=%s,%s", username, cfg.BaseDN), "", password))
	if err != nil {
		log.Printf("fail to set password: %s.", err)
		return err
	}

	return nil
}
