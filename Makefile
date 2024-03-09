
.PHONY: all
all:
	go run src/main.go -cpuprofile=billion.prof
	go tool pprof -raw billion.prof > bil.raw
	../FlameGraph/stackcollapse-go.pl bil.raw > bil.collapsed
	../FlameGraph/flamegraph.pl bil.collapsed > bil.svg
	open bil.svg
