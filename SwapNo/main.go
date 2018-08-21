package main

import (
    "io"
    "os"
    "fmt"
    "time"
    "flag"
    "bufio"
    "strconv"
    "strings"
    "os/exec"
    "net/http"
    "io/ioutil"
    "math/rand"
    "path/filepath"
    urlstr "net/url"
    "net/http/cookiejar"
    "github.com/json-iterator/go"
)

var Success = 0
var FileLogPath string
var ProxyPath string
var UsersPath string
var check_time = 0
var TProcesss = 0
var ProxyList = make(map[int]map[string]string)
var UserList = make(map[int]map[string]string)

var (
    url = flag.String("url", "", "The Register Url of SWAPIDC which you want to add users to~~~")
    TicketUrl = flag.String("ticketurl", "", "The Ticket Url of SWAPIDC which you want to add tickets to~~~")
    FileLog = flag.Bool("log", false, "Log the outputs")
    FileLogLimit = flag.Int("loglimit", 10240, "Log Limit")
    Proxy = flag.Bool("proxy", false, "Enable proxy mode")
    SaveUsers = flag.Bool("saveusers", false, "Save the registered Users(Cann't know success or not)")
    ProxyUpdate = flag.Bool("proxyupdate", false, "Update the Proxy list")
    CheckRates = flag.Int("rate", 1, "The rate of the Import Process")
    ShowResult = flag.Bool("debug", false, "Show Results")
    OverClock = flag.Bool("overclock", false, "Run faster")
    TicketMode = flag.Bool("tickets", false, "I love Tickets")
    TicketProcess = flag.Int("ticketprocess", 10, "Tickets Process")
)

func main() {
    flag.Parse()
    FileLogPath = getCurrentPath() + "/log.txt"
    ProxyPath = getCurrentPath() + "/proxy.txt"
    UsersPath = getCurrentPath() + "/users.txt"
    serviceLogger("Now Loading...", 0)
    if(*FileLog){
        serviceLogger(fmt.Sprintf("Saving Logs to %s", FileLogPath), 0)
    }
    if(*ProxyUpdate){
        updateProxy()
        os.Exit(0)
    }
    if(*TicketMode){
        serviceLogger(fmt.Sprintf("Running in Ticket Mode, Process: %d", *TicketProcess), 0)
        loadUsers()
    }
    if(*Proxy){
        serviceLogger(fmt.Sprintf("Enabled Proxy Mode"), 0)
        loadProxy()
    }
    if(*url == ""){
        serviceLogger("Please input URL!", 31)
        os.Exit(0)
    }else{
        serviceLogger(fmt.Sprintf("Target: %s", *url), 32)
        startProcess(*url)
    }
}

func updateProxy(){
    serviceLogger("Updating Proxy List", 33)
    client := http.Client{}
    resp, err := client.Get("https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list")
    if err != nil {
        serviceLogger("Update Failed", 31)
        os.Exit(0)
    }
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        serviceLogger("Update Failed, IO Error", 31)
        os.Exit(0)
    }
    err = os.Remove(ProxyPath)
    if err != nil {
        serviceLogger(fmt.Sprintf("Old Proxy File Removed Failed, Error: %s", err), 31)
    } else {
        serviceLogger("Old File Removed!", 32)
    }
    fd, err := os.OpenFile(ProxyPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)  
    if(err != nil){
        serviceLogger(fmt.Sprintf("[%s]Load Error: %s", ProxyPath, err), 31)
    }else{
        buf := []byte(string(data))  
        fd.Write(buf)  
        fd.Close()
    }
    serviceLogger(fmt.Sprintf("Proxy Updated"), 32)
}

func startProcess(url string){
    var CheckRate = int(*CheckRates)
    if(*OverClock){
        serviceLogger(fmt.Sprintf("Loaded ImportRate : %d Millisecond", int(CheckRate)), 32)
    }else{
        serviceLogger(fmt.Sprintf("Loaded ImportRate : %d Second", int(CheckRate)), 32)
    }
    if(*TicketMode){
        for {
            serviceLogger(fmt.Sprintf("[%d]Start Opening!", int(check_time) + 1), 0)   
            if(TProcesss <= int(*TicketProcess)){
                go func() {
                    ImportTicket(url, int(check_time))
                    TProcesss = TProcesss - 1
                }()
                check_time = check_time + 1
            }
            timeSleep(CheckRate)
        }
    }else{
        for {
            serviceLogger(fmt.Sprintf("[%d]Start Importing!", int(check_time) + 1), 0)   
            go func() {
                ImportUser(url, int(check_time))
            }()
            check_time = check_time + 1
            timeSleep(CheckRate)
        }
    }
}

func ImportTicket(url string, round int) error{
    TProcesss = TProcesss + 1
    var client http.Client
    var hostV map[string]string
    if(*Proxy){
        hostV = getRandomProxy()
        if(hostV["status"] != "Success"){
            serviceLogger(fmt.Sprintf("[%d]Error: Get Proxy Failed", round), 31)
        }else{
            serviceLogger(fmt.Sprintf("[%d]Using Proxy: %s", round, hostV["host"]), 31)
            urlc := urlstr.URL{}
            urlproxy, _ := urlc.Parse(hostV["host"])
            client = http.Client{
                Transport: &http.Transport{
                    Proxy: http.ProxyURL(urlproxy),
                },
            }
        }
    }else{
        client = http.Client{}
    }
    jar, err := cookiejar.New(nil)
    if err != nil {
        serviceLogger(fmt.Sprintf("[%d]Get CookieJar Failed. Error: %s", round, err), 31)
        return nil
    }
    client.Jar = jar
    userV := getRandomUser()
    uuu := urlstr.Values{
        "swapname" : {userV["username"]},
        "swappass" : {userV["password"]},
    }
    resp, err := client.PostForm(url, uuu)
    if err != nil {
        serviceLogger(fmt.Sprintf("[%d]Error: %s", round, err), 31)
        if(*Proxy){
            vint, err := strconv.Atoi(hostV["id"]) 
            if(err != nil){
                serviceLogger(fmt.Sprintf("[%d]Int Error: %s", round, err), 31)
            }else{
                serviceLogger(fmt.Sprintf("[%d]Removed Proxy(%s): %s", round, hostV["id"], hostV["host"]), 31)
                ProxyList[vint]["available"] = "false"
            }
        }
        return nil
    }
    defer resp.Body.Close()
    if(*ShowResult){
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            serviceLogger(fmt.Sprintf("Round %d, Error: %s", round, err), 31)
            return nil
        }
        serviceLogger(fmt.Sprintf("[%d]Result: %s", round, string(body)), 32)
    }
    for {
        serviceLogger(fmt.Sprintf("[%d]Start Ticketing!", int(check_time) + 1), 0)  
        rttket := getRandomTicket() 
        uuu := urlstr.Values{
            "name": {rttket["name"]},
            "email" : {rttket["email"]},
            "subject" : {rttket["subject"]},
            "message" : {rttket["message"]},
        }
        resp1, err := client.PostForm(*TicketUrl, uuu)
        if err != nil {
            serviceLogger(fmt.Sprintf("[%d]Error: %s", round, err), 31)
            if(*Proxy){
                vint, err := strconv.Atoi(hostV["id"]) 
                if(err != nil){
                    serviceLogger(fmt.Sprintf("[%d]Int Error: %s", round, err), 31)
                }else{
                    serviceLogger(fmt.Sprintf("[%d]Removed Proxy(%s): %s", round, hostV["id"], hostV["host"]), 31)
                    ProxyList[vint]["available"] = "false"
                }
            }
            return nil
        }
        defer resp1.Body.Close()
        serviceLogger(fmt.Sprintf("[%d]Done~ Random Ticket Title: %s %s, Email: %s", round, rttket["name"], rttket["lastname"], rttket["email"]), 32)
        if(*ShowResult){
            body, err := ioutil.ReadAll(resp1.Body)
            if err != nil {
                serviceLogger(fmt.Sprintf("Round %d, Error: %s", round, err), 31)
                return nil
            }
            serviceLogger(fmt.Sprintf("[%d]Result: %s", round, string(body)), 32)
        }
        addSuccess()
        timeSleep(int(*CheckRates))
    }
    return nil
}

func loadUsers(){
    fi, err := os.Open(UsersPath)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    defer fi.Close()
    br := bufio.NewReader(fi)
    for {
        a, _, c := br.ReadLine()
        if c == io.EOF {
            break
        }
        initSingleUser(string(a))
    }
    serviceLogger(fmt.Sprintf("Loaded %d Users", len(UserList)), 32)
}

func initSingleUser(str string){
    strb := []byte(str)
    username := jsoniter.Get(strb, "username").ToString()
    password := jsoniter.Get(strb, "password").ToString()
    usermap := make(map[string]string)
    usermap["username"] = username
    usermap["password"] = password
    UserList[len(UserList) + 1] = usermap
}

func getRandomUser() map[string]string{
    returnV := make(map[string]string)
    for i := 1; i <= len(UserList); i++ {
        vid := rand.Intn(len(UserList))
        sproxy := UserList[vid]
        returnV["status"] = "Success"
        returnV["username"] = sproxy["username"]
        returnV["password"] = sproxy["password"]
        returnV["id"] = strconv.Itoa(vid)
        return returnV
    }
    returnV["status"] = "Error"
    return returnV
}

func getRandomTicket()map[string]string{
    strR := make(map[string]string)
    strR["name"] = getRandomString(5)
    strR["email"] = getRandomString(10) + "@" + getRandomString(3) + ".com"
    strR["subject"] = getRandomString(5)
    strR["message"] = getRandomString(10)
    return strR
}

func ImportUser(url string, round int) error{
    resstrR := createRandomUser(round)
    uuu := urlstr.Values{
        "username": {resstrR["username"]},
        "password" : {resstrR["password"]},
        "repassword" : {resstrR["password"]},
        "email" : {resstrR["email"]},
        "name" : {resstrR["name"]},
        "country" : {"China"},
        "address" : {resstrR["address"]},
        "zip" : {resstrR["zip"]},
        "phone" : {resstrR["phone"]},
    }
    var client http.Client
    var hostV map[string]string
    if(*Proxy){
        hostV = getRandomProxy()
        if(hostV["status"] != "Success"){
            serviceLogger(fmt.Sprintf("[%d]Error: Get Proxy Failed", round), 31)
        }else{
            serviceLogger(fmt.Sprintf("[%d]Using Proxy: %s", round, hostV["host"]), 31)
            urlc := urlstr.URL{}
            urlproxy, _ := urlc.Parse(hostV["host"])
            client = http.Client{
                Transport: &http.Transport{
                    Proxy: http.ProxyURL(urlproxy),
                },
            }
        }
    }else{
        client = http.Client{}
    }
    resp, err := client.PostForm(url, uuu)
    if err != nil {
        serviceLogger(fmt.Sprintf("[%d]Error: %s", round, err), 31)
        if(*Proxy){
            vint, err := strconv.Atoi(hostV["id"]) 
            if(err != nil){
                serviceLogger(fmt.Sprintf("[%d]Int Error: %s", round, err), 31)
            }else{
                serviceLogger(fmt.Sprintf("[%d]Removed Proxy(%s): %s", round, hostV["id"], hostV["host"]), 31)
                ProxyList[vint]["available"] = "false"
            }
        }
        return nil
    }
    defer resp.Body.Close()
    serviceLogger(fmt.Sprintf("[%d]Done~ Random Username: %s, Email: %s, Password %s, Phone: %s", round, resstrR["username"], resstrR["email"], resstrR["password"], resstrR["phone"]), 32)
    if(*SaveUsers){
        fd, err := os.OpenFile(UsersPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)  
        if(err != nil){
            fmt.Printf("%c[1;0;%dm[%s], %s%c[0m\n", 0x1B, 31, UsersPath, err, 0x1B)
        }else{
            var userstr = `{"username": "` + resstrR["username"] + `", "password": "` + resstrR["password"] + `"}`
            fd_content := strings.Join([]string{userstr, "\n"}, "")  
            buf := []byte(fd_content)  
            fd.Write(buf)  
            fd.Close()
        }
    }
    if(*ShowResult){
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            serviceLogger(fmt.Sprintf("Round %d, Error: %s", round, err), 31)
            return nil
        }
        serviceLogger(fmt.Sprintf("[%d]Result: %s", round, string(body)), 32)
    }
    addSuccess()
    return nil
}

func addSuccess(){
    Success = Success + 1
    serviceLogger(fmt.Sprintf("[Count]Success: %d", int(Success)), 32)
}

func loadProxy(){
    fi, err := os.Open(ProxyPath)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    defer fi.Close()
    br := bufio.NewReader(fi)
    for {
        a, _, c := br.ReadLine()
        if c == io.EOF {
            break
        }
        initSingleProxy(string(a))
    }
    serviceLogger(fmt.Sprintf("Loaded %d Proxys", len(ProxyList)), 32)
}

func initSingleProxy(str string){
    strb := []byte(str)
    anonymity := jsoniter.Get(strb, "anonymity").ToString()
    host := jsoniter.Get(strb, "host").ToString()
    port := jsoniter.Get(strb, "port").ToString()
    from := jsoniter.Get(strb, "from").ToString()
    vtype := jsoniter.Get(strb, "type").ToString()
    response_time := jsoniter.Get(strb, "response_time").ToString()
    proxymap := make(map[string]string)
    proxymap["host"] = vtype + "://" + host + ":" + port
    proxymap["available"] = "true"
    ProxyList[len(ProxyList) + 1] = proxymap
    //serviceLogger(fmt.Sprintf("Loaded Proxy: %s:%s(%s), Anonymity: %s, From: %s", host, port, response_time, anonymity, from), 32)
    fmt.Sprintf("Loaded Proxy: %s:%s(%s), Anonymity: %s, From: %s", host, port, response_time, anonymity, from)
}

func getRandomProxy() map[string]string{
    returnV := make(map[string]string)
    for i := 1; i <= len(ProxyList); i++ {
        vid := rand.Intn(len(ProxyList))
        sproxy := ProxyList[vid]
        if(sproxy["available"] == "true"){
            returnV["status"] = "Success"
            returnV["host"] = sproxy["host"]
            returnV["id"] = strconv.Itoa(vid)
            return returnV
        }
    }
    returnV["status"] = "Error"
    return returnV
}

func createRandomUser(round int) map[string]string{
    strR := make(map[string]string)
    strR["username"] = getRandomString(10)
    strR["password"] = getRandomString(10)
    strR["email"] = getRandomString(10) + "@" + getRandomString(3) + ".com"
    strR["name"] = getRandomString(5)
    strR["address"] = getRandomString(10)
    strR["zip"] = getRandomStringInt(10)
    strR["phone"] = "156" + getRandomStringInt(8)
    return strR
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
    if(*OverClock){
        time.Sleep(time.Duration(second) * time.Millisecond)
    }else{
        time.Sleep(time.Duration(second) * time.Second)
    }
}

func serviceLogger(log string, color int){
    checkLogOverSized()
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

func checkLogOverSized(){
    logInfo := getSize(FileLogPath)
    if((int(logInfo) / 1024) >= int(*FileLogLimit) && int(*FileLogLimit) > 0){
        fmt.Printf("%c[1;0;%dm[Error]%s%c[0m\n", 0x1B, 31, "[Log]Log Oversized, Cleaning!", 0x1B)
        err := os.Remove(FileLogPath)
        if err != nil {
            fmt.Printf("%c[1;0;%dm%s%c[0m\n", 0x1B, 31, fmt.Sprintf("[Log]Log Remove Error: %s !", err), 0x1B)
        } else {
            fmt.Printf("%c[1;0;%dm%s%c[0m\n", 0x1B, 32, "[Log]Log Cleaned!", 0x1B)
        }
    }
}

func getSize(path string) int64 {
    fileInfo, err := os.Stat(path)
    if err != nil {
        fmt.Printf("%c[1;0;%dm%s%c[0m\n", 0x1B, 31, fmt.Sprintf("Error: %v !", err), 0x1B)
        return 0
    }
    fileSize := fileInfo.Size()
    return fileSize
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