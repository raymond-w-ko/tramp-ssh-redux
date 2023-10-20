all:
	go build -o bin/tramp-ssh-redux-client ./cmd/tramp-ssh-redux-client	
	go build -o bin/tramp-ssh-redux-proxy ./cmd/tramp-ssh-redux-proxy
	go build -o bin/tramp-ssh-redux-server ./cmd/tramp-ssh-redux-server
