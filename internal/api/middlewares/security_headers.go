package middlewares

import (
	"fmt"
	"net/http"
)

func SecurityHandler(next http.Handler) http.Handler {
	fmt.Println("SECURITY HEADERS MIDDLEWARE STARTED")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off")                                                    // Avoid DNS prefetching in background by browser; Adds more privacy and prevents DNS attacks;
		w.Header().Set("X-Frame-Options", "DENY")                                                          // Prevents displaying of website using iframe in other websites;
		w.Header().Set("X-XSS-Protection", "1;mode=block")                                                 // Enables cross site scripting filters;
		w.Header().Set("X-Content-Type-Options", "nosniff")                                                // This prevents browsers from mime sniffing, ensuring the files are served with proper mime types;
		w.Header().Set("Strict-Transport-Security", "max-age=31536000;includeSubDomains;preload")          // Enforces HTTPS (for the specified duration);
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'") // Controls which resources can be loaded;
		w.Header().Set("Referred-Policy", "no-referrer")                                                   // How much referrer information should be included with requests made from the site;
		w.Header().Set("X-Powered-By", "something private")                                                // Announcing the backend tech used;
		next.ServeHTTP(w, r)
		fmt.Println("SECURITY HEADERS MIDDLEWARE ENDED")

	})
}

/**

******** Basic Middleware Skeleton ********

func securityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
*/
