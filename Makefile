deploy:
	go run .

test-curl:
	curl -v -w "\n" http://localhost:8080/
	curl -v -w "\n" http://localhost:8080/version
	curl -v -w "\n" http://localhost:8080/duration
