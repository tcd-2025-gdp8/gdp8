// App.tsx
import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./firebase/authContext";
import Login from "./loginPage/login";
import StudyGroupsPage from "./StudyGroupsPage/StudyGroupsPage";

const App: React.FC = () => {
  const { user } = useAuth();

  return (
    <BrowserRouter>
      <Routes>
        {/* If not logged in, redirect to /login, else show StudyGroupsPage */}
        <Route path="/" element={user ? <StudyGroupsPage /> : <Navigate to="/login" />} />
        {/* If already authenticated, redirect /login to /study-groups */}
        <Route path="/login" element={user ? <Navigate to="/study-groups" /> : <Login />} />
        <Route path="/study-groups" element={<StudyGroupsPage />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
