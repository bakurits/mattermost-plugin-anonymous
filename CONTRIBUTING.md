Fork the repository to your own account and then clone it to a directory outside of $GOPATH matching your plugin name:

git clone https://github.com/owner/mattermost-plugin-anonymous
Note that this project uses Go modules. Be sure to locate the project outside of $GOPATH, or allow the use of Go modules within your $GOPATH with an export GO111MODULE=on.

To build your plugin use make

Use make check-style to check the style.

Use make debug-dist and make debug-deploy in place of make dist and make deploy to configure webpack to generate unminified Javascript.

make will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

dist/com.example.my-plugin.tar.gz
Alternatively you can deploy a plugin automatically to your server, but it requires login credentials:

export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
or configuration of a personal access token:

export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
In production, deploy and upload your plugin via the System Console.
