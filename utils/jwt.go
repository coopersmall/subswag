package utils

func MustGetJWTSigningSecret() []byte {
	found, err := GetEnvVar(JWT_SECRET)
	if err != nil {
		panic("PUBLIC_SIGNING_KEY not found")
	}
	return []byte(found)
}
