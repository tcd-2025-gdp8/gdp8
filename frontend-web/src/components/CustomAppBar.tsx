import { useEffect, useState } from "react";
import {
    Box,
    Typography,
    Drawer,
    IconButton,
    AppBar,
    Toolbar,
    Badge,
} from "@mui/material";
import NotificationsIcon from "@mui/icons-material/Notifications";
import CloseIcon from "@mui/icons-material/Close";

import { Notification, getNotifications, markNotificationAsRead } from "../api/notifications";
import { useAuth } from "../auth/useAuth";


export default function CustomAppBar() {
    const { token, user } = useAuth();

    const [openNotifications, setOpenNotifications] = useState(false);
    const [notifications, setNotifications] = useState<Notification[]>([]);


    useEffect(() => {
        const fetchNotifications = async () => {
            if(!user) return;
            try {
                const fetchedNotifications = await getNotifications(token, user.uid);
                setNotifications(fetchedNotifications);
            } catch (error) {
                console.error("Error fetching notifications:", error);
                alert("Error fetching notifications");
            }
        };
        void fetchNotifications();
    }, [token, user]);

    const handleDeleteNotification = (notificationId: number) => {
        void (async () => {
            try {
                await markNotificationAsRead(token, notificationId);
                setNotifications((prev) =>
                    prev.filter((notification) => notification.id !== notificationId)
                );
            } catch (error) {
                console.error("Error deleting notification:", error);
                alert("Error deleting notification");
            }
        })();
    };


    return (
        <>
            <AppBar
                position="fixed"
                sx={{
                    backgroundColor: "#ffffff",
                    color: "#000000",
                    zIndex: 1201,
                    width: "calc(100% - 240px)",
                    marginLeft: "240px",
                }}
            >
                <Toolbar sx={{ flexDirection: "row", alignItems: "center" }}>
                    <Box sx={{ flexGrow: 1, display: "flex", alignItems: "center" }}>
                        <Typography variant="h6">
                            Blackboard + StudyWise
                        </Typography>
                    </Box>
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
                        zIndex: 1300,
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
                                <Typography>{notification.content}</Typography>
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
        </>
    )
}
