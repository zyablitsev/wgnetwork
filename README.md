# WGNetwork. Managing a Private Secured Network

Tool for creating, configuring and managing [Wireguard](https://www.wireguard.com/) network and [NFTables](https://netfilter.org/projects/nftables/) traffic filtering system using web-interface.

* **Wireguard** — modern, fast and secure communication protocol.
* **NFTables** — Linux kernel subsystem providing filtering network connections. 

Access to the web management interface is provided within the created network. Access control of the device (wireguard peer) and the user with the role of "manager" associated with this device (time-based one-time password authentication).

---

- [Control interfaces](#control-interfaces)
    - [Web](#web)
    - [Cli-Tool](#cli-tool)
- [Getting started](#getting-started)
    - [Run in Docker](#run-in-docker)
    - [Installing manually into the system](#installing-manually-into-the-system)
    - [Using script for automating installation](#using-script-for-automating-installation)
    - [Build from sources](#build-from-sources)
- [Testing](#testing)

## Control interfaces

Create, modify, delete users, their devices, manage the list of ip-addresses allowed to access the server via ssh.

### Web

Access protocol: `http over tcp/ip`, connection encryption is provided by the Wireguard protocol.

Available by default at `http://172.16.0.1`, can be changed during installation.

#### Screenshots of web-interface pages

##### Authentication
![authentication](/stuff/screenshots/01-auth.png "auth")

##### Manage the list of ip-addresses allowed for remote access to the server via ssh
![server](/stuff/screenshots/02-server.png "server")

##### List of users
![users](/stuff/screenshots/03-users.png "users")

##### Creating a user
![user-create](/stuff/screenshots/04-user-create.png "user-create")

##### User information
![user](/stuff/screenshots/05-user.png "user")

##### List of devices
![devices](/stuff/screenshots/06-devices.png "devices")

##### Creating a device
![device-create](/stuff/screenshots/07-device-create.png "device-create")

##### Device information
![device](/stuff/screenshots/08-device.png "device")


### Cli-Tool

Access protocol: `http over unix-socket`

After installation in the system, the `wgn_managercli` command is available in the server console.

#### Arguments and parameters

##### Display the parameters of the wireguard server
```bash
~$ wgn_managercli wgcfg
  -unix-socket
    	path (default "/tmp/wgmanager.sock")
```

##### Adding a new user *(if you pass the parameter is_manager=true, the user will be created with the role of manager and a qr-code will be displayed to quickly import the totp key into the mobile device)*
```bash
~$ wgn_managercli user-create
  -is_manager
    	is manager flag (default "false")
  -name
    	name
  -unix-socket
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Editing user by `uuid`
```bash
~$ wgn_managercli user-edit
  -is_manager
    	is manager flag
  -name string
    	name
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
  -uuid string
    	user uuid
```

##### Deleting user and all associated devices by `uuid`
```bash
~$ wgn_managercli user-remove
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
  -uuid string
    	user uuid
```

##### Getting user information by `uuid`
```bash
~$ wgn_managercli user
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
  -uuid string
    	user uuid
```

##### Display complete list of users
```bash
~$ wgn_managercli users
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Adding a new user device *(if you do not pass the wg_pubkey parameter, the keys will be generated and a qr-code will be displayed to quickly import the wireguard-configuration into the mobile device)*
```bash
~$ wgn_managercli device-create
  -label string
    	label
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
  -user_uuid string
    	user uuid
  -wan_forward
    	allow ip forwarding (default "false")
  -wg_pubkey string
    	wireguard public key (optional)
```

##### Editing device by `ip`
```bash
~$ wgn_managercli device-edit
  -ip string
    	device ip
  -label string
    	label
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
  -wan_forward
    	wan_forward
  -wg_pubkey string
    	wireguard public key
```

##### Deleting device by `ip`
```bash
~$ wgn_managercli device-remove
  -ip string
    	device ip
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Getting device information by `ip`
```bash
~$ wgn_managercli device
  -ip string
    	device ip
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Display complete list of devices
```bash
~$ wgn_managercli devices
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Adding an ip-address to the list of permissions for remote access to the server via ssh
```bash
~$ wgn_managercli trust-ipset-add
  -ip string
    	device ip
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Removing an ip-address from the list of permissions for remote access to the server via ssh
```bash
~$ wgn_managercli trust-ipset-remove
  -ip string
    	device ip
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Display full list of ip-addresses allowed for remote access to the server via ssh
```bash
~$ wgn_managercli trust-ipset
  -unix-socket string
    	unix-socket (default "/tmp/wgmanager.sock")
```

##### Database initialization with the list of ip-addresses allowed for remote access to the server via ssh. *(used in the server installation process before the first start of the service)*
```bash
~$ wgn_bootstrap-trust-ipset
  -dbpath string
    	dbpath
  -trustip value
    	device ip
```

## Getting started

**Requirements:**
* Wireguard
* NFTables

the service can be run in a docker container or installed to run on the system.

### Run in Docker

1. download tool to initialize the service database
```bash
~$ curl -L -o ./wgn_bootstrap-trust-ipset "https://github.com/zyablitsev/wgnetwork/releases/download/v0.0.1/wgn-bootstrap-trust-ipset_linux_amd64"
~$ chmod +x ./wgn_bootstrap-trust-ipset
```

2. initialize the database with your ip-address, which will be added to the list of allowed remote access via ssh protocol when you start the service
```bash
~$ mkdir /usr/local/boltdb/
~$ TRUSTIP=`last -1w | grep $USER | awk '{ print $3 }'`
~$ ./wgn_bootstrap-trust-ipset -dbpath="/usr/local/boltdb/wgnetwork.db" -trustip="$TRUSTIP"
```

3. turn on ip_forward
```bash
~$ sed -i -e '/^#net.ipv4.ip_forward/s/^.*$/net.ipv4.ip_forward=1/' /etc/sysctl.conf
~$ sysctl -p
```

4. start the service container
```bash
~$ SESSION_SECRET=`cat /dev/urandom | tr -dc '[:alpha:]' | fold -w ${1:-20} | head -n 1`
~$ docker run \
    -e LOG_LEVEL="info" \
    -e DB_PATH="/wgnetwork.db" \
    -e WG_BINARY="/usr/bin/wg" \
    -e WG_PORT="51820" \
    -e WG_CIDR="172.16.0.1/24" \
    -e FE_HTTP_PORT="80" \
    -e API_HTTP_PORT="8080" \
    -e OTP_ISSUER="wgnetwork" \
    -e SESSION_SECRET="$SESSION_SECRET" \
    -e SESSION_TTL="5m" \
    -e NFT_ENABLED="true" \
    -e NFT_DEFAULT_POLICY="drop" \
    -e NFT_TRUST_PORTS="22" \
    --network host \
    --cap-add NET_ADMIN \
    --volume /usr/bin/wg:/usr/bin/wg \
    --volume /usr/local/boltdb/wgnetwork.db:/wgnetwork.db \
    --restart always \
    --name wgnetwork \
    -d zyablitsev/wgnetwork
```

5. create the first user with the role of "manager" and register the device

---

**IMPORTANT**: access to the management web-interface is possible only from the devices of users with the role of "manager"

---
```bash
~$ docker exec wgnetwork \
    /wgn_managercli user-create -name="admin" -is_manager="true"
```
scan the qr-code into your authentication application (e.g. Google Authenticator), the totp code is required to authenticate the user in the management interface.

```bash
~$ docker exec wgnetwork \
    /wgn_managercli device-create --label="mobile" --user_uuid="INSERT_VALUE"
```
the configuration for your device will be generated, add it to your Wireguard client.

Activate the tunnel created in wireguard and you will be able to access the management web interface using totp code from the authentication program to authorize at `http://172.16.0.1`

### Installing manually into the system

1. install required system packages
```bash
~$ apt-get update -y
~$ apt-get upgrade -y
~$ apt-get install -y ca-certificates curl nftables wireguard
```

2. turn on ip_forward
```bash
~$ sed -i -e '/^#net.ipv4.ip_forward/s/^.*$/net.ipv4.ip_forward=1/' /etc/sysctl.conf
~$ sysctl -p
```

3. add system user wgnetwork
```bash
~$ useradd --system \
    -M \
    --user-group \
    --shell /sbin/nologin \
    wgnetwork
```

4. download binaries
```bash
~$ curl -L -o /usr/local/bin/wgn_bootstrap-trust-ipset "https://github.com/zyablitsev/wgnetwork/releases/download/v0.0.1/wgn-bootstrap-trust-ipset_linux_amd64"
~$ chmod +x /usr/local/bin/wgn_bootstrap-trust-ipset
~$ chown wgnetwork:wgnetwork /usr/local/bin/wgn_bootstrap-trust-ipset

~$ curl -L -o /usr/local/bin/wgn_managercli "https://github.com/zyablitsev/wgnetwork/releases/download/v0.0.1/wgn-managercli_linux_amd64"
~$ chmod +x /usr/local/bin/wgn_managercli
~$ chown wgnetwork:wgnetwork /usr/local/bin/wgn_managercli

~$ curl -L -o /usr/local/bin/wgnetwork "https://github.com/zyablitsev/wgnetwork/releases/download/v0.0.1/wgnetwork_linux_amd64"
~$ chmod +x /usr/local/bin/wgnetwork
~$ chown wgnetwork:wgnetwork /usr/local/bin/wgnetwork
~$ setcap cap_net_admin,cap_net_bind_service+eip /usr/local/bin/wgnetwork
```

5. create service environment variables configuration
```bash
~$ WG_BINARY=`which wg`
~$ SESSION_SECRET=`cat /dev/urandom | tr -dc '[:alpha:]' | fold -w ${1:-20} | head -n 1`
~$ cat <<EOF > /etc/default/wgnetwork
LOG_LEVEL="info"
DB_PATH="/usr/local/boltdb/wgnetwork.db"
WG_BINARY="$WG_BINARY"
WG_PORT="51820"
WG_CIDR="172.16.0.1/24"
FE_HTTP_PORT="80"
API_HTTP_PORT="8080"
API_UNIX_SOCKET="/tmp/wgmanager.sock"
OTP_ISSUER="wgnetwork"
SESSION_SECRET="$SESSION_SECRET"
SESSION_TTL="5m"
NFT_ENABLED="true"
NFT_DEFAULT_POLICY="drop"
NFT_TRUST_PORTS="22"
EOF

~$ chown root:root /etc/default/wgnetwork
~$ chmod 0644 /etc/default/wgnetwork
```

6. create a systemd service configuration description
```bash
~$ cat <<EOF > /lib/systemd/system/wgnetwork.service
[Unit]
Description=WGNetworkService
Wants=network-online.target
After=network-online.target
AssertFileIsExecutable=/usr/local/bin/wgnetwork

[Service]
WorkingDirectory=/usr/local/

User=wgnetwork
Group=wgnetwork

EnvironmentFile=/etc/default/wgnetwork

ExecStart=/usr/local/bin/wgnetwork

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

~$ chown root:root /lib/systemd/system/wgnetwork.service
~$ chmod 0644 /lib/systemd/system/wgnetwork.service
```

7. initialize the database with your ip-address, which will be added to the list of allowed remote access via ssh protocol when you start the service
```bash
~$ mkdir /usr/local/boltdb/
~$ TRUSTIP=`last -1w | grep $USER | awk '{ print $3 }'`
~$ wgn_bootstrap-trust-ipset -dbpath="/usr/local/boltdb/wgnetwork.db" -trustip="$TRUSTIP"
~$ chown wgnetwork:wgnetwork /usr/local/boltdb
~$ chmod 0700 /usr/local/boltdb
~$ chown wgnetwork:wgnetwork /usr/local/boltdb/wgnetwork.db
~$ chmod 0600 /usr/local/boltdb/wgnetwork.db
```

8. run service
```bash
~$ systemctl enable nftables.service
~$ systemctl enable wgnetwork
~$ systemctl start nftables.service
~$ systemctl start wgnetwork
```

9. create the first user with the role of "manager" and register the device

---

**IMPORTANT**: access to the management web-interface is possible only from the devices of users with the role of "manager"

---
```bash
~$ wgn_managercli user-create -name="admin" -is_manager="true"
```
scan the qr-code into your authentication application (e.g. Google Authenticator), the totp code is required to authenticate the user in the management interface.

```bash
~$ wgn_managercli device-create --label="mobile" --user_uuid="INSERT_VALUE"
```
the configuration for your device will be generated, add it to your Wireguard client.

Activate the tunnel created in wireguard and you will be able to access the management web interface using totp code from the authentication program to authorize at `http://172.16.0.1`

### Using script for automating installation

**Limitations: Debian 11 (bullseye)**

1. open terminal and run:
```bash
~$ apt-get install -y ca-certificates curl
~$ bash <(curl -s "https://raw.githubusercontent.com/zyablitsev/wgnetwork/main/stuff/install.sh")
```

2. create the first user with the role of "manager" and register the device

---

**IMPORTANT**: access to the management web-interface is possible only from the devices of users with the role of "manager"

---
```bash
~$ wgn_managercli user-create -name="admin" -is_manager="true"
```
scan the qr-code into your authentication application (e.g. Google Authenticator), the totp code is required to authenticate the user in the management interface.

```bash
~$ wgn_managercli device-create -label="laptop" -user_uuid="INSERT_VALUE" -wan_forward="false"
```
the configuration for your device will be generated, add it to your Wireguard client.

Activate the tunnel created in wireguard and you will be able to access the management web interface using totp code from the authentication program to authorize at `http://172.16.0.1`

### Build from sources

**Requirements:**
* go 1.19+
* node 16.14+

1. clone repository
```bash
~$ git clone git@github.com:zyablitsev/wgnetwork.git
```

2. get dependencies
```bash
~$ make install-dependencies-fe
```

3. run build
```bash
~$ BIN_DIR=./bin/ make build
```

4. build docker-image
```bash
~$ make docker-build
```

## Testing

```bash
~$ make test
```
