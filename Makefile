all:
	echo "build and push httpserver"
	GOOS=linux GOARCH=amd64 go build -o bin/httpserver httpserver/main.go
	docker build -t abhijaju/httpserver -f docker/Dockerfile.httpserver .
	docker push abhijaju/httpserver

	echo "build and push operator"
	GOOS=linux GOARCH=amd64 go build -o bin/operator operator/main.go
	docker build -t abhijaju/operator -f docker/Dockerfile.operator .
	docker push abhijaju/operator
