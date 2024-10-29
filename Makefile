all:
	docker build -t abhijaju/httpserver -f docker/Dockerfile.httpserver .
	docker push abhijaju/httpserver
