# Anonymous Plugin [![CircleCI branch](https://img.shields.io/circleci/project/github/bakurits/mattermost-plugin-anonymous/key-pair-generation.svg)](https://circleci.com/gh/bakurits/mattermost-plugin-anonymous) [![codecov](https://codecov.io/gh/bakurits/mattermost-plugin-anonymous/branch/key-pair-generation/graph/badge.svg)](https://codecov.io/gh/bakurits/mattermost-plugin-anonymous)

This plugin serves as a starting point for writing a Mattermost plugin. Feel free to base your own plugin off this repository.

To learn more about plugins, see [our plugin documentation](https://developers.mattermost.com/extend/plugins/).

## Getting Started
Use GitHub's template feature to make a copy of this repository by clicking the "Use this template" button then clone outside of `$GOPATH`.

Alternatively shallow clone the repository to a directory outside of `$GOPATH` matching your plugin name:
```
git clone --depth 1 https://github.com/mattermost/mattermost-plugin-starter-template com.example.my-plugin
```

Note that this project uses [Go modules](https://github.com/golang/go/wiki/Modules). Be sure to locate the project outside of `$GOPATH`, or allow the use of Go modules within your `$GOPATH` with an `export GO111MODULE=on`.

Edit `plugin.json` with your `id`, `name`, and `description`:
```
{
    "id": "com.example.my-plugin",
    "name": "My Plugin",
    "description": "A plugin to enhance Mattermost."
}
```

Build your plugin:
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.example.my-plugin.tar.gz
```

There is a build target to automate deploying and enabling the plugin to your server, but it requires login credentials:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or configuration of a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```

Alternatively, if you are running your `mattermost-server` out of a sibling directory by the same name, use the `deploy` target alone to  unpack the files into the right directory. You will need to restart your server and manually enable your plugin.

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

## Q&A

### How do I make a server-only or web app-only plugin?

Simply delete the `server` or `webapp` folders and remove the corresponding sections from `plugin.json`. The build scripts will skip the missing portions automatically.

### How do I include assets in the plugin bundle?

Place them into the `assets` directory. To use an asset at runtime, build the path to your asset and open as a regular file:

```go
bundlePath, err := p.API.GetBundlePath()
if err != nil {
    return errors.Wrap(err, "failed to get bundle path")
}

profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "profile_image.png"))
if err != nil {
    return errors.Wrap(err, "failed to read profile image")
}

if appErr := p.API.SetProfileImage(userID, profileImage); appErr != nil {
    return errors.Wrap(err, "failed to set profile image")
}
```

### How do I build the plugin with unminified JavaScript?
Use `make debug-dist` and `make debug-deploy` in place of `make dist` and `make deploy` to configure webpack to generate unminified Javascript.
