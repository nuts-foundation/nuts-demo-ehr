echo Generating CA
openssl ecparam -genkey -name prime256v1 -noout -out ca.key
openssl req -x509 -new -nodes -key ca.key -sha256 -days 1825 -out ca.pem -subj "/CN=Nuts demo CA/O=Nuts/C=NL"

# Generate key/certs for all domains by calling issue-cert.sh
./issue-cert.sh left.local
./issue-cert.sh node.left.local
./issue-cert.sh right.local
./issue-cert.sh node.right.local

