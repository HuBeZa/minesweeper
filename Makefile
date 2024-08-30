prompt-build:
	go build -o bin/minesweeper-prompt.exe minesweeper-prompt/main.go

prompt-run: prompt-build
	./bin/minesweeper-prompt.exe
