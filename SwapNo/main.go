package main

import (
    "os"
    "fmt"
    "time"
    "flag"
    "strings"
    "os/exec"
    "net/http"
    "io/ioutil"
    "math/rand"
    "path/filepath"
    urlstr "net/url"
)


var FileLogPath string
var check_time = 1
var (
    url = flag.String("url", "", "The Register Url of SWAPIDC which you want to add users to~~~")
    FileLog = flag.Bool("log", false, "Log the outputs")
    CheckRates = flag.Int("rate", 1, "The rate of the Import Process")
)

func main() {
    flag.Parse()
    FileLogPath = getCurrentPath() + "/log.txt"
    serviceLogger("Now Loading...", 0)
    if(*FileLog){
        serviceLogger(fmt.Sprintf("Saving Logs to %s", FileLogPath), 0)
    }
    if(*url == ""){
        serviceLogger("Please input URL!", 31)
        os.Exit(0)
    }else{
        serviceLogger(fmt.Sprintf("Target: %s", *url), 32)
        startProcess(*url)
    }
}

func startProcess(url string){
    var CheckRate = int(*CheckRates)
    ch := make(chan string, 1)
    serviceLogger(fmt.Sprintf("Loaded ImportRate : %d Second", int(CheckRate)), 32)
    for {
        serviceLogger(fmt.Sprintf("Start Importing, Round %d", int(check_time)), 0)   
        go func() {
            ImportUser(url, int(check_time))
            ch <- "done"
        }()
        select {
        case <-ch:
            serviceLogger(fmt.Sprintf("Task(%d) Is Done!!!!!", int(check_time)), 32)
        case <-time.After(time.Duration(CheckRate - 1) * time.Second):
            serviceLogger(fmt.Sprintf("Task(%d) Is Timeout!!!!!", int(check_time)), 31)
        }
        check_time = check_time + 1
        timeSleep(CheckRate)
    }
}

func ImportUser(url string, round int) error{
    resstr := createRandomUser()
    resp, err := http.PostForm(url, resstr)
    if err != nil {
        serviceLogger(fmt.Sprintf("Round %d, Error: %s", round, err), 31)
        return nil
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        serviceLogger(fmt.Sprintf("Round %d, Error: %s", round, err), 31)
        return nil
    }
    body = body
    serviceLogger(fmt.Sprintf("Done~"), 32)
    return nil
}

func createRandomUser() urlstr.Values{
    username := getRandomString(10)
    password := getRandomString(10)
    email := getRandomString(10) + "@" + getRandomString(3) + ".com"
    name := getRandomString(5)
    address := getRandomString(10)
    zip := getRandomStringInt(10)
    phone := getRandomStringInt(11)
    uuu := urlstr.Values{
        "username": {username},
        "password" : {password},
        "repassword" : {password},
        "email" : {email},
        "name" : {name},
        "country" : {"China"},
        "address" : {address},
        "zip" : {zip},
        "phone" : {phone},
    }
    serviceLogger(fmt.Sprintf("Random Username: %s, Email: %s, Password %s, Phone: %s", username, email, password, phone), 33)
    return uuu
}

func getRandomString(length int) string{
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    result := []byte{}
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

func getRandomStringInt(length int) string{
    str := "0123456789"
    bytes := []byte(str)
    result := []byte{}
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

func timeSleep(second int){
    time.Sleep(time.Duration(second) * time.Second)
}

func serviceLogger(log string, color int){
    log = strings.Replace(log, "\n", "", -1)
    if(color == 0){
        fmt.Printf("%s\n", log)
    }else{
        fmt.Printf("%c[1;0;%dm%s%c[0m\n", 0x1B, color, log, 0x1B)
    }
    if(*FileLog){
        fd, err := os.OpenFile(FileLogPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)  
        if(err != nil){
            fmt.Printf("%c[1;0;%dm[%s], %s%c[0m\n", 0x1B, 31, FileLogPath, err, 0x1B)
        }else{
            fd_time := time.Now().Format("2006/01/02-15:04:05");  
            fd_content := strings.Join([]string{fd_time, "  ", log, "\n"}, "")  
            buf := []byte(fd_content)  
            fd.Write(buf)  
            fd.Close()
        }
    }
}

func getCurrentPath() string {  
    file, _ := exec.LookPath(os.Args[0])  
    path, _ := filepath.Abs(file)  
    path = substr(path, 0, strings.LastIndex(path, "/"))
    return path  
}  

func substr(s string, pos, length int) string {
    runes := []rune(s)
    l := pos + length
    if l > len(runes) {
        l = len(runes)
    }
    return string(runes[pos:l])
}