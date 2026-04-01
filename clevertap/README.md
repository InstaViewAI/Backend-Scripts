# CleverTap Script

Before running the script, set `CLEVERTAP_ACCOUNT` to the account key name from [accounts.json](/Users/shaharyaralam/Desktop/Instavision/Backend-Scripts/clevertap/accounts.json).

The account key decides two things:

1. Which CleverTap credentials are used.
2. Which input JSON file is loaded.

Example account config:

```json
{
  "accounts": {
    "instaview": {
      "base_url": "https://in1.api.clevertap.com",
      "account_id": "TEST-756-49R-867Z",
      "passcode": "UTC-AUC-CEEL"
    },
    "luna": {
      "base_url": "https://us1.api.clevertap.com",
      "account_id": "YOUR-SECOND-ACCOUNT-ID",
      "passcode": "YOUR-SECOND-PASSCODE"
    },
    "homiQ": {
      "base_url": "https://us1.api.clevertap.com",
      "account_id": "YOUR-THIRD-ACCOUNT-ID",
      "passcode": "YOUR-THIRD-PASSCODE"
    }
  }
}
```

Input file mapping:

- `instaview` -> `instaview_user.json`
- `luna` -> `luna_user.json`
- `homiQ` -> `homiQ_user.json`

## Run

Go to the CleverTap folder:

```bash
cd /Users/shaharyaralam/Desktop/Instavision/Backend-Scripts/clevertap
```

Run for `instaview`:

```bash
CLEVERTAP_ACCOUNT=instaview go run ./cmd
```

Run for `luna`:

```bash
CLEVERTAP_ACCOUNT=luna go run ./cmd
```

Run for `homiQ`:

```bash
CLEVERTAP_ACCOUNT=homiQ go run ./cmd
```

You can also export it first:

```bash
export CLEVERTAP_ACCOUNT=instaview
go run ./cmd
```

## Important

- `CLEVERTAP_ACCOUNT` must exactly match the key name inside `accounts.json`.
- The input file name must follow this format: `<account_key>_user.json`
- Example: if account key is `luna`, the file must be `luna_user.json`
