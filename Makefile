
run:
	go run main.go
pushmain:
	git add . && git commit -m "$(MSG)" && git push 

pushv:
	git tag $(TAG) && git push origin $(TAG)
cleancli:
	rm -rf ~/devkitcli/ ~/.config/devkitcli/ ~/release/
build : 
	go build -o devkit
clean:
	cd example && rm -rf app api proto sup*
init:
	make build && rm -rf new_fork/devkit.env && cd new_fork && ../devkit init
init_dirs:
	cd example && mkdir -p app api proto/devkit/v1 supabase/queries && cp ../service.proto proto/devkit/v1/devkit_service.proto && cp ../api_template.txt api/api.go
new_api:
	rm -rf new_fork && make build &&  ./devkit new api new_fork --org esolveeg  
new_endpoint:
	 make build && cd example && ../devkit new endpoint accounts roles create_update
new_feature:
	 make build && cd example && ../devkit new feature accounts roles
new_domain:
	make build && cd example && rm -rf app/accounts && ../devkit new domain accounts
download_cli:
	curl -sSL https://raw.githubusercontent.com/darwishdev/devkit-cli/refs/heads/main/install.sh | bash

endpoint_test:
	make build && make clean init_dirs new_domain new_feature new_endpoint
