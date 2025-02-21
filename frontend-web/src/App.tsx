// App.tsx
import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./auth/useAuth";
import Login from "./loginPage/login";
import StudyGroupsPage from "./StudyGroupsPage/StudyGroupsPage";
import ModuleSettings from "./moduleSettings/moduleSettings";
import LandingPage from "./landingPage/landingPage";
import Chat from "./chat/Chat";

const App: React.FC = () => {
  const { user } = useAuth();

  return (
    <BrowserRouter>
      <Routes>
        {/* Always show login page on "/" and "/login" */}
        <Route path="/" element={<Login />} />
        <Route path="/login" element={<Login />} />
        {/* Example: Only allow landing page if signed in */}
        <Route
          path="/landing"
          element={user ? <LandingPage /> : <Navigate to="/login" />}
        />
        <Route
          path="/study-groups"
          element={user ? <StudyGroupsPage /> : <Navigate to="/login" />}
        />
        <Route
          path="/module"
          element={<ModuleSettings />}
        />
        <Route path="/chat/:groupId" element={<Chat />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
