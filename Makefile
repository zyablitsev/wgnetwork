BUILDCMD=CGO_ENABLED=0 go build -mod=vendor
TESTCMD=go test -v -cover -race -mod=vendor

# load config parameters
ifneq (,$(wildcard ./config))
	include config
endif

v_bin_dir=$(or ${BIN_DIR},${bin_dir})

v_log_level=$(or ${LOG_LEVEL},${log_level})
v_db_path=$(or ${DB_PATH},${db_path})
v_wg_binary=$(or ${WG_BINARY},${wg_binary})
v_wg_iface=$(or ${WG_IFACE},${wg_iface})
v_wg_port=$(or ${WG_PORT},${wg_port})
v_wg_cidr=$(or ${WG_CIDR},${wg_cidr})
v_dns_tcp_port=$(or ${DNS_TCP_PORT},${dns_tcp_port})
v_dns_udp_port=$(or ${DNS_UDP_PORT},${dns_udp_port})
v_dns_resolver_addrs=$(or ${DNS_RESOLVER_ADDRS},${dns_resolver_addrs})
v_dns_zone=$(or ${DNS_ZONE},${dns_zone})
v_fe_http_port=$(or ${FE_HTTP_PORT},${fe_http_port})
v_api_http_port=$(or ${API_HTTP_PORT},${api_http_port})
v_api_unix_socket=$(or ${API_UNIX_SOCKET},${api_unix_socket})
v_otp_issuer=$(or ${OTP_ISSUER},${otp_issuer})
v_httporigin=$(or ${HTTPORIGIN},${httporigin})
v_session_secret=$(or ${SESSION_SECRET},${session_secret})
v_session_ttl=$(or ${SESSION_TTL},${session_ttl})

v_nft_enabled=$(or ${NFT_ENABLED},${nft_enabled})
v_nft_default_policy=$(or ${NFT_DEFAULT_POLICY},${nft_default_policy})

v_dev_hostname=$(or ${DEV_HOSTNAME},${dev_hostname})
v_dev_authip=$(or ${DEV_AUTHIP},${dev_authip})

v_env=LOG_LEVEL="${v_log_level}" \
	DB_PATH="${v_db_path}" \
	WG_BINARY="${v_wg_binary}" \
	WG_IFACE="${v_wg_iface}" \
	WG_PORT="${v_wg_port}" \
	WG_CIDR="${v_wg_cidr}" \
	DNS_TCP_PORT="${v_dns_tcp_port}" \
	DNS_UDP_PORT="${v_dns_udp_port}" \
	DNS_RESOLVER_ADDRS="${v_dns_resolver_addrs}" \
	DNS_ZONE="${v_dns_zone}" \
	FE_HTTP_PORT="${v_fe_http_port}" \
	API_HTTP_PORT="${v_api_http_port}" \
	API_UNIX_SOCKET="${v_api_unix_socket}" \
	OTP_ISSUER="${v_otp_issuer}" \
	HTTPORIGIN="${v_httporigin}" \
	SESSION_SECRET="${v_session_secret}" \
	SESSION_TTL="${v_session_ttl}" \
	NFT_ENABLED="${v_nft_enabled}" \
	NFT_DEFAULT_POLICY="${v_nft_default_policy}" \
	DEV_HOSTNAME="${v_dev_hostname}" \
	DEV_AUTHIP="${v_dev_authip}"

install-dependencies: install-dependencies-fe install-dependencies-be

install-dependencies-fe:
	cd femanager/spa/ && npm i

install-dependencies-be:
	go get .

build: build-fe build-bootstrap_trust_ipset build-managercli build-service

build-fe:
	cd femanager/spa/ && \
		npm run build

build-bootstrap_trust_ipset:
	GOOS=linux GOARCH=amd64 \
		${BUILDCMD} -o ${v_bin_dir}/wgn-bootstrap-trust-ipset_linux_amd64 cmd/bootstrap-trust-ipset/main.go

build-managercli:
	GOOS=linux GOARCH=amd64 \
		${BUILDCMD} -o ${v_bin_dir}/wgn-managercli_linux_amd64 cmd/managercli/main.go

build-service:
	GOOS=linux GOARCH=amd64 \
		${BUILDCMD} -o ${v_bin_dir}/wgnetwork_linux_amd64 cmd/service/main.go

docker-build:
	docker build -t zyablitsev/wgnetwork .

build-dev: build-fe build-bootstrap_trust_ipset-dev build-managercli-dev build-service-dev

build-bootstrap_trust_ipset-dev:
	${BUILDCMD} -o ${v_bin_dir}/wgn-bootstrap-trust-ipset cmd/bootstrap-trust-ipset/main.go

build-managercli-dev:
	${BUILDCMD} -o ${v_bin_dir}/wgn-managercli cmd/managercli/main.go

build-service-dev:
	${BUILDCMD} -o ${v_bin_dir}/wgnetwork cmd/service/main.go

docker-run-dev: docker-rm-dev
	docker run \
		-e LOG_LEVEL=${v_log_level} \
		-e DB_PATH=wgnetwork.db \
		-e WG_BINARY=${v_wg_binary} \
		-e WG_IFACE=${v_wg_iface} \
		-e WG_PORT=${v_wg_port} \
		-e WG_CIDR=${v_wg_cidr} \
		-e DNS_TCP_PORT=${v_dns_tcp_port} \
		-e DNS_UDP_PORT=${v_dns_udp_port} \
		-e DNS_RESOLVER_ADDRS=${v_dns_resolver_addrs} \
		-e DNS_ZONE=${v_dns_zone} \
		-e FE_HTTP_PORT=${v_fe_http_port} \
		-e API_HTTP_PORT=${v_api_http_port} \
		-e OTP_ISSUER=${v_otp_issuer} \
		-e HTTPORIGIN=${v_httporigin} \
		-e SESSION_SECRET=${v_session_secret} \
		-e SESSION_TTL=${v_session_ttl} \
		-e NFT_ENABLED=${v_nft_enabled} \
		-e NFT_DEFAULT_POLICY=${v_nft_default_policy} \
		-e DEV_HOSTNAME=${v_dev_hostname} \
		-e DEV_AUTHIP=${v_dev_authip} \
		--network host \
		--cap-add NET_ADMIN \
		--volume /usr/bin/wg:/usr/bin/wg \
		--volume $(shell pwd)/${v_db_path}:/wgnetwork.db \
		--name wgnetwork \
		-d zyablitsev/wgnetwork

docker-rm-dev:
	docker rm -f wgnetwork 2>/dev/null || true

set-cap-dev:
	setcap cap_net_admin,cap_net_bind_service+eip ${v_bin_dir}/wgnetwork

run-service-dev:
	${v_env} ./bin/wgnetwork

test-envconfig:
	go clean -testcache && \
		${TESTCMD} ./pkg/envconfig/ -run Test

test-model:
	go clean -testcache && \
		${TESTCMD} ./model/ -run Test

test-pkg-otp:
	go clean -testcache && \
		${TESTCMD} ./pkg/otp/ -run Test

test-ipset:
	go clean -testcache && \
		${TESTCMD} ./pkg/ipset/ -run Test

test-wgmngr:
	go clean -testcache && \
		${TESTCMD} ./pkg/wgmngr/ -run Test

test-pretty:
	go clean -testcache && \
		${TESTCMD} ./pkg/pretty/ -run Test

test: test-envconfig \
	test-model \
	test-pkg-otp \
	test-ipset \
	test-wgmngr
