package gux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igoracmelo/gux"
)

func Test(t *testing.T) {
	t.Run("must return not found", func(t *testing.T) {
		mux := gux.New()
		mux.Get("/some/path", func(c *gux.Ctx) {
			c.W.WriteHeader(http.StatusAccepted)
		})
		mux.Get("/some/other/path", func(c *gux.Ctx) {
			c.W.WriteHeader(http.StatusAccepted)
		})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/some/path/but", nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("want: status %d, got: status %d", http.StatusNotFound, rec.Code)
		}

		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/some", nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("want: status %d, got: status %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("should return not allowed", func(t *testing.T) {
		mux := gux.New()
		mux.Post("/user", func(c *gux.Ctx) {
			c.W.WriteHeader(http.StatusAccepted)
		})
		mux.Put("/user", func(c *gux.Ctx) {
			c.W.WriteHeader(http.StatusAccepted)
		})
		mux.Get("/user/:id", func(c *gux.Ctx) {
			c.W.WriteHeader(http.StatusAccepted)
		})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/user", nil))
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("want: status %d, got: status %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("should match equal paths", func(t *testing.T) {
		mux := gux.New()

		ok := false
		mux.Head("/courses", func(c *gux.Ctx) {
			ok = true
		})
		mux.ServeHTTP(nil, httptest.NewRequest("HEAD", "/courses", nil))
		if !ok {
			t.Error("HEAD /courses not reached")
		}

		ok = false
		mux.Delete("/users", func(c *gux.Ctx) {
			ok = true
		})
		mux.ServeHTTP(nil, httptest.NewRequest("DELETE", "/users", nil))
		if !ok {
			t.Error("DELETE /users not reached")
		}
	})

	t.Run("should parse path variables", func(t *testing.T) {
		mux := gux.New()

		ok := false
		mux.Put("/courses/:id", func(c *gux.Ctx) {
			id := c.Vars["id"]
			if id != "10" {
				t.Errorf("id - want: 10, got: %s", id)
				return
			}
			ok = true
		})

		mux.ServeHTTP(nil, httptest.NewRequest("PUT", "/courses/10", nil))
		if !ok {
			t.Error("PUT /courses/:id not reached")
		}

		ok = false
		mux.Patch("/course/:courseID/lesson/:lessonID", func(c *gux.Ctx) {
			courseID := c.Vars["courseID"]
			lessonID := c.Vars["lessonID"]

			if courseID != "1" {
				t.Errorf("courseID - want: 1, got: %s", courseID)
			}
			if lessonID != "2" {
				t.Errorf("lessonID - want: 2, got: %s", lessonID)
			}
			ok = true
		})

		mux.ServeHTTP(nil, httptest.NewRequest("PATCH", "/course/1/lesson/2", nil))
		if !ok {
			t.Error("PATCH /course/:courseID/lesson/:lessonID not reached")
		}
	})

	t.Run("should call GET handler when HEAD request is done", func(t *testing.T) {
		mux := gux.New()

		ok := false
		mux.Get("/head-me", func(c *gux.Ctx) {
			ok = true
		})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("HEAD", "/head-me", nil))
		if !ok {
			t.Error("HEAD /head-me not reached")
		}
	})
}
