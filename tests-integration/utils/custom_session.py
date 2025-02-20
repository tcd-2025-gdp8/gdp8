import os
import requests
from urllib.parse import urljoin


BACKEND_HOST = os.getenv("BACKEND_HOST", "127.0.0.1:8080")
BASE_URL = f"http://{BACKEND_HOST}/api"


class CustomSession(requests.Session):
    def __init__(self, access_token=None):
        super().__init__()

        self.headers.update({"Authorization": f"Bearer {access_token}"} if access_token else {})


    def request(self, method, url, **kwargs):
        full_url = urljoin(BASE_URL + "/", url.lstrip("/"))
        
        headers = kwargs.pop("headers", {})
        headers.update(self.headers)

        return super().request(method, full_url, headers=headers, **kwargs)
    