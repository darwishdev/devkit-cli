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


