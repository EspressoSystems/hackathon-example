# Hotshot Query Example

This is a simple Go project that fetches and prints namespace transactions from the Hotshot query service. The application periodically polls for the latest block height, then retrieves and pretty-prints the namespace transactions in your terminal. It skips fallback responses (e.g. those containing `"FetchBlock"`).

## Overview

- **Block Height Endpoint:**  
  Retrieves the latest block height from `/status/block-height`.

- **Namespace Transactions:**  
  Uses the configured HotShot URL along with the block height and namespace to build the full endpoint:
  ```
  <hotshot_url>/availability/block/<blockHeight>/namespace/<namespace>
  ```
  The returned JSON is filtered and pretty-printed to display only the relevant transaction data. Fallback responses (such as those containing `"FetchBlock"`) are omitted.

## Requirements

- Go 1.23.2 or higher

## Setup

This project uses the OffchainLabs fork of go-ethereum as a Git submodule. To clone and set up the repository:

### **Option 1: Clone with submodules**

`git clone --recurse-submodules https://github.com/your-username/hackathon-example.git`
`cd hackathon-example`

### **Option 2: Clone and initialize submodules separately**

`git clone https://github.com/your-username/hackathon-example.git`
`cd hackathon-example`
`git submodule init`
`git submodule update`

### **Ensure the correct go-ethereum commit**

`cd go-ethereum`
`git checkout a437823d683b7392d3513331c7ca90705b548ef7`
`cd ..`

### **Update dependencies**

`go mod tidy`

## Getting Started

**Run the Application**

   ```bash
   go run .
   ```

## Configuration

The application uses a configuration file located at `config/config.json`. This file allows you to customize the following parameters:

- **Query Service Endpoint:** Set your Hotshot query service URLs.
- **Namespace (Chain ID):** Change the namespace identifier. **Note:** In this project, the namespace is the same as the chain id used by the network. This means that when you configure the namespace, you are effectively specifying the chain identifier.
- **Polling Interval:** Adjust the time interval for how frequently the project polls the API.

Below is an example configuration file:

```json
{
  "hotshot_url": "https://query.decaf.testnet.espresso.network/v0",
  "namespace": 20250115,
  "polling_interval": 10,
  "from": "0x0000000000000000000000000000000000000000",
  "value": 100
}
```


## Extending the Project

You can easily extend this project to interact with additional endpoints (such as block headers or explorer transactions) by following the existing structure. Simply update the configuration and add new functions to handle the additional API calls.

## Additional Resources

Below are some common API endpoints along with example responses used in this project:

- **/status/block-height**  
  *Example Response:*  
  ```json
  {"blockHeight": 12345}
  ```

- **/availability/block/"blockHeight"/namespace/"chainId"**  
  *Example Response:*  
  ```json
  {
    "proof": null,
    "transactions": []
  }
  ```

**Example Query Responses:**  
For additional reference, please see [`example-logs/transactions-data.json`](example-logs/transactions-data.json), which contains valid API responses to help you know what to expect.

For a complete overview of all supported endpoints and details on each module, please refer to the [Espresso API Guide](https://docs.espressosys.com/network/api-reference/sequencer-api).

## Espresso API Modules

The Espresso API is divided into several modules:
 - [Status API](https://docs.espressosys.com/network/api-reference/sequencer-api/status-api) – Provides node-specific state and consensus metrics.
 - [Catchup API](https://docs.espressosys.com/network/api-reference/sequencer-api/catchup-api) – Serves recent consensus state to allow peers to catch up with the network.
 - [Availability API](https://docs.espressosys.com/network/api-reference/sequencer-api/availability-api) – Serves data recorded by the [Tiramisu DA](https://docs.espressosys.com/network/learn/the-espresso-network/properties-of-hotshot/espresso-data-availability-layer) layer (EspressoDA), such as committed blocks. 
 - [Node API](https://docs.espressosys.com/network/api-reference/sequencer-api/node-api) – Complements the availability API by serving eventually consistent data that is not (yet) necessarily agreed upon by all nodes.
 - [State API](https://docs.espressosys.com/network/api-reference/sequencer-api/state-api) – Serves consensus state derived from finalized blocks.
 - [Events API](https://docs.espressosys.com/network/api-reference/sequencer-api/events-api) – Streams events from HotShot.
 - [Submit API](https://docs.espressosys.com/network/api-reference/sequencer-api/submit-api) – Allows users to submit transactions to the public mempool.

## License

This project is provided as-is with no warranty.