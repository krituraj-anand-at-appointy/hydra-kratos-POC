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



# Multiple Domian Single SignOn Demo

To verify if Single Sign-On (SSO) works across multiple domains using Hydra and Kratos integration, we will set up two different projects from the current repository named POC1 and POC2. Each project will have its own Hydra client.

## Step 1 : Create Hydra Clients

We will start by creating two OAuth2 clients for Hydra using the following requests:

```bash
curl --location 'http://127.0.0.1:4445/admin/clients' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Cookie: csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=coe3Z55OqL2b94fCaXvUXYnl5sPb7QiJu8gEdazIYJk=' \
--data '{
    "client_name": "Test OAuth2 Client 2",
    "client_secret": "secret",
    "grant_types": [
        "authorization_code",
        "refresh_token"
    ],
    "redirect_uris": [
        "http://app2.local:8081/callback"
    ],
    "post_logout_redirect_uris": [
        "http://app2.local:8081"
    ],
    "response_types": [
        "code",
        "id_token"
    ],
    "scope": "openid offline",
    "token_endpoint_auth_method": "client_secret_post"
}'

curl --location 'http://127.0.0.1:4445/admin/clients' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Cookie: csrf_token_be481debe9e1ebcf14d99f6f631d9a520ca6701ba0f3e4398508af30ebb1f509=coe3Z55OqL2b94fCaXvUXYnl5sPb7QiJu8gEdazIYJk=' \
--data '{
    "client_name": "Test OAuth2 Client 1",
    "client_secret": "secret",
    "grant_types": [
        "authorization_code",
        "refresh_token"
    ],
    "redirect_uris": [
        "http://app2.local:8080/callback"
    ],
    "post_logout_redirect_uris": [
        "http://app2.local:8080"
    ],
    "response_types": [
        "code",
        "id_token"
    ],
    "scope": "openid offline",
    "token_endpoint_auth_method": "client_secret_post"
}'
```

## Step 2: Update Client IDs

Once both clients are created, you will receive unique client IDs for each Hydra client. Update the corresponding client IDs in the oauthConfig section of the main.go files for each project:

Test OAuth2 Client 1 → POC1 main.go (oauthConfig)
Test OAuth2 Client 2 → POC2 main.go (oauthConfig)

## Step 3: Start Both Servers

Start the servers for both POC1 and POC2.

## Step4 : Configure Hosts File

Map the following domains to localhost in your system's hosts file (typically located at /etc/hosts):

127.0.0.1   app1.local
127.0.0.1   app2.local

## Step5 : Test SSO

Open app1.local:8080 and app2.local:8081 in your browser (e.g., Chrome). After logging in on the first domain, you will not need to log in again on the second domain, as Kratos will recognize the user session and handle the SSO seamlessly.

