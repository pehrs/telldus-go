
all:	build

build:
	go build github.com/pehrs/telldus-go/td


c:
	cd td && 	go tool cgo -objdir ../_obj sensorevent.go
