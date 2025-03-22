import { useState, useEffect } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  Button,
  FormControl,
  Select,
  MenuItem,
  Box,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  IconButton,
} from "@mui/material";
import Calendar from "react-calendar";
import "react-calendar/dist/Calendar.css";
import { useNavigate, useParams } from "react-router-dom";
import { createStudySession, fetchStudySessions, deleteStudySession } from "../api/studySessions"; // Still used for "Create Study Session"
import { useAuth } from "../auth/useAuth";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

interface Session {
  id: number;
  name: string;
  date: Date | null;
  earliestTime: string;
  latestTime: string;
}

export default function SchedulingUI() {
  const [sessions, setSessions] = useState<Session[]>([]);

  // --- State for Request Availability dialog ---
  const [openAvailability, setOpenAvailability] = useState(false);
  const [availabilityName, setAvailabilityName] = useState("");
  const [availabilityRequests, setAvailabilityRequests] = useState<Session[]>([]);
  const [availabilityDate, setAvailabilityDate] = useState<Date | null>(null);
  const [earliestTime, setEarliestTime] = useState("09:00");
  const [latestTime, setLatestTime] = useState("17:00");

  // --- State for Create Study Session dialog ---
  const [openSessionDialog, setOpenSessionDialog] = useState(false);
  const [sessionTitle, setSessionTitle] = useState("");
  const [sessionDate, setSessionDate] = useState<Date | null>(null);
  const [startTime, setStartTime] = useState("09:00");
  const [endTime, setEndTime] = useState("17:00");

  const navigate = useNavigate();
  const { token } = useAuth();
  const { groupId } = useParams();
  const numericGroupId = groupId ? Number(groupId) : NaN;

  function formatTime(dateObj: Date): string {
    return dateObj.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }

  async function loadSessions() {
    if (!token || isNaN(numericGroupId)) return;
    try {
      const backendSessions = await fetchStudySessions(token, numericGroupId);

      // The backend returns objects with `title`, `startTime`, `endTime` (and more).
      // Convert them into our local structure:
      const mapped: Session[] = backendSessions.map((bs) => ({
        id: bs.id,
        name: bs.title,
        date: bs.startTime, // convert "string" or "Date" from the backend as needed
        earliestTime: formatTime(bs.startTime),
        latestTime: formatTime(bs.endTime),
      }));

      setSessions(mapped);
    } catch (err) {
      console.error("Error fetching sessions:", err);
    }
  }

  useEffect(() => {
    loadSessions();
  }, [token, numericGroupId]);

  /**
   * Utility to calculate duration minutes from start/end times on the same date.
   */
  function calculateDurationMinutes(dateObj: Date, start: string, end: string): number {
    const [startHour, startMin] = start.split(":").map(Number);
    const [endHour, endMin] = end.split(":").map(Number);

    const startDateTime = new Date(dateObj.getTime());
    startDateTime.setHours(startHour, startMin, 0, 0);

    const endDateTime = new Date(dateObj.getTime());
    endDateTime.setHours(endHour, endMin, 0, 0);

    return Math.max(0, (endDateTime.getTime() - startDateTime.getTime()) / 60000);
  }

  /**
   * Generate hourly time slots between a start and end time (same pattern as File 1).
   */
  const generateTimeSlots = (start: string, end: string): string[] => {
    const slots: string[] = [];
    const current = new Date();
    const [startHour, startMin] = start.split(":").map(Number);
    const [endHour, endMin] = end.split(":").map(Number);

    current.setHours(startHour, startMin, 0, 0);
    const endTime = new Date();
    endTime.setHours(endHour, endMin, 0, 0);

    while (current <= endTime) {
      slots.push(
        current.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
      );
      current.setMinutes(current.getMinutes() + 60); // 1-hour intervals
    }
    return slots;
  };

  // Delete or edit a session in the "Scheduled Sessions" box
  const handleDeleteSession = async (index: number) => {
    if (!token) return;
    const sessionToDelete = sessions[index];
    try {
      // Call backend to delete
      await deleteStudySession(token, sessionToDelete.id);

      // Remove from local state
      setSessions((prev) => prev.filter((_, i) => i !== index));
    } catch (error) {
      console.error("Error deleting study session:", error);
    }
  };
  const handleEditSession = (index: number) => {
    //TO DO implement Backend for updating sessions
    console.log("Edit session:", sessions[index]);
  };

  /**
   * 1. Handler for the "Request Availability" form
   *    (now mimics File 1's createSession behavior).
   */
  const handleAvailabilityRequest = () => {
    // Only proceed if we have an event name and date
    if (!availabilityName || !availabilityDate) return;

    // Create new availability request
    const newRequest: Session = {
      id: Date.now(), //update this to a real ID!!
      name: availabilityName,
      date: availabilityDate,
      earliestTime,
      latestTime,
    };

    // Add to local state
    setAvailabilityRequests((prev) => [...prev, newRequest]);

    // Store the generated time slots in localStorage
    localStorage.setItem(
      "selectedTimeSlots",
      JSON.stringify(generateTimeSlots(earliestTime, latestTime))
    );

    // Reset fields
    setAvailabilityName("");
    setAvailabilityDate(null);
    setEarliestTime("09:00");
    setLatestTime("17:00");
    setOpenAvailability(false);

    // Finally navigate to /availability
    navigate(`/study-groups/${groupId}/availability`);
  };

  /**
   * 2. Handler for the "Create Study Session" form
   *    (uses createStudySession API).
   */
  const handleCreateStudySession = async () => {
    if (!sessionTitle || !sessionDate || !token || isNaN(numericGroupId)) return;

    try {
      // Calculate total duration in minutes
      const durationMinutes = calculateDurationMinutes(
        sessionDate,
        startTime,
        endTime
      );

      // Build a startTime Date from sessionDate + startTime
      const [sHour, sMin] = startTime.split(":").map(Number);
      const combinedStart = new Date(sessionDate.getTime());
      combinedStart.setHours(sHour, sMin, 0, 0);

      const sessionDetails = {
        title: sessionTitle,
        startTime: combinedStart,
        durationMinutes,
      };

      // Create the study session in the backend
      const newBackendSession = await createStudySession(
        token,
        numericGroupId,
        sessionDetails
      );

      // Update local state
      setSessions((prev) => [
        ...prev,
        {
          id: newBackendSession.id,
          name: newBackendSession.title,
          date: newBackendSession.startTime,
          earliestTime: startTime,
          latestTime: endTime,
        },
      ]);

      // Reset form
      setSessionTitle("");
      setSessionDate(null);
      setStartTime("09:00");
      setEndTime("17:00");
      setOpenSessionDialog(false);

      // Optionally navigate after creation
      // navigate(`/study-groups/${groupId}/details`);
    } catch (error) {
      console.error("Error creating study session:", error);
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Scheduled Sessions</h1>
      <Box sx={{ display: "flex", flexDirection: "column", alignItems: "flex-start", gap: 2, mt: 4 }}>
        {/* Create Session */}
        <Button variant="contained" color="primary" onClick={() => setOpenSessionDialog(true)}>
          Create Session
        </Button>

        {/* Dark-blue box for SESSIONS */}
        <Paper
          sx={{
            p: 2,
            bgcolor: "#3b5998",
            borderRadius: "8px",
            maxHeight: "220px",
            overflowY: "auto",
            width: "100%",
            maxWidth: "1200px",
          }}
        >
          <List>
            {sessions.map((session, index) => (
              <ListItem
                key={session.id}
                sx={{
                  backgroundColor: "white",
                  borderRadius: "4px",
                  mb: 1,
                  padding: 1.5,
                  "&:hover": {
                    backgroundColor: "#f5f5f5",
                  },
                }}
              >
                <ListItemText
                  primary={`${session.name}, ${session.date ? session.date.toLocaleDateString() : ""
                    } ${session.earliestTime} - ${session.latestTime}`}
                />
                <Box sx={{ display: "flex", gap: 1 }}>
                  <IconButton onClick={() => handleEditSession(index)} sx={{ color: "primary.main" }}>
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDeleteSession(index)} sx={{ color: "error.main" }}>
                    <DeleteIcon />
                  </IconButton>
                </Box>
              </ListItem>
            ))}
          </List>
        </Paper>

        {/* Title for availability requests */}
        <h1 className="text-2xl font-bold">Availability Requested</h1>

        {/* Request Availability */}
        <Button variant="contained" color="primary" onClick={() => setOpenAvailability(true)}>
          Request Availability
        </Button>
        <Paper
          sx={{
            p: 2,
            bgcolor: "#3b5998",
            borderRadius: "8px",
            maxHeight: "220px",
            overflowY: "auto",
            width: "100%",
            maxWidth: "1200px",
          }}
        >
          <List>
            {availabilityRequests.map((req) => (
              <ListItem
                key={req.id}
                sx={{
                  backgroundColor: "white",
                  borderRadius: "4px",
                  mb: 1,
                  padding: 1.5,
                  "&:hover": {
                    backgroundColor: "#f5f5f5",
                  },
                }}
              >
                <ListItemText
                  primary={req.name}
                  secondary={
                    <>
                      <Typography variant="body2" color="text.primary">
                        Date: {req.date?.toLocaleDateString()}
                      </Typography>
                      <Typography variant="body2" color="text.primary">
                        Time: {req.earliestTime} - {req.latestTime}
                      </Typography>
                    </>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Paper>
      </Box>

      {/* DIALOG 1: Request Study Session Availability (No API call) */}
      <Dialog open={openAvailability} onClose={() => setOpenAvailability(false)}>
        <DialogContent className="max-w-lg p-6">
          <DialogTitle>Request Study Session Availability</DialogTitle>
          <div className="space-y-4">
            <TextField
              fullWidth
              label="Event Name"
              value={availabilityName}
              onChange={(e) => setAvailabilityName(e.target.value)}
            />
            <div>
              <Typography variant="subtitle1" gutterBottom>
                Select Dates
              </Typography>
              <Calendar
                onChange={(newDate) => setAvailabilityDate(newDate as Date | null)}
                value={availabilityDate}
                tileDisabled={({ date }) => date < new Date()}
              />
            </div>
            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle1" gutterBottom>
                Available Time Range
              </Typography>
              <Box sx={{ display: "flex", gap: 2 }}>
                <FormControl fullWidth variant="outlined">
                  <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                    Earliest Time
                  </Typography>
                  <Select
                    value={earliestTime}
                    onChange={(e) => setEarliestTime(e.target.value)}
                    sx={{
                      height: "56px",
                      "& .MuiOutlinedInput-notchedOutline": { borderColor: "#1976d2" },
                      "&:hover .MuiOutlinedInput-notchedOutline": {
                        borderColor: "#1976d2",
                      },
                    }}
                  >
                    {[...Array(24).keys()].map((hour) => (
                      <MenuItem key={hour} value={`${hour.toString().padStart(2, "0")}:00`}>
                        {`${hour.toString().padStart(2, "0")}:00`}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                <FormControl fullWidth variant="outlined">
                  <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                    Latest Time
                  </Typography>
                  <Select
                    value={latestTime}
                    onChange={(e) => setLatestTime(e.target.value)}
                    sx={{
                      height: "56px",
                      "& .MuiOutlinedInput-notchedOutline": { borderColor: "#1976d2" },
                      "&:hover .MuiOutlinedInput-notchedOutline": {
                        borderColor: "#1976d2",
                      },
                    }}
                  >
                    {[...Array(24).keys()].map((hour) => (
                      <MenuItem key={hour} value={`${hour.toString().padStart(2, "0")}:00`}>
                        {`${hour.toString().padStart(2, "0")}:00`}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Box>
            </Box>
            <Button
              variant="contained"
              color="primary"
              fullWidth
              onClick={handleAvailabilityRequest}
              disabled={!availabilityName || !availabilityDate}
              sx={{ mt: 3, py: 1.5, fontSize: "1rem" }}
            >
              CREATE AVAILABILITY REQUEST
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* DIALOG 2: Create Study Session (With API Call) */}
      <Dialog open={openSessionDialog} onClose={() => setOpenSessionDialog(false)}>
        <DialogContent className="max-w-lg p-6">
          <DialogTitle>Create Study Session</DialogTitle>
          <div className="space-y-4">
            <TextField
              fullWidth
              label="Session Title"
              value={sessionTitle}
              onChange={(e) => setSessionTitle(e.target.value)}
            />
            <div>
              <Typography variant="subtitle1" gutterBottom>
                Select Date
              </Typography>
              <Calendar
                onChange={(newDate) => setSessionDate(newDate as Date | null)}
                value={sessionDate}
                tileDisabled={({ date }) => date < new Date()}
              />
            </div>
            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle1" gutterBottom>
                Study Session Duration
              </Typography>
              <Box sx={{ display: "flex", gap: 2 }}>
                <FormControl fullWidth variant="outlined">
                  <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                    Start Time
                  </Typography>
                  <Select
                    value={startTime}
                    onChange={(e) => setStartTime(e.target.value)}
                    sx={{
                      height: "56px",
                      "& .MuiOutlinedInput-notchedOutline": { borderColor: "#1976d2" },
                      "&:hover .MuiOutlinedInput-notchedOutline": {
                        borderColor: "#1976d2",
                      },
                    }}
                  >
                    {[...Array(24).keys()].map((hour) => (
                      <MenuItem key={hour} value={`${hour.toString().padStart(2, "0")}:00`}>
                        {`${hour.toString().padStart(2, "0")}:00`}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                <FormControl fullWidth variant="outlined">
                  <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                    End Time
                  </Typography>
                  <Select
                    value={endTime}
                    onChange={(e) => setEndTime(e.target.value)}
                    sx={{
                      height: "56px",
                      "& .MuiOutlinedInput-notchedOutline": { borderColor: "#1976d2" },
                      "&:hover .MuiOutlinedInput-notchedOutline": {
                        borderColor: "#1976d2",
                      },
                    }}
                  >
                    {[...Array(24).keys()].map((hour) => (
                      <MenuItem key={hour} value={`${hour.toString().padStart(2, "0")}:00`}>
                        {`${hour.toString().padStart(2, "0")}:00`}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Box>
            </Box>
            <Button
              variant="contained"
              color="primary"
              fullWidth
              onClick={handleCreateStudySession}
              disabled={!sessionTitle || !sessionDate}
              sx={{ mt: 3, py: 1.5, fontSize: "1rem" }}
            >
              CREATE STUDY SESSION
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
