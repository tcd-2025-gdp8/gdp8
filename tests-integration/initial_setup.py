import time
import firebase_admin
from firebase_admin import auth


TEST_USER_EMAIL = "testuser@example.com"
TEST_USER_PASSWORD = "password123"
TEST_USER_NAME = "Mary Jane Holland"

TEST_USER2_EMAIL = "user2@example.com"
TEST_USER2_PASSWORD = "password123"

TEST_USER3_EMAIL = "user3@example.com"
TEST_USER3_PASSWORD = "password123"


if __name__ == "__main__":
    print("Initial setup started. Waiting for Firebase emulator.", flush=True)

    time.sleep(30)

    firebase_admin.initialize_app(None, { "projectId": "demo-backend" })

    # Create Test User 1
    try:
        user1 = auth.create_user(
            email=TEST_USER_EMAIL,
            password=TEST_USER_PASSWORD,
            display_name="Test User"
        )
        print("Created Test User 1")
    except Exception as e:
        print("Test User 1 already exists or error:", e)

    # Create Test User 2
    try:
        user2 = auth.create_user(
            email=TEST_USER2_EMAIL,
            password=TEST_USER2_PASSWORD,
            display_name="Test User 2"
        )
        print("Created Test User 2")
    except Exception as e:
        print("Test User 2 already exists or error:", e)

    # Create Test User 3
    try:
        user3 = auth.create_user(
            email=TEST_USER3_EMAIL,
            password=TEST_USER3_PASSWORD,
            display_name="Test User 3"
        )
        print("Created Test User 3")
    except Exception as e:
        print("Test User 3 already exists or error:", e)

    print("Initial setup completed.", flush=True)
