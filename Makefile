
run:
	go run main.go
pushmain:
	git add . && git commit -m "$(MSG)" && git push 

pushv:
	git tag $(TAG) && git push origin $(TAG)
cleancli:
	rm -rf ~/devkitcli/ ~/.config/devkitcli/ ~/release/
buld_mac:
	GOOS=darwin GOARCH=amd64 go build -o devkit-mac-amd64 && \
	GOOS=darwin GOARCH=arm64 go build -o devkit-mac-arm64

build: 
	go build -o devkit && cp devkit ~/devkitcli/
clean:
	cd example && rm -rf app api proto sup*
init:
	make build && rm -rf new_fork/devkit.env && cd new_fork && ../devkit init
init_dirs:
	cd example && mkdir -p app api proto/devkit/v1 supabase/queries && cp ../service.proto proto/devkit/v1/devkit_service.proto && cp ../api_template.txt api/api.go
new_api:
	rm -rf new_fork && make build &&  ./devkit new api new_fork --org esolveeg  
new_endpoint:
	 make build && cd new_fork && ../devkit new endpoint list -d products -f categories
new_feature:
	 make build && cd new_fork && ../devkit new feature categories -d products
new_domain:
	make build && cd new_fork && rm -rf app/products && ../devkit new domain products
download_cli:
	curl -sSL https://raw.githubusercontent.com/darwishdev/devkit-cli/refs/heads/main/install.sh | bash
seed_storage:
	make build && cd new_fork && ../devkit seed storage -f assets -i icons
seed:
	make build && cd new_fork && ../devkit seed accounts_schema --file-path accounts.xlsx --out-file q.sql
endpoint_test:
	make build && make clean init_dirs new_domain new_feature new_endpoint
