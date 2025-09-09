# Titan
Wee CLI to help working with micro-frontends and/or micro-services. Its two main functions are:

- It allows to have your apps under a reverse proxy
- It has some functionalities to work with the repositories, like pulling/fetching or building
the apps

## Roadmap

Repository actions
[x] - Fetch & Pull from remote on configured repositories
[x] - Build - run build command on configured repositories
[x] - Clean - remove node_modules and dist folders on configured repositories
[x] - Install - run install command on configured repositories

Proxy
[x] - Start a reverse proxy to serve several apps under the same host
[x] - Start tasks in parallel to the reverse proxy (like apps we want proxied)

## Configuration
You can find info about how to configure titan [here](./docs/configuration.md).

## Usage

### Repository actions
Example usages for repository actions. The explanation takes into consideration the default action
if none is provided via config. Via Titan's [configuration](./docs/configuration.md) we can
change the defaul action on each action if desired.

**fetch**
Runs `git fetch -p && git pull` and `git fetch --tags --force && git fetch --prune --prune-tags`
on the configured repositories

```bash
./titan fetch -c /path/to/config/file.yaml
```

**clean**
Runs `rm -rf`, recursively, to remove `dist` and `node_modules` folders on the configured

repositories
```bash
./titan clean -c /path/to/config/file.yaml
```

**install**
Runs `pnpm install --frozen-lockfile --prefer-offline` on the configured repositories

```bash
./titan install -c /path/to/config/file.yaml
```

**build**
Runs `pnpm run build:local` on the configured repositories

```bash
./titan build -c /path/to/config/file.yaml
```

### Proxy
Example usage to use the proxy

**serve**
It starts the proxy using the profile passed as argument. Via Titan's [configuration](./docs/configuration.md)
we determine the profile, what routes Titan must proxy and to where and, if desired, tasks that
would be run alongside the proxy server. Example tasks could be running the apps we want to proxy.

```bash
./titan serve -c /path/to/config/file.yaml -p local:all
```
