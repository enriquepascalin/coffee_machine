# Coffee Machine

A console coffee machine simulator written in Go.

## Features

- User mode and hidden admin mode
- Product catalog with typed product IDs
- Typed milk and flavor enums
- Maintenance gauge and service workflow
- Inventory management
- Money handling in integer cents
- Accepted denominations and change-making
- Colored terminal output
- Simple brewing progress bar
- Admin price updates
- Realistic refill units

## Project Layout

```text
cmd/coffee-machine          application entrypoint
internal/app                application flow and mode handling
internal/auth               authentication
internal/catalog            product catalog and recipes
internal/config             defaults and bootstrap values
internal/machine            coffee machine domain
internal/money              money, denominations, change
internal/ui                 terminal I/O helpers