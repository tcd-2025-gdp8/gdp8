import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import Layout from "./Layout";
import LoginPage from "./pages/login";
import LandingPage from "./pages/landing";
import StudyGroupsPage from "./pages/StudyGroupsPage";
import IndividualStudyGroup from "./pages/IndividualStudyGroup";
import ModuleSettingsPage from "./pages/moduleSettings";
import ChatPage from "./pages/ChatPage";
import FilesPage from "./pages/FilesPage";
import GroupDetailsPage from "./pages/groupDetailsPage";
import SchedulePage from "./pages/SchedulePage";
import AvailabilitySelection from "./pages/Availability";
import Leaderboard from "./pages/LeadershipPage";

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>

        <Route
          path="/login"
          element={<LoginPage />}
        />

        <Route path="/" element={<Layout />}>
          <Route
            index
            element={<LandingPage />}
          />
          <Route
            path="/landing"
            element={<LandingPage />}
          />
          <Route
            path="/study-groups"
            element={<StudyGroupsPage />}
          />
          <Route
            path="/study-groups/:groupId/details"
            element={<GroupDetailsPage />}
          />
          <Route path="/study-groups/:groupId"
            element={<IndividualStudyGroup />}
          />
          <Route
            path="/module"
            element={<ModuleSettingsPage />}
          />
          <Route
            path="/chat/:groupId"
            element={<ChatPage />}
          />
          <Route path="/study-groups/:groupId/files"
            element={<FilesPage />}
          />
          <Route path="/study-groups/:groupId/schedule"
            element={<SchedulePage />}
          />
        </Route>
        <Route path="/study-groups/:groupId/availability"
          element={<AvailabilitySelection />}
        />
        
        <Route path="/study-groups/:groupId/leadership-board"
            element={<Leaderboard />}
          />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
