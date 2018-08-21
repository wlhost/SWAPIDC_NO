# SWAPIDC_NO
Add random users to SWAPIDC

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

### Simple Sample
```./SwapNo -log -proxy -rate 2 -url https://site.com/index/register/```