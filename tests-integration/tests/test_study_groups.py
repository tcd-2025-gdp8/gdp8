from utils.custom_session import CustomSession


class TestStudyGroups:
    def test_get_study_group_does_not_exist(self, access_token):
        session = CustomSession(access_token)

        response = session.get("/study-groups/12345678")
        assert response.status_code == 404


    def test_create_study_group(self, access_token):
        session = CustomSession(access_token)

        response = session.post("/study-groups", json={
            "name": "New study group",
            "description": "Description of the new study group.",
            "type": "public"
        })

        assert response.status_code == 200

        study_group_id = response.json().get("id")

        response = session.get(f"/study-groups/{study_group_id}")
        assert response.status_code == 200
