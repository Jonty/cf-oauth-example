#cf-oauth-example
----------------

This is an example cloudfoundry UAAC oAuth client written to test out authentication and use of the CloudFoundry API. It was based off [this example](https://github.com/cloudfoundry-community/cf-go-client-example) but heavily modified and updated to autoconfigure and use the actual CF API once the token has been fetched.

Usage
-----

The [maifest.yml](manifest.yml) contains all the app configuration. For the app to function you'll need to create a client on the UAAC server and ensure that the manifest has the correct CLIENT_ID and CLIENT_SECRET set.

An example of creating the correct UAAC client for the example defaults to work is:

````uaac client add --scope openid,cloud_controller.read,cloud_controller.write,password.write,console.admin,console.support --authorized_grant_types authorization_code,client_credentials --refresh_token_validity 1209600 --secret c1oudc0w cf-oauth-example````

You should be able to simply ````cf push```` this app and visit the endpoint. It'll bounce you to the UAAC login page, at which point you can login and approve the app. Once authenticated you'll be bounced back to the app and shown a list of CF spaces you have permission to operate in (fetched from the API).
