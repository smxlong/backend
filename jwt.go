package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// JWT is a middleware that verifies the JWT in the Authorization header.
func JWT(issuer, audience, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.Next()
			return
		}
		token, err := jwt.Parse([]byte(bearer),
			jwt.WithIssuer(issuer),
			jwt.WithAudience(audience),
			jwt.WithKey(jwa.HS256(), []byte(secret)),
			jwt.WithValidate(true),
			jwt.WithVerify(true),
		)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("token", token)
		c.Next()
	}
}

// RequirePermissionsClaim is a middleware that requires the token to have the given
// permissions in the claim.
func RequirePermissionsClaim(claim string, assertion PermissionsAssertion) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, _ := c.Get("token")
		token := t.(jwt.Token)
		var permissions []string
		if err := token.Get(claim, &permissions); err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		tokenPermissions := map[string]bool{}
		for _, p := range permissions {
			tokenPermissions[p] = true
		}
		if !assertion(tokenPermissions) {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequirePermissions is a middleware that requires the token to have the given
// permissions in the "permissions" claim.
func RequirePermissions(assertion PermissionsAssertion) gin.HandlerFunc {
	return RequirePermissionsClaim("permissions", assertion)
}

// Token returns the token from the gin.Context.
func Token(c *gin.Context) jwt.Token {
	t, _ := c.Get("token")
	return t.(jwt.Token)
}

// PermissionsAssertion is an assertion about permissions
type PermissionsAssertion func(map[string]bool) bool

// HasAny returns an assertion that requires the token to have any of the given permissions.
func HasAny(permissions ...string) PermissionsAssertion {
	return func(token map[string]bool) bool {
		for _, p := range permissions {
			if token[p] {
				return true
			}
		}
		return false
	}
}

// HasAll returns an assertion that requires the token to have all of the given permissions.
func HasAll(permissions ...string) PermissionsAssertion {
	return func(token map[string]bool) bool {
		for _, p := range permissions {
			if !token[p] {
				return false
			}
		}
		return true
	}
}

// Or returns an assertion that requires any of the given assertions to be true.
func Or(assertions ...PermissionsAssertion) PermissionsAssertion {
	return func(token map[string]bool) bool {
		for _, a := range assertions {
			if a(token) {
				return true
			}
		}
		return false
	}
}

// And returns an assertion that requires all of the given assertions to be true.
func And(assertions ...PermissionsAssertion) PermissionsAssertion {
	return func(token map[string]bool) bool {
		for _, a := range assertions {
			if !a(token) {
				return false
			}
		}
		return true
	}
}
