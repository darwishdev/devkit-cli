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

```go
// getDomainCmd returns the cobra command for 'devkit new domain'.
// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *Command) getDomainCmd() *cobra.Command {
	apiCmd := &cobra.Command{
		Use:   "domain [domain_name]",
		Short: "Create a new domain in your application.",
		Long: `The 'domain' command generates a new domain within your Go backend application.
It creates the necessary directory structure (adapter, repo, usecase) and 
populates them with boilerplate code for the given domain name.`,
		Args: cobra.ExactArgs(1), // Ensure exactly one argument (domain_name) is provided
		Run: func(cmd *cobra.Command, args []string) {
			c.newCmd.NewDomain(args, cmd.Flags())
		},
	}

	apiCmd.Flags().StringP("org", "o", "", "GitHub organization to fork the repository to. If not provided, defaults to the configured Git user.")
	apiCmd.Flags().StringP("buf-user", "b", "", "Buf.build username for pushing API documentation. If not provided, defaults to the configured Buf user.")

	return apiCmd
}
```

## README Update

### `devkit new domain`

This command generates a new domain within your Go backend application.

**Usage:**

```bash
devkit new domain [domain-name]
```

*   `domain-name`: The name of the domain you want to create.

**Functionality:**

1.  **Directory Creation:** Creates a new directory named under app/`domain-name`.
2.  **Subdirectory Generation:** Creates three subdirectories within the domain directory:
    *   `adapter`: For code related to external interactions (e.g., database, external APIs).
    *   `repo`: For data access and persistence logic.
    *   `usecase`: For business logic and application-specific rules.
3.  **Boilerplate Code:** Populates the subdirectories with basic Go code files:
    *   `adapter.go`, `repo.go`, `usecase.go`
4.  **API Integration:** Updates the `api/api.go` file to include the new domain's use case:
    *   Adds necessary imports, fields, instantiations, and dependency injections.

This command streamlines the process of adding new domains to your application, providing a structured foundation and reducing boilerplate code. You can then customize the generated code to implement the specific logic for your domain.


---


## README Update

### `devkit new feature`

This command generates a new feature within a specified domain in your Go backend application.

**Usage:**

```bash
devkit new feature [feature-name] --domain [domain-name]
```

*   `feature-name`: The name of the feature you want to create.
*   `--domain`: The name of the existing domain where this feature belongs.

**Functionality:**

1.  **File Generation:** Creates new files within the specified domain's subdirectories:
    *   `adapter`: `[feature-name]_adapter.go`
    *   `repo`: `[feature-name]_repo.go`
    *   `usecase`: `[feature-name]_usecase.go`
    *   `proto`: `[domain-name]_[feature-name].proto` (in the `proto` directory)
    *   `query`: `[domain-name]_[feature-name].sql` (in the `supabase/queries` directory)
    *   `api`: `[domain-name]_[feature-name]_rpc.go` (in the `api` directory)
2.  **Boilerplate Code:** Populates the generated files with basic Go code relevant to the feature.
3.  **Proto Integration:** Updates the domain's service proto file (`[domain-name]_service.proto`) to import the new feature's proto definition.

This command automates the process of adding new features to your application's domains, ensuring consistency and reducing manual effort. You can then customize the generated code to implement the specific logic for your feature.

---

### `devkit new endpoint`

This command generates a new endpoint within a specified feature and domain in your Go backend application.

**Usage:**

```bash
devkit new endpoint [endpoint-name] --domain [domain-name] --feature [feature-name]
```

*   `endpoint-name`: The name of the endpoint you want to create.
*   `--domain`: The name of the existing domain where this endpoint belongs.
*   `--feature`: The name of the existing feature within the domain where this endpoint belongs.
*   `--get`: (Optional) If enabled, this endpoint will be marked as having no side effects, and the REST endpoint will be of type GET.
*   `--empty-response`: (Optional) If enabled, this endpoint will return an empty response.
*   `--empty-request`: (Optional) If enabled, this endpoint will accept an empty request.
*   `--list`: (Optional) If enabled, this endpoint will return a set of records; if disabled, it will use a single record, except for endpoints with the term 'list' in their name.

**Functionality:**

1.  **Code Generation:** Updates various files with the necessary code for the new endpoint:
    *   Appends code to existing `adapter`, `repo`, `usecase`, `proto`, `query`, and `api` files within the specified domain and feature.
    *   Injects interfaces and methods into `adapter/adapter.go`, `usecase/usecase.go`, `proto/[api-service-name]_service.proto`, and `repo/repo.go`.
2.  **Boilerplate Code:** Provides basic Go code relevant to the endpoint, including request/response handling, database interaction, and business logic.
3.  **Proto Integration:** Adds the new endpoint to the domain's service proto file (`[domain-name]_service.proto`).
4.  **Query Generation:** Creates a SQL query file (`[domain-name]_[feature-name].sql`) for the endpoint's database interaction.

This command further automates the development process by generating endpoints with predefined structures and best practices. You can then customize the generated code to implement the specific logic for your endpoint.

---


### `devkit seed`

This command automates the process of seeding your database tables with data from an Excel file and optionally creates users in Supabase Auth.

**Usage:**

```bash
devkit seed [schema_name] --file-path [excel_file_path]
```

*   `schema_name`: The name of the database schema you want to seed.
*   `--file-path`: The path to the Excel file containing the data.
*   `--out-file`: (Optional) The path to the output SQL file. If not provided, the SQL will be printed to the console.
*   `--execute`: (Optional) If enabled, the generated SQL will be executed against the configured database.
*   `--skip-supabase`: (Optional) If enabled, the command will skip creating users in Supabase Auth.

**Functionality:**

1.  **Excel Parsing:** Reads the Excel file and extracts data from each sheet.
2.  **SQL Generation:** Generates SQL insert statements for each sheet, mapping the sheet name to the table name in the specified schema.
3.  **Relationship Handling:** Automatically handles one-to-many and many-to-many relationships between tables.
4.  **Hashing:** Hashes password columns automatically (columns with names ending in `(#)`).
5.  **Supabase Auth:** If the Excel file contains a sheet named "users", the command will create users in Supabase Auth with the data from that sheet, unless the `--skip-supabase` flag is enabled.
6.  **Output:** Prints the generated SQL to the console or saves it to a file.
7.  **Execution:** Optionally executes the SQL against the configured database.

**Column Mapping:**

*   **One-to-many:** `<primary_key_column><OneToManyDelimiter><table_name><OneToManyDelimiter><search_key_column>`
*   **Many-to-many:** `<joining_table_primary_key><ManyToManyDelimiter><joining_table_name><ManyToManyDelimiter><second_table_name><ManyToManyDelimiter><second_table_search_column><ManyToManyDelimiter><first_table_search_column>`
*   **Hashing:** `column_name+(#)`

**Example:**

To seed the `accounts_schema` from an Excel file named `accounts.xlsx`, create users in Supabase Auth, and save the generated SQL to a file named `q.sql`:

```bash
devkit seed accounts_schema --file-path accounts.xlsx --out-file q.sql
```

To seed the same schema, skip Supabase Auth, and execute the SQL directly:

```bash
devkit seed accounts_schema --file-path accounts.xlsx --execute --skip-supabase
```

This command simplifies database seeding and integrates with Supabase Auth for easier user management in your projects.his command significantly simplifies the process of seeding your database, especially when dealing with complex relationships and data transformations. It leverages the (`sqlseeder`)[https://github.com/darwishdev/sqlseeder] package  to handle the intricacies of SQL generation and relationship mapping.

---

### `devkit seed storage`

This command seeds Supabase storage with files and icons from specified paths.

**Usage:**

```bash
devkit seed storage --files-path [files_path] --icons-path [icons_path]
```

*   `--files-path`: The path to the folder containing your images and files. This is required if `--icons-path` is not provided.
*   `--icons-path`: The path to the folder containing your SVG icons. This is required if `--files-path` is not provided.

**Functionality:**

1.  **Bucket Creation:** Creates buckets in Supabase storage based on the subfolder structure in the `files-path`. For example, if `files-path` is `assets`, and `assets` has subfolders `images` and `documents`, the command will create buckets named `images` and `documents`.
2.  **File Upload:** Uploads the files from each subfolder in `files-path` to their corresponding buckets in Supabase storage.
3.  **Icon Insertion:** If `--icons-path` is provided, the command reads all SVG files from that path and inserts their content into the `icons` table in your database. Each icon's name is derived from its filename.

**Example:**

To seed Supabase storage with files from the `assets` folder:

```bash
devkit seed storage --files-path assets
```

To seed Supabase storage with files from the `assets` folder and icons from the `svg-icons` folder:

```bash
devkit seed storage --files-path assets --icons-path svg-icons
```

This command provides a convenient way to populate your Supabase storage with files and icons, automating the process of bucket creation and file uploads.

---
