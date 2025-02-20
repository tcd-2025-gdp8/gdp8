// App.tsx
import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./firebase/useAuth";
import Login from "./loginPage/login";
import StudyGroupsPage from "./StudyGroupsPage/StudyGroupsPage";
import ModuleSettings from "./moduleSettings/moduleSettings";
import Chat from "./chat/Chat";
const App: React.FC = () => {
  const { user } = useAuth();

  return (
    <BrowserRouter>
      <Routes>
        {/* Always show login page on "/" and "/login" */}
        <Route path="/" element={<Login />} />
        <Route path="/login" element={<Login />} />
        {/* Only allow study groups if signed in */}
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
