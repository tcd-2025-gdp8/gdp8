from initial_setup import TEST_USER
from utils.custom_session import CustomSession
from utils.firebase_client_utils import get_uid


class TestStudyGroups:
    def test_get_study_group_does_not_exist(self, access_token):
        session = CustomSession(access_token)

        response = session.get("/study-groups/12345678")
        assert response.status_code == 404


    def test_create_study_group(self, access_token):
        session = CustomSession(access_token)
        uid = get_uid(access_token)

        study_group_data = {
            "name": "New study group",
            "description": "Description of the new study group.",
            "type": "public",
            "moduleId": 1,
            "maxMembers": 10
        }

        expected_member_data = {
            "id": uid,
            "name": TEST_USER["name"],
            "role": "admin"
        }

        response = session.post("/study-groups", json=study_group_data)

        assert response.status_code == 200
        data = response.json()

        study_group_id = data["id"]

        assert_study_group_data(data, study_group_data)
        assert len(data["members"]) == 1
        assert_study_group_has_member(data["members"], expected_member_data)

        response = session.get(f"/study-groups/{study_group_id}")
        assert response.status_code == 200

        data = response.json()
        assert_study_group_data(data, study_group_data)
        assert len(data["members"]) == 1
        assert_study_group_has_member(data["members"], expected_member_data)
    

    def test_get_relevant_study_groups(self, access_token):
        """Tests whether the `GET /study-groups` endpoint retrieves only the groups relevant to the user's chosen modules. """
        session = CustomSession(access_token)
        uid = get_uid(access_token)

        # set module preferences to include modules 3 and 4
        response = session.post(f"/user/{uid}/modules", json={
            "selectedModules": [3, 4]
        })
        assert response.status_code == 200

        # add two new groups for modules 3 and 4
        response = session.post("/study-groups", json={
            "name": "New study group 3",
            "description": "Description of the new study group.",
            "type": "public",
            "moduleId": 3,
            "maxMembers": 10
        })
        assert response.status_code == 200
        group3_id = response.json().get("id")
        response = session.post("/study-groups", json={
            "name": "New study group 4",
            "description": "Description of the new study group.",
            "type": "public",
            "moduleId": 4,
            "maxMembers": 10
        })
        assert response.status_code == 200
        group4_id = response.json().get("id")

        # assert that both groups are in the response
        response = session.get("/study-groups")
        assert response.status_code == 200
        data = response.json()
        assert any(group["id"] == group3_id for group in data)
        assert any(group["id"] == group4_id for group in data)

        # modify the preferences to only include group 3
        response = session.post(f"/user/{uid}/modules", json={
            "selectedModules": [3]
        })
        assert response.status_code == 200

        # assert that the group for module 3 is in the response, but the group for module 4 is not
        response = session.get("/study-groups")
        assert response.status_code == 200
        data = response.json()
        assert any(group["id"] == group3_id for group in data)
        assert not any(group["id"] == group4_id for group in data)


def assert_study_group_data(data, expected_data):
    assert data["name"] == expected_data["name"]
    assert data["description"] == expected_data["description"]
    assert data["type"] == expected_data["type"]
    assert data["moduleId"] == expected_data["moduleId"]
    assert data["maxMembers"] == expected_data["maxMembers"]


def assert_study_group_has_member(members, expected_member_data):
    return any(
        member["id"] == expected_member_data["id"] and 
        member["name"] == expected_member_data["name"] and
        member["role"] == expected_member_data["role"]
        for member in members
    )
