# Antresol addon template

A starting point for an Antresol addon: a Go program the server runs on a
Unix socket, and a web bundle the app loads into its slots. The contract
between them and the server is `docs/addons/contract.md` in the core
repository; read it first.

## Layout

| Path | What |
| --- | --- |
| `antresol-addon.json` | the manifest without `platforms` and `web` (the build fills them in) |
| `main.go`, `addon.go` | the program on `addonkit`: `GET /health` (`addonkit.Ready`), `GET /api/hello` (`/api/addons/<id>/x/hello` in the app), an admin-only `GET /api/admin/hello` (`addonkit.RequireAdmin`) and the hooks (`addonkit.HookMux` at `/hooks/`) |
| `go.mod` | requires `github.com/antresol-app/antresol/addonkit`, the core's Go library for addons (MIT, standard library only): the environment, a client of the core API v1, hooks, the acting profile |
| `web/` | the web part: `src/index.tsx` calls `registerAddon`; built by Vite with the SDK's plugin |
| `tools/manifest` | writes the release's `antresol-addon.json` and checks each archive the way the server will unpack it |
| `.github/workflows/release.yml` | a tag `v<version>` builds and publishes the release |

## Make it yours

1. Pick an id (`^[a-z][a-z0-9-]{1,39}$`, not `web`, `preview`, `x` or `ui`) and
   put it in `antresol-addon.json` and in `ID` of `web/src/index.tsx`.
2. Set the name, description, author, homepage and the permissions you need
   (`events`, `library.read`).
3. Rename the module in `go.mod`.

## Build

The program builds on `addonkit` (`go.mod`'s require). The web part builds on
the SDK package `@antresol/addon-sdk`: `ANTRESOL_ADDON_SDK=<its packed .tgz>`,
or else `ANTRESOL_CORE=<a core checkout>`, whose `web/` has run `npm ci`; the
build packs the SDK from its `web/addon-sdk` and copies it into
`web/node_modules`.

```sh
make test                                    # go vet + go test
make build ANTRESOL_CORE=<core checkout>     # dist/: linux/amd64, linux/arm64, darwin/arm64, web, manifest
make local ANTRESOL_CORE=<core checkout>     # dist-local/<id>/: this machine only
```

Archives are made with plain `tar` (ustar, no extended attributes, no `._`
files). The server refuses links, pax global headers (so never `git archive`)
and anything but regular files and directories.

## Try it on a dev server

```sh
make local
VAULT_ADDONS_DEV=$PWD/dist-local bin/vault    # in the core checkout
```

Then in the app: «Управление» → «Расширения» → «Добавить расширение», enter
`local:<id>`. A rebuilt local release with the same version is not an update:
raise the version or remove the addon first.

## Playing videos

An addon that serves videos of its own registers them as inputs
(`Client.RegisterInput`, permission `play`) and sends the profile to
`/play/input/<id>`; see `addons/hello` for a working example (a sample
video route, `playerNotes` and the `play.progress` hook) and the contract's
section on inputs.

## Release

Push a tag `v<version>` matching `version` in `antresol-addon.json`; the
workflow attaches `dist/*` to the release. Installing takes the repository's
URL, `https://github.com/<owner>/<repo>`.

## Licence

Licence: MIT

The MIT licence (see `LICENSE`), Copyright (c) 2026 Antresol Contributors. The addon
SDK it builds on is MIT too; the Antresol server stays AGPL-3.0-or-later.
