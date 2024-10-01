package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
)

var (
	BASE_URL = "http://app1.local:8080"
)

var oauthConfig = &oauth2.Config{
	ClientID:     "f79560a2-4bef-40de-bcb9-b011df232e6f",
	ClientSecret: "secret",
	RedirectURL:  BASE_URL + "/callback",
	Scopes:       []string{"openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://127.0.0.1:4444/oauth2/auth",
		TokenURL: "http://127.0.0.1:4444/oauth2/token",
	},
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/logout", handleLogout)

	log.Println("Server starting on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {

	_, err := r.Cookie("custom-auth-cookie")
	if err == nil {
		w.Header().Set("content-type", "text/html")
		fmt.Fprintf(w, "Already Logged In <a href='/logout'>Log out</a>")
		return
	}

	w.Header().Set("content-type", "text/html")
	w.Header().Set("content-type", "text/html")
	fmt.Fprintf(w, "Welcome! <a href='/login'>Login with Hydra</a>")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("custom-auth-cookie")
	if err == nil {
		w.Header().Set("content-type", "text/html")
		fmt.Fprintf(w, "Already Logged In <a href='/logout'>Log out</a>")
		return
	}

	url := oauthConfig.AuthCodeURL("state123445", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	// Extract the ID Token (included in token.Extra)
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token found in token response", http.StatusInternalServerError)
		return
	}

	// Set the custom authentication cookie
	cookie := http.Cookie{
		Name:     "custom-auth-cookie",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
	}

	// Set a cookie for the id_token
	idTokenCookie := http.Cookie{
		Name:     "id-token-cookie",
		Value:    idToken,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.SetCookie(w, &idTokenCookie)

	w.Header().Set("content-type", "text/html")
	fmt.Fprintf(w, "Login successful! You can now access protected resources.")
	fmt.Fprintf(w, "</br> Already Logged In <a href='/logout'>Log out</a>")
	return
}

func handleLogout(w http.ResponseWriter, r *http.Request) {

	// Retrieve the id_token from the cookie
	idTokenCookie, err := r.Cookie("id-token-cookie")
	if err != nil {
		http.Error(w, "ID token not found, user might not be logged in", http.StatusBadRequest)
		return
	}
	idToken := idTokenCookie.Value

	// Hydra logout endpoint
	logoutURL := "http://localhost:4444/oauth2/sessions/logout"

	// Clear the custom authentication cookie
	cookie := http.Cookie{
		Name:   "custom-auth-cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Delete the cookie
	}
	http.SetCookie(w, &cookie)

	// Clear the id_token cookie
	idTokenClear := http.Cookie{
		Name:   "id-token-cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Delete the cookie
	}
	http.SetCookie(w, &idTokenClear)

	// Build the logout request URL with id_token_hint and post-logout redirect URI
	postLogoutRedirect := BASE_URL + "/"
	logoutRequestURL := fmt.Sprintf("%s?id_token_hint=%s&post_logout_redirect_uri=%s",
		logoutURL, url.QueryEscape(idToken), url.QueryEscape(postLogoutRedirect))

	// Redirect the user to Hydra's logout endpoint
	http.Redirect(w, r, logoutRequestURL, http.StatusTemporaryRedirect)

}
