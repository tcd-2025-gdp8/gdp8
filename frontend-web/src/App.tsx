import { useState, useEffect } from "react";
import "./App.css";

interface StudyGroup {
  id: string;
  name: string;
  module: string;
}

function App() {
  const [studyGroups, setStudyGroups] = useState<StudyGroup[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStudyGroups = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/study-groups");
        if (!response.ok) throw new Error("Failed to fetch study groups");
        const data = await response.json();
        setStudyGroups(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchStudyGroups();
  }, []);

  return (
    <div className="container">
      <h1>Study Groups</h1>
      {loading && <p>Loading...</p>}
      {error && <p className="error">{error}</p>}
      <ul>
        {studyGroups.map((group) => (
          <li key={group.id}>
            <strong>{group.name}</strong> - {group.module}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
