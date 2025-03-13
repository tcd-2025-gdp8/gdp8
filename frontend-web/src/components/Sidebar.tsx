import { useState } from "react";
import { Link } from "react-router-dom";
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
} from "@mui/material";

import profileLogo from "../assets/profileLogo.png";
import { useAuth } from "../auth/useAuth"; // Assuming you have an auth hook

export default function Sidebar() {
    const { logout } = useAuth(); // Assuming your auth hook provides a logout function
    const [logoutDialogOpen, setLogoutDialogOpen] = useState(false);

    const handleLogout = () => {
        logout(); // Perform logout
        setLogoutDialogOpen(false); // Close the dialog
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
                    },
                }}
            >
                <Box sx={{ padding: "1rem", textAlign: "center" }}>
                    <img
                        src={profileLogo}
                        alt="Profile Logo"
                        style={{ width: '48px', height: '48px', marginRight: '1rem' }}
                    />
                </Box>
                <List>
                    <ListItemButton component={Link} to="/study-groups">
                        <ListItemText primary="Study Groups" />
                    </ListItemButton>
                    <ListItemButton component={Link} to="/module">
                        <ListItemText primary="Modules" />
                    </ListItemButton>
                    <ListItemButton onClick={() => setLogoutDialogOpen(true)}>
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
                    <Button onClick={() => setLogoutDialogOpen(false)} color="primary">
                        Cancel
                    </Button>
                    <Button onClick={handleLogout} color="secondary" variant="contained">
                        Log Out
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
