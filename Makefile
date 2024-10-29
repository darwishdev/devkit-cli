build : 
	go build 
clean:
	cd example && rm -rf app api proto sup*
init:
	make build && cd example && rm -rf devkit.yaml && ../devkit-cli init
init_dirs:
	cd example && mkdir -p app api proto/devkit/v1 supabase/queries && cp ../service.proto proto/devkit/v1/devkit_service.proto && cp ../api_template.txt api/api.go
new_api:
	rm -rf new_fork &&	make build &&  ./devkit-cli new api new_fork esolveeg buf.build/ahmeddarwish && cp devkit.yaml new_fork/
new_endpoint:
	 make build && cd example && ../devkit-cli new endpoint accounts roles create_update
new_feature:
	 make build && cd example && ../devkit-cli new feature accounts roles
new_domain:
	make build && cd example && rm -rf app/accounts && ../devkit-cli new domain accounts
endpoint_test:
	make build && make clean init_dirs new_domain new_feature new_endpoint
