#! /bin/bash

if [ -z "$KHULNASOFT_REGISTRY_USERNAME" ]; then echo "KHULNASOFT_REGISTRY_USERNAME env is unset" && exit 1; fi
if [ -z "$KHULNASOFT_REGISTRY_PASSWORD" ]; then echo "KHULNASOFT_REGISTRY_PASSWORD env is unset" && exit 1; fi
if [ -z "$KHULNASOFT_VERSION" ]; then echo "KHULNASOFT_VERSION env is unset" && exit 1; else echo "KHULNASOFT_VERSION env is set to '$KHULNASOFT_VERSION'"; fi

HARBOR_HOME="/opt/harbor"
HARBOR_PKI_DIR="/etc/harbor/pki/internal"
HARBOR_SCANNER_KHULNASOFT_VERSION="0.14.0"
SCANNER_UID=1000
SCANNER_GID=1000

mkdir -p $HARBOR_HOME/common/config/khulnasoft-adapter
mkdir -p /data/khulnasoft-adapter/reports
mkdir -p /data/khulnasoft-adapter/opt
mkdir -p /var/lib/khulnasoft-db/data

# Login to Khulnasoft registry.
echo $KHULNASOFT_REGISTRY_PASSWORD | docker login registry.khulnasoft.com \
  --username $KHULNASOFT_REGISTRY_USERNAME \
  --password-stdin

# Copy the scannercli binary from the registry.khulnasoft.com/scanner image.
docker run --rm --entrypoint "" \
  --volume $HARBOR_HOME/common/config/khulnasoft-adapter:/out registry.khulnasoft.com/scanner:$KHULNASOFT_VERSION \
  cp /opt/khulnasoft/scannercli /out

# Generate a private key.
openssl genrsa -out $HARBOR_PKI_DIR/khulnasoft_adapter.key 4096

# Generate a certificate signing request (CSR).
openssl req -sha512 -new \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=khulnasoft-adapter" \
  -key $HARBOR_PKI_DIR/khulnasoft_adapter.key \
  -out $HARBOR_PKI_DIR/khulnasoft_adapter.csr

# Generate an x509 v3 extension file.
cat > $HARBOR_PKI_DIR/khulnasoft_adapter_v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=khulnasoft-adapter
EOF

# Use the v3.ext file to generate a certificate for your Harbor host.
openssl x509 -req -sha512 -days 365 \
  -extfile $HARBOR_PKI_DIR/khulnasoft_adapter_v3.ext \
  -CA $HARBOR_PKI_DIR/harbor_internal_ca.crt \
  -CAkey $HARBOR_PKI_DIR/harbor_internal_ca.key \
  -CAcreateserial \
  -in $HARBOR_PKI_DIR/khulnasoft_adapter.csr \
  -out $HARBOR_PKI_DIR/khulnasoft_adapter.crt

chown $SCANNER_UID:$SCANNER_GID /data/khulnasoft-adapter/reports
chown $SCANNER_UID:$SCANNER_GID /data/khulnasoft-adapter/opt
chown $SCANNER_UID:$SCANNER_GID $HARBOR_HOME/common/config/khulnasoft-adapter/scannercli
chown $SCANNER_UID:$SCANNER_GID $HARBOR_PKI_DIR/khulnasoft_adapter.key
chown $SCANNER_UID:$SCANNER_GID $HARBOR_PKI_DIR/khulnasoft_adapter.crt

cat << EOF > $HARBOR_HOME/common/config/khulnasoft-adapter/env
SCANNER_LOG_LEVEL=debug
SCANNER_API_ADDR=:8443
SCANNER_API_TLS_KEY=/etc/pki/khulnasoft_adapter.key
SCANNER_API_TLS_CERTIFICATE=/etc/pki/khulnasoft_adapter.crt
SCANNER_KHULNASOFT_USERNAME=administrator
SCANNER_KHULNASOFT_PASSWORD=@Khulnasoft12345
SCANNER_KHULNASOFT_HOST=https://khulnasoft-console:8443
SCANNER_CLI_NO_VERIFY=true
SCANNER_KHULNASOFT_REGISTRY=Harbor
SCANNER_KHULNASOFT_USE_IMAGE_TAG=false
SCANNER_KHULNASOFT_REPORTS_DIR=/var/lib/scanner/reports
SCANNER_REDIS_URL=redis://redis:6379
SCANNER_CLI_OVERRIDE_REGISTRY_CREDENTIALS=false
EOF

cat << EOF > $HARBOR_HOME/docker-compose.override.yml
version: '2.3'
services:
  khulnasoft-adapter:
    networks:
      - harbor
    container_name: khulnasoft-adapter
    # image: docker.io/khulnasoft/harbor-scanner-khulnasoft:dev
    # image: docker.io/khulnasoft/harbor-scanner-khulnasoft:$HARBOR_SCANNER_KHULNASOFT_VERSION
    image: public.ecr.aws/khulnasoft/harbor-scanner-khulnasoft:$HARBOR_SCANNER_KHULNASOFT_VERSION
    restart: always
    cap_drop:
      - ALL
    depends_on:
      - redis
    volumes:
      - type: bind
        source: $HARBOR_PKI_DIR/khulnasoft_adapter.key
        target: /etc/pki/khulnasoft_adapter.key
      - type: bind
        source: $HARBOR_PKI_DIR/khulnasoft_adapter.crt
        target: /etc/pki/khulnasoft_adapter.crt
      - type: bind
        source: $HARBOR_HOME/common/config/khulnasoft-adapter/scannercli
        target: /usr/local/bin/scannercli
      - type: bind
        source: /data/khulnasoft-adapter/reports
        target: /var/lib/scanner/reports
      - type: bind
        source: /data/khulnasoft-adapter/opt
        target: /opt/khulnasoftscans
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "khulnasoft-adapter"
    env_file:
      $HARBOR_HOME/common/config/khulnasoft-adapter/env
  khulnasoft-db:
    networks:
      - harbor
    image: registry.khulnasoft.com/database:$KHULNASOFT_VERSION
    container_name: khulnasoft-db
    environment:
      - POSTGRES_PASSWORD=lunatic0
    volumes:
      - /var/lib/khulnasoft-db/data:/var/lib/postgresql/data
  khulnasoft-console:
    networks:
      - harbor
    ports:
      - 9080:8080
    image: registry.khulnasoft.com/console:$KHULNASOFT_VERSION
    container_name: khulnasoft-console
    environment:
      - ADMIN_PASSWORD=@Khulnasoft12345
      - SCALOCK_DBHOST=khulnasoft-db
      - SCALOCK_DBNAME=scalock
      - SCALOCK_DBUSER=postgres
      - SCALOCK_DBPASSWORD=lunatic0
      - SCALOCK_AUDIT_DBHOST=khulnasoft-db
      - SCALOCK_AUDIT_DBNAME=slk_audit
      - SCALOCK_AUDIT_DBUSER=postgres
      - SCALOCK_AUDIT_DBPASSWORD=lunatic0
      - KHULNASOFT_DOCKERLESS_SCANNING=1
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    depends_on:
      - khulnasoft-db
  khulnasoft-gateway:
    image: registry.khulnasoft.com/gateway:$KHULNASOFT_VERSION
    container_name: khulnasoft-gateway
    environment:
      - SCALCOK_LOG_LEVEL=DEBUG
      - KHULNASOFT_CONSOLE_SECURE_ADDRESS=khulnasoft-console:8443
      - SCALOCK_DBHOST=khulnasoft-db
      - SCALOCK_DBNAME=scalock
      - SCALOCK_DBUSER=postgres
      - SCALOCK_DBPASSWORD=lunatic0
      - SCALOCK_AUDIT_DBHOST=khulnasoft-db
      - SCALOCK_AUDIT_DBNAME=slk_audit
      - SCALOCK_AUDIT_DBUSER=postgres
      - SCALOCK_AUDIT_DBPASSWORD=lunatic0
    networks:
      - harbor
    depends_on:
      - khulnasoft-db
      - khulnasoft-console
EOF

cd /opt/harbor
docker-compose up --detach

# Use Harbor 2.0 REST API to register khulnasoft-adapter as an Interrogation Service.
cat << EOF > /tmp/khulnasoft-adapter.registration.json
{
  "name": "Khulnasoft Enterprise $KHULNASOFT_VERSION",
  "url": "https://khulnasoft-adapter:8443",
  "description": "Khulnasoft Enterprise $KHULNASOFT_VERSION vulnerability scanner."
}
EOF

curl --include \
  --user admin:Harbor12345 \
  --request POST \
  --header "accept: application/json" \
  --header "Content-Type: application/json" \
  --data-binary "@/tmp/khulnasoft-adapter.registration.json" \
  "http://localhost:8080/api/v2.0/scanners"
