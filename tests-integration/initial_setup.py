import time
import firebase_admin
from firebase_admin import auth


TEST_USER = {
    "email": "testuser@example.com",
    "password": "password123",
    "name": "Mary Jane Holland"
}

ADDITIONAL_TEST_USERS = [
    {
        "email": "user2@example.com",
        "password": "password123",
        "name": "User Two"
    },
    {
        "email": "user3@example.com",
        "password": "password123",
        "name": "User Three"
    }
]


if __name__ == "__main__":
    print("Initial setup started. Waiting for Firebase emulator.", flush=True)

    time.sleep(30)

    firebase_admin.initialize_app(None, { "projectId": "demo-backend" })

    auth.create_user(
        email=TEST_USER["email"],
        password=TEST_USER["password"],
        display_name="Test User"
    )

    for user in ADDITIONAL_TEST_USERS:
        auth.create_user(
            email=user["email"],
            password=user["password"],
            display_name=user["name"]
        )

    print("Initial setup completed.", flush=True)
