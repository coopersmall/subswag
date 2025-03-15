package domain

import (
	"fmt"
	"strings"

	"github.com/coopersmall/subswag/domain/user"
)

const NamespaceKeyPrefix = "subswag"
const UserIDKeyPrefix = "users/%d"

func NamespaceKey(key string) string {
	return createPath(NamespaceKeyPrefix, key)
}

func UserIDKey(userId user.UserID, paths ...string) string {
	path := createPath(NamespaceKeyPrefix, fmt.Sprintf(UserIDKeyPrefix, userId))
	for _, p := range paths {
		path = createPath(path, p)
	}
	return path
}

func createPath(strs ...string) string {
	return strings.Join(strs, "/")
}
