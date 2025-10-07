package auth

import (
	"crypto/tls"
	"fmt"

	ldap "github.com/go-ldap/ldap/v3"
)

// AuthenticateAD authenticates a user against Active Directory via LDAP
func AuthenticateAD(ldapURL, baseDN, username, password string) (bool, error) {
	conn, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return false, fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer conn.Close()

	userDN := fmt.Sprintf("cn=%s,%s", username, baseDN)
	if err := conn.Bind(userDN, password); err != nil {
		return false, nil // Invalid credentials
	}
	return true, nil // Authenticated
}
