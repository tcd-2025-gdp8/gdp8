// frontend-web/pages/IndividualStudyGroup.tsx

import { useState } from "react";
import { Box, Button, Menu, MenuItem } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";

import Sidebar from "../components/Sidebar";
import CustomAppBar from "../components/CustomAppBar";

export default function IndividualStudyGroup() {
    const navigate = useNavigate();
    const { groupId } = useParams();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    // const [, setShowSchedule] = useState(false);

    const handleMenuOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
        setAnchorEl(null);
    };

    return (
        <Box sx={{ display: "flex" }}>
            <Sidebar />
            <Box sx={{ flexGrow: 1, ml: "300px" }}>
                <CustomAppBar />
                <Box sx={{ p: 2, mt: 10, position: "relative", paddingLeft: "80px", marginTop: "-30px" }}>
                    {/* "Sections" dropdown button */}
                    <Button
                        variant="contained"
                        onClick={handleMenuOpen}
                        sx={{
                            position: "absolute",
                            left: "-550px", // do not change
                            top: "24px",
                            zIndex: 1,
                            backgroundColor: "#3b5998",
                            "&:hover": {
                                backgroundColor: "#2d4373",
                            },
                        }}
                    >
                        Sections
                    </Button>

                    <Menu
                        anchorEl={anchorEl}
                        open={Boolean(anchorEl)}
                        onClose={handleMenuClose}
                        anchorOrigin={{
                            vertical: "bottom",
                            horizontal: "left",
                        }}
                        transformOrigin={{
                            vertical: "top",
                            horizontal: "left",
                        }}
                    >
                        <MenuItem
                            onClick={() => {
                                handleMenuClose();
                                if (groupId) {
                                    void navigate(`/chat/${groupId}`);
                                }
                            }}
                        >
                            Chats
                        </MenuItem>
                        <MenuItem
                            onClick={() => {
                                handleMenuClose();
                                if (groupId) {
                                    // Now navigate to the real schedule route:
                                    navigate(`/study-groups/${groupId}/schedule`);
                                }
                            }}
                        >
                            Schedule
                        </MenuItem>
                        <MenuItem onClick={handleMenuClose}>Leadership Board</MenuItem>

                        {/* New: Navigate to our Group Details page */}
                        <MenuItem
                            onClick={() => {
                                handleMenuClose();
                                if (groupId) {
                                    void navigate(`/study-groups/${groupId}/details`);
                                }
                            }}
                        >
                            Group Details
                        </MenuItem>
                    </Menu>

                    {/* 
            You can leave this blank or display something else.
            The GroupDetailsPage is now a separate route.
          */}
                </Box>
            </Box>
        </Box>
    );
}
