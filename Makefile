test:
	go run main.go
all:
	go build main.go
	ls example/wavs/*wav | xargs basename | parallel --progress --bar ./main -in example/wavs/{} -out example/processed/{}
all-drums:
	go build main.go
	ls example/wavs/drums/*.wav | xargs basename | parallel --progress --bar ./main -b 24 -r 44100 -in example/wavs/drums/{} -out example/processed/drums/{}
all-speech:
	go build main.go
	ls example/wavs/speech/*.wav | xargs basename | parallel --progress --bar ./main -b 16 -r 16000 -in example/wavs/speech/{} -out example/processed/speech/{}
all-instrumental:
	go build main.go
	ls example/wavs/instrumental/*.wav | xargs basename | parallel --progress --bar ./main -b 16 -r 44100 -in example/wavs/instrumental/{} -out example/processed/instrumental/{}
listen:
	for file in example/processed/*.wav; do \
		echo $$file; \
		sox $$file -n stats ; \
		afplay $$file ; \
	done
listen-speech:
	for file in example/processed/speech/*.wav; do \
		echo $$file; \
		sox $$file -n stats ; \
		afplay $$file ; \
	done
listen-drums:
	for file in example/processed/drums/*.wav; do \
		echo $$file; \
		sox $$file -n stats ; \
		afplay $$file ; \
	done
listen-instrumental:
	for file in example/processed/instrumental/*.wav; do \
		echo $$file; \
		sox $$file -n stats ; \
		afplay $$file ; \
	done
listen-clean:
	for file in example/wavs/*.wav; do \
		echo $$file ; \
		sox $$file -n stats ; \
		afplay $$file ; \
	done