# devkit-cli


## Boost your Go development with devkit-cli: The ultimate tool for bootstrapping robust and scalable backend applications.

devkit-cli is a command-line interface designed to streamline the creation of new Go projects. It provides a solid foundation with essential backend functionalities out-of-the-box, allowing you to focus on building core features instead of writing repetitive boilerplate code.
By leveraging a modular architecture and powerful code generation capabilities, devkit-cli accelerates your development workflow and helps you deliver high-quality applications faster.


## Getting Started

### Dependencies

Before you begin, ensure you have the following dependencies installed on your system:

1.  **Buf CLI:**  Buf is used for managing Protobuf files and generating code. Install it by following the instructions on the [Buf website](https://buf.build/docs/installation).

2.  **Supabase CLI:** Supabase CLI is required for interacting with the Supabase database. You can find installation instructions on the Supabase CLI GitHub repository [invalid URL removed].

3.  **SQLc CLI:** SQLc generates type-safe Go code from SQL queries. Get it from the SQLc website [invalid URL removed].

4.  **Connect-Go:** Connect-Go is used for building gRPC servers and clients. Install it using 
```
go get github.com/bufbuild/connect-go
```

5.  **protoc-gen-go:** This is the official Go code generator for Protocol Buffers. Install it with
```
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

6.  **Resend Account:**  You'll need a Resend account and API key to utilize the email sending functionality. Sign up for an account on the [Resend website](https://resend.com).

7.  **Google Cloud Project:**  Create a project on Google Cloud Platform and obtain a client ID and secret for enabling Google authentication [Google Tutorial](https://support.google.com/cloud/answer/6158849?hl=en).

8.  **GitHub Personal Access Token:** devkit-cli requires a GitHub personal access token to automatically fork the base project and create new repositories for your applications. You can generate a token in your GitHub account settings under "Developer settings" -> "Personal access tokens".  Make sure to grant the token the necessary permissions for repository creation and management. 

Once you have these dependencies set up, you're ready to install devkit-cli!


---

## Installation

To install `devkit-cli`, run the following command in your terminal:

```bash
curl -sSL https://raw.githubusercontent.com/darwishdev/devkit-cli/refs/heads/main/install.sh | bash
```

This command downloads and runs the installation script directly from the repository, setting up `devkit-cli` on your system.

---


## Configuration

devkit-cli relies on a configuration file to store essential settings and API keys. This file is automatically created during installation at `~/.config/devkitcli/devkit`.

Here's a breakdown of the fields you need to populate:

*   **GIT_USER**: Your GitHub username.
*   **DOCKER_HUB_USER**: Your Docker Hub username.
*   **BUF_USER**: Your Buf username.
*   **GITHUB_TOKEN**:  Your GitHub personal access token. This token requires permissions for repository creation and management.
*   **API_SERVICE_NAME**: The desired name for your API service (e.g., "devkit").
*   **API_VERSION**: The initial version of your API (e.g., "v1").
*   **GOOGLE_CLIENT_ID**: The client ID from your Google Cloud project.
*   **GOOGLE_CLIENT_SECRET**: The client secret from your Google Cloud project.
*   **RESEND_API_KEY**: Your Resend API key.

**Important:**

*   Open the `~/.config/devkitcli/devkit` file in a text editor and replace the placeholder values with your actual credentials.
*   Ensure the `GITHUB_TOKEN` has the necessary permissions for repository creation and management within your GitHub account.
*   Keep your configuration file secure and do not share sensitive information like API keys publicly.

Once you've configured these settings, devkit-cli will be able to leverage these credentials to streamline your development workflow and automate various tasks.

---

### `devkit new api`

This command bootstraps a new Go backend API application by forking a base repository, cloning it locally, and performing the necessary configurations.

**Usage:**

```bash
devkit new api [app-name] --org [github-org] --buf-user [buf-user]
```

*   `app-name`: The name of your new application. This will also be used for the repository name.
*   `--org`: (Optional) The GitHub organization to fork the base repository to. If not provided, it defaults to the `GIT_USER` specified in your configuration file.
*   `--buf-user`: (Optional) Your Buf.build username. This is used for pushing API documentation. If not provided, it defaults to the `BUF_USER` in your configuration file.

**Functionality:**

1.  **Forking:** Forks the base API repository to the specified GitHub organization.
2.  **Cloning:** Clones the forked repository to your local machine.
3.  **Configuration:** Copies configuration files and replaces placeholders with your application-specific settings.
4.  **Code Generation:** Executes `make buf` and `make sqlc` to generate code from Protobuf definitions and SQL queries.
5.  **Dependency Management:** Runs `go mod tidy` to manage Go module dependencies.
6.  **Project Initialization:** Runs `devkit init` to initialize the project and create a configuration file specific to the new project.

This command automates the initial setup of a new Go backend API project, providing a solid foundation with pre-configured functionalities. You can then focus on building your core application logic without worrying about repetitive boilerplate code.

---


### `devkit init`

This command initializes a new project configuration file in the current directory.

**Usage:**

```bash
devkit init
```

**Functionality:**

*   Creates a file named `devkit.env` in the current directory.
*   Populates the file with default values based on your global CLI configuration (`~/.config/devkitcli/devkit`).
*   Includes fields for:
    *   `GIT_USER`, `DOCKERHUB_USER`, `BUF_USER`
    *   `GOOGLE_CLIENT_ID`, `RESEND_API_KEY`, `GOOGLE_CLIENT_SECRET`
    *   `GITHUB_TOKEN`, `API_SERVICE_NAME`, `API_VERSION`
    *   `ENVIRONMENT`, `DB_PROJECT_REF`, `SUPABASE_SERVICE_ROLE_KEY`
    *   `SUPABASE_API_KEY`, `DB_SOURCE`, `APP_NAME`

**Note:** If you created your project using `devkit new api`, this command is already executed as part of the setup. You only need to run `devkit init` manually if you don't have a `.devkit.env` file in your project directory.

---


