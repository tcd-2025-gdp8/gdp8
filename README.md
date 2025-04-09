# PeerSphere - Group 8 - CSU44098 Group Design Project (2024/25)

## Overview

PeerSphere is an intelligent academic collaboration platform that enables students to discover, join, and engage with study groups based on their academic modules. The platform supports real-time communication, smart scheduling, file-based collaboration, and even conversational AI for contextual document support.

### Key Features

**Module-Based Study Groups**  
Students can create, discover, and join study groups based on the available modules.

**Real-Time Communication**  
WebSocket-based implementation allows for group chats with persistent history to encourage academic collaboration.

**Smart Scheduling with Google Calendar**  
Users can effortlessly schedule study sessions through the app, with automatic Google Calendar synchronisation, and including support for meeting cancellations and reminders.

**AI-Powered Chatbot (Gemini API)**  
An integrated chatbot enables users to ask contextual questions about group-specific uploaded PDF files and get relevant responses without leaving the app, or having to reupload the files each time they want to ask questions.

**File Upload & Management**  
Users can safely upload and remove pdf files, other documents, spreadsheets, and images within their study groups for collaborative use and enhanced learning.

**Study Group Leaderboard**  
A motivational leaderboard displays study groups ranked by cumulative study hours against different study groups in the same modules.

**Secure Authentication**  
Firebase Authentication is used for OAuth-based user identity management, ensuring seamless and secure user login.

**Robust Backend Architecture**  
Built in Go, utilising a MySQL database for efficient data storage, and following scalable and modular design patterns for performance and maintainability.

**Containerized Deployment**  
Docker and Docker Compose ensure consistent local and production environments.

**CI Pipeline**  
GitHub Actions and custom scripts automate testing, code quality checks, and deployment.

## Step by Step

TLDR: to run the app locally it is enough to run `./scripts/run-localprod` from the main project directory. Once the server has fully started, the app can then be accessed through the browser by navigating to `localhost`.

Details about the provided run scripts are outlined [below](#detailed-description-of-the-run-scripts). Details regarding the integration with external services are provided in the [*Supplying credentials for external services* section](#supplying-credentials-for-external-services).

### Detailed description of the run scripts

The app server can be easily run using the scripts provided in the `scripts` directory. The `scripts` directory contains the necessary cross-platform implementations of the run scripts with support for Windows (PowerShell) and Unix-based systems (macOS/Linux via shell scripts). All the necessary Docker compose files are also located in the `scripts` directory.

For authentication, the server relies on Firebase, which affects how it can be run locally and in production. To support both scenarios, two run scripts are provided: `run-localprod` and `run-prod`. The `localprod` script uses a Firebase emulator, allowing the server to run entirely locally without requiring real Firebase credentials. Conversely, the `prod` script connects to the actual Firebase services, requiring valid Firebase credentials to be configured.

In order to run the app in `localprod` mode it is enough to run `./scripts/run-localprod` from the main project directory. The script will then run all the necessary services and expose the web-based frontend on port `80`. The app can then be accessed through the browser by navigating to `localhost`.

Similarly, the app can be run in `prod` mode by running `./scripts/run-prod` from the main project directory. Since it requires connection to the actual Firebase service, however, it is also needed to provide Firebase credentials files. A `serviceAccountKey.json` file must be placed in `backend/credentials` directory. The environmental variables (as specified in the `frontend/src/auth/firebase.tsx`) must also be defined in a `.env.production` file placed in the `frontend-web` directory. With the credentials files in place, the script will similarly run all the necessary services and expose the web-based frontend on port `80`. The app can then be accessed through the browser by navigating to `localhost`.

Docker compose is needed to run the application using the provided scripts.

Additional scripts are also provided for running the backend in `dev` mode (`./scripts/run-dev`), as well as for executing integration tests (`./scripts/run-tests-integration`).

### Supplying credentials for external services

In all cases, some features relying on external integrations (Gemini API, Google Calendar) require additional environmental variables to be set. This can be done by supplying the necessary values in a root directory `.env` file. In the absence of these environmental variables, the server will still run correctly, but with the particular features disabled.
The supported environmental variables are outlined below:

- `GEMINI_API_KEY` - necessary for the Gemini API integration for the intelligent chatbot.
- `GOOGLE_CALENDAR_HOST_CLIENT_ID` (OAuth 2.0 Client ID associated with a Google Cloud project), `GOOGLE_CALENDAR_HOST_CLIENT_SECRET` (OAuth 2.0 Client secret), `GOOGLE_CALENDAR_HOST_REFRESH_TOKEN` (a refresh token generated for the OAuth 2.0 Client) - necessary for the Google Calendar integration for study sessions.

The OAuth 2.0 Client ID and secret can be obtained from the Google Cloud console under the *APIs and Services->Credentials* section within the project.

## Backends

The application backend utilises a MySQL database and Firebase for authentication. If the app is run using the provided scripts, any ports for backend services used will also be exposed by default:

- the main backend server (REST and WebSocket API) can be accessed on port `8080`;
- the Firebase emulator (if used) can be accessed on port `4000` (emulator UI) and `9099` (auth emulator).

If external integrations are enabled by providing the necessary credentials as described in the [*Supplying credentials for external services* section](#supplying-credentials-for-external-services), the application backend also accesses the following external services:

- Google Calendar API
- Gemini API
