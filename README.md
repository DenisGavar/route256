
# Route256

## Introduction

**Route256** is a project designed to manage and monitor different services related to e-commerce operations. It includes services like `checkout` and `loms` (Logistics and Order Management System), and integrates monitoring tools such as Prometheus and Alertmanager.

## Services

- **Checkout Service**: Handles customer orders and payment processing.
- **LOMS Service**: Manages logistics and order fulfillment.

## Contracts

- **[Russian](./contracts.ru.md)**
- **[English](./contracts.en.md)**

## Getting Started

To run the project:

1. Clone the repository:
    ```sh
    git clone https://github.com/DenisGavar/route256.git
    cd route256
    ```

2. Copy configuration files:
    ```sh
    cp checkout/config.yml.example checkout/config.yml
    cp loms/config.yml.example loms/config.yml
    ```
    Fill in the required values in `config.yml`.

3. Run all services:
    ```sh
    make run-all
    ```

4. Apply migrations:
    ```sh
    cd checkout
    ./migration.sh
    cd ..
    cd loms
    ./migration.sh
    cd ..
    ```

## Configuration

- **Docker Compose**: The project uses `docker-compose.yml` to set up services like databases and monitoring tools.
- **Environment Variables**: Set required environment variables in the configuration files.

## Usage

The application includes various handlers and daemons for managing checkout processes, logistics, and monitoring. Currently, more detailed usage instructions and application flow are available in [contracts.ru.md](./contracts.ru.md) (in Russian) and [contracts.en.md](./contracts.en.md) (in English).

## Technologies

- **Go**: The primary programming language for services.
- **Prometheus & Alertmanager**: For monitoring and alerting.
- **Docker**: To containerize the services.
- **Makefile**: To automate build and deployment processes.

## TODO

- Document addresses for services like Prometheus, Alertmanager, Kafka, Grafana, Jaeger, GRPC and databases.

