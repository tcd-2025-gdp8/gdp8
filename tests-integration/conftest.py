import pytest

from utils.custom_session import CustomSession
from utils.firebase_client_utils import get_firebase_access_token, get_uid
from initial_setup import TEST_USER_EMAIL, TEST_USER_PASSWORD, TEST_USER_NAME, TEST_USER2_EMAIL, TEST_USER2_PASSWORD, TEST_USER3_EMAIL, TEST_USER3_PASSWORD


@pytest.fixture(scope="session", autouse=True)
def access_token() -> str:
    token = get_firebase_access_token(TEST_USER_EMAIL, TEST_USER_PASSWORD)

    uid = get_uid(token)

    session = CustomSession(token)
    response = session.post("/user", json={
        "id": uid,
        "name": TEST_USER_NAME,
    })
    assert response.status_code == 201

    return token

@pytest.fixture(scope="session")
def user2_access_token() -> str:
    token = get_firebase_access_token(TEST_USER2_EMAIL, TEST_USER2_PASSWORD)
    decoded_token = jwt.decode(token, options={"verify_signature": False})
    uid = decoded_token.get("user_id")

    session = CustomSession(token)
    response = session.post("/user", json={
        "id": uid,
        "name": "User Two",
    })
    assert response.status_code == 201
    return token

@pytest.fixture(scope="session")
def user3_access_token() -> str:
    token = get_firebase_access_token(TEST_USER3_EMAIL, TEST_USER3_PASSWORD)
    decoded_token = jwt.decode(token, options={"verify_signature": False})
    uid = decoded_token.get("user_id")

    session = CustomSession(token)
    response = session.post("/user", json={
        "id": uid,
        "name": "User Three",
    })
    assert response.status_code == 201
    return token
