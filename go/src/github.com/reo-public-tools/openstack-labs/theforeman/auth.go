package theforeman

import (
    "os"
    "fmt"
    "log"
    "bufio"
//    "errors"
    "regexp"
    "unsafe"
    "syscall"
    "strings"
    "net/http"
    "io/ioutil"
    "crypto/tls"
)

// Allow for disabling of terminal echo while asking for password
func terminalEcho(show bool) (error) {
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
        return err
    }

    return nil
}

// Prompt user for credentials
func promptUserForCredentials() (string, string, error) {

    // Request username and password from the user
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("\n\nTheForeman username: ")
    username, _ := reader.ReadString('\n')
    username = strings.Replace(username, "\n", "", -1)

    fmt.Print("TheForeman password: ")
    err := terminalEcho(false)
    if err != nil {
        return "", "", err
    }
    password, _ := reader.ReadString('\n')
    password = strings.Replace(password, "\n", "", -1)
    err = terminalEcho(true)
    if err != nil {
        return "", "", err
    }
    fmt.Printf("\n\n")

    return username, password, nil
}

// Save theforeman session to ~/.theforeman-session for future runs
func saveTheForemanSession(session string) error {

    // Get the home directory and set the full path for the file
    var home string = os.Getenv("HOME")
    var sessionFile = fmt.Sprintf("%s/.theforeman-session", home)

    // Make sure the file exists
    err := ioutil.WriteFile(sessionFile, []byte(session), 0600)
    if err != nil {
        return err
    }

    return nil
}

// Check for and test an existing session
func checkExistingTheForemanSession() (string, error) {

    // Get the home directory and set the full path for the file
    var home string = os.Getenv("HOME")
    var sessionFile = fmt.Sprintf("%s/.theforeman-session", home)

    // See if token file exists
    _, err := os.Stat(sessionFile)
    if os.IsNotExist(err) {
        return "notfound", nil
    }

    // Pull the token from the file if we got this far
    sessionbyte, err := ioutil.ReadFile(sessionFile)
    if err != nil {
        return "notfound", err
    }
    return string(sessionbyte), nil
}

// Check the if the session is valid
func isSessionValid(url string, session string) (bool, error) {

    // Rackspace Identity URL for token validation
    var testurl string = fmt.Sprintf("%s/api/v2/status", url)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", testurl, nil)
    if err != nil {
        return false, err
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
        return false, err
    }
    defer resp.Body.Close()

    // Print out the results
    return resp.Status == "200 OK", nil

}

func TheForemanLogin(url string) (string, error) {

    // Check for and pull the existing token from the token file
    session, err := checkExistingTheForemanSession()
    if err != nil {
        return "invalid", err
    } else {
        isvalid, err := isSessionValid(url, session)
        if err != nil {
            return "invalid", err
        } else {
            if isvalid {
                return session, nil
            }
        }
    }


    // Prompt the user for some credentials 
    username, password, err := promptUserForCredentials()
    if err != nil {
        return "invalid", err
    }

    // Set the login url
    var testurl string = fmt.Sprintf("%s/api/v2/status", url)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", testurl, nil)
    if err != nil {
        return "invalid", err
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
        return "invalid", err
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
    err = saveTheForemanSession(session)
    if err != nil {
        return "invalid", err
    }

    return session, nil
}
