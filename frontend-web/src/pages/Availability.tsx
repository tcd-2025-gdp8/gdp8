import { useState, useEffect, useRef } from "react";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@mui/material";
import { useNavigate, useParams, useLocation } from "react-router-dom";
import {
  fetchCurrentAvailabilityRequestsByGroupId,
  upsertAvailabilityEntries,
  createStudySession,
  StudySessionAvailabilityEntry
} from "../api/studySessions";
import { fetchStudyGroupById } from "../api/studyGroups";
import { useAuth } from "../auth/useAuth";

export default function AvailabilitySelection() {
  const navigate = useNavigate();
  const location = useLocation();
  const { token, user } = useAuth();
  const { groupId } = useParams();
  const numericGroupId = groupId ? Number(groupId) : NaN;

  // Store the array of groupMembers as { id, name }
  const [groupMembers, setGroupMembers] = useState<{ id: string; name: string }[]>([]);

  // Use the actual backend user ID for the current user
  const currentUserId = user?.uid ?? "UnknownUid";

  const [availabilityRequestId, setAvailabilityRequestId] = useState<number | null>(null);
  const [requestStartDate, setRequestStartDate] = useState<Date | null>(null);
  const [requestEndDate, setRequestEndDate] = useState<Date | null>(null);
  const [timeSlots, setTimeSlots] = useState<string[]>([]);

  // Store the backend availability entries
  const [backendAvailabilityEntries, setBackendAvailabilityEntries] = useState<StudySessionAvailabilityEntry[]>([]);

  // A helper to format a Date into "hh:mm AM/PM"
  const formatTime = (date: Date): string =>
    date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });

  // Return userIds (string[]) for everyone who has availability matching the given time slot
  const getBackendAvailableUsers = (timeSlot: string): string[] => {
    return Array.from(
      new Set(
        backendAvailabilityEntries
          .filter((entry) => formatTime(entry.availabilityStart) === timeSlot)
          .map((entry) => entry.userId)
      )
    );
  };

  // Return the array of { id, name } that are available for a given time slot
  function getAvailableUsers(time: string): { id: string; name: string }[] {
    const availableUserIds = getBackendAvailableUsers(time);
    return groupMembers.filter((member) => availableUserIds.includes(member.id));
  }

  // Generate hourly time slots between two dates
  function generateHourlyTimeSlots(start: Date, end: Date): string[] {
    const slots: string[] = [];
    const current = new Date(start.getTime());
    // Round down to the nearest hour
    current.setMinutes(0, 0, 0);

    while (current < end) {
      slots.push(
        current.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
      );
      current.setHours(current.getHours() + 1);
    }
    return slots;
  }

  // Generate time slots whenever the start/end changes
  useEffect(() => {
    if (requestStartDate && requestEndDate && requestStartDate < requestEndDate) {
      const newSlots = generateHourlyTimeSlots(requestStartDate, requestEndDate);
      setTimeSlots(newSlots);
    }
  }, [requestStartDate, requestEndDate]);

  useEffect(() => {
    console.log("availabilityRequestId:", availabilityRequestId);
  }, [availabilityRequestId]);

  // Load group members
  useEffect(() => {
    if (!token || isNaN(numericGroupId)) return;
    fetchStudyGroupById(token, numericGroupId)
      .then((group) => {
        setGroupMembers(group.members);
      })
      .catch((err) => console.error("Error fetching group details", err));
  }, [token, numericGroupId]);

  // Check for query parameters: reqId, start, end
  useEffect(() => {
    if (!token || isNaN(numericGroupId)) return;

    const queryParams = new URLSearchParams(location.search);
    const reqIdParam = queryParams.get("reqId");
    const startParam = queryParams.get("start");
    const endParam = queryParams.get("end");

    // 1) If we have a reqId in the query string, fetch that request
    if (reqIdParam) {
      const wantedId = Number(reqIdParam);
      if (!isNaN(wantedId)) {
        fetchCurrentAvailabilityRequestsByGroupId(token, numericGroupId)
          .then((requests) => {
            const found = requests.find((r) => r.id === wantedId);
            if (found) {
              setAvailabilityRequestId(found.id);
              setRequestStartDate(new Date(found.availabilityPeriodStart));
              setRequestEndDate(new Date(found.availabilityPeriodEnd));
              setBackendAvailabilityEntries(found.availabilityEntries);
            } else {
              console.warn("No availability request found for reqId=", wantedId);
            }
          })
          .catch((err) => console.error("Error fetching requests:", err));
        return;
      }
    }

    // 2) If we have start/end in the URL, use them
    if (startParam && endParam) {
      const startDate = new Date(startParam);
      const endDate = new Date(endParam);
      setRequestStartDate(startDate);
      setRequestEndDate(endDate);

      fetchCurrentAvailabilityRequestsByGroupId(token, numericGroupId)
        .then((requests) => {
          if (requests.length > 0) {
            setAvailabilityRequestId(requests[0].id);
            setBackendAvailabilityEntries(requests[0].availabilityEntries);
          }
        })
        .catch((err) => console.error("Error fetching requests:", err));
      return;
    }

    // 3) Otherwise, just fetch the first request
    fetchCurrentAvailabilityRequestsByGroupId(token, numericGroupId)
      .then((requests) => {
        if (requests.length > 0) {
          setAvailabilityRequestId(requests[0].id);
          setRequestStartDate(new Date(requests[0].availabilityPeriodStart));
          setRequestEndDate(new Date(requests[0].availabilityPeriodEnd));
          setBackendAvailabilityEntries(requests[0].availabilityEntries);
        } else {
          // Possibly create a new request if needed
        }
      })
      .catch((err) => console.error("Error:", err));
  }, [token, numericGroupId, location.search]);

  // Persist local availability in localStorage
  const [availability, setAvailability] = useState<Record<string, Set<string>>>(() => {
    const stored = localStorage.getItem("userAvailability");
    if (stored) {
      try {
        const parsed = JSON.parse(stored) as Record<string, string[]>;
        const result: Record<string, Set<string>> = {};
        for (const key in parsed) {
          if (Object.prototype.hasOwnProperty.call(parsed, key)) {
            result[key] = new Set(parsed[key]);
          }
        }
        return result;
      } catch (err) {
        console.error("Error parsing local availability", err);
        return {};
      }
    }
    return {};
  });

  const [open, setOpen] = useState(false);
  const [isSelecting, setIsSelecting] = useState(false);
  const [selectionMode, setSelectionMode] = useState<"add" | "remove">("add");
  const gridRef = useRef<HTMLDivElement>(null);

  // Mouse handling for click-and-drag selection
  useEffect(() => {
    const handleMouseUp = () => setIsSelecting(false);
    document.addEventListener("mouseup", handleMouseUp);
    return () => document.removeEventListener("mouseup", handleMouseUp);
  }, []);

  // Toggle a single time slot for the current user
  const toggleAvailability = (time: string, mode?: "add" | "remove") => {
    const action = mode ?? selectionMode;
    setAvailability((prev) => {
      const newAvail = { ...prev };
      const currentAvail = new Set(newAvail[currentUserId]);
      if (action === "add") {
        currentAvail.add(time);
      } else {
        currentAvail.delete(time);
      }
      newAvail[currentUserId] = currentAvail;
      return newAvail;
    });
  };

  const handleMouseDown = (time: string) => {
    const isAvail = availability[currentUserId]?.has(time);
    setSelectionMode(isAvail ? "remove" : "add");
    toggleAvailability(time);
    setIsSelecting(true);
  };

  const handleMouseEnter = (time: string) => {
    if (isSelecting) {
      toggleAvailability(time, selectionMode);
    }
  };

  const handleMouseUp = () => {
    setIsSelecting(false);
  };

  // Compute "best times" based on how many members are available
  const getBestTimeSlots = () => {
    const timesWithCounts = timeSlots.map((time) => ({
      time,
      count: getAvailableUsers(time).length,
    }));

    return timesWithCounts
      .sort((a, b) => b.count - a.count)
      .filter((item) => item.count > 0)
      .slice(0, 3);
  };

  // Save local availability to the backend
  const saveAvailabilityToBackend = async () => {
    if (!token || !availabilityRequestId || !requestStartDate || !requestEndDate) {
      console.error("Missing token, availabilityRequestId, requestStartDate, or requestEndDate");
      return;
    }

    const selectedTimes = Array.from(availability[currentUserId] || []);

    const entries = selectedTimes
      .map((slotString) => {
        const [timePart, amPm] = slotString.split(" ");
        const [rawHour, rawMin] = timePart.split(":");
        let hour = parseInt(rawHour, 10);
        const minute = parseInt(rawMin, 10);

        if (amPm === "PM" && hour < 12) hour += 12;
        if (amPm === "AM" && hour === 12) hour = 0;

        // Use the poll’s actual start date from DB
        const entryStart = new Date(requestStartDate);
        entryStart.setHours(hour, minute, 0, 0);

        const entryEnd = new Date(entryStart.getTime() + 60 * 60 * 1000);

        // Check if it's in range
        if (entryStart < requestStartDate || entryEnd > requestEndDate) {
          console.warn(`Entry ${slotString} is out of range. Skipping.`);
          return null;
        }
        return { availabilityStart: entryStart, availabilityEnd: entryEnd };
      })
      .filter((entry) => entry !== null);

    try {
      await upsertAvailabilityEntries(token, availabilityRequestId, entries);
      alert("Your availability has been saved to the backend!");

      // Re-fetch the updated availability requests
      const requests = await fetchCurrentAvailabilityRequestsByGroupId(token, numericGroupId);
      const updatedRequest = requests.find((r) => r.id === availabilityRequestId);
      if (updatedRequest) {
        setBackendAvailabilityEntries(updatedRequest.availabilityEntries);
      }
    } catch (err) {
      console.error("Failed to upsert availability:", err);
      alert("Failed to save availability.");
    }
  };

  // Additional mouse-up listener
  useEffect(() => {
    document.addEventListener("mouseup", handleMouseUp);
    return () => {
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  async function handleCreateSessionFromBestTime(time: string) {
    if (!token || !requestStartDate || isNaN(numericGroupId)) {
      console.error("Missing token or group ID or requestStartDate");
      return;
    }

    // Parse the "time" string (e.g. "02:00 PM") into a JS Date
    const [timePart, amPm] = time.split(" ");
    const [rawHour, rawMin] = timePart.split(":");
    let hour = parseInt(rawHour, 10);
    const minute = parseInt(rawMin, 10);

    if (amPm === "PM" && hour < 12) hour += 12;
    if (amPm === "AM" && hour === 12) hour = 0;

    // Build the start time from requestStartDate
    const sessionStart = new Date(requestStartDate);
    sessionStart.setHours(hour, minute, 0, 0);

    // For demonstration, let's do a 60-minute duration
    const durationMinutes = 60;

    try {
      // Make a default title or do something like "Study Session"
      const newSession = await createStudySession(token, numericGroupId, {
        title: `Study Session @ ${time}`,
        startTime: sessionStart,
        durationMinutes,
      });
      alert(`Session created! ID=${newSession.id}`);

      // Optionally navigate to schedule page or re-fetch sessions
      // navigate(`/study-groups/${groupId}/schedule`);
    } catch (err) {
      console.error("Failed to create session:", err);
      alert("Failed to create session.");
    }
  }

  return (
    <div style={{ padding: "24px 64px 24px 64px", maxWidth: "1200px", margin: "0 auto 0 100px", position: "relative" }}>
      <h1 style={{ fontSize: "24px", fontWeight: "bold", marginBottom: "16px", textAlign: "center" }}>
        Availability Selection
      </h1>

      {/* Current user indicator */}
      <div style={{
        marginBottom: "24px",
        backgroundColor: "#f5f5f5",
        padding: "12px",
        borderRadius: "4px",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between"
      }}>
        <div>
          <span style={{ fontWeight: "bold" }}>Logged in as:</span> {currentUserId}
        </div>
        <div style={{ color: "#666", fontSize: "14px" }}>
          You can only modify your own availability
        </div>
      </div>

      {/* Instructions */}
      <div style={{
        marginBottom: "16px",
        padding: "16px",
        backgroundColor: "#f0f8ff",
        border: "1px solid #cce5ff",
        borderRadius: "4px"
      }}>
        <p style={{ fontWeight: "bold" }}>How to use:</p>
        <p>
          Click and drag across the grid to mark when you are available. Click again to unmark.
          Your selections are automatically saved locally. You can sync them to the backend by
          clicking &quot;Save Availability to Server.&quot;
        </p>
      </div>

      {/* Main content with grid and side stats */}
      <div style={{ display: "flex", gap: "24px" }}>
        {/* Time grid – Left side */}
        <div
          ref={gridRef}
          style={{
            flex: "1 1 auto",
            border: "1px solid #ddd",
            borderRadius: "4px",
            overflow: "auto",
            boxShadow: "0 2px 4px rgba(0,0,0,0.1)"
          }}
        >
          {timeSlots.map((time) => (
            <div
              key={time}
              style={{
                display: "flex",
                borderBottom: time === timeSlots[timeSlots.length - 1] ? "none" : "1px solid #ddd",
              }}
            >
              {/* Time label */}
              <div style={{
                padding: "12px",
                backgroundColor: "#f5f5f5",
                borderRight: "1px solid #ddd",
                fontWeight: 500,
                width: "100px",
                textAlign: "right"
              }}>
                {time}
              </div>
              {/* Selection cell */}
              <div
                style={{
                  height: "40px",
                  flexGrow: 1,
                  cursor: "pointer",
                  backgroundColor: availability[currentUserId]?.has(time) ? "#3f51b5" : "#e0e0e0",
                  transition: "background-color 0.2s"
                }}
                onMouseDown={() => handleMouseDown(time)}
                onMouseEnter={() => handleMouseEnter(time)}
                onTouchStart={() => handleMouseDown(time)}
                onTouchMove={(e) => {
                  if (isSelecting && gridRef.current) {
                    const touch = e.touches[0];
                    const elements = document.elementsFromPoint(touch.clientX, touch.clientY);
                    const timeElement = elements.find(el => el.getAttribute("data-time"));
                    const timeValue = timeElement?.getAttribute("data-time");
                    if (timeValue) handleMouseEnter(timeValue);
                  }
                }}
                data-time={time}
              />
            </div>
          ))}
        </div>

        {/* Group availability stats – Right side */}
        <div style={{
          width: "300px",
          flexShrink: 0,
          display: "flex",
          flexDirection: "column",
          gap: "16px"
        }}>

          {/* "Group Availability" bar chart */}
          <div style={{
            border: "1px solid #ddd",
            borderRadius: "4px",
            overflow: "hidden",
            boxShadow: "0 2px 4px rgba(0,0,0,0.1)"
          }}>
            <div style={{
              padding: "12px",
              backgroundColor: "#f5f5f5",
              borderBottom: "1px solid #ddd",
              fontWeight: "bold"
            }}>
              Group Availability
            </div>
            <div style={{ padding: "12px" }}>
              {timeSlots.map((time) => {
                const availableUserIds = getBackendAvailableUsers(time);
                const memberCount = groupMembers.length;
                const percentage = memberCount > 0
                  ? Math.round((availableUserIds.length / memberCount) * 100)
                  : 0;

                return (
                  <div key={time} style={{ marginBottom: "8px", display: "flex", alignItems: "center" }}>
                    <div style={{ width: "100px", fontSize: "14px" }}>{time}</div>
                    <div style={{
                      flex: "1",
                      backgroundColor: "#eee",
                      height: "16px",
                      borderRadius: "8px",
                      overflow: "hidden",
                      marginRight: "8px"
                    }}>
                      <div style={{
                        width: `${percentage}%`,
                        backgroundColor: percentage > 0 ? "#3f51b5" : "#eee",
                        height: "100%"
                      }} />
                    </div>
                    <div style={{ width: "40px", fontSize: "14px", textAlign: "right" }}>
                      {availableUserIds.length}/{memberCount}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Best times */}
          <div style={{
            border: "1px solid #ddd",
            borderRadius: "4px",
            overflow: "hidden",
            boxShadow: "0 2px 4px rgba(0,0,0,0.1)"
          }}>
            <div style={{
              padding: "12px",
              backgroundColor: "#f5f5f5",
              borderBottom: "1px solid #ddd",
              fontWeight: "bold"
            }}>
              Best Times
            </div>
            <div style={{ padding: "12px" }}>
              {getBestTimeSlots().length > 0 ? (
                // Only map over the top 3 slots
                getBestTimeSlots().slice(0, 3).map(({ time, count }, index) => {
                  // Percentage of total members
                  const percentage = groupMembers.length > 0
                    ? Math.round((count / groupMembers.length) * 100)
                    : 0;

                  return (
                    <div
                      key={time}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "space-between",
                        marginBottom: index < 2 ? "8px" : 0
                      }}
                    >
                      <span style={{ fontSize: "12px" }}>
                        <strong>{time}</strong> — {percentage}% available
                      </span>
                      <Button
                        variant="outlined"
                        color="primary"
                        sx={{
                          fontSize: "12px",
                          padding: "4px 8px",
                          textTransform: "none",
                          whiteSpace: "nowrap",
                        }}
                        onClick={() => void handleCreateSessionFromBestTime(time)}
                      >
                        Create Session
                      </Button>
                    </div>
                  );
                })
              ) : (
                <p style={{ fontStyle: "italic", color: "#666" }}>
                  No availability selected yet
                </p>
              )}
            </div>
          </div>


          {/* Others' Availability */}
          <div style={{
            border: "1px solid #ddd",
            borderRadius: "4px",
            overflow: "hidden",
            boxShadow: "0 2px 4px rgba(0,0,0,0.1)"
          }}>
            <div style={{
              padding: "12px",
              backgroundColor: "#f5f5f5",
              borderBottom: "1px solid #ddd",
              fontWeight: "bold"
            }}>
              Others&apos; Availability
            </div>
            <div style={{ padding: "12px" }}>
              {groupMembers
                .filter((member) => member.id !== currentUserId)
                .map((member) => {
                  // Count how many time slots the backend says they are available for:
                  let count = 0;
                  for (const slot of timeSlots) {
                    const users = getBackendAvailableUsers(slot);
                    if (users.includes(member.id)) {
                      count++;
                    }
                  }

                  const totalSlots = timeSlots.length;
                  const percentage = totalSlots > 0
                    ? Math.round((count / totalSlots) * 100)
                    : 0;

                  return (
                    <div key={member.id} style={{ marginBottom: "8px" }}>
                      <div style={{ display: "flex", justifyContent: "space-between" }}>
                        <strong>{member.name}</strong>
                        <span>{count}/{totalSlots} slots</span>
                      </div>
                      <div
                        style={{
                          backgroundColor: "#eee",
                          height: "8px",
                          borderRadius: "4px",
                          overflow: "hidden",
                          marginTop: "4px"
                        }}
                      >
                        <div
                          style={{
                            width: `${percentage}%`,
                            backgroundColor: percentage > 0 ? "#3f51b5" : "#eee",
                            height: "100%",
                          }}
                        />
                      </div>
                    </div>
                  );
                })}
            </div>
          </div>
        </div>
      </div>

      {/* Action buttons */}
      <div style={{ display: "flex", justifyContent: "center", gap: "16px", marginTop: "32px" }}>
        <Button variant="contained" color="primary" onClick={() => setOpen(true)}>
          View Summary
        </Button>
        <Button variant="contained" color="success" onClick={() => void saveAvailabilityToBackend()}>
          Save Availability to Server
        </Button>
        <Button variant="outlined" color="secondary" onClick={() => void navigate(`/study-groups/${groupId}/schedule`)}>
          Back to Schedule
        </Button>
      </div>

      {/* Summary Dialog */}
      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Availability Summary</DialogTitle>
        <DialogContent>
          <div>
            {/* By Person */}
            <h3 style={{ fontWeight: "bold", marginBottom: "8px" }}>By Person</h3>
            {groupMembers.map((member) => {
              const userTimes = Array.from(availability[member.id] || []);
              return (
                <div
                  key={member.id}
                  style={{
                    borderBottom: "1px solid #eee",
                    paddingBottom: "8px",
                    marginBottom: "8px",
                    backgroundColor: member.id === currentUserId ? "#f0f8ff" : "transparent"
                  }}
                >
                  <strong>
                    {member.name}
                    {member.id === currentUserId ? " (You)" : ""}
                  </strong>{" "}
                  {userTimes.length > 0 ? (
                    <span>{userTimes.join(", ")}</span>
                  ) : (
                    <span style={{ fontStyle: "italic", color: "#666" }}>
                      No availability selected
                    </span>
                  )}
                </div>
              );
            })}

            {/* By Time Slot */}
            <h3 style={{ fontWeight: "bold", marginTop: "16px", marginBottom: "8px" }}>By Time Slot</h3>
            {timeSlots.map((time) => {
              const slotUsers = getAvailableUsers(time);
              const displayNames = slotUsers.map((m) => m.name);

              return (
                <div
                  key={time}
                  style={{
                    borderBottom: "1px solid #eee",
                    paddingBottom: "8px",
                    marginBottom: "8px",
                  }}
                >
                  <strong>{time}:</strong>{" "}
                  {slotUsers.length > 0 ? (
                    <span>{displayNames.join(", ")}</span>
                  ) : (
                    <span style={{ fontStyle: "italic", color: "#666" }}>
                      No one available
                    </span>
                  )}
                </div>
              );
            })}

            {/* Best Options */}
            <h3 style={{ fontWeight: "bold", marginTop: "16px", marginBottom: "8px" }}>
              Best Options
            </h3>
            {getBestTimeSlots().length > 0 ? (
              getBestTimeSlots().map(({ time, count }) => {
                const slotUsers = getAvailableUsers(time);
                const displayNames = slotUsers.map((m) => m.name);
                const percentage = groupMembers.length > 0
                  ? Math.round((count / groupMembers.length) * 100)
                  : 0;

                return (
                  <div
                    key={time}
                    style={{
                      backgroundColor: "#f0f8ff",
                      padding: "8px",
                      marginBottom: "8px",
                      borderRadius: "4px",
                      border: "1px solid #cce5ff"
                    }}
                  >
                    <div style={{ display: "flex", justifyContent: "space-between" }}>
                      <strong>{time}</strong>
                      <span>
                        {count}/{groupMembers.length} people available ({percentage}%)
                      </span>
                    </div>
                    <div>{displayNames.join(", ")}</div>
                  </div>
                );
              })
            ) : (
              <div style={{ fontStyle: "italic", color: "#666" }}>
                No common availability found yet
              </div>
            )}
          </div>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
