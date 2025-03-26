import { useState, useEffect, useCallback } from "react";
import {
  Container,
  Typography,
  Grid,
  CircularProgress,
  Alert,
} from "@mui/material";
import { fetchUserStudyGroups, StudyGroup } from "../api/studyGroups"; // Import the new function
import { getUserModules } from "../api/users";
import { useAuth } from "../auth/useAuth";
import StudyGroupCard from "../components/StudyGroupCard";
import { Module } from "../api/modules";

export default function LandingPage() {
  const { token, user } = useAuth();
  const currentUserId = user?.uid;

  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [modulesList, setModulesList] = useState<Module[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!token || !currentUserId) return;

    try {
      const [userStudyGroups, userModules] = await Promise.all([
        fetchUserStudyGroups(token, currentUserId), // Use the new function
        getUserModules(token, currentUserId),
      ]);

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

      {/* Study Sessions Section (Placeholder for now) */}
      <Typography variant="h5" style={{ marginBottom: "1.5rem", marginTop: "3rem" }}>
        Upcoming Study Sessions
      </Typography>
      <Typography variant="body1" style={{ marginBottom: "2rem" }}>
        This section will display upcoming study sessions in the future.
      </Typography>
    </Container>
  );
}