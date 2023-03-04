# Golang based REST API

Expected Key-Value struct:
{
  "key": "{keyValue}",
  "value": "{value}"
}

## Testing
To simplify the build, deploy, and test process, a docker-compose.yml file and Makefile were provied.

Aside from a docker installation being present on the system, there should be very little required from the user.

Run `make all` to perform testing. If specific Key-Value pairs to test are desired, there are jsons present in the testingJsons directory at the repository root.

test_set.json contains the Key-Value json structures to be created on the endpoint first and validated that the POST endpoint at `/api/keys` is functioning as expected.

test_overwrite.json uses a similar structure as test_set.json. However, it verifies first that the keys provided are already present on the server and that the value overwriting the old is different. To test this, build a key-value struct which uses a key which was created during by test_set.json, but with a different value.

test_delete.json is used to test the delete function against the endpoing. Inside is a json array of key values which are to be deleted from the server. The delete process first validates that the key is present before the delete call and that it is removed afterwards.

`make all` should build the server and client, prepare the build directory, and call `docker-compose up` to spin up two containers.
The first container, https-server, runs the server in the background. The second, https-client, runs the client tests against the server.

### Assumptions

The API should be lightweight, secure, and easy to test. While it is currently only intended to run on a single Docker container locally, it should ideally be able to deploy and scale horizontally if necessary, while handling all kinds of user data. In that case, TLS should be used to encrypt and decrypt request and responses.

For the scope of this, I've implemented a basic set of scripts to generate a server self-signed certificate and have it certified for the client. In the long term, implementing a certificate authority on the server and building a cert pool on the client should be considered, but are out of scope for this initial implementation.

Given that this API is able to read, write, and delete data on the server, it should be account protected, so a basic login endpoint was implemented to generate and return a JSON Web Token. JWTs do have the risk of being reused and not being tied to users, but the secret and expiration time. However, they allow for shared authentication among any services that share the same secret and wouldn't require an additional Account endpoint added.

Separately, there is only the admin account {"admin", "admin"}, but the Key-Value endpoint can be repurposed to maintain accounts and access leves.


### Design decisions

Chose to use golang, as the java frameworks I've worked with tend to be much more complicated in setup, and significantly less flexible. 

While spring boot would be able to handle both without significant issue, I've found that spring boot and spring boot framework's reliance on components and services would make initial development (such as for POC work) too inefficient.

After some work with golang's net/http library on the server, I changed over to gin-gonic to build the server's endpoints, the reason being the quality of life it grants in writing simple to understand endpoint handlers, as opposed to net/http, which is notably more complicated.

This is due to gin-gonic acting as a shim for many calls to the net/http library, but the abstraction helps significantly with development, although the performance impact hasn't been determined yet.

Shared constants and structures were put into src/helper/common.go to prevent repetition and allow for easy modification of common use variables if needed.

All JWT processing is in the jwtMiddleware package. However, as all authentication side code is in the middleware package, it would be possible to expand to other forms of authentication (such as basic authentication) without needing significant rewrites.

With the client side, I still chose to use net/http as client functions as a test service and is less likely to see significant and frequent overhauls, but additions, as expected behaviour should remain constant. However, there was an effort to build standard methods to perform the GET, POST, and DELETE calls to prevent unnecessary repetition.
