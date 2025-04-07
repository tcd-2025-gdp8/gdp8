import { useState, useEffect, useRef, useCallback, Fragment } from "react";
import {
  Container,
  Typography,
  Button,
  Grid,
  MenuItem,
  Select,
  SelectChangeEvent,
  FormControl,
  InputLabel,
} from "@mui/material";
//import { Link } from "react-router-dom";
import { Module } from "../api/modules";
import { fetchStudyGroups, StudyGroup } from "../api/studyGroups";
import { getUserModules } from "../api/users";
import { useAuth } from "../auth/useAuth";
import StudyGroupCard from "../components/StudyGroupCard";
import DialogCreateStudyGroup from "../components/DialogCreateStudyGroup";
import { DialogRef } from "../utils/useDialog";


export default function StudyGroupsPage() {
  const { token, user } = useAuth();
  const currentUserId = user?.uid;

  const createStudyGroupDialogRef = useRef<DialogRef>(null);

  const [modulesList, setModulesList] = useState<Module[]>([]);
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);

  const [selectedFilterModuleId, setSelectedFilterModuleId] = useState<number>(-1);


  const fetchData = useCallback(async () => {
    const fetchRelevantStudyGroups = async () => {
      const fetchedStudyGroups = await fetchStudyGroups(token);
      setStudyGroups(fetchedStudyGroups);
    };

    const fetchModules = async () => {
      if (!currentUserId) return;
      const data = await getUserModules(token, currentUserId);
      if (data) {
        setModulesList([{ id: -1, code: "All", name: "All" }, ...data]);
      } else {
        setModulesList([{ id: -1, code: "All", name: "All" }]);
      }
    }

    try {
      await Promise.all([
        fetchRelevantStudyGroups(),
        fetchModules()
      ]);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  }, [token, currentUserId]);

  useEffect(() => {
    void fetchData();
  }, [fetchData]);


  let filteredGroups = studyGroups;
  if (selectedFilterModuleId !== -1) {
    filteredGroups = studyGroups.filter((group) => group.moduleId === selectedFilterModuleId);
  }


  return (
    <Container maxWidth="md" style={{ marginTop: "20px", position: "relative", textAlign: "center" }}>
      <Typography variant="h4" style={{ marginBottom: "2rem" }}>
        Study Groups
      </Typography>
      <FormControl fullWidth style={{ marginBottom: "2rem", maxWidth: "600px", margin: "0 auto 2em auto" }}>
        <InputLabel>Filter by Module</InputLabel>
        <Select
          value={selectedFilterModuleId}
          onChange={(e: SelectChangeEvent<number | string>) => setSelectedFilterModuleId(Number(e.target.value))}
          label="Filter by Module"
        >
          {modulesList.map((module) => (
            <MenuItem key={module.id} value={module.id}>
              {module.code === "All" ? module.name : `${module.code}: ${module.name}`}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <Grid container spacing={2} style={{ marginBottom: "2rem" }}>
        {filteredGroups.map((group) => {
          const module = modulesList.find((module) => module.id === group.moduleId);
          if (!module) return <Fragment key={group.id}></Fragment>;
          //const isMember = group.members?.some(
          //(member) => member.id === currentUserId
          //);
          return (
            <Grid item xs={12} sm={6} md={4} key={group.id} style={{ minWidth: "280px" }}>
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

      <Button
        variant="contained"
        color="primary"
        onClick={() => createStudyGroupDialogRef.current?.openDialog()}
        sx={{ maxWidth: "600px", margin: "2rem auto", display: "block", width: "100%" }}
      >
        Create a Study Group
      </Button>
      
      
      <DialogCreateStudyGroup
        modules={modulesList.filter((module) => module.id !== -1)}
        ref={createStudyGroupDialogRef}
        onUpdate={() => void fetchData()}
      />
      <Button
  variant="contained"
  color="primary"
  onClick={() => createStudyGroupDialogRef.current?.openDialog()}
  sx={{ maxWidth: "600px", margin: "2rem auto", display: "block", width: "100%" }}
>
  Create a Study Group
</Button>

<Button
  variant="outlined"
  color="secondary"
  onClick={async () => {
    if (!token) return;
    try {
      await navigator.clipboard.writeText(token);
      alert("Firebase token copied to clipboard!");
    } catch (err) {
      console.error("Failed to copy token", err);
    }
  }}
  sx={{ marginBottom: "1.5rem" }}
>
  Copy Firebase Token
</Button>


    </Container>
  );
}