# Hydra and Kratos Self-Hosted Servers

This guide provides instructions to start the Hydra and Kratos self-hosted servers along with all required services using Docker Compose.

## Getting Started

To start all the services on your local Docker machine, run the following command:

```bash
cd contrib/hydra/
docker-compose up --build
```

This will start all necessary services for Hydra and Kratos.

Creating a Hydra Client for OAuth Flow
Once all services are running, you will need to create a new Hydra client to perform the OAuth flow using Kratos as an Identity Provider (IDP). Follow the steps below to set up the client and initiate the authorization flow.

## Step 1: Create a Hydra Client
Run the following command to create a new client using Hydra:

```bash
cd contrib/hydra/
code_client=$(docker-compose exec hydra \
    hydra create client \
    --endpoint http://127.0.0.1:4445 \
    --grant-type authorization_code,refresh_token \
    --response-type code,id_token \
    --format json \
    --scope openid --scope offline \
    --redirect-uri http://127.0.0.1:5555/callback)
```


## Step 2: Extract Client ID and Secret
After the client is created, extract the client_id and client_secret using the following commands:

```bash
cd contrib/hydra/
code_client_id=$(echo $code_client | jq -r '.client_id')
code_client_secret=$(echo $code_client | jq -r '.client_secret')
```


## Step 3: Perform Authorization Code Flow
To initiate the authorization code flow, run the following command:

```bash
cd contrib/hydra/
docker-compose exec hydra \
    hydra perform authorization-code \
    --client-id $code_client_id \
    --client-secret $code_client_secret \
    --endpoint http://127.0.0.1:4444/ \
    --port 5555 \
    --scope openid --scope offline
```

## Step 4: Complete Authentication Flow

Open a browser and go to http://localhost:5555. This will begin the authentication process using Kratos as the IDP.




