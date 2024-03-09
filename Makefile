
.PHONY: all
all:
	go run src/main.go -cpuprofile=billion.prof
	go tool pprof -raw billion.prof > bil.raw
	../FlameGraph/stackcollapse-go.pl bil.raw > bil.collapsed
	../FlameGraph/flamegraph.pl bil.collapsed > bil.svg
	open bil.svg

cheat:
	go run cheat.go measurements.txt
	go tool pprof -raw cheat.prof > cheat.raw
	../FlameGraph/stackcollapse-go.pl cheat.raw > cheat.collapsed
	../FlameGraph/flamegraph.pl cheat.collapsed > cheat.svg
	open cheat.svg