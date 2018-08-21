# SWAPIDC_NO
Add random Users and Tickets to SWAPIDC

### Usage(You can get the build version in Releases)
* go build
* ./SwapNo -log -rate 5 -url https://www.site.com/index/register/

### Infos
* -log Log the outputs
* -loglimit Log Limit(MB)(default 10240M)
* -proxy Enable proxy mode
* -proxyupdate Update the Proxy list(From https://github.com/fate0/proxylist)
* -rate The rate of the Import Process(s)
* -url The Register Url of SWAPIDC which you want to add users to~~~
* -debug Show the Post results
* -overclock Change the rate to Millisecond(1s = 1000ms)
* -tickets Enable Ticket Mode
* -ticketurl The Ticket Url of SWAPIDC which you want to add tickets to~~~
* -ticketprocess The Amount of Process to open tickets

### Simple Sample
```./SwapNo -log -proxy -rate 2 -url https://site.com/index/register/```

### OverClock Mode
```./SwapNo -log -proxy -rate 20 -url https://site.com/index/register/ -overclock```

### Update Proxys
```./SwapNo -proxyupdate```

### TicketMode
* 1. Add Users into users.txt like:
* ({"username": "xxxxx", "password": "123456"})
* 2. Run the Script
```./SwapNo -log -rate 2 -url http://site.com/index/login/ -proxy -tickets -ticketurl http://site.com/ticket/submit/```
