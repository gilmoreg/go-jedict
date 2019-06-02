run: 
	go run cmd/main.go -import data/JMdict_e.xml

read:
	go run tools/read.go

ipa:
	go run tools/ipacsv.go