import pytest

from utils.custom_session import CustomSession
from utils.firebase_client_utils import get_firebase_access_token, get_uid
from initial_setup import TEST_USER_EMAIL, TEST_USER_PASSWORD, TEST_USER_NAME


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
