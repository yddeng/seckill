build_windows:
	GOOS=windows go build -o bin/windows main/seckill.go

build_linux:
	GOOS=linux go build -o bin/linux main/seckill.go