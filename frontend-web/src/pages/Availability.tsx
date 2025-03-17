import { useState, useEffect, useRef } from "react";
import { Button, Dialog, DialogTitle, DialogContent, DialogActions } from "@mui/material";
import { useNavigate } from "react-router-dom";

// Example users - in a real app, this would come from authentication
const users = ["Alice", "Bob", "Charlie", "David"];
// Example current user - in a real app, this would come from auth context
const currentUser = "Alice"; 

export default function AvailabilitySelection() {
  const navigate = useNavigate();
  const [timeSlots, setTimeSlots] = useState<string[]>([
    "09:00 AM", "10:00 AM", "11:00 AM", "12:00 PM", 
    "01:00 PM", "02:00 PM", "03:00 PM", "04:00 PM", "05:00 PM"
  ]);
  const [availability, setAvailability] = useState<Record<string, Set<string>>>(
    () => {
      // Try to load from localStorage first
      // Define a type for your stored data structure
      type UserAvailability = Record<string, string[]>;
      const storedAvailability = localStorage.getItem("userAvailability");
      
      if (storedAvailability) {
        try {
          // Parse with type assertion
          const parsed = JSON.parse(storedAvailability) as UserAvailability;
          
          // Convert arrays back to Sets
          const result: Record<string, Set<string>> = {};
          
          for (const user in parsed) {
            // Make sure property actually belongs to the object
            if (Object.prototype.hasOwnProperty.call(parsed, user)) {
              result[user] = new Set(parsed[user]);
            }
          }
          
          return result;
        } catch (error) {
          console.error("Error parsing availability from localStorage:", error);
          return {}; // Return default value on error
        }
      }
      
      // Return default value if no stored data
      return {};  // Note the semicolon here was missing
    }
  );
  const [open, setOpen] = useState(false);
  const [isSelecting, setIsSelecting] = useState(false);
  const [selectionMode, setSelectionMode] = useState<"add" | "remove">("add");
  const gridRef = useRef<HTMLDivElement>(null);

  // Load time slots from localStorage
  useEffect(() => {
    const storedTimeSlots = localStorage.getItem("selectedTimeSlots");
    if (storedTimeSlots) {
      try {
        const parsedSlots: unknown = JSON.parse(storedTimeSlots);
        if (Array.isArray(parsedSlots) && parsedSlots.every((slot) => typeof slot === "string")) {
          setTimeSlots(parsedSlots);
        }
      } catch (error) {
        console.error("Error parsing time slots from localStorage:", error);
      }
    }
  }, []);

  // Save availability to localStorage whenever it changes
  useEffect(() => {
    // Convert Sets to arrays for JSON serialization
    const availabilityForStorage: Record<string, string[]> = {};
    for (const user in availability) {
      availabilityForStorage[user] = Array.from(availability[user]);
    }
    localStorage.setItem("userAvailability", JSON.stringify(availabilityForStorage));
  }, [availability]);

  // Handle toggle for single time slot
  const toggleAvailability = (time: string, mode?: "add" | "remove") => {
    const actionMode = mode ?? selectionMode;
    
    setAvailability((prev) => {
      const newAvailability = { ...prev };
      const userAvailability = new Set(newAvailability[currentUser]);
      
      if (actionMode === "add") {
        userAvailability.add(time);
      } else {
        userAvailability.delete(time);
      }
      
      newAvailability[currentUser] = userAvailability;
      return newAvailability;
    });
  };

  // Handle mouse events for drag selection
  const handleMouseDown = (time: string) => {
    const isAvailable = availability[currentUser].has(time);
    setSelectionMode(isAvailable ? "remove" : "add");
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

  // Calculate who's available at each time slot
  const getAvailableUsers = (time: string) => {
    return users.filter(user => availability[user].has(time));
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

  useEffect(() => {
    // Add global mouse up handler
    document.addEventListener("mouseup", handleMouseUp);
    return () => {
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  return (
    <div style={{ padding: "24px 64px 24px 64px", maxWidth: "1200px", margin: "0 auto 0 100px", position: "relative"}}>
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
          <span style={{ fontWeight: "bold" }}>Logged in as:</span> {currentUser}
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
        <p>Click and drag across the grid to mark when you are available. Click again to unmark. Your selections are automatically saved. You can view other users availability but cannot modify it.</p>
      </div>
      
      {/* Main content with grid and side stats */}
      <div style={{ display: "flex", gap: "24px" }}>
        {/* Time grid - Left side */}
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
                  backgroundColor: availability[currentUser].has(time) ? "#3f51b5" : "#e0e0e0",
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

        {/* Group availability stats - Right side */}
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
                const availableUsers = getAvailableUsers(time);
                const percentage = Math.round((availableUsers.length / users.length) * 100);
                
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
                      {availableUsers.length}/{users.length}
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
                getBestTimeSlots().map(({time, count}, index) => (
                  <div key={time} style={{ marginBottom: index < 2 ? "8px" : 0 }}>
                    <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "4px" }}>
                      <strong>{time}</strong> 
                      <span>{count}/{users.length} people</span>
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

          {/* Other users' availability */}
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
              Others Availability
            </div>
            <div style={{ padding: "12px" }}>
              {users.filter(user => user !== currentUser).map(user => {
                const userTimes = Array.from(availability[user]);
                return (
                  <div key={user} style={{ marginBottom: "8px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between" }}>
                      <strong>{user}</strong>
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
        <Button 
          variant="contained" 
          color="primary" 
          onClick={() => setOpen(true)}
        >
          View Summary
        </Button>
        <Button 
          variant="outlined" 
          color="secondary" 
          onClick={() => void navigate("/study-groups/:groupId/schedule")}
        >
          Back to Schedule
        </Button>
      </div>
      
      {/* Summary Dialog */}
      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Availability Summary</DialogTitle>
        <DialogContent>
          <div>
            <h3 style={{ fontWeight: "bold", marginBottom: "8px" }}>By Person</h3>
            {users.map((user) => (
              <div key={user} style={{ 
                borderBottom: "1px solid #eee", 
                paddingBottom: "8px", 
                marginBottom: "8px",
                backgroundColor: user === currentUser ? "#f0f8ff" : "transparent"
              }}>
                <strong>{user}{user === currentUser ? " (You)" : ""}:</strong>{" "}
                {Array.from(availability[user]).length > 0 ? (
                  <span>{Array.from(availability[user]).join(", ")}</span>
                ) : (
                  <span style={{ fontStyle: "italic", color: "#666" }}>No availability selected</span>
                )}
              </div>
            ))}
            
            <h3 style={{ fontWeight: "bold", marginTop: "16px", marginBottom: "8px" }}>By Time Slot</h3>
            {timeSlots.map((time) => {
              const availableUsers = getAvailableUsers(time);
              return (
                <div key={time} style={{ borderBottom: "1px solid #eee", paddingBottom: "8px", marginBottom: "8px" }}>
                  <strong>{time}:</strong>{" "}
                  {availableUsers.length > 0 ? (
                    <span>{availableUsers.join(", ")}</span>
                  ) : (
                    <span style={{ fontStyle: "italic", color: "#666" }}>No one available</span>
                  )}
                </div>
              );
            })}
            
            <h3 style={{ fontWeight: "bold", marginTop: "16px", marginBottom: "8px" }}>Best Options</h3>
            {getBestTimeSlots().length > 0 ? (
              getBestTimeSlots().map(({time, count}) => (
                <div key={time} style={{ 
                  backgroundColor: "#f0f8ff", 
                  padding: "8px", 
                  marginBottom: "8px", 
                  borderRadius: "4px",
                  border: "1px solid #cce5ff"
                }}>
                  <div style={{ display: "flex", justifyContent: "space-between" }}>
                    <strong>{time}</strong>
                    <span>{count}/{users.length} people available ({Math.round((count/users.length)*100)}%)</span>
                  </div>
                  <div>
                    {getAvailableUsers(time).join(", ")}
                  </div>
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