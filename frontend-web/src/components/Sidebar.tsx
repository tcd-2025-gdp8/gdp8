import { useState } from "react";
import { useLocation, useNavigate, matchPath, Link } from "react-router-dom";
import {
    Box,
    List,
    ListItemButton,
    ListItemText,
    Drawer,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Paper,
    ListItemIcon,
} from "@mui/material";
import { Logout } from "@mui/icons-material";
import profileLogo from "../assets/profileLogo.png";
import { useAuth } from "../auth/useAuth";

export default function Sidebar() {
    const { logout } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();
    const [logoutDialogOpen, setLogoutDialogOpen] = useState(false);

    // Check if we're on a group details page or chat page
    const groupDetailsMatch = matchPath(
        { path: "/study-groups/:groupId/details" },
        location.pathname
    );
    const chatMatch = matchPath(
        { path: "/chat/:groupId" },
        location.pathname
    );
    const groupId = groupDetailsMatch?.params?.groupId ?? chatMatch?.params?.groupId;    // Use whichever match finds a groupId


    const handleLogout = () => {
        void logout();
        setLogoutDialogOpen(false);
    };

    return (
        <>
            <Drawer
                variant="permanent"
                sx={{
                    width: 240,
                    flexShrink: 0,
                    "& .MuiDrawer-paper": {
                        width: 240,
                        boxSizing: "border-box",
                        backgroundColor: "#f5f5f5",
                        position: "fixed",
                        zIndex: 1200,
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "flex-start",
                        height: "100vh",
                        paddingBottom: "1rem"
                    },
                }}
            >
                {/* Profile Section */}
                <Box sx={{ padding: "1rem", textAlign: "center" }}>
                    <img
                        src={profileLogo}
                        alt="Profile Logo"
                        style={{ width: "48px", height: "48px", marginRight: "1rem" }}
                    />
                </Box>

                {/* Permanent Menu Items */}
                <List sx={{ flexGrow: 0 }}>
                    <ListItemButton component={Link} to="/study-groups">
                        <ListItemText primary="Study Groups" />
                    </ListItemButton>
                    <ListItemButton component={Link} to="/module">
                        <ListItemText primary="Modules" />
                    </ListItemButton>
                </List>

                {/* White box below Modules, only when we have a groupId */}
                {groupId && (
                    <Paper
                        elevation={3}
                        sx={{
                            width: "100%",
                            // No top margin so it sits directly under Modules
                            mt: 0,
                        }}
                    >
                        <List
                            disablePadding
                            sx={{
                                // Remove extra padding so items are flush at the top
                                paddingTop: 0,
                                paddingBottom: 0,
                            }}
                        >
                            <ListItemButton onClick={() => { void navigate(`/chat/${groupId}`); }}>
                                <ListItemText primary="Chats" />
                            </ListItemButton>
                            <ListItemButton onClick={() => { void navigate(`/study-groups/${groupId}/files`); }}>
                                <ListItemText primary="Files" />
                            </ListItemButton>
                            <ListItemButton onClick={() => { void navigate(`/study-groups/${groupId}/schedule`); }}>
                                <ListItemText primary="Scheduling" />
                            </ListItemButton>
                            <ListItemButton onClick={() => { void navigate(`/study-groups/${groupId}/leadership-board`); }}>
                                <ListItemText primary="Leadership Board" />
                            </ListItemButton>
                            <ListItemButton onClick={() => { void navigate(`/study-groups/${groupId}/details`); }}>
                                <ListItemText primary="Group Details" />
                            </ListItemButton>
                        </List>
                    </Paper>
                )}
                <Box sx={{ flexGrow: 1 }} />

                {/* Logout Button at the Bottom */}
                <List>
                    <ListItemButton
                        onClick={() => void logout()}
                        sx={{
                            backgroundColor: "#d32f2f",
                            color: "white",
                            "&:hover": { backgroundColor: "#e57373" },
                        }}
                    >
                        <ListItemIcon>
                            <Logout sx={{ color: "white" }} />
                        </ListItemIcon>
                        <ListItemText primary="Log Out" />
                    </ListItemButton>
                </List>
            </Drawer>

            {/* Logout Confirmation Dialog */}
            <Dialog
                open={logoutDialogOpen}
                onClose={() => setLogoutDialogOpen(false)}
            >
                <DialogTitle>Confirm Logout</DialogTitle>
                <DialogContent>
                    Are you sure you want to log out?
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setLogoutDialogOpen(false)} color="primary" variant="contained">
                        Cancel
                    </Button>
                    <Button onClick={handleLogout} color="error" variant="contained">
                        Log Out
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
