build_windows:
	test -d bin/windows || mkdir -p bin/windows;
	cd bin/windows; GOOS=windows go build  ../../main/seckill.go

build_linux:
	test -d bin/linux || mkdir -p bin/linux;
	cd bin/linux;GOOS=linux go build ../../main/seckill.go