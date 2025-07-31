# Onyxia Onboarding Service

<p align="center">
    <a href="https://github.com/onyxia-datalab/onyxia-onboarding/actions/workflows/release.yml">
      <img src="https://github.com/onyxia-datalab/onyxia-onboarding/actions/workflows/release.yml/badge.svg?branch=main">
    </a>
    <a href="https://join.slack.com/t/3innovation/shared_invite/zt-2skhjkavr-xO~uTRLgoNOCm6ubLpKG7Q">
      <img src="https://img.shields.io/badge/slack-550_Members-brightgreen.svg?logo=slack">
    </a>
</p>

Onyxia Onboarding is a **Go-based REST API** responsible for managing user onboarding in **Onyxia** by provisioning **Kubernetes namespaces** with associated **quotas, and annotations**.

This project is a **refactor** of the [Java-based Onyxia API](https://github.com/inseeFrlab/onyxia-api) and is designed to provide a lightweight and efficient onboarding mechanism for **data science and cloud workspaces**. The original Java API will be split into three Go REST API modules: Onboarding, Bootstrap, and Feature API.

## 🌍 Project Website

Visit [onyxia.sh](https://onyxia.sh) to learn more about the Onyxia ecosystem.

## 🚀 Features

- **Automated namespace creation**: Ensures users have their own dedicated Kubernetes namespace.
- **Resource quotas**: Enforces limits on CPU, GPU, memory, and storage usage.
- **Namespace annotations**: Allows additional metadata if enabled via environment variables.
- **REST API**: Simple and efficient API for managing onboarding operations.

## 🏗️ Installation & Setup

This module is designed to be used within the Onyxia ecosystem. Detailed installation guides are coming soon. In the meantime, you can explore the local development setup below.

### Local Development

Prerequisites

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerized deployment)
- Kubernetes cluster with appropriate RBAC setup

1. **Clone the repository**:

   ```sh
   git clone https://github.com/onyxia-datalab/onyxia-onboarding.git
   cd onyxia-onboarding
   make install
   ```

2. **Setup environment variables**:

   Developers can create a local environment configuration by copying the default template:

   ```sh
   cp internal/bootstrap/env.default.yaml env.yaml
   ```

   Modify `env.yaml` as needed to configure the service. Feel free, env.yaml is git ignored.

3. **Run locally**:

   ```sh
   make run
   ```

4. **Run tests**:
   ```sh
   make test
   ```

## 🐳 Docker Deployment

You can build and run the service using **Docker** with the following commands:

```sh
make docker-build
make docker-run
```

This will:

- **Build the Docker image** using the appropriate platform for your system.
- **Run the container** and expose it on port **8080**.

### 🖀 Multi-Architecture Build

By default, the Docker image is built for **the local system architecture** that is running the command. If you need to support **both AMD64 and ARM64** (e.g., Apple M1/M2 chips), enable multi-architecture builds by setting the `MULTIARCH` flag:

```sh
MULTIARCH=1 make docker-build
```

### 🚀 Pushing to Docker Hub

First, set credentials in your environment.

Make sure there is a repository called: `$(DOCKER_REGISTRY)/$(PROJECTNAME)`.

```sh
DOCKER_REGISTRY=<Your registry> make docker-push
```

This will:

- **Build the image**.
- **Push** it to your Docker registry.

## 🛠️ Environnement Values

The service is configurable via **environment variables**.

### Configuration Structure

The configuration is loaded using **Viper** and can be provided via:

- Embedded `env.default.yaml`
- An external `env.yaml` file in the root directory (overrides defaults)
- Direct environment variables

### Available Configuration Options

#### **General**

| Variable             | Description                      | Default |
| -------------------- | -------------------------------- | ------- |
| `authenticationMode` | Authentication mode (none, oidc) | `none`  |

#### **Server**

| Variable | Description | Default |
| -------- | ----------- | ------- |
| `port`   | Server port | `8080`  |

#### **Security**

| Variable             | Description                  | Default |
| -------------------- | ---------------------------- | ------- |
| `corsAllowedOrigins` | List of allowed CORS origins | `[]`    |

#### **OIDC Authentication**

| Variable        | Description           | Default              |
| --------------- | --------------------- | -------------------- |
| `issuerURI`     | OIDC Issuer URI       | `""`                 |
| `skipTLSVerify` | Skip TLS verification | `false`              |
| `clientID`      | OIDC Client ID        | `""`                 |
| `audience`      | OIDC Audience         | `""`                 |
| `usernameClaim` | Claim for username    | `preferred_username` |
| `groupsClaim`   | Claim for groups      | `groups`             |
| `rolesClaim`    | Claim for roles       | `roles`              |

#### **Onboarding Configuration**

| Variable               | Description                                                                    | Default                      |
| ---------------------- | ------------------------------------------------------------------------------ | ---------------------------- |
| `namespacePrefix`      | Prefix for user namespaces                                                     | `user-`                      |
| `groupNamespacePrefix` | Prefix for group namespaces                                                    | `projet-`                    |
| `namespaceLabels`      | Static labels to add to the namespace (at creation and subsequent user logins) | `{ "created-by": "onyxia" }` |
| `annotations`          | See [Annotations](#annotations)                                                |                              |
| `quotas`               | See [Quotas](#quotas)                                                          |                              |

##### **Annotations**

| Variable                     | Description                                                                                     | Default |
| ---------------------------- | ----------------------------------------------------------------------------------------------- | ------- |
| `enabled`                    | Enable annotations                                                                              | `false` |
| `static`                     | Static annotations key-value pairs                                                              | `{}`    |
| `dynamic.lastLoginTimestamp` | Track last login timestamp by adding `onyxia_last_login_timestamp: <unix time in milliseconds>` | `false` |
| `dynamic.userAttributes`     | List of user attributes                                                                         | `[]`    |

##### **Quotas**

| Variable       | Description                                                                                                                                                                                      | Default |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| `enabled`      | Enable quotas                                                                                                                                                                                    | `false` |
| `default`      | Default quotas values [See](#quotas-values)                                                                                                                                                      |         |
| `userEnabled`  | Enable user-specific quotas                                                                                                                                                                      | `false` |
| `user`         | User quotas values [See](#quotas-values)                                                                                                                                                         |         |
| `groupEnabled` | Enable group-specific quotas                                                                                                                                                                     | `false` |
| `group`        | Group quotas values [See](#quotas-values)                                                                                                                                                        |         |
| `roles`        | Map of quotas corresponding to user roles. In case the user has multiple of those roles, only the first one will be applied. If user has no role from this list then user quota will be applied. | `{}`    |

##### **Quotas Values**

| Variable                     | Description                       | Default |
| ---------------------------- | --------------------------------- | ------- |
| `requests.memory`            | Default requested memory limit    | `10Gi`  |
| `requests.cpu`               | Default requested CPU limit       | `10`    |
| `limits.memory`              | Default memory limit              | `10Gi`  |
| `limits.cpu`                 | Default CPU limit                 | `10`    |
| `requests.storage`           | Default storage request           | `100Gi` |
| `count.pods`                 | Default max pods count            | `50`    |
| `requests.ephemeral-storage` | Default ephemeral storage request | `10Gi`  |
| `limits.ephemeral-storage`   | Default ephemeral storage limit   | `20Gi`  |
| `requests.nvidia.com/gpu`    | Default GPU requests              | `0`     |
| `limits.nvidia.com/gpu`      | Default GPU limits                | `0`     |

This is a subset of the configuration options available. The full configuration structure can be found in `env.default.yaml`.

## 📖 Contributing

We welcome contributions! To get started:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature-name`)
3. Commit changes (`git commit -m "Add new feature"`)
4. Push the branch (`git push origin feature-name`)
5. Open a pull request 🚀
