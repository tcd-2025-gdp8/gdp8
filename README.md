# PeerSphere - Group 8 - CSU44098 Group Design Project (2024/25)

## Overview

PeerSphere is an intelligent academic collaboration platform that enables students to discover, join, and engage with study groups based on their academic modules. The platform supports real-time communication, smart scheduling, file-based collaboration, and even conversational AI for contextual document support.

## Key Features

- **Module-Based Study Groups**  
  Students can create, discover, and join study groups based on their enrolled academic modules.

- **Real-Time Communication**  
  WebSocket-powered group chat with persistent history for seamless academic collaboration.

- **Smart Scheduling with Google Calendar**  
  Users can schedule study sessions that are automatically synced with Google Calendar, including support for meeting cancellations and reminders.

- **AI-Powered Chatbot (Gemini API)**  
  Integrated chatbot enables users to ask contextual questions about uploaded PDF files and get relevant responses without leaving the app.

- **File Upload & Management**  
  Users can upload and remove PDF files within their study groups for collaborative use.

- **Study Group Leaderboard**  
  A motivational leaderboard displays study groups ranked by cumulative study hours.

- **Secure Authentication**  
  Firebase Authentication with OAuth2 ensures safe and seamless user login.

- **Robust Backend Architecture**  
  Built in Go with MySQL and modular design patterns for performance and maintainability.

- **Containerized Deployment**  
  Docker and Docker Compose ensure consistent local and production environments.

- **CI/CD Pipeline**  
  GitHub Actions and custom scripts automate testing, code quality checks, and deployment.


## Step by Step

TLDR: to run the app locally it is enough to run `./scripts/run-localprod` from the main project directory. Details are provided below.

The app server can be easily run using the scripts provided in the `scripts` directory.

For authentication, the server relies on Firebase, which affects how it can be run locally and in production. To support both scenarios, two run scripts are provided: `run-localprod` and `run-prod`. The `localprod` script uses a Firebase emulator, allowing the server to run entirely locally without requiring real Firebase credentials. Conversely, the `prod` script connects to the actual Firebase services, requiring valid Firebase credentials to be configured.

In order to run the app in `localprod` mode it is enough to run `./scripts/run-localprod` from the main project directory. The script will then run all the necessary services and expose the web-based frontend on port `80`. The app can then be accessed through the browser by navigating to `localhost`.

Similarly, the app can be run in `prod` mode by running `./scripts/run-prod` from the main project directory. Since it requires connection to the actual Firebase service, however, it is also needed to provide Firebase credentials files. A `serviceAccountKey.json` file must be placed in `backend/credentials` directory. The environmental variables (as specified in the `frontend/src/auth/firebase.tsx`) must also be defined in a `.env.production` file placed in the `frontend-web` directory. With the credentials files in place, the script will similarly run all the necessary services and expose the web-based frontend on port `80`. The app can then be accessed through the browser by navigating to `localhost`.

Docker compose is needed to run the application using the provided scripts.

Additional scripts are also provided for running the backend in `dev` mode (`./scripts/run-dev`), as well as for executing integration tests (`./scripts/run-tests-integration`).

## Backends

The application backend utilises a MySQL database and Firebase for authentication. If the app is run using the provided scripts, any ports for backend services used will also be exposed by default:

- the main backend server (REST and WebSocket API) can be accessed on port `8080`;
- the Firebase emulator (if used) can be accessed on port `4000` (emulator UI) and `9099` (auth emulator).
