go build -o main cmd/sync/main.go
./main
rm -rf testing/b/*
time ./main
rm main