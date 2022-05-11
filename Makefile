build-prometheus-server:
	docker build -f Dockerfile.prometheus -t my-prom-server . 
run-prometheus-server:
	docker run -p 9090:9090 my-prom-server
run-app-server:
	go run main.go