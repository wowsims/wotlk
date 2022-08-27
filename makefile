OUT_DIR := dist/wotlk
TS_CORE_SRC := $(shell find ui/core -name '*.ts' -type f)
ASSETS_INPUT := $(shell find assets/ -type f)
ASSETS := $(patsubst assets/%,$(OUT_DIR)/assets/%,$(ASSETS_INPUT))
# Recursive wildcard function
rwildcard := $(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))
GOROOT := $(shell go env GOROOT)
UI_SRC := $(shell find ui -name '*.ts' -o -name '*.scss' -o -name '*.html')
HTML_INDECIES := ui/balance_druid/index.html \
				 ui/feral_druid/index.html \
				 ui/feral_tank_druid/index.html \
				 ui/elemental_shaman/index.html \
				 ui/enhancement_shaman/index.html \
				 ui/hunter/index.html \
				 ui/mage/index.html \
				 ui/rogue/index.html \
				 ui/retribution_paladin/index.html \
				 ui/protection_paladin/index.html \
				 ui/healing_priest/index.html \
				 ui/shadow_priest/index.html \
				 ui/smite_priest/index.html \
				 ui/warlock/index.html \
				 ui/warrior/index.html \
				 ui/protection_warrior/index.html \
				 ui/deathknight/index.html \
				 ui/tank_deathknight/index.html \
				 ui/raid/index.html \
				 ui/detailed_results/index.html

$(OUT_DIR)/.dirstamp: \
  $(OUT_DIR)/lib.wasm \
  ui/core/proto/api.ts \
  $(ASSETS) \
  $(OUT_DIR)/bundle/.dirstamp
	touch $@

$(OUT_DIR)/bundle/.dirstamp: \
  $(UI_SRC) \
  $(HTML_INDECIES) \
  vite.config.js \
  node_modules \
  tsconfig.json \
  ui/core/index.ts \
  ui/core/proto/api.ts \
  $(OUT_DIR)/net_worker.js \
  $(OUT_DIR)/sim_worker.js
	npx tsc --noEmit
	npx vite build
	touch $@

$(OUT_DIR)/sim_worker.js: ui/worker/sim_worker.js
	cat '$(GOROOT)/misc/wasm/wasm_exec.js' > $(OUT_DIR)/sim_worker.js
	cat ui/worker/sim_worker.js >> $(OUT_DIR)/sim_worker.js

$(OUT_DIR)/net_worker.js: ui/worker/net_worker.js
	cp ui/worker/net_worker.js $(OUT_DIR)

ui/core/index.ts: $(TS_CORE_SRC)
	find ui/core -name '*.ts' | \
	  awk -F 'ui/core/' '{ print "import \x22./" $$2 "\x22;" }' | \
	  sed 's/\.ts";$$/";/' | \
	  grep -v 'import "./index";' > $@

.PHONY: clean
clean:
	rm -rf ui/core/proto/*.ts \
	  sim/core/proto/*.pb.go \
	  wowsimwotlk \
	  wowsimwotlk-windows.exe \
	  wowsimwotlk-amd64-darwin \
	  wowsimwotlk-amd64-linux \
	  dist \
	  binary_dist \
	  ui/core/index.ts \
	  ui/core/proto/*.ts \
	  node_modules \
	  $(HTML_INDECIES)
	find . -name "*.results.tmp" -type f -delete


ui/core/proto/api.ts: proto/*.proto node_modules
	npx protoc --ts_opt generate_dependencies --ts_out ui/core/proto --proto_path proto proto/api.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/test.proto
	npx protoc --ts_out ui/core/proto --proto_path proto proto/ui.proto

ui/%/index.html: ui/index_template.html
	$(eval title := $(shell echo $(shell basename $(@D)) | sed -r 's/(^|_)([a-z])/\U \2/g' | cut -c 2-))
	cat ui/index_template.html | sed 's/@@TITLE@@/WOTLK $(title) Simulator/g' > $@

package-lock.json:
	npm install

node_modules: package-lock.json
	npm ci

# Generic rule for hosting any class directory
.PHONY: host_%
host_%: $(OUT_DIR) node_modules
	npx http-server $(OUT_DIR)/..

# Generic rule for building index.html for any class directory
$(OUT_DIR)/%/index.html: ui/index_template.html $(OUT_DIR)/assets
	$(eval title := $(shell echo $(shell basename $(@D)) | sed -r 's/(^|_)([a-z])/\U \2/g' | cut -c 2-))
	echo $(title)
	mkdir -p $(@D)
	cat ui/index_template.html | sed 's/@@TITLE@@/WOTLK $(title) Simulator/g' > $@

.PHONY: wasm
wasm: $(OUT_DIR)/lib.wasm

# Builds the generic .wasm, with all items included.
$(OUT_DIR)/lib.wasm: sim/wasm/* sim/core/proto/api.pb.go $(filter-out sim/core/items/all_items.go, $(call rwildcard,sim,*.go))
	@echo "Starting webassembly compile now..."
	@if GOOS=js GOARCH=wasm go build -o ./$(OUT_DIR)/lib.wasm ./sim/wasm/; then \
		printf "\033[1;32mWASM compile successful.\033[0m\n"; \
	else \
		printf "\033[1;31mWASM COMPILE FAILED\033[0m\n"; \
		exit 1; \
	fi
	
$(OUT_DIR)/assets/%: assets/%
	mkdir -p $(@D)
	cp $< $@

binary_dist/dist.go: sim/web/dist.go.tmpl
	mkdir -p binary_dist/wotlk
	touch binary_dist/wotlk/embedded
	cp sim/web/dist.go.tmpl binary_dist/dist.go

binary_dist: $(OUT_DIR)
	rm -rf binary_dist
	mkdir -p binary_dist
	cp -r $(OUT_DIR) binary_dist/
	rm binary_dist/wotlk/lib.wasm
	rm -rf binary_dist/wotlk/assets/item_data
	mkdir -p binary_dist/wotlk/assets/item_data
	cp $(OUT_DIR)/assets/item_data/all_items_db.json ./binary_dist/wotlk/assets/item_data/all_items_db.json
	rm -rf binary_dist/wotlk/assets/spell_data
	mkdir -p binary_dist/wotlk/assets/spell_data
	cp $(OUT_DIR)/assets/spell_data/all_spells_db.json ./binary_dist/wotlk/assets/spell_data/all_spells_db.json

# Builds the web server with the compiled client.
.PHONY: wowsimwotlk
wowsimwotlk: binary_dist devserver

.PHONY: devserver
devserver: sim/core/proto/api.pb.go sim/web/main.go binary_dist/dist.go
	@echo "Starting server compile now..."
	@if go build -o wowsimwotlk ./sim/web/main.go; then \
		printf "\033[1;32mBuild Completed Succeessfully\033[0m\n"; \
	else \
		printf "\033[1;31mBUILD FAILED\033[0m\n"; \
		exit 1; \
	fi

rundevserver: devserver
	./wowsimwotlk --usefs=true --launch=false

release: wowsimwotlk
	GOOS=windows GOARCH=amd64 go build -o wowsimwotlk-windows.exe -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=darwin GOARCH=amd64 go build -o wowsimwotlk-amd64-darwin -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 go build -o wowsimwotlk-amd64-linux   -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
# Now compress into a zip because the files are getting large.
	zip wowsimwotlk-windows.exe.zip wowsimwotlk-windows.exe
	zip wowsimwotlk-amd64-darwin.zip wowsimwotlk-amd64-darwin
	zip wowsimwotlk-amd64-linux.zip wowsimwotlk-amd64-linux

sim/core/proto/api.pb.go: proto/*.proto
	protoc -I=./proto --go_out=./sim/core ./proto/*.proto

.PHONY: items
items: sim/core/items/all_items.go sim/core/proto/api.pb.go

sim/core/items/all_items.go: generate_items/*.go $(call rwildcard,sim/core/proto,*.go)
	go run generate_items/*.go -outDir=sim/core/items
	gofmt -w ./sim/core/items

.PHONY: test
test: $(OUT_DIR)/lib.wasm binary_dist/dist.go
	go test ./...

.PHONY: update-tests
update-tests:
	find . -name "*.results" -type f -delete
	find . -name "*.results.tmp" -exec bash -c 'cp "$$1" "$${1%.results.tmp}".results' _ {} \;

.PHONY: fmt
fmt: tsfmt
	gofmt -w ./sim
	gofmt -w ./generate_items

.PHONY: tsfmt
tsfmt:
	for dir in $$(find ./ui -maxdepth 1 -type d -not -path "./ui" -not -path "./ui/worker"); do \
		echo $$dir; \
		npx tsfmt -r --useTsfmt ./tsfmt.json --baseDir $$dir; \
	done

# one time setup to install pre-commit hook for gofmt and npm install needed packages
setup:
	cp pre-commit .git/hooks
	chmod +x .git/hooks/pre-commit

# Host a local server, for dev testing
.PHONY: host
host: $(OUT_DIR)/.dirstamp node_modules
	# Intentionally serve one level up, so the local site has 'wotlk' as the first
	# directory just like github pages.
	npx http-server $(OUT_DIR)/..
