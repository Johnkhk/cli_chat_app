# CLI Chat App

- [CLI Chat App](#cli-chat-app)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
      - [Client Installation](#client-installation)
  - [Demo](#demo)
  - [Usage](#usage)
  - [Testing](#testing)
  - [Contributing](#contributing)
  - [License](#license)
  - [Acknowledgments](#acknowledgments)


## Introduction

Welcome to the CLI Chat App, a end-to-end encrypted chat application that runs on your terminal!

It is designed to be lightweight, fast, and easy to use, making it perfect for developers and tech enthusiasts who prefer working in a terminal environment.

## Features

- **Secure Messaging**: Utilizes end-to-end encryption ([The Signal Protocol](https://signal.org/docs/)) to ensure that your messages remain private and secure. This means chat history is stored locally on your device and is not accessible by the server or any third parties.
- **User Authentication**: Register and log in with a username and password. JWTs are used to keep you signed in between sessions.
- **Friend Management**: Send and receive friend requests, and manage your friend list.
- **Real-time Communication**: Chat with your friends in real-time using a simple and intuitive interface.
- **Multi-media support**: Send and receive images, videos, and files.
- **Cross-Platform**: Available on Linux, macOS (Intel and ARM), and Windows.

## Getting Started

### Installation

#### Client Installation

1. **Install the binary**: Download the binary from the [releases](https://github.com/Johnkhk/cli_chat_app/releases) page for your operating system.

2. **Configure the binary**: If on Linux or MacOS, you can run `chmod +x cli_chat_app` to make the binary executable. If on MacOS, you may need to run `xattr -cr cli_chat_app` to remove the quarantine attribute.
3. **Run the binary**: Run the binary by typing `./cli_chat_app` in your terminal.
4. **Add server address**: Add the following line to your `.env` file in the same path as the binary:
   ```
   SERVER_ADDRESS=clichatapp.click:50051
   ```

   Alternatively, you can set the `SERVER_ADDRESS` environment variable by running `export SERVER_ADDRESS=clichatapp.click:50051` in your terminal.
5. **Add an Alias (Optional)**: Add an alias to your `.bashrc` or `.zshrc` file to easily access the binary. For example, `alias cli_chat="~/path/to/cli_chat_app"`.

## Demo

[![Watch the demo video](https://img.youtube.com/vi/E5gffV7ap5g/0.jpg)](https://youtu.be/E5gffV7ap5g?si=Mz2KRdsPKwo6KU22)

Click the image above to watch a demo video of the CLI Chat App in action!


## Usage

- **Register**: Create a new account by selecting the "Register" option.
- **Login**: Log in with your credentials to access the chat features.
- **Send Friend Requests**: Add friends by sending them a request.
- **Chat**: Start a conversation with your friends. (Send text or files)

## Testing

Testing is done using [gotestsum](https://github.com/gotestyourself/gotestsum). Tests set up and teardown a single server, a specified number of clients, and the necessary (local client and server) databases.

To run the tests, set the following environment variables:

```
CLI_CHAT_APP_JWT_SECRET_KEY=your-generated-key
TEST_DATABASE_URL="root@tcp(127.0.0.1:3306)/
```

Optionally, you can set the `TEST_LOG_DIR` environment variable to specify the directory for the test logs. If not set, the logs will be stored in the app directory. (See `GetAppDirPath()` in `client/app/utils.go`)

The jwt key is a random 32-byte key encoded in Base64. You can generate one using the following command:

```
openssl rand -base64 32
```

Then run the following command:

```
gotestsum --format=short-verbose ./test/...
```

Tests are also ran in github actions.


## Contributing

We welcome contributions! Please fork the repository and submit a pull request with your changes. Ensure that your code follows the project's coding standards and includes appropriate tests.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the contributors and the open-source community for their support and contributions.

---

*Note: This README is a work in progress. More detailed instructions, including GIFs and images, will be added soon.*
