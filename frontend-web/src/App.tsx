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
        </Route>
        <Route
            path="/files"
            element={<FilesPage />}
          />

      </Routes>
    </BrowserRouter>
  );
};

export default App;
