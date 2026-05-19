package service

func isUnsafeHeader(headerKey string) bool {
	var unsfafeHeaders = map[string]bool{
		"Host":                true,
		"Content-Length":      true,
		"Transfer-Encoding":   true,
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"TE":                  true,
		"Trailer":             true,
		"Upgrade":             true,
		"Cookie":              true,
		"Set-Cookie":          true,
		"Authorization":       true,
		"User-Agent":          true,
	}
	return unsfafeHeaders[headerKey]
}
