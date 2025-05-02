## MingW-w64
https://github.com/go-vgo/Mingw

##build
```
# rsrc.exe -manifest taobao_links.exe.manifest -arch amd64 -ico ./assets/img/icon.ico -o rsrc.syso
# go build -ldflags="-H windowsgui -buildmode=exe -X github.com/wppzxc/taobao_links/pkg/version.version='$version'"
``` 
