import jwt
import pytest

from utils.firebase_client_utils import get_firebase_access_token
from utils.custom_session import CustomSession
from initial_setup import TEST_USER_EMAIL, TEST_USER_PASSWORD


@pytest.fixture(scope="session", autouse=True)
def access_token() -> str:
    token = get_firebase_access_token(TEST_USER_EMAIL, TEST_USER_PASSWORD)

    decoded_token = jwt.decode(token, options={"verify_signature": False})
    uid = decoded_token.get("user_id")

    session = CustomSession(token)
    response = session.post("/user", json={
        "id": uid,
        "name": "Mary Jane Holland",
    })
    assert response.status_code == 201

    return token
