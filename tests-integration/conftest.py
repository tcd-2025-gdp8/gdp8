import pytest

from utils.firebase_client_utils import get_firebase_access_token
from initial_setup import TEST_USER_EMAIL, TEST_USER_PASSWORD


@pytest.fixture(scope="session", autouse=True)
def access_token() -> str:
    return get_firebase_access_token(TEST_USER_EMAIL, TEST_USER_PASSWORD)
