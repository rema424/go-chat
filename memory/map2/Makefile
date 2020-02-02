profile-mem: run
	go tool pprof -http=":8081" main mem.pprof

profile-cpu: run
	go tool pprof -http=":8081" main cpu.pprof

run: build
	./main

build: clean
	go build main.go

clean:
	rm -f main *.pprof