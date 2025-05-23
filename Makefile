

all:
	cd expParser/src && go build -o ../../exp cmd/main.go && cd ../.. && \
	cd assembler && go build -o ../asmp1 cmd/main.go && cd .. && \
	cd neanderExecutor && go build -o ../neander cmd/main.go && cd ..

test:
	./exp programa.lpn output.asm
	./asmp1 output.asm output.bin
	./neander output.bin

clear:
	rm exp asmp1 neander
