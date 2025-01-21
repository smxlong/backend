package backend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func Test_that_NewRouter_fails_if_gin_Engine_is_not_registered(t *testing.T) {
	di := NewInjector()
	_, err := NewRouter(di)
	require.NotNil(t, err)
}

func Test_that_Router_can_dispatch_a_GET_request(t *testing.T) {
	g := gin.New()
	di := NewInjector()
	di.AddSingleton(g)
	r, err := NewRouter(di)
	require.Nil(t, err)
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test")
	})

	server := httptest.NewServer(g)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/test")
	require.Nil(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// intYielder is a dependency that yields a constant integer value when its
// Get method is called.
type intYielder struct {
	value int
}

// Get returns the constant integer value.
func (i *intYielder) Get() int {
	return i.value
}

func Test_that_Router_can_dispatch_requests_with_dependencies(t *testing.T) {
	g := gin.New()
	di := NewInjector()
	di.AddSingleton(g)
	di.AddSingleton(&intYielder{42})
	r, err := NewRouter(di)
	require.Nil(t, err)

	r.GET("/get", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})
	r.POST("/post", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})
	r.PUT("/put", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})
	r.DELETE("/delete", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})
	r.PATCH("/patch", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})
	r.OPTIONS("/options", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})
	r.HEAD("/head", func(c *gin.Context, i *intYielder) {
		c.String(200, "%d", i.Get())
	})

	server := httptest.NewServer(g)
	defer server.Close()

	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"} {
		methodLower := strings.ToLower(method)
		req, _ := http.NewRequest(method, server.URL+"/"+methodLower, nil)
		resp, err := server.Client().Do(req)
		require.Nil(t, err)
		require.Equal(t, 200, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		if method == "HEAD" {
			require.Equal(t, "", string(body))
		} else {
			require.Equal(t, "42", string(body))
		}
	}
}

func Test_that_Router_implements_ServeHTTP(t *testing.T) {
	g := gin.New()
	di := NewInjector()
	di.AddSingleton(g)
	r, err := NewRouter(di)
	require.Nil(t, err)

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test")
	})

	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/test")
	require.Nil(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
