OUT_DIR := dist/wotlk
TS_CORE_SRC := $(shell find ui/core -name '*.ts' -type f)
ASSETS_INPUT := $(shell find assets/ -type f)
ASSETS := $(patsubst assets/%,$(OUT_DIR)/assets/%,$(ASSETS_INPUT))
# Recursive wildcard function. Needs to be '=' instead of ':=' because of recursion.
rwildcard = $(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))
GOROOT := $(shell go env GOROOT)
UI_SRC := $(shell find ui -name '*.ts' -o -name '*.tsx' -o -name '*.scss' -o -name '*.html')
HTML_INDECIES := ui/balance_druid/index.html \
				 ui/feral_druid/index.html \
				 ui/feral_tank_druid/index.html \
				 ui/restoration_druid/index.html \
				 ui/elemental_shaman/index.html \
				 ui/enhancement_shaman/index.html \
				 ui/restoration_shaman/index.html \
				 ui/hunter/index.html \
				 ui/mage/index.html \
				 ui/rogue/index.html \
				 ui/holy_paladin/index.html \
				 ui/protection_paladin/index.html \
				 ui/retribution_paladin/index.html \
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
	  wowsimwotlk-arm64-darwin \
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
	cat ui/index_template.html | sed -e 's/@@TITLE@@/WOTLK $(title) Simulator/g' -e 's/@@SPEC@@/$(shell basename $(@D))/g' > $@

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
	cat ui/index_template.html | sed -e 's/@@TITLE@@/WOTLK $(title) Simulator/g' -e 's/@@SPEC@@/$(shell basename $(@D))/g' > $@

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

binary_dist: $(OUT_DIR)/.dirstamp
	rm -rf binary_dist
	mkdir -p binary_dist
	cp -r $(OUT_DIR) binary_dist/
	rm binary_dist/wotlk/lib.wasm
	rm -rf binary_dist/wotlk/assets/db_inputs
	rm binary_dist/wotlk/assets/database/db.bin
	rm binary_dist/wotlk/assets/database/leftover_db.bin

# Rebuild the protobuf generated code.
.PHONY: proto
proto: sim/core/proto/api.pb.go ui/core/proto/api.ts

# Builds the web server with the compiled client.
.PHONY: wowsimwotlk
wowsimwotlk: binary_dist devserver

.PHONY: devserver
devserver: sim/core/proto/api.pb.go sim/web/main.go binary_dist/dist.go
	@echo "Starting server compile now..."
	@if go build -o wowsimwotlk ./sim/web/main.go; then \
		printf "\033[1;32mBuild Completed Successfully\033[0m\n"; \
	else \
		printf "\033[1;31mBUILD FAILED\033[0m\n"; \
		exit 1; \
	fi

.PHONY: air
air:
ifeq ($(WATCH), 1)
	@if ! command -v air; then \
		echo "Missing air dependency. Please run \`make setup\`"; \
		exit 1; \
	fi
endif

rundevserver: air devserver
ifeq ($(WATCH), 1)
	npx vite build -m development --watch &
	ulimit -n 10240 && air -tmp_dir "/tmp" -build.include_ext "go,proto" -build.args_bin "--usefs=true --launch=false" -build.bin "./wowsimwotlk" -build.cmd "make devserver" -build.exclude_dir "assets,dist,node_modules,ui,tools"
else
	./wowsimwotlk --usefs=true --launch=false --host=":3333"
endif

wowsimwotlk-windows.exe: wowsimwotlk
# go build only considers syso files when invoked without specifying .go files: https://github.com/golang/go/issues/16090
	cp ./assets/favicon_io/icon-windows_amd64.syso ./sim/web/icon-windows_amd64.syso
	cd ./sim/web/ && GOOS=windows GOARCH=amd64 GOAMD64=v2 go build -o wowsimwotlk-windows.exe -ldflags="-X 'main.Version=$(VERSION)' -s -w"
	cd ./cmd/wowsimcli && GOOS=windows GOARCH=amd64 GOAMD64=v2 go build -o wowsimcli-windows.exe --tags=with_db -ldflags="-X 'main.Version=$(VERSION)' -s -w"
	rm ./sim/web/icon-windows_amd64.syso
	mv ./sim/web/wowsimwotlk-windows.exe ./wowsimwotlk-windows.exe
	mv ./cmd/wowsimcli/wowsimcli-windows.exe ./wowsimcli-windows.exe

release: wowsimwotlk wowsimwotlk-windows.exe
	GOOS=darwin GOARCH=amd64 GOAMD64=v2 go build -o wowsimwotlk-amd64-darwin -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=darwin GOARCH=arm64 go build -o wowsimwotlk-arm64-darwin -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 GOAMD64=v2 go build -o wowsimwotlk-amd64-linux   -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./sim/web/main.go
	GOOS=linux GOARCH=amd64 GOAMD64=v2 go build -o wowsimcli-amd64-linux --tags=with_db -ldflags="-X 'main.Version=$(VERSION)' -s -w" ./cmd/wowsimcli/cli_main.go
# Now compress into a zip because the files are getting large.
	zip wowsimwotlk-windows.exe.zip wowsimwotlk-windows.exe
	zip wowsimwotlk-amd64-darwin.zip wowsimwotlk-amd64-darwin
	zip wowsimwotlk-arm64-darwin.zip wowsimwotlk-arm64-darwin
	zip wowsimwotlk-amd64-linux.zip wowsimwotlk-amd64-linux
	zip wowsimcli-amd64-linux.zip wowsimcli-amd64-linux
	zip wowsimcli-windows.exe.zip wowsimcli-windows.exe

sim/core/proto/api.pb.go: proto/*.proto
	protoc -I=./proto --go_out=./sim/core ./proto/*.proto

# Only useful for building the lib on a host platform that matches the target platform
.PHONY: locallib
locallib: sim/core/proto/api.pb.go
	go build -buildmode=c-shared -o wowsimwotlk.so --tags=with_db ./sim/lib/library.go

.PHONY: nixlib
nixlib: sim/core/proto/api.pb.go
	GOOS=linux GOARCH=amd64 GOAMD64=v2 go build -buildmode=c-shared -o wowsimwotlk-linux.so --tags=with_db ./sim/lib/library.go

.PHONY: winlib
winlib: sim/core/proto/api.pb.go
	GOOS=windows GOARCH=amd64 GOAMD64=v2 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -o wowsimwotlk-windows.dll --tags=with_db ./sim/lib/library.go

.PHONY: items
items: sim/core/items/all_items.go sim/core/proto/api.pb.go

sim/core/items/all_items.go: $(call rwildcard,tools/database,*.go) $(call rwildcard,sim/core/proto,*.go)
	go run tools/database/gen_db/*.go -outDir=./assets -gen=db

.PHONY: test
test: $(OUT_DIR)/lib.wasm binary_dist/dist.go
	go test --tags=with_db ./sim/...

.PHONY: update-tests
update-tests:
	find . -name "*.results" -type f -delete
	find . -name "*.results.tmp" -exec bash -c 'cp "$$1" "$${1%.results.tmp}".results' _ {} \;

.PHONY: fmt
fmt: tsfmt
	gofmt -w ./sim
	gofmt -w ./tools

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
	! command -v air && curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin || true

# Host a local server, for dev testing
.PHONY: host
host: air $(OUT_DIR)/.dirstamp node_modules
ifeq ($(WATCH), 1)
	ulimit -n 10240 && air -tmp_dir "/tmp" -build.include_ext "go,ts,js,html" -build.bin "npx" -build.args_bin "http-server $(OUT_DIR)/.." -build.cmd "make" -build.exclude_dir "dist,node_modules,tools"
else
	# Intentionally serve one level up, so the local site has 'wotlk' as the first
	# directory just like github pages.
	npx http-server $(OUT_DIR)/..
endif
