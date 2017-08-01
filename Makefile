test:
	go run main.go
all:
	go build main.go
	ls example/wavs/*wav | xargs basename | parallel --progress ./main -in example/wavs/{} -out example/processed/{}
listen:
	for file in example/processed/*.wav; do \
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