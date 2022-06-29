# Recursive wildcard function
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

OUT_DIR=dist/tbc
GOROOT:=$(shell go env GOROOT)

ifeq ($(shell go env GOOS),darwin)
    SED:=sed -i "" -E -e
else
    SED:=sed -i -E -e
endif

# Make everything. Keep this first so it's the default rule.
$(OUT_DIR): ui_shared balance_druid feral_druid feral_tank_druid elemental_shaman enhancement_shaman hunter mage rogue retribution_paladin protection_paladin shadow_priest smite_priest warlock warrior protection_warrior raid

# Add new sim rules here! Don't forget to add it as a dependency to the default rule above.
balance_druid: $(OUT_DIR)/balance_druid/index.js $(OUT_DIR)/balance_druid/index.css $(OUT_DIR)/balance_druid/index.html
feral_druid: $(OUT_DIR)/feral_druid/index.js $(OUT_DIR)/feral_druid/index.css $(OUT_DIR)/feral_druid/index.html
feral_tank_druid: $(OUT_DIR)/feral_tank_druid/index.js $(OUT_DIR)/feral_tank_druid/index.css $(OUT_DIR)/feral_tank_druid/index.html
elemental_shaman: $(OUT_DIR)/elemental_shaman/index.js $(OUT_DIR)/elemental_shaman/index.css $(OUT_DIR)/elemental_shaman/index.html
enhancement_shaman: $(OUT_DIR)/enhancement_shaman/index.js $(OUT_DIR)/enhancement_shaman/index.css $(OUT_DIR)/enhancement_shaman/index.html
hunter: $(OUT_DIR)/hunter/index.js $(OUT_DIR)/hunter/index.css $(OUT_DIR)/hunter/index.html
mage: $(OUT_DIR)/mage/index.js $(OUT_DIR)/mage/index.css $(OUT_DIR)/mage/index.html
rogue: $(OUT_DIR)/rogue/index.js $(OUT_DIR)/rogue/index.css $(OUT_DIR)/rogue/index.html
retribution_paladin: $(OUT_DIR)/retribution_paladin/index.js $(OUT_DIR)/retribution_paladin/index.css $(OUT_DIR)/retribution_paladin/index.html
protection_paladin: $(OUT_DIR)/protection_paladin/index.js $(OUT_DIR)/protection_paladin/index.css $(OUT_DIR)/protection_paladin/index.html
shadow_priest: $(OUT_DIR)/shadow_priest/index.js $(OUT_DIR)/shadow_priest/index.css $(OUT_DIR)/shadow_priest/index.html
smite_priest: $(OUT_DIR)/smite_priest/index.js $(OUT_DIR)/smite_priest/index.css $(OUT_DIR)/smite_priest/index.html
warlock: $(OUT_DIR)/warlock/index.js $(OUT_DIR)/warlock/index.css $(OUT_DIR)/warlock/index.html
warrior: $(OUT_DIR)/warrior/index.js $(OUT_DIR)/warrior/index.css $(OUT_DIR)/warrior/index.html
protection_warrior: $(OUT_DIR)/protection_warrior/index.js $(OUT_DIR)/protection_warrior/index.css $(OUT_DIR)/protection_warrior/index.html


raid: $(OUT_DIR)/raid/index.js $(OUT_DIR)/raid/index.css $(OUT_DIR)/raid/index.html

ui_shared: $(OUT_DIR)/lib.wasm $(OUT_DIR)/sim_worker.js $(OUT_DIR)/net_worker.js detailed_results $(OUT_DIR)/index.html
detailed_results: $(OUT_DIR)/detailed_results/index.js $(OUT_DIR)/detailed_results/index.css $(OUT_DIR)/detailed_results/index.html

$(OUT_DIR)/index.html:
	cp ui/index.html $(OUT_DIR)

clean:
	rm -f ui/core/proto/*.ts
	rm -f sim/core/proto/*.pb.go
	rm -f wowsimtbc
	rm -f wowsimtbc-windows.exe
	rm -f wowsimtbc-amd64-darwin
	rm -f wowsimtbc-amd64-linux
	rm -rf dist
	rm -rf binary_dist
	find . -name "*.results.tmp" -type f -delete

# Host a local server, for dev testing
host: $(OUT_DIR)
	# Intentionally serve one level up, so the local site has 'tbc' as the first
	# directory just like github pages.
	npx http-server $(OUT_DIR)/..

ui/core/proto/api.ts: proto/*.proto node_modules
	mkdir -p $(OUT_DIR)/protobuf-ts
	cp -r node_modules/@protobuf-ts/runtime/build/es2015/* $(OUT_DIR)/protobuf-ts
	$(SED) "s/from '(.*)';/from '\1.js';/g" $(OUT_DIR)/protobuf-ts/*.js
	$(SED) "s/from \"(.*)\";/from '\1.js';/g" $(OUT_DIR)/protobuf-ts/*.js
	npx protoc --ts_opt generate_dependencies --ts_out ui/core/proto --proto_path proto proto/api.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/test.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/ui.proto

node_modules: package-lock.json
	npm install

$(OUT_DIR)/core/tsconfig.tsbuildinfo: $(call rwildcard,ui/core,*.ts) ui/core/proto/api.ts
	npx tsc -p ui/core
	$(SED) 's#@protobuf-ts/runtime#/tbc/protobuf-ts/index#g' $(OUT_DIR)/core/proto/*.js
	$(SED) "s/from \"(.*)\";/from '\1.js';/g" $(OUT_DIR)/core/proto/*.js

# Generic rule for hosting any class directory
host_%: ui_shared %
	npx http-server $(OUT_DIR)/..

# Generic rule for building index.js for any class directory
$(OUT_DIR)/%/index.js: ui/%/index.ts ui/%/*.ts $(OUT_DIR)/core/tsconfig.tsbuildinfo
	npx tsc -p $(<D) 

# Generic rule for building index.css for any class directory
$(OUT_DIR)/%/index.css: ui/%/index.scss ui/%/*.scss $(call rwildcard,ui/core,*.scss)
	mkdir -p $(@D)
	npx sass $< $@

# Generic rule for building index.html for any class directory
$(OUT_DIR)/%/index.html: ui/index_template.html $(OUT_DIR)/assets
	$(eval title := $(shell echo $(shell basename $(@D)) | sed -r 's/(^|_)([a-z])/\U \2/g' | cut -c 2-))
	echo $(title)
	mkdir -p $(@D)
	cat ui/index_template.html | sed 's/@@TITLE@@/TBC $(title) Simulator/g' > $@

.PHONY: wasm
wasm: $(OUT_DIR)/lib.wasm

# Builds the generic .wasm, with all items included.
$(OUT_DIR)/lib.wasm: sim/wasm/* sim/core/proto/api.pb.go $(filter-out sim/core/items/all_items.go, $(call rwildcard,sim,*.go))
	@echo "Starting webassembly compile now..."
	@if GOOS=js GOARCH=wasm go build -o ./$(OUT_DIR)/lib.wasm ./sim/wasm/; then \
		echo "\033[1;32mWASM compile successful.\033[0m"; \
	else \
		echo "\033[1;31mWASM COMPILE FAILED\033[0m"; \
		exit 1; \
	fi
	

# Generic sim_worker that uses the generic lib.wasm
$(OUT_DIR)/sim_worker.js: ui/worker/sim_worker.js
	cat $(GOROOT)/misc/wasm/wasm_exec.js > $(OUT_DIR)/sim_worker.js
	cat ui/worker/sim_worker.js >> $(OUT_DIR)/sim_worker.js

$(OUT_DIR)/net_worker.js: ui/worker/net_worker.js
	cp ui/worker/net_worker.js $(OUT_DIR)

$(OUT_DIR)/assets: assets/*
	cp -r assets $(OUT_DIR)

binary_dist/dist.go: sim/web/dist.go.tmpl
	mkdir -p binary_dist/tbc
	touch binary_dist/tbc/embedded
	cp sim/web/dist.go.tmpl binary_dist/dist.go

binary_dist: $(OUT_DIR)
	rm -rf binary_dist
	mkdir -p binary_dist
	cp -r $(OUT_DIR) binary_dist/
	rm binary_dist/tbc/lib.wasm
	rm -rf binary_dist/tbc/assets/item_data

# Builds the web server with the compiled client.
wowsimtbc: binary_dist devserver

devserver: sim/web/main.go binary_dist/dist.go
	@echo "Starting server compile now..."
	@if go build -o wowsimtbc ./sim/web/main.go; then \
		echo "\033[1;32mBuild Completed Succeessfully\033[0m"; \
	else \
		echo "\033[1;31mBUILD FAILED\033[0m"; \
		exit 1; \
	fi

rundevserver: devserver
	./wowsimtbc --usefs=true --launch=false

release: wowsimtbc
	GOOS=windows GOARCH=amd64 go build -o wowsimtbc-windows.exe -ldflags="-X 'main.Version=$(VERSION)'" ./sim/web/main.go
	GOOS=darwin GOARCH=amd64 go build -o wowsimtbc-amd64-darwin -ldflags="-X 'main.Version=$(VERSION)'" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 go build -o wowsimtbc-amd64-linux   -ldflags="-X 'main.Version=$(VERSION)'" ./sim/web/main.go

sim/core/proto/api.pb.go: proto/*.proto
	protoc -I=./proto --go_out=./sim/core ./proto/*.proto

.PHONY: items
items: sim/core/items/all_items.go sim/core/proto/api.pb.go

sim/core/items/all_items.go: generate_items/*.go $(call rwildcard,sim/core/proto,*.go)
	go run generate_items/*.go -outDir=sim/core/items
	gofmt -w ./sim/core/items

test: $(OUT_DIR)/lib.wasm binary_dist/dist.go
	go test ./...

update-tests:
	find . -name "*.results" -type f -delete
	find . -name "*.results.tmp" -exec bash -c 'cp "$$1" "$${1%.results.tmp}".results' _ {} \;

fmt: tsfmt
	gofmt -w ./sim
	gofmt -w ./generate_items

tsfmt:
	for dir in $$(find ./ui -maxdepth 1 -type d -not -path "./ui" -not -path "./ui/worker"); do \
		echo $$dir; \
		npx tsfmt -r --useTsfmt ./tsfmt.json --baseDir $$dir; \
	done

# one time setup to install pre-commit hook for gofmt and npm install needed packages
setup:
	cp pre-commit .git/hooks
	chmod +x .git/hooks/pre-commit
