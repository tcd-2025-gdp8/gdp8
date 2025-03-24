from utils.custom_session import CustomSession


class TestGoogleCalendarInvite:
    def test_create_calendar_invite(self, access_token):
        session = CustomSession(access_token)
        response = session.post("/api/calendar/invite", json={
            "summary": "Integration Test Meeting",
            "location": "Virtual",
            "description": "Test event for Google Calendar integration.",
            "start": {
                "dateTime": "2025-03-30T10:00:00Z",
                "timeZone": "UTC"
            },
            "end": {
                "dateTime": "2025-03-30T11:00:00Z",
                "timeZone": "UTC"
            },
            "attendees": [
                {"email": "test1@example.com"},
                {"email": "test2@example.com"}
            ]
        })

        assert response.status_code == 200, f"Expected 200, got {response.status_code}: {response.text}"
        data = response.json()
        assert "id" in data, f"No event ID returned: {data}"
        assert data.get("summary") == "Integration Test Meeting", "Unexpected event summary"
