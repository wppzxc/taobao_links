##build
```
# rsrc.exe -manifest taobao_links.exe.manifest -arch amd64 -ico ./assets/img/icon.ico -o rsrc.syso
# go build -ldflags="-H windowsgui -X github.com/wppzxc/taobao_links/pkg/version.version='$version'"
``` 
