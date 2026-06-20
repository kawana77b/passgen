# passgen

A cross-platform CLI tool for generating cryptographically secure random passwords and unique IDs.

- Passwords are generated using the OS-provided CSPRNG (`crypto/rand`)
- Each character group (lowercase, uppercase, digits, symbols) is guaranteed to appear at least once
- Strength evaluated by entropy (bits) based on RFC 4086 and NIST SP 800-63B
- Single binary, no runtime dependencies
- Supports Windows / macOS / Linux

## Why did you create this?

It's for personal use. There are plenty of tools like this out there.

## Installation

### Download

Download the latest binary from the [Releases](../../releases) page.

### Build from source

```sh
git clone https://github.com/kawana77b/passgen.git
cd passgen
make install
```

## Usage

### Password generation

Passwords are sorted by entropy (strongest first).

```sh
# Default: 10 passwords, 16 characters, letters + digits + common symbols
passgen

# Specify length
passgen -l 24

# Specify count
passgen -n 5

# Letters and digits only
passgen -r alnum

# With symbols — safe set (-_.)
passgen -r full -s safe

# With symbols — full set, 32 characters
passgen -r full -s full -l 32

# Show entropy-based strength level next to each password
passgen --show-strength
passgen -S

# Copy the strongest password to the clipboard
passgen --clip
passgen -c
```

> [!NOTE]
> On Linux, `xclip` or `xsel` must be installed for clipboard support.

#### Flags

| Flag              | Short | Default  | Description                                             |
| ----------------- | ----- | -------- | ------------------------------------------------------- |
| `--length`        | `-l`  | `16`     | Password length                                         |
| `--rule`          | `-r`  | `full`   | Generation rule (see below)                             |
| `--symbol-set`    | `-s`  | `common` | Symbol set when using `--rule=full`                     |
| `--count`         | `-n`  | `10`     | Number of passwords to generate                         |
| `--show-strength` | `-S`  | `false`  | Show entropy-based strength level next to each password |
| `--clip`          | `-c`  | `false`  | Copy the strongest password to the clipboard            |

#### Rules (`--rule`)

| Value   | Characters             | Description                              |
| ------- | ---------------------- | ---------------------------------------- |
| `lower` | a-z                    | Lowercase letters only                   |
| `mixed` | a-z, A-Z               | Upper and lowercase letters              |
| `alnum` | a-z, A-Z, 0-9          | Letters and digits                       |
| `full`  | a-z, A-Z, 0-9, symbols | Letters, digits, and symbols _(default)_ |

#### Symbol sets (`--symbol-set`, used with `--rule=full`)

| Value    | Characters                  | Description                                |
| -------- | --------------------------- | ------------------------------------------ |
| `safe`   | `-_.`                       | Safe — unlikely to be rejected by websites |
| `common` | `!@#$%^&*`                  | Common _(default)_                         |
| `full`   | `!@#$%^&*()-_=+[]{};:,./?'` | Full — all supported symbols               |

> [!TIP]
> Run `passgen list-rules` to display this information in the terminal.

---

### Password strength check

Check the entropy and strength level of any password string. Output is JSON.

```sh
passgen check "myPassword123!"
```

```json
{
  "password": "myPassword123!",
  "length": 14,
  "entropy_bits": 91.76,
  "level": "very strong"
}
```

Strength levels are based on entropy (bits):

| Level       | Entropy      |
| ----------- | ------------ |
| weak        | < 40 bits    |
| fair        | 40 – 59 bits |
| strong      | 60 – 79 bits |
| very strong | ≥ 80 bits    |

Please note that this is a kind of score and is not strictly based on professional information engineering.

---

### Unique ID generation

This tool includes a feature for generating unique IDs.

```sh
# UUID v4 (random)
passgen uid v4

# UUID v7 (time-sortable)
passgen uid v7

# NanoID
passgen uid nanoid

# Specify count
passgen uid v4 -n 5
```

#### Flags

| Flag      | Short | Default | Description               |
| --------- | ----- | ------- | ------------------------- |
| `--count` | `-n`  | `10`    | Number of IDs to generate |

---
