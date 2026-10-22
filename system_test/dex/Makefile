CERT_DIR := ../cert

.PHONY: dex
dex: dex.yaml
	# Wait for kind network.
	until docker network inspect kind; do sleep 1; done

	# Build dex with upgraded ca-certificates (see Dockerfile).
	docker build -q -t dex-server:local .

	# Create a container.
	docker create -q --name dex-server -p 10443:10443 --network kind dex-server:local dex serve /dex.yaml

	# Deploy the config.
	docker cp $(CERT_DIR)/server.crt dex-server:/
	docker cp $(CERT_DIR)/server.key dex-server:/
	docker cp dex.yaml dex-server:/

	# Start the container.
	docker start dex-server
	docker logs dex-server

.PHONY: terminate
terminate:
	docker stop dex-server
	docker rm dex-server
