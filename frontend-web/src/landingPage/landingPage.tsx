// src/LandingPage.tsx
import React, { useState } from "react";
import { Link } from "react-router-dom"; // Removed useNavigate since it's unused
import { useAuth } from "../firebase/useAuth";
import {
    Box,
    Typography,
    List,
    ListItemButton,
    ListItemText,
    Drawer,
    IconButton,
    AppBar,
    Toolbar,
    Badge,
} from "@mui/material";
import NotificationsIcon from "@mui/icons-material/Notifications";
import CloseIcon from "@mui/icons-material/Close";
import styles from "./LandingPage.module.css"; // Ensure TypeScript recognizes this correctly

interface Notification {
    id: number;
    message: string;
}

const LandingPage: React.FC = () => {
    const { user } = useAuth();

    // Notifications state
    const [notifications, setNotifications] = useState<Notification[]>([
        { id: 1, message: 'Your request to join "The Prefects" has been accepted.' },
        { id: 2, message: 'New study group "CS Wizards" has been created for CSU44051: Human Factors.' },
        { id: 3, message: 'New study group "The Elites" has been created for CSU44052: Computer Graphics.' },
    ]);
    const [openNotifications, setOpenNotifications] = useState(false);

    // Delete a notification
    const handleDeleteNotification = (notificationId: number) => {
        setNotifications((prev) =>
            prev.filter((notification) => notification.id !== notificationId)
        );
    };

    return (
        <Box sx={{ display: "flex" }}>
            {/* Sidebar */}
            <Drawer
                variant="permanent"
                sx={{
                    width: 240,
                    flexShrink: 0,
                    "& .MuiDrawer-paper": {
                        width: 240,
                        boxSizing: "border-box",
                        backgroundColor: "#f5f5f5",
                        position: "fixed", // Fix the sidebar position
                        zIndex: 1200,     // Ensure sidebar is below AppBar
                    },
                }}
            >
                <Box sx={{ padding: "1rem", textAlign: "center" }}>
                    <img
                        src="/src/assets/profileLogo.png"
                        alt="Profile Logo"
                        className={(styles as Record<string, string>).logo ?? ""}
                    />
                </Box>
                <List>
                    <ListItemButton component={Link} to="/module">
                        <ListItemText primary="Modules" />
                    </ListItemButton>
                    <ListItemButton component={Link} to="/study-groups">
                        <ListItemText primary="Study Groups" />
                    </ListItemButton>
                </List>
            </Drawer>

            {/* Main Content */}
            <Box sx={{ flexGrow: 1, marginLeft: "240px" }}>
                {/* AppBar at the top */}
                <AppBar
                    position="fixed"
                    sx={{
                        backgroundColor: "#ffffff",
                        color: "#000000",
                        zIndex: 1201, // Ensure AppBar is above sidebar
                        width: "calc(100% - 240px)", // Adjust width to account for sidebar
                        marginLeft: "240px", // Shift AppBar to the right
                    }}
                >
                    <Toolbar sx={{ flexDirection: "row", alignItems: "center" }}>
                        {/* Left section: Title */}
                        <Box sx={{ flexGrow: 1, display: "flex", alignItems: "center" }}>
                            <Typography variant="h6">
                                Blackboard + StudyWise
                            </Typography>
                        </Box>

                        {/* Right section: User Email & Notifications */}
                        {user && (
                            <Typography variant="subtitle1" sx={{ mr: 2 }}>
                                {user.email}
                            </Typography>
                        )}
                        <IconButton color="inherit" onClick={() => setOpenNotifications(true)}>
                            <Badge badgeContent={notifications.length} color="error">
                                <NotificationsIcon />
                            </Badge>
                        </IconButton>
                    </Toolbar>
                </AppBar>

                {/* Notifications Drawer */}
                <Drawer
                    anchor="right"
                    open={openNotifications}
                    onClose={() => setOpenNotifications(false)}
                    sx={{
                        "& .MuiDrawer-paper": {
                            zIndex: 1300, // Ensure notifications drawer is above AppBar
                        },
                    }}
                >
                    <Box sx={{ width: 300, padding: "1rem" }}>
                        <Box
                            sx={{
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                                marginBottom: "1rem",
                            }}
                        >
                            <Typography variant="h6">Notifications</Typography>
                            <IconButton onClick={() => setOpenNotifications(false)}>
                                <CloseIcon />
                            </IconButton>
                        </Box>
                        {notifications.length === 0 ? (
                            <Typography
                                variant="body2"
                                color="textSecondary"
                                sx={{ textAlign: "center" }}
                            >
                                No notifications
                            </Typography>
                        ) : (
                            notifications.map((notification) => (
                                <Box
                                    key={notification.id}
                                    sx={{
                                        padding: "0.5rem",
                                        borderBottom: "1px solid #ccc",
                                        display: "flex",
                                        justifyContent: "space-between",
                                        alignItems: "center",
                                    }}
                                >
                                    <Typography>{notification.message}</Typography>
                                    <IconButton
                                        size="small"
                                        onClick={() => handleDeleteNotification(notification.id)}
                                    >
                                        <CloseIcon fontSize="small" />
                                    </IconButton>
                                </Box>
                            ))
                        )}
                    </Box>
                </Drawer>

                {/* Main Content Area (Below the AppBar) */}
                <Box sx={{ padding: "2rem", marginTop: "64px" }}>
                    {/* Removed the instructions text */}
                </Box>
            </Box>
        </Box>
    );
};

export default LandingPage;
