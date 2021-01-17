module github.com/wppzxc/taobao_links

go 1.13

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/atotto/clipboard v0.1.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-vgo/robotgo v0.0.0-20191128163956-6b94d024dc37
	github.com/hpcloud/tail v1.0.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/lxn/win v0.0.0-20201111105847-2a20daff6a55
	github.com/mvdan/xurls v1.1.0
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4
	github.com/tealeg/xlsx v1.0.5
	github.com/zserge/lorca v0.1.8
	github.com/zserge/webview v0.0.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/klog v1.0.0
	sigs.k8s.io/yaml v1.1.0
)

replace github.com/zserge/webview v0.0.0 => ./libs/webview
