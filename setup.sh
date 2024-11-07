#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Update packages
echo "Updating packages..."
sudo yum update -y

# Install Git
echo "Installing Git..."
sudo yum install -y git

# Install Docker
echo "Installing Docker..."
sudo yum install -y docker

# Start Docker service and enable it on boot
sudo service docker start
sudo chkconfig docker on  # Enables Docker to start on boot

# Add the current user to the docker group (to avoid using sudo with Docker commands)
echo "Adding $(whoami) to the docker group..."
sudo usermod -aG docker $USER

# Install Docker Compose
echo "Installing Docker Compose..."
LATEST_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep -Po '"tag_name": "\K.*\d')
sudo curl -L "https://github.com/docker/compose/releases/download/${LATEST_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify Docker and Docker Compose installation
docker --version
docker-compose --version

# Clone the repository
echo "Cloning repository..."
REPO_URL="https://github.com/Johnkhk/cli_chat_app"  # Replace with your repository URL
APP_DIR="/home/ec2-user/cli_chat_app"
if [ -d "$APP_DIR" ]; then
  sudo rm -rf "$APP_DIR"
fi
git clone "$REPO_URL" "$APP_DIR"

# Navigate to the app directory
cd "$APP_DIR"

# Start Docker Compose
echo "Starting application with Docker Compose..."
docker-compose up -d

echo "Setup completed successfully!"

# Reminder to logout and log back in for Docker group permissions to take effect
echo "Please log out and back in to apply Docker group permissions."
