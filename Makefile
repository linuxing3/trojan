test:
	go test trojan/... -v

run:
	go run main.go

build:
	packr2
	go build -o "result/xxxray" .
	echo "清理临时文件"
	packr2 clean
	rm -rf result

upload:
	scp result/xxxray root@dongxishijie.xyz:xxxray
