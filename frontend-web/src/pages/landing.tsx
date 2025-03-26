import { useState, useEffect, useCallback } from "react";
import {
  Container,
  Typography,
  Grid,
  CircularProgress,
  Alert,
} from "@mui/material";
import { fetchUserStudyGroups, StudyGroup } from "../api/studyGroups";
import { getUserModules } from "../api/users";
import { useAuth } from "../auth/useAuth";
import StudyGroupCard from "../components/StudyGroupCard";
import { Module } from "../api/modules";
import { fetchAllUserStudySessions, Session } from "../api/studySessions";

export default function LandingPage() {
  const { token, user } = useAuth();
  const currentUserId = user?.uid;

  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [modulesList, setModulesList] = useState<Module[]>([]);
  const [studySessions, setStudySessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!token || !currentUserId) return;

    try {
      const [userStudyGroups, userModules] = await Promise.all([
        fetchUserStudyGroups(token, currentUserId),
        getUserModules(token, currentUserId),
      ]);

      const userStudySessions = await fetchAllUserStudySessions(token, userStudyGroups);
      setStudySessions(userStudySessions); // Set the fetched study sessions
      setStudyGroups(userStudyGroups);
      setModulesList(userModules || []);
    } catch (err) {
      console.error("Error fetching study groups or modules:", err);
      setError("Failed to load study groups.");
    } finally {
      setLoading(false);
    }
  }, [token, currentUserId]);

  useEffect(() => {
    void fetchData();
  }, [fetchData]);

  return (
    <Container maxWidth="md" style={{ marginTop: "20px", textAlign: "center" }}>
      <Typography variant="h4" style={{ marginBottom: "2rem" }}>
        Welcome to Your Study Hub
      </Typography>

      {/* Study Groups Section */}
      <Typography variant="h5" style={{ marginBottom: "1.5rem" }}>
        Your Study Groups
      </Typography>
      {loading && <CircularProgress />}
      {error && <Alert severity="error">{error}</Alert>}
      {!loading && studyGroups.length === 0 && (
        <Typography>No study groups found. Join or create one!</Typography>
      )}

      <Grid container spacing={2} style={{ marginBottom: "2rem" }}>
        {studyGroups.map((group) => {
          const module = modulesList.find((mod) => mod.id === group.moduleId);
          if (!module) return null;
          return (
            <Grid item xs={12} sm={6} md={4} key={group.id}>
              <StudyGroupCard
                group={group}
                module={module}
                currentUserId={currentUserId!}
                onUpdate={() => void fetchData()}
              />
            </Grid>
          );
        })}
      </Grid>

      {/* Study Sessions Section */}
      <Typography variant="h5" style={{ marginBottom: "1.5rem", marginTop: "3rem" }}>
        Upcoming Study Sessions
      </Typography>
      {loading && <CircularProgress />}
      {error && <Alert severity="error">{error}</Alert>}
      {!loading && studySessions.length === 0 && (
        <Typography>No upcoming study sessions found.</Typography>
      )}

      <Grid container spacing={2}>
        {studySessions.map((session) => (
          <Grid item xs={12} sm={6} md={4} key={session.id}>
            <Typography variant="h6">{session.name}</Typography>
            <Typography variant="body2">Date: {session.date?.toLocaleString()}</Typography>
            <Typography variant="body2">Start: {session.earliestTime}</Typography>
            <Typography variant="body2">End: {session.latestTime}</Typography>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
}
