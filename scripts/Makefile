
mac_amd64 :
	go build -o ../build/n2a/mac/amd64 -v ../main.go
mac_arm64 :
	go build GOARCH=arm64 -o ../build/n2a/mac/arm64 -v ../main.go

win_64 :
	go build GOOS=windows GOARCH=386 -o ../build/n2a/windows/win64 -v ../main.go
win_32 :
	go build GOOS=windows GOARCH=amd64 -o ../build/n2a/windows/win32 -v ../main.go

linux_amd64 :
	go build GOOS=linux GOARCH=amd64 -o ../build/n2a/linux/amd64 -v ../main.go
linux_arm64 :
	go build GOOS=linux GOARCH=arm64 -o ../build/n2a/linux/arm64 -v ../main.go