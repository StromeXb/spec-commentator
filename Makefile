# https://stackoverflow.com/questions/2214575/passing-arguments-to-make-run
# If the first argument is "run"...
ifeq (gen,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  SERVICE_NAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(SERVICE_NAME):;@:)
endif

ifeq (copy_proto,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  SERVICE_NAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(SERVICE_NAME):;@:)
endif

generate_orm:
	# генерация ent файлов
	go run entgo.io/ent/cmd/ent generate --target ./ms/$(SERVICE_NAME)/repository/ent ./ms/$(SERVICE_NAME)/repository/schema

generate_gateway:
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -package generated -generate types api/openapi/external/speccomentor.yaml > ms/gateway/generated/types.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -package generated -generate chi-server api/openapi/external/speccomentor.yaml > ms/gateway/generated/server.gen.go
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -package generated -generate spec api/openapi/external/speccomentor.yaml > ms/gateway/generated/spec.gen.go