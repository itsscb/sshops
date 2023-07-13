package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	remoteComputer = "192.168.128.99"
	commandFormat  = "ssh root@%s -i cert '%s'"
)

func isAuthorized(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	password := os.Getenv("AUTH_PASSWORD")
	if password == "" {
		return false
	}

	providedPassword := strings.TrimPrefix(authHeader, "Bearer ")
	return providedPassword == password
}

func enableHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cmd := fmt.Sprintf(commandFormat, remoteComputer, "uci set firewall.@rule[13].enabled=1 && uci commit firewall && fw3 reload &>/dev/null")
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		log.Println("Failed to enable:", err)
		http.Error(w, "Failed to enable", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Enabled")
}

func disableHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cmd := fmt.Sprintf(commandFormat, remoteComputer, "uci set firewall.@rule[13].enabled=0 && uci commit firewall && fw3 reload &>/dev/null")
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		log.Println("Failed to disable:", err)
		http.Error(w, "Failed to disable", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Disabled")
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cmd := fmt.Sprintf(commandFormat, remoteComputer, "uci get firewall.@rule[13].enabled")
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Println("Failed to get state:", err)
		http.Error(w, "Failed to get state", http.StatusInternalServerError)
		return
	}

	state := string(out)
	state = state[:len(state)-1]

	// w.WriteHeader(http.StatusOK)
	if state == "0" {
		fmt.Fprintln(w, 0) // Rule is enabled, return 0
	} else {
		fmt.Fprintln(w, 1) // Rule is disabled, return 1
	}
}