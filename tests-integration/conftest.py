import pytest

from initial_setup import TEST_USER, ADDITIONAL_TEST_USERS
from utils.custom_session import CustomSession
from utils.firebase_client_utils import get_firebase_access_token, get_uid


@pytest.fixture(scope="session", autouse=True)
def access_token() -> str:
    token = register_user(TEST_USER)
    return token

@pytest.fixture(scope="session")
def user2_access_token() -> str:
    token = register_user(ADDITIONAL_TEST_USERS[0])
    return token

@pytest.fixture(scope="session")
def user3_access_token() -> str:
    token = register_user(ADDITIONAL_TEST_USERS[1])
    return token


def register_user(user) -> str:
    token = get_firebase_access_token(user["email"], user["password"])

    uid = get_uid(token)

    session = CustomSession(token)
    response = session.post("/user", json={
        "id": uid,
        "name": user["name"],
    })
    assert response.status_code == 201

    return token
