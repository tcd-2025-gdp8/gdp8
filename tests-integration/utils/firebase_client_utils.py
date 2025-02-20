import os
import requests


FIREBASE_AUTH_EMULATOR_HOST = os.getenv("FIREBASE_AUTH_EMULATOR_HOST", "127.0.0.1:9099")
FIREBASE_EMULATOR_SIGNIN_URI = f"http://{FIREBASE_AUTH_EMULATOR_HOST}/identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=foo"


def get_firebase_access_token(email: str, password: str) -> str:
    payload = {
        "email": email,
        "password": password,
        "returnSecureToken": True
    }
    
    response = requests.post(FIREBASE_EMULATOR_SIGNIN_URI, json=payload)
    data = response.json()
    
    if "idToken" in data:
        return data["idToken"]
    else:
        raise Exception(f"Error retrieving token: {data.get('error', 'Unknown error')}")
