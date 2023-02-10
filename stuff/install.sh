#!/bin/bash

set -eu

VERSION="v0.0.1"
BASE_URL="https://github.com/zyablitsev/wgnetwork/releases/download/${VERSION}"

v_wgn_bootstrap_trust_ipset_url=${WGN_BOOTSTRAP_TRUST_IPSET_URL:-"${BASE_URL}/bootstrap-trust-ipset_linux_amd64"}
v_wgn_managercli_url=${WGN_MANAGERCLI_URL:-"${BASE_URL}/managercli_linux_amd64"}
v_wgn_service_url=${WGN_SERVICE_URL:-"${BASE_URL}/service_linux_amd64"}
v_wgn_bootstrap_trust_ipset_binary="/usr/local/bin/wgn_bootstrap-trust-ipset"
v_wgn_managercli_binary="/usr/local/bin/wgn_managercli"
v_wgn_service_binary="/usr/local/bin/wgnetwork"
v_wgn_service_working_dir="/usr/local/"
v_wgn_service_systemd_env_file="/etc/default/wgnetwork"
v_wgn_service_file="/lib/systemd/system/wgnetwork.service"

v_wgn_db_dir="/usr/local/boltdb/"
v_wgn_db_path="/usr/local/boltdb/wgnetwork.db"
v_wgn_fe_http_port="80"
v_wgn_api_http_port="8080"
v_wgn_name="wgnetwork"

v_wgn_user=${WGN_USER:-"wgnetwork"}

v_sysctl_conf_file="/etc/sysctl.conf"

v_trust_ip="`last -1w | grep $USER | awk '{ print $3 }'`"
v_ssh_port="`cat /etc/ssh/sshd_config 2>/dev/null | grep Port | awk '{ print $2 }' | head -n 1`"
v_wgn_session_secret="`cat /dev/urandom | tr -dc '[:alpha:]' | fold -w ${1:-20} | head -n 1`"

v_wg_port=${WG_PORT:-"51820"}
v_wg_cidr=${WG_CIDR:-"172.16.0.1/24"}

read -p "Enter IP to allow ssh access (default '$v_trust_ip'): " -i $v_trust_ip -e v_trust_ip
read -p "Enter Network Name (default '$v_wgn_name'): " -i $v_wgn_name -e v_wgn_name
read -p "Enter Wireguard Port (default '$v_wg_port'): " -i $v_wg_port -e v_wg_port
read -p "Enter Wireguard CIDR (default '$v_wg_cidr'): " -i $v_wg_cidr -e v_wg_cidr
read -p "Enter Web Management Interface Port (default '$v_wgn_fe_http_port'): " -i $v_wgn_fe_http_port -e v_wgn_fe_http_port
read -p "Enter API Management Interface Port (default '$v_wgn_api_http_port'): " -i $v_wgn_api_http_port -e v_wgn_api_http_port

v_program="${0##*/}"
v_self="$(readlink -f "${BASH_SOURCE[0]}")"
[[ $UID == 0 ]] || exec sudo -p "$v_program must be run as root. Please enter the password for %u to continue: " -- "$BASH" -- "$v_self"

# update system packages
apt-get update -y
apt-get upgrade -y

# install required system packages
apt-get install -y ca-certificates curl nftables wireguard

# configure ip forward
sed -i -e '/^#net.ipv4.ip_forward/s/^.*$/net.ipv4.ip_forward=1/' \
    $v_sysctl_conf_file
sysctl -p

# add system user
useradd --system \
    -M \
    --user-group \
    --shell /sbin/nologin \
    $v_wgn_user

# download binaries
curl -L -o $v_wgn_bootstrap_trust_ipset_binary "$v_wgn_bootstrap_trust_ipset_url"
chmod 700 $v_wgn_bootstrap_trust_ipset_binary
chown $v_wgn_user:$v_wgn_user $v_wgn_bootstrap_trust_ipset_binary

curl -L -o $v_wgn_managercli_binary "$v_wgn_managercli_url"
chmod 700 $v_wgn_managercli_binary
chown $v_wgn_user:$v_wgn_user $v_wgn_managercli_binary

curl -L -o $v_wgn_service_binary "$v_wgn_service_url"
chmod 700 $v_wgn_service_binary
chown $v_wgn_user:$v_wgn_user $v_wgn_service_binary
setcap cap_net_admin,cap_net_bind_service+eip $v_wgn_service_binary

# configure systemd
cat <<EOF > $v_wgn_service_systemd_env_file
LOG_LEVEL="info"
DB_PATH="$v_wgn_db_path"
WG_BINARY="`which wg`"
WG_PORT="$v_wg_port"
WG_CIDR="$v_wg_cidr"
FE_HTTP_PORT="$v_wgn_fe_http_port"
API_HTTP_PORT="$v_wgn_api_http_port"
API_UNIX_SOCKET="/tmp/wgmanager.sock"
OTP_ISSUER="$v_wgn_name"
SESSION_SECRET="$v_wgn_session_secret"
SESSION_TTL="5m"
NFT_ENABLED="true"
NFT_DEFAULT_POLICY="drop"
NFT_TRUST_PORTS="${v_ssh_port}"
EOF

chown root:root $v_wgn_service_systemd_env_file
chmod 0644 $v_wgn_service_systemd_env_file

cat <<EOF > $v_wgn_service_file
[Unit]
Description=WGNetworkService
Wants=network-online.target
After=network-online.target
AssertFileIsExecutable=$v_wgn_service_binary

[Service]
WorkingDirectory=$v_wgn_service_working_dir

User=$v_wgn_user
Group=$v_wgn_user

EnvironmentFile=$v_wgn_service_systemd_env_file

ExecStart=$v_wgn_service_binary

# Let systemd restart this service always
Restart=always

# Specifies the maximum file descriptor number that can be opened by this process
LimitNOFILE=65536

# Disable timeout logic and wait until process is stopped
TimeoutStopSec=infinity
SendSIGKILL=no

[Install]
WantedBy=multi-user.target
EOF

chown root:root $v_wgn_service_file
chmod 0644 $v_wgn_service_file

# bootstrap db
mkdir $v_wgn_db_dir
chown $v_wgn_user:$v_wgn_user $v_wgn_db_dir
chmod 700 $v_wgn_db_dir

$v_wgn_bootstrap_trust_ipset_binary -dbpath="$v_wgn_db_path" -trustip="$v_trust_ip"
chown $v_wgn_user:$v_wgn_user $v_wgn_db_path

# run
systemctl enable nftables.service
systemctl enable wgnetwork
systemctl start nftables.service
systemctl start wgnetwork
