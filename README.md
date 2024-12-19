# CLI Chat App

- [CLI Chat App](#cli-chat-app)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Examples](#examples)
    - [App usage between two users](#app-usage-between-two-users)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
      - [Client Installation](#client-installation)
  - [Usage](#usage)
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
- **Cross-Platform**: Available on Linux, macOS (Intel and ARM), and Windows.

## Examples

### App usage between two users

[![Example](https://img.youtube.com/vi/auK8GqmSAlw/3.jpg)](https://youtu.be/auK8GqmSAlw)


## Getting Started

### Installation

#### Client Installation

1. **Install the binary**: Download the binary from the [releases](https://github.com/Johnkhk/cli_chat_app/releases) page for your operating system.

2. **Configure the binary**: If on Linux or MacOS, you can run `chmod +x cli_chat_client` to make the binary executable. If on MacOS, you may need to run `xattr -cr cli_chat_client` to remove the quarantine attribute.
3. **Run the binary**: Run the binary by typing `./cli_chat_client` in your terminal.
4. **Add server address**: Add the following line to your `.env` file in the same path as the binary:
   ```
   SERVER_ADDRESS=clichatapp.click:50051
   ```

   Alternatively, you can set the `SERVER_ADDRESS` environment variable by running `export SERVER_ADDRESS=clichatapp.click:50051` in your terminal.
5. **Add an Alias (Optional)**: Add an alias to your `.bashrc` or `.zshrc` file to easily access the binary. For example, `alias cli_chat="~/path/to/cli_chat_client"`.


## Usage

- **Register**: Create a new account by selecting the "Register" option.
- **Login**: Log in with your credentials to access the chat features.
- **Send Friend Requests**: Add friends by sending them a request.
- **Chat**: Start a conversation with your friends.

## Contributing

We welcome contributions! Please fork the repository and submit a pull request with your changes. Ensure that your code follows the project's coding standards and includes appropriate tests.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the contributors and the open-source community for their support and contributions.

---

*Note: This README is a work in progress. More detailed instructions, including GIFs and images, will be added soon.*
