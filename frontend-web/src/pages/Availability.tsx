import { useState, useEffect, useRef } from "react";
import {
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import {
  fetchCurrentAvailabilityRequestsByGroupId,
  createAvailabilityRequest,
  upsertAvailabilityEntries
} from "../api/studySessions";
import { fetchStudyGroupById } from "../api/studyGroups";
import { useAuth } from "../auth/useAuth";

export default function AvailabilitySelection() {
  const navigate = useNavigate();
  const { token, user } = useAuth();
  const { groupId } = useParams();
  const numericGroupId = groupId ? Number(groupId) : NaN;

  const [groupMembers, setGroupMembers] = useState<{ id: string; name: string }[]>([]);
  const memberNames = groupMembers.map(member => member.name);

  const currentUserName = user?.displayName ?? user?.email ?? "Unknown";

  const [availabilityRequestId, setAvailabilityRequestId] = useState<number | null>(null);
  useEffect(() => {
    console.log("availabilityRequestId:", availabilityRequestId);
  }, [availabilityRequestId]);

  const [timeSlots] = useState<string[]>([
    "09:00 AM", "10:00 AM", "11:00 AM", "12:00 PM",
    "01:00 PM", "02:00 PM", "03:00 PM", "04:00 PM", "05:00 PM"
  ]);
  // Local availability is stored as a record mapping user names to sets of time strings.
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

  useEffect(() => {
    if (!token || isNaN(numericGroupId)) return;
    fetchStudyGroupById(token, numericGroupId)
      .then((group) => {
        setGroupMembers(group.members);
      })
      .catch((err) => console.error("Error fetching group details", err));
  }, [token, numericGroupId]);

  useEffect(() => {
    if (!token || isNaN(numericGroupId)) {
      console.error("Invalid group ID or missing token");
      return;
    }
    fetchCurrentAvailabilityRequestsByGroupId(token, numericGroupId)
      .then((requests) => {
        if (requests.length > 0) {
          setAvailabilityRequestId(requests[0].id);
        } else {
          const now = new Date();
          const oneWeekLater = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);
          const newRequest = {
            title: "Availability Poll",
            availabilityPeriodStart: now,
            availabilityPeriodEnd: oneWeekLater,
          };
          return createAvailabilityRequest(token, numericGroupId, newRequest)
            .then(() => fetchCurrentAvailabilityRequestsByGroupId(token, numericGroupId));
        }
      })
      .then((updatedRequests) => {
        if (updatedRequests && updatedRequests.length > 0) {
          setAvailabilityRequestId(updatedRequests[0].id);
        }
      })
      .catch((err) => {
        console.error("Error fetching or creating availability request:", err);
      });
  }, [token, numericGroupId]);

  // Persist local availability changes to localStorage.
  useEffect(() => {
    const availForStorage: Record<string, string[]> = {};
    for (const key in availability) {
      availForStorage[key] = Array.from(availability[key]);
    }
    localStorage.setItem("userAvailability", JSON.stringify(availForStorage));
  }, [availability]);

  // Global mouseup handler to end drag selection.
  useEffect(() => {
    const handleMouseUp = () => setIsSelecting(false);
    document.addEventListener("mouseup", handleMouseUp);
    return () => document.removeEventListener("mouseup", handleMouseUp);
  }, []);

  // Toggle a single time slot for the current user.
  const toggleAvailability = (time: string, mode?: "add" | "remove") => {
    const action = mode ?? selectionMode;
    setAvailability(prev => {
      const newAvail = { ...prev };
      const currentAvail = new Set(newAvail[currentUserName]);
      if (action === "add") {
        currentAvail.add(time);
      } else {
        currentAvail.delete(time);
      }
      newAvail[currentUserName] = currentAvail;
      return newAvail;
    });
  };

  const handleMouseDown = (time: string) => {
    const isAvail = availability[currentUserName]?.has(time);
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

  // Returns an array of member names who have selected the given time.
  const getAvailableUsers = (time: string) => {
    return memberNames.filter(name => availability[name]?.has(time));
  };

  // Get the best time slots (most availability)
  const getBestTimeSlots = () => {
    const timesWithCounts = timeSlots.map(time => ({
      time,
      count: getAvailableUsers(time).length
    }));

    return timesWithCounts
      .sort((a, b) => b.count - a.count)
      .filter(item => item.count > 0)
      .slice(0, 3);
  };

  const saveAvailabilityToBackend = async () => {
    if (!token || !availabilityRequestId) {
      console.error("No token or availabilityRequestId available");
      return;
    }
    // Get the current user's selected time strings.
    const selectedTimes = Array.from(availability[currentUserName] || []);
    // Convert each selected time into a date range (assuming 1-hour blocks).
    const entries = selectedTimes.map((slotString) => {
      const [timePart, amPm] = slotString.split(" ");
      const [rawHour, rawMin] = timePart.split(":");
      let hour = parseInt(rawHour, 10);
      const minute = parseInt(rawMin, 10);
      if (amPm === "PM" && hour < 12) hour += 12;
      if (amPm === "AM" && hour === 12) hour = 0;
      const startDate = new Date();
      startDate.setHours(hour, minute, 0, 0);
      const endDate = new Date(startDate.getTime() + 60 * 60 * 1000);
      return { availabilityStart: startDate, availabilityEnd: endDate };
    });
    try {
      await upsertAvailabilityEntries(token, availabilityRequestId, entries);
      alert("Your availability has been saved to the backend!");
    } catch (err) {
      console.error("Failed to upsert availability:", err);
      alert("Failed to save availability.");
    }
  };
  useEffect(() => {
    // Add global mouse up handler
    document.addEventListener("mouseup", handleMouseUp);
    return () => {
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

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
          <span style={{ fontWeight: "bold" }}>Logged in as:</span> {currentUserName}
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
                  backgroundColor: availability[currentUserName]?.has(time) ? "#3f51b5" : "#e0e0e0",
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
                const availableNames = getAvailableUsers(time);
                const percentage = Math.round((availableNames.length / memberNames.length) * 100);
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
                      {availableNames.length}/{memberNames.length}
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
                getBestTimeSlots().map(({ time, count }, index) => (
                  <div key={time} style={{ marginBottom: index < 2 ? "8px" : 0 }}>
                    <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "4px" }}>
                      <strong>{time}</strong>
                      <span>{count}/{memberNames.length} people available ({Math.round((count / memberNames.length) * 100)}%)</span>
                    </div>
                    <div style={{ fontSize: "14px" }}>
                      {getAvailableUsers(time).join(", ")}
                    </div>
                  </div>
                ))
              ) : (
                <p style={{ fontStyle: "italic", color: "#666" }}>
                  No availability selected yet
                </p>
              )}
            </div>
          </div>

          {/* Others' availability */}
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
              {memberNames.filter(name => name !== currentUserName).map(name => {
                const userTimes = Array.from(availability[name] || []);
                return (
                  <div key={name} style={{ marginBottom: "8px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between" }}>
                      <strong>{name}</strong>
                      <span>{userTimes.length}/{timeSlots.length} slots</span>
                    </div>
                    <div style={{
                      backgroundColor: "#eee",
                      height: "8px",
                      borderRadius: "4px",
                      overflow: "hidden",
                      marginTop: "4px"
                    }}>
                      <div style={{
                        width: `${(userTimes.length / timeSlots.length) * 100}%`,
                        backgroundColor: userTimes.length > 0 ? "#3f51b5" : "#eee",
                        height: "100%"
                      }} />
                    </div>
                    {userTimes.length > 0 && (
                      <div style={{ fontSize: "12px", marginTop: "4px", color: "#666" }}>
                        {userTimes.slice(0, 3).join(", ")}
                        {userTimes.length > 3 ? ` and ${userTimes.length - 3} more...` : ""}
                      </div>
                    )}
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
            <h3 style={{ fontWeight: "bold", marginBottom: "8px" }}>By Person</h3>
            {memberNames.map((name) => (
              <div key={name} style={{
                borderBottom: "1px solid #eee",
                paddingBottom: "8px",
                marginBottom: "8px",
                backgroundColor: name === currentUserName ? "#f0f8ff" : "transparent"
              }}>
                <strong>{name}{name === currentUserName ? " (You)" : ""}:</strong>{" "}
                {Array.from(availability[name] || []).length > 0 ? (
                  <span>{Array.from(availability[name] || []).join(", ")}</span>
                ) : (
                  <span style={{ fontStyle: "italic", color: "#666" }}>
                    No availability selected
                  </span>
                )}
              </div>
            ))}

            <h3 style={{ fontWeight: "bold", marginTop: "16px", marginBottom: "8px" }}>By Time Slot</h3>
            {timeSlots.map((time) => {
              const availableNames = getAvailableUsers(time);
              return (
                <div key={time} style={{ borderBottom: "1px solid #eee", paddingBottom: "8px", marginBottom: "8px" }}>
                  <strong>{time}:</strong>{" "}
                  {availableNames.length > 0 ? (
                    <span>{availableNames.join(", ")}</span>
                  ) : (
                    <span style={{ fontStyle: "italic", color: "#666" }}>No one available</span>
                  )}
                </div>
              );
            })}

            <h3 style={{ fontWeight: "bold", marginTop: "16px", marginBottom: "8px" }}>Best Options</h3>
            {getBestTimeSlots().length > 0 ? (
              getBestTimeSlots().map(({ time, count }) => (
                <div key={time} style={{
                  backgroundColor: "#f0f8ff",
                  padding: "8px",
                  marginBottom: "8px",
                  borderRadius: "4px",
                  border: "1px solid #cce5ff"
                }}>
                  <div style={{ display: "flex", justifyContent: "space-between" }}>
                    <strong>{time}</strong>
                    <span>{count}/{memberNames.length} people available ({Math.round((count / memberNames.length) * 100)}%)</span>
                  </div>
                  <div>{getAvailableUsers(time).join(", ")}</div>
                </div>
              ))
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