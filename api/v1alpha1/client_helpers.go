package v1alpha1

import "strings"

func (c *Client) InternalSubject() string {
	return strings.Join([]string{"client", c.Name}, ":")
}

func (c *Client) Usernames(prefix string) []string {
	usernames := []string{
		prefix + strings.Join([]string{"client", c.Name}, ":"),                             // New portable format
		prefix + strings.Join([]string{"client", c.Namespace, c.Name, string(c.UID)}, ":"), // Legacy format
	}

	if c.Spec.Username != nil {
		usernames = append(usernames, *c.Spec.Username)
	}

	return usernames
}
