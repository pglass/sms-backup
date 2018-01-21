.PHONY: run

INPUT = ~/Downloads/sms-20180119165444.xml
OUTFILE = out.png
MY_NUMBER ?=

main: main.go parse/*.go analyze/*.go
	go build $<

run: main
	rm -f $(OUTFILE)
	./main -f $(INPUT) -o $(OUTFILE) -n $(MY_NUMBER) -t messagesTimeOfDay
	open $(OUTFILE)

clean:
	rm -f main
