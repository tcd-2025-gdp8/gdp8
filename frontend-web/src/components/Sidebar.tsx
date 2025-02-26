import { Link } from "react-router-dom";
import {
    Box,
    List,
    ListItemButton,
    ListItemText,
    Drawer,
} from "@mui/material";


export default function Sidebar() {
    return (
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
                    src="/src/assets/profileLogo.png"
                    alt="Profile Logo"
                    style={{ width: '48px', height: '48px', marginRight: '1rem' }}
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
    );
}
