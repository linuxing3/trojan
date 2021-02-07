test:
	go test trojan/... -v

run:
	go run main.go

dist:
	make pack
	make build

pack:
	packr2

build:
	go build -ldflags "-s -w -X 'trojan/xray.MVersion=`git describe --tags $(git rev-list --tags --max-count=1)`' -X 'trojan/xray.BuildDate=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`' -X 'trojan/xray.GoVersion=`go version|awk '{print $3,$4}'`' -X 'trojan/xray.GitVersion=`git rev-parse HEAD`'" -o "result/xxxray" .

clean:	
	packr2 clean
	rm -rf result

upload:
	scp result/xxxray root@dongxishijie.xyz:xxxray
