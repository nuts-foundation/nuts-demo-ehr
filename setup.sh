#!/usr/bin/env bash

echo "----------------------------------------"
echo "Removing data and restarting nodes..."
echo "----------------------------------------"
# stop all containers so we can delete all data
docker compose down
sleep 0.5 # If containers fail to restart below, make this longer

# delete all data
rm -r ./docker-compose/{left,right}/data/*/*
rm -r ./docker-compose/{left,right}/config/demo/customers.json
touch ./docker-compose/{left,right}/config/demo/customers.json # or docker will create directories for these mounts during startup

docker compose up --wait

echo "----------------------------------------"
echo "Creating DIDs..."
echo "----------------------------------------"
DID_LEFT=$(docker exec nuts-demo-ehr-node-left-1 curl -sS -X POST "http://localhost:8081/internal/vdr/v2/did" | jq -r .id)
DID_RIGHT=$(docker exec nuts-demo-ehr-node-right-1 curl -sS -X POST "http://localhost:8081/internal/vdr/v2/did" | jq -r .id)
echo "DID_LEFT: $DID_LEFT"
echo "DID_RIGHT: $DID_RIGHT"

echo "----------------------------------------"
echo "Issuing NutsOrganizationCredentials..."
echo "----------------------------------------"
# issue Left
REQUEST="{\"type\":\"NutsOrganizationCredential\",\"issuer\":\"${DID_LEFT}\", \"credentialSubject\": {\"id\":\"${DID_LEFT}\", \"organization\":{\"name\":\"Left\", \"city\":\"Enske\"}},\"withStatusList2021Revocation\": false}"
RESPONSE=$(echo $REQUEST | docker exec -i nuts-demo-ehr-node-left-1 curl -sS -X POST --data-binary @- http://localhost:8081/internal/vcr/v2/issuer/vc -H "Content-Type:application/json")
if echo $RESPONSE | grep -q "VerifiableCredential"; then
  echo "NutsOrganizationCredential issued for Left"
else
  echo "FAILED: Could not issue NutsOrganizationCredential for Left" 1>&2
  echo $RESPONSE
fi

# add to wallet Left
RESPONSE=$(echo $RESPONSE | docker exec -i nuts-demo-ehr-node-left-1 curl -sS -X POST --data-binary @- http://localhost:8081/internal/vcr/v2/holder/${DID_LEFT}/vc -H "Content-Type:application/json")
if [[ $RESPONSE -eq "" ]]; then
  echo "VC stored in wallet"
else
  echo "FAILED: Could not load NutsOrganizationCredential for Left" 1>&2
  echo $RESPONSE
fi

# issue Right
REQUEST="{\"type\":\"NutsOrganizationCredential\",\"issuer\":\"${DID_RIGHT}\", \"credentialSubject\": {\"id\":\"${DID_RIGHT}\", \"organization\":{\"name\":\"Right\", \"city\":\"Enske\"}},\"withStatusList2021Revocation\": false}"
RESPONSE=$(echo $REQUEST | docker exec -i nuts-demo-ehr-node-right-1 curl -sS -X POST --data-binary @- http://localhost:8081/internal/vcr/v2/issuer/vc -H "Content-Type:application/json")
if echo $RESPONSE | grep -q "VerifiableCredential"; then
  echo "NutsOrganizationCredential issued for Right"
else
  echo "FAILED: Could not issue NutsOrganizationCredential for Right" 1>&2
  echo $RESPONSE
fi

# add to wallet Right
RESPONSE=$(echo $RESPONSE | docker exec -i nuts-demo-ehr-node-right-1 curl -sS -X POST --data-binary @- http://localhost:8081/internal/vcr/v2/holder/${DID_RIGHT}/vc -H "Content-Type:application/json")
if [[ $RESPONSE -eq "" ]]; then
  echo "VC stored in wallet"
else
  echo "FAILED: Could not load NutsOrganizationCredential for Left" 1>&2
  echo $RESPONSE
fi

echo "----------------------------------------"
echo "Registering DIDs on Discovery Service..."
echo "----------------------------------------"
SERVICE="urn:nuts.nl:usecase:eOverdrachtDemo2024"
docker exec nuts-demo-ehr-node-left-1 curl -sS -X POST http://localhost:8081/internal/discovery/v1/${SERVICE}/${DID_LEFT}
docker exec nuts-demo-ehr-node-right-1 curl -sS -X POST http://localhost:8081/internal/discovery/v1/${SERVICE}/${DID_RIGHT}

echo "----------------------------------------"
echo "Adding services to DIDs..."
echo "----------------------------------------"
docker exec nuts-demo-ehr-node-left-1 curl -sS -X POST "http://localhost:8081/internal/vdr/v2/did/$DID_LEFT/service" -H  "Content-Type: application/json" -d "{\"type\": \"eOverdracht-sender\",\"serviceEndpoint\": {\"auth\": \"https://node.left.local/oauth2/$DID_LEFT/authorize\",\"fhir\": \"https://left.local/fhir/1\"}}"
docker exec nuts-demo-ehr-node-left-1 curl -sS -X POST "http://localhost:8081/internal/vdr/v2/did/$DID_LEFT/service" -H  "Content-Type: application/json" -d "{\"type\": \"eOverdracht-receiver\",\"serviceEndpoint\": {\"auth\": \"https://node.left.local/oauth2/$DID_LEFT/authorize\",\"notification\": \"https://left.local/web/external/transfer/notify\"}}"
docker exec nuts-demo-ehr-node-right-1 curl -sS -X POST "http://localhost:8081/internal/vdr/v2/did/$DID_RIGHT/service" -H  "Content-Type: application/json" -d "{\"type\": \"eOverdracht-sender\",\"serviceEndpoint\": {\"auth\": \"https://node.right.local/oauth2/$DID_RIGHT/authorize\",\"fhir\": \"https://right.local/fhir/1\"}}"
docker exec nuts-demo-ehr-node-right-1 curl -sS -X POST "http://localhost:8081/internal/vdr/v2/did/$DID_RIGHT/service" -H  "Content-Type: application/json" -d "{\"type\": \"eOverdracht-receiver\",\"serviceEndpoint\": {\"auth\": \"https://node.right.local/oauth2/$DID_RIGHT/authorize\",\"notification\": \"https://right.local/web/external/transfer/notify\"}}"

echo "----------------------------------------"
echo "Creating customers.json for demo-ehr..."
echo "----------------------------------------"
printf "{\n\t\"1\":{\"active\":false,\"city\":\"Enske\",\"did\":\"$DID_LEFT\",\"domain\":\"\",\"id\":1,\"name\":\"Left\"}\n}\n" > ./docker-compose/left/config/demo/customers.json
printf "{\n\t\"1\":{\"active\":false,\"city\":\"Enske\",\"did\":\"$DID_RIGHT\",\"domain\":\"\",\"id\":1,\"name\":\"Right\"}\n}\n" > ./docker-compose/right/config/demo/customers.json

docker compose down # at the minimum a restart is needed to load the new customers.json file