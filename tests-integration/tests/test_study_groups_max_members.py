from utils.custom_session import CustomSession

class TestStudyGroupsMaxMembers:
    def test_max_members_limit(self, access_token, user2_access_token, user3_access_token):
        # Step 1: Create a study group with maxMembers=2 using the admin token (access_token)
        admin_session = CustomSession(access_token)
        create_response = admin_session.post("/study-groups", json={
            "name": "Max Members Test Group",
            "description": "Testing maxMembers enforcement",
            "type": "public",
            "moduleId": 1,
            "maxMembers": 2
        })
        assert create_response.status_code == 200, f"Group creation failed: {create_response.text}"
        group_data = create_response.json()
        group_id = group_data.get("id")
        assert group_id, "No group id returned"

        # Step 2: User2 sends a join request; should succeed.
        user2_session = CustomSession(user2_access_token)
        join_response_user2 = user2_session.post(f"/study-groups/{group_id}/request-to-join")
        assert join_response_user2.status_code == 200, f"User2 join failed: {join_response_user2.text}"

        # Step 3: User3 sends a join request; should fail (group is full).
        user3_session = CustomSession(user3_access_token)
        join_response_user3 = user3_session.post(f"/study-groups/{group_id}/request-to-join")
        assert join_response_user3.status_code == 400, (
            f"Expected failure joining full group, but got {join_response_user3.status_code}: {join_response_user3.text}"
        )
        # Optionally, check that the error message mentions "full"
        assert "full" in join_response_user3.text.lower(), (
            f"Error message does not indicate group is full: {join_response_user3.text}"
        )
