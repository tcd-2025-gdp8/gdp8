import { useState } from "react";
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
} from "@mui/material";
import Calendar from "react-calendar";
import "react-calendar/dist/Calendar.css";
import { useNavigate } from "react-router-dom";

interface Session {
  name: string;
  date: Date | null;
  earliestTime: string;
  latestTime: string;
}

export default function SchedulingUI() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [open, setOpen] = useState(false);
  const [sessionName, setSessionName] = useState("");
  const [date, setDate] = useState<Date | null>(null);
  const [earliestTime, setEarliestTime] = useState("09:00");
  const [latestTime, setLatestTime] = useState("17:00");
  const navigate = useNavigate();

  const generateTimeSlots = (start: string, end: string) => {
    const slots = [];
    const current = new Date();
    const [startHour, startMin] = start.split(":").map(Number);
    const [endHour, endMin] = end.split(":").map(Number);

    current.setHours(startHour, startMin, 0, 0);
    const endTime = new Date();
    endTime.setHours(endHour, endMin, 0, 0);

    while (current <= endTime) {
      slots.push(current.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }));
      current.setMinutes(current.getMinutes() + 60); // 1-hour intervals
    }
    return slots;
  };

  const createSession = () => {
    if (sessionName && date) {
      const newSession: Session = {
        name: sessionName,
        date,
        earliestTime,
        latestTime,
      };

      setSessions([...sessions, newSession]);
      localStorage.setItem("selectedTimeSlots", JSON.stringify(generateTimeSlots(earliestTime, latestTime)));
      setOpen(false);
      setSessionName("");
      setDate(null);
      setEarliestTime("09:00");
      setLatestTime("17:00");
      void navigate("/availability");
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Scheduled Sessions</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
        {sessions.map((session, index) => (
          <div key={index} className="p-4 border rounded-lg shadow-md">
            <h2 className="text-lg font-semibold">{session.name}</h2>
            <p>Date: {session.date?.toLocaleDateString()}</p>
            <p>Time: {session.earliestTime} - {session.latestTime}</p>
          </div>
        ))}
      </div>
      <Button variant="contained" color="primary" onClick={() => setOpen(true)} className="mt-4">
        Create Session
      </Button>

      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogContent className="max-w-lg p-6">
          <DialogTitle>Create a Study Session</DialogTitle>
          <div className="space-y-4">
            <TextField
              fullWidth
              label="Event Name"
              value={sessionName}
              onChange={(e) => setSessionName(e.target.value)}
            />
            <div>
              <Typography variant="subtitle1" gutterBottom>
                Select Dates
              </Typography>
              <Calendar
                onChange={(newDate) => setDate(newDate as Date | null)}
                value={date}
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
                      "&:hover .MuiOutlinedInput-notchedOutline": { borderColor: "#1976d2" }
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
                      "&:hover .MuiOutlinedInput-notchedOutline": { borderColor: "#1976d2" }
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
              onClick={createSession}
              disabled={!sessionName || !date}
              sx={{ mt: 3, py: 1.5, fontSize: "1rem" }}
            >
              NEXT: SCHEDULE AVAILABILITY
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
