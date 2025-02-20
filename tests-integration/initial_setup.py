import time
import firebase_admin
from firebase_admin import auth


TEST_USER_EMAIL = "testuser@example.com"
TEST_USER_PASSWORD = "password123"


if __name__ == "__main__":
    print("Initial setup started. Waiting for Firebase emulator.", flush=True)

    time.sleep(30)

    firebase_admin.initialize_app(None, { "projectId": "demo-backend" })

    user = auth.create_user(
        email=TEST_USER_EMAIL,
        password=TEST_USER_PASSWORD,
        display_name="Test User"
    )

    print("Initial setup completed.", flush=True)
