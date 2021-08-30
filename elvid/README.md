# ElvID

## Test

- Setup following environment variables.

```text
source_env ${HOME}
export ELVID_ACCESS_TOKEN="ey...."
export ELVID_BASE_URL="https://elvid.test-elvia.io"
export ELVID_CLIENT_ID="00000000-0000-4000-8000-000000000000"
export ELVID_ID_TOKEN="ey...."
export ELVID_MACHINE_CLIENT_ID="00000000-0000-4000-8000-000000000000"
export ELVID_MACHINE_CLIENT_SECRET="...."
```

```bash
go test -v -tags=integration -timeout 0
```
