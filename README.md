# Rollup Transaction Monitor

This is a Go application that monitors an espresso integrated rollup for specific transactions. The application acts as a security bot that watches for outgoing transactions from your wallet address that exceed a configured value threshold, helping you detect potentially unauthorized transfers or drained funds.

## Overview

- **Block Monitoring:**  
  Retrieves the latest processed block from your Caff node and inspects transactions.

- **Transaction Filtering:**  
  Filters transactions based on sender address and minimum value. The application will log the transaction details to the console.

## Requirements

- Go 1.23.2 or higher
- go-ethereum submodule
- Caff node (local or remote)

## Getting Started

1. **Clone the Repository**

   ```bash
   git clone https://github.com/EspressoSystems/hackathon-example --recursive
   cd hackathon-example
   ```

   If you've already cloned the repository without the `--recursive` flag, you can fetch the go-ethereum submodule with:

   ```bash
   git submodule update --init --recursive
   ```

2. **Install Dependencies**

   ```bash
   go mod tidy
   ```
3. **Configure the Application**

   Edit the configuration file at `config/config.json` to set your Caff node URL, polling interval, and transaction filtering criteria.

3. **Run the Application**

   ```bash
   go run .
   ```

## Configuration

The application uses a configuration file located at `config/config.json`. This file allows you to customize the following parameters:

- **Caff Node URL:** The URL of your Caff node (e.g., http://127.0.0.1:8550)
- **Polling Interval:** How frequently (in seconds) the application checks for new blocks
- **Value:** Minimum transaction value to monitor (in wei)
- **From:** Ethereum address to monitor for outgoing transactions

Below is an example configuration file:

```json
{
  "caff_node_url": "http://YOUR_HOST:8550",
  "polling_interval": 1,
  "value": VALUE_IN_WEI,
  "from": "YOUR_WALLET_ADDRESS"
}
```

## Project Structure

- `main.go` - Application entry point
- `transactions.go` - Transaction monitoring logic
- `config/` - Configuration handling
- `go-ethereum/` - Ethereum client library (submodule)

## Extending the Project

You can easily extend this project to:
- Monitor multiple addresses
- Filter by recipient address
- Track specific contract interactions
- Implement notifications when matching transactions are found

## License

This project is provided as-is with no warranty.
