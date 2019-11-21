package theforeman

import (
     "os"
     "fmt"
     "log"
     "bufio"
     "unsafe"
     "regexp"
     "syscall"
     "strings"
     "net/http"
     "io/ioutil"
     "crypto/tls"
)

// Allow for disabling of terminal echo while asking for password
func terminalEcho(show bool) {
    var termios = &syscall.Termios{}
    var fd = os.Stdout.Fd()

    _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)))
    if err != 0 {
        log.Fatal(err)
    }

    if show {
        termios.Lflag |= syscall.ECHO
    } else {
        termios.Lflag &^= syscall.ECHO
    }

    _, _, err = syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios)))
    if err != 0 {
        log.Fatal(err)
    }
}

// Prompt user for credentials
func promptUserForCredentials() (string, string) {

    // Request username and password from the user
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("\n\nRackspace SSO username: ")
    username, _ := reader.ReadString('\n')
    username = strings.Replace(username, "\n", "", -1)

    fmt.Print("Rackspace SSO password: ")
    terminalEcho(false)
    password, _ := reader.ReadString('\n')
    password = strings.Replace(password, "\n", "", -1)
    terminalEcho(true)
    fmt.Printf("\n\n")

    return username, password
}

// Save theforeman session to ~/.theforeman-session for future runs
func saveTheForemanSession(session string) {

    // Get the home directory and set the full path for the file
    var home string = os.Getenv("HOME")
    var sessionFile = fmt.Sprintf("%s/.theforeman-session", home)

    // Make sure the file exists
    err := ioutil.WriteFile(sessionFile, []byte(session), 0600)
    if err != nil {
        log.Fatal(err)
    }

}

// Check for and test an existing session
func checkExistingTheForemanSession() (bool, string) {

    // Get the home directory and set the full path for the file
    var home string = os.Getenv("HOME")
    var sessionFile = fmt.Sprintf("%s/.theforeman-session", home)

    // See if token file exists
    _, err := os.Stat(sessionFile)
    if os.IsNotExist(err) {
        return false, "notfound"
    }

    // Pull the token from the file if we got this far
    sessionbyte, err := ioutil.ReadFile(sessionFile)
    if err != nil {
        log.Fatal(err)
    }
    return true, string(sessionbyte)
}

// Check the if the session is valid
func isSessionValid(url string, session string) bool {

    // Rackspace Identity URL for token validation
    var testurl string = fmt.Sprintf("%s/api/v2/status", url)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", testurl, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    // Set the session to test 
    req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s",session))

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    // Print out the results
    return resp.Status == "200 OK"

}

func TheForemanLogin(url string) (string) {

    // Check for and pull the existing token from the token file
    sessionfound, session := checkExistingTheForemanSession()

    if sessionfound {
        if isSessionValid(url, session) {
            return session
        }
    }

    // Prompt the user for some credentials 
    username, password := promptUserForCredentials()

    // Set the login url
    var testurl string = fmt.Sprintf("%s/api/v2/status", url)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", testurl, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Set up basic auth
    req.SetBasicAuth(username, password)

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json,version=2")

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
            log.Fatal(err)
    }
    defer resp.Body.Close()

    // Pull out the Set-Cookie header
    cookieheader := resp.Header.Get("Set-Cookie")

    // Use regex to pull out just the session id
    re := regexp.MustCompile("_session_id=([^;]*)")
    match := re.FindStringSubmatch(cookieheader)
    if len(match) != 0 { session = match[1] }
    fmt.Println("Cookie: ", session)

    // Print out the body
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("resp Body: ", string(body))

    // Save the session after pulling a new one
    saveTheForemanSession(session)

    return session
}
