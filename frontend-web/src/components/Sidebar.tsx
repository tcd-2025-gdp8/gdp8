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
    ListItemIcon,
} from "@mui/material";
import { Logout } from "@mui/icons-material";
import profileLogo from "../assets/profileLogo.png";
import { useAuth } from "../auth/useAuth";

export default function Sidebar() {
    const { logout } = useAuth();
    const [logoutDialogOpen, setLogoutDialogOpen] = useState(false);

    const handleLogout = () => {
        logout();
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
                        justifyContent: "space-between",
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

                {/* Menu Items */}
                <List sx={{ flexGrow: 1 }}>
                    <ListItemButton component={Link} to="/study-groups">
                        <ListItemText primary="Study Groups" />
                    </ListItemButton>
                    <ListItemButton component={Link} to="/module">
                        <ListItemText primary="Modules" />
                    </ListItemButton>
                </List>

                {/* Logout Button at the Bottom */}
                <ListItemButton
                    onClick={() => setLogoutDialogOpen(true)}
                    sx={{
                        position: "absolute",
                        bottom: 0,
                        width: "100%",
                        backgroundColor: "#e57373",
                        color: "white",
                        "&:hover": { backgroundColor: "#d32f2f" },
                    }}
                >
                    <ListItemIcon>
                        <Logout sx={{ color: "white" }} />
                    </ListItemIcon>
                    <ListItemText primary="Log Out" />
                </ListItemButton>
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
