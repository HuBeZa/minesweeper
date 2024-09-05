prompt-build:
	go build -o bin/minesweeper-prompt.exe minesweeper-prompt/main.go

prompt-run: prompt-build
	./bin/minesweeper-prompt.exe

bubbletea-build:
	go build -o bin/minesweeper-bubbletea.exe minesweeper-bubbletea/main.go

bubbletea-run: bubbletea-build
	./bin/minesweeper-bubbletea.exe
