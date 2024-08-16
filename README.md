# Veritas

"Veritas" means "truth" in Latin, reflecting its purpose to discover the true security posture of container images.

Veritas offers an HTTP API for scanning public container images for security vulnerabilities. It also includes a dashboard for visualizing vulnerability data and downloading a Software Bill of Materials (SBOM).

## Getting Started

### Prerequisites

- Go
- Make

Ensure these tools are installed and available in your system's PATH before building the project.

### 1. Clone the repository

```shell
git clone https://github.com/lucasrod16/veritas.git
```

Navigate to the veritas directory:

```shell
cd veritas
```

### 2. Build the binary

```shell
make
```

The binary will be created at `./bin/veritas`.

### 3. Start veritas

```shell
./bin/veritas
```

Veritas will start and run in the foreground.

## Usage

### Scan container images

#### Dashboard

To access the dashboard, navigate to http://localhost:8080 in your browser. Enter the container image you want to scan into the text box, and either press Enter or click the Scan button.

#### API

For programmatic access, use curl or any other HTTP client to interact with veritas’ API endpoints.

- Scan a container image and get an SBOM in CycloneDX JSON format:

```shell
curl "http://localhost:8080/scan/report?image=<your-container-image>"
```

- Scan a container image and get detailed information about vulnerabilities:

```shell
curl "http://localhost:8080/scan/details?image=<your-container-image>"
```

### Stop veritas

To stop the program, press Ctrl+C.
