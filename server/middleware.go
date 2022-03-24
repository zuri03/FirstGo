package server

//Add auth middleware
import (
	"net/http"
	"strings"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		currentTime := time.Now()
		logger.Printf("%s : %s : %s \n", currentTime.Format("2006-01-02 15:04:05"), req.Method, req.URL.String())
		next.ServeHTTP(writer, req)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		if _, ok := values["error"]; ok {
			writer.Write([]byte("Internal Server Error"))
			return
		}

		if _, ok := values["code"]; !ok {
			path := strings.Split(req.URL.Path, "/")
			var url string
			switch path[0] {
			case "Analysis":
				url = generateAuthorizationCodeUrl("http://localhost:8080/Analysis", "user-top-read")
				break
			case "MakePlaylist":
				url = generateAuthorizationCodeUrl("http://localhost:8080/Playlist", "user-top-read",
					"playlist-modify-public", "user-library-read")
				break
			default:
				writer.Write([]byte("Internal Server Error"))
				return
			}
			http.Redirect(writer, req, url, 301)
			return
		}
		next.ServeHTTP(writer, req)
	})
}
