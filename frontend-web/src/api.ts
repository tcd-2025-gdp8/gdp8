export const fetchStudyGroups = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/study-groups");
      if (!response.ok) throw new Error("Failed to fetch study groups");
      return await response.json();
    } catch (error) {
      console.error("Error fetching study groups:", error);
      return [];
    }
  };
  