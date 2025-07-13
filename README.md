# HomeKit Proxy

HomeKit Proxy is a lightweight, self-hosted bridge that exposes your shell commands as HomeKit accessories. It allows you to control your custom devices and scripts through Apple's Home app and Siri.

## Features

*   **HomeKit Bridge:** Exposes shell commands as HomeKit accessories.
*   **Device Configuration:** Configure devices and their characteristics using a simple YAML file (`device.yaml`).
*   **Automations:** Schedule shell commands to run at specific times using cron expressions (`automation.yaml`).
*   **Web UI:** A simple web interface to view the status of your accessories and automations, and to manually trigger automations.
*   **Action Log:** Logs all actions to a SQLite database for auditing and debugging.
*   **Configuration Hot-Reload:** Automatically reloads the server when configuration files are changed.
*   **Customizable:** Customize network interface, address, and PIN code.

## How it Works

HomeKit Proxy works by creating a virtual HomeKit bridge. You define your accessories and their characteristics in a `device.yaml` file. For each characteristic, you specify shell commands to get or set its value. The proxy then exposes these accessories to your HomeKit network.

When you interact with an accessory in the Home app (e.g., turn on a light), HomeKit Proxy executes the corresponding shell command. It can also periodically poll the get command to keep the state of the accessory up-to-date.

## Getting Started

### Prerequisites

*   Go 1.18 or later

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/waynezhang/homekit-proxy.git
    cd homekit-proxy
    ```

2.  Build the binary:
    ```bash
    make build
    ```
    The binary will be located in the `bin` directory.

### Configuration

1.  Create a `config` directory:
    ```bash
    mkdir config
    ```

2.  Create a `device.yaml` file in the `config` directory. See the [device.sample.yaml](config/device.sample.yaml) for an example.

3.  Create an `automation.yaml` file in the `config` directory. See the [automation.sample.yaml](config/automation.sample.yaml) for an example.

### Running the Proxy

```bash
./bin/homekit-proxy serve --config ./config --db ./db
```

This will start the HomeKit Proxy server. You can then add the bridge to your Home app by scanning the QR code displayed in the console or by entering the PIN code.

## Configuration

See `config/device.sample.yaml` and `config/automation.sample.yaml`.

## Environment Variables

| Name                     | Description                                                                 |
| ------------------------ | --------------------------------------------------------------------------- |
| `HOMEKIT_PROXY_USER`     | The username for the web UI.                                                |
| `HOMEKIT_PROXY_PASSWORD` | The password for the web UI.                                                |
| `HOMEKIT_PROXY_BINDADDR` | The address and port to bind the HomeKit proxy to.                                   |
| `HOMEKIT_PROXY_IFACE`    | The network interface to use for the HomeKit proxy.                         |
| `HOMEKIT_PROXY_NTFY_TOPIC` | The ntfy.sh topic to send notifications to when an automation is triggered. |

## Web UI

HomeKit Proxy provides a simple web UI to view the status of your accessories and automations. You can access it at `http://bindaddr`.

The web UI also allows you to manually trigger automations.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
