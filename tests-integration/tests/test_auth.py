from utils.custom_session import CustomSession


class TestAuth:
    def test_unauthorized(self):
        unauthorized_session = CustomSession()

        response = unauthorized_session.get("/study-groups")
        assert response.status_code == 401


    def test_authorized(self, access_token):
        authorized_session = CustomSession(access_token)

        response = authorized_session.get("/study-groups")
        assert response.status_code == 200
