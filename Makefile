build:
	GOOS=linux go build parsepacket.go

clean:
	rm -f parsepacket parsepacket.exe
