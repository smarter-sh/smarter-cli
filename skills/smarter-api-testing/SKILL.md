---
name: smarter-api-testing
description: Use when you need to run the smarter cli, or curl, against a real Smarter platform by hand from the smarter-cli side - e.g. to see a route's real response before changing a command, to confirm a fix end to end, or to reproduce a user's bug. Explains why only the local platform, the throwaway api key, the throwaway $HOME and config.yaml, the cli's routes, what never to run, and how to clean up.
---

# Running the CLI Against the Local Platform

The unit tests use a mock server (`smarter-cli-testing`), so they can't tell
whether the platform really serves a route or returns the shape the cli
formats. Check that against the **local** platform in `../smarter`, whose dev
server is `http://localhost:9357` (the `smarter-app` container), and which the
cli reaches with `--environment local`.

Start it with `make run` in `../smarter` if
`curl -s -o /dev/null -w '%{http_code}' localhost:9357/` doesn't answer 302.
See the platform's `smarter-docker-environment` skill.

**Never** run the cli against `alpha`, `beta`, `next` or `prod`, and never with
the user's own `~/.smarter/config.yaml`.

## 1. Create a throwaway key

```console
docker exec smarter-app bash -c 'cd /home/smarter_user/smarter && python manage.py shell -c "
from smarter.lib.drf.models import SmarterAuthToken
from smarter.apps.account.utils import smarter_cached_objects as s
from smarter.apps.account.models import UserProfile
u = s.smarter_admin
up = UserProfile.objects.filter(user=u).first()
_, key = SmarterAuthToken.objects.create(user=u, user_profile=up, name=\"cli_probe\", description=\"DELETE ME\", is_active=True)
print(\"KEY=\" + key)
"' 2>/dev/null | grep '^KEY=' | cut -d= -f2 > "$SCRATCHPAD/key.txt"
```

Keep the key in your scratchpad, never in the repository. The admin user is a
superuser; for what an ordinary user sees, create the key for a non-staff user.

## 2. A throwaway home

```console
mkdir -p "$SCRATCHPAD/home/.smarter"
printf 'config:\n  environment: local\n  root_domain: platform.smarter.sh\n  output_format: yaml\nlocal:\n  api_key: "%s"\n' \
  "$(cat "$SCRATCHPAD/key.txt")" > "$SCRATCHPAD/home/.smarter/config.yaml"
go build -o "$SCRATCHPAD/smarter" .
```

Writing the file directly keeps the key out of your shell history. On a fresh
`$HOME`, `HOME=... smarter configure --api_key <key> --environment local` works
too.

## 3. Run commands

```console
HOME="$SCRATCHPAD/home" "$SCRATCHPAD/smarter" whoami
HOME="$SCRATCHPAD/home" "$SCRATCHPAD/smarter" get llmclients
HOME="$SCRATCHPAD/home" "$SCRATCHPAD/smarter" manifest secret -o yaml > "$SCRATCHPAD/secret.yaml"
HOME="$SCRATCHPAD/home" "$SCRATCHPAD/smarter" apply -f "$SCRATCHPAD/secret.yaml" --dry-run   # validates only
HOME="$SCRATCHPAD/home" "$SCRATCHPAD/smarter" describe secret <name>
```

Add `-v` to see the request and response, but don't paste `-v` output into a
summary, an issue or a file: it contains the api key (ADR-007).

To see the platform's raw response without the cli's formatting, use curl. All
cli routes are `POST`, with arguments in the query string:

```console
KEY=$(cat "$SCRATCHPAD/key.txt")
curl -s -X POST -H "Authorization: Token $KEY" "http://localhost:9357/api/v1/cli/get/secret/?i=5"
curl -s -X POST -H "Authorization: Token $KEY" "http://localhost:9357/api/v1/cli/describe/secret/?name=<name>"
curl -s -X POST -H "Authorization: Token $KEY" -H "Content-Type: application/json" \
  --data-binary @"$SCRATCHPAD/secret.yaml" "http://localhost:9357/api/v1/cli/validate/"
```

| Route (POST, under `/api/v1/cli/`)                     | cli command                                   |
| ------------------------------------------------------ | --------------------------------------------- |
| `apply/`, `validate/`                                  | `apply -f`, `validate -f`                     |
| `get/<Kind>/`, `describe/<Kind>/`                      | `get <kinds>`, `describe <kind>`              |
| `delete/<Kind>/`, `deploy/<Kind>/`, `undeploy/<Kind>/` | `delete`, `deploy`, `undeploy`                |
| `logs/<Kind>/`                                         | `logs`                                        |
| `example-manifest/<Kind>/`, `json-schema/<Kind>/`      | `manifest`, `json-schema`                     |
| `prompt/<name>/`, `prompt/config/<name>/`              | `chat prompt`, `chat config`                  |
| `whoami/`, `status/`, `version/`, `resources/`         | `whoami`, `status`, `version -v`, `resources` |

The authority is the platform's `openapi.json`, which the cli pins in
`testdata/contracts/openapi.json` (ADR-022), and
`../smarter/smarter/smarter/apps/api/v1/cli/urls.py`. Errors come back as
`{"error": {"description", "context", "errorClass", "stacktrace", ...}}`.

`chat prompt` sends a live, billable LLM request; that is accepted.

**Never run `deploy`** while probing unless the user asks: it creates real,
billable infrastructure (DNS, certificates, Kubernetes ingresses), even from
the local platform. Prefer `apply --dry-run`.

## 4. Clean up

Delete what you applied (`smarter delete ...`, or the web console), then the key:

```console
docker exec smarter-app bash -c 'cd /home/smarter_user/smarter && python manage.py shell -c "
from smarter.lib.drf.models import SmarterAuthToken
print(SmarterAuthToken.objects.filter(name=\"cli_probe\", description=\"DELETE ME\").delete())
"' 2>/dev/null
rm -rf "$SCRATCHPAD/home" "$SCRATCHPAD/key.txt"
```
