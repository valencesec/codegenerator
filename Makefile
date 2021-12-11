all:
	echo "Choose make target"

.PHONY: test
test:
	rm -fr test/workdir
	cp -a test/original test/workdir
	go run cmd/codegenerator/main.go -dir test/workdir
	diff test/workdir test/expected

accept_new_test_result:
	cp -a test/workdir/* test/expected/
