import React, { createContext, useContext, useState, useEffect } from "react";

// Define context type
interface ModuleContextType {
  activeModules: string[];
  setActiveModules: (modules: string[]) => void;
}

// Create context with default values
const ModuleContext = createContext<ModuleContextType>({
  activeModules: [],
  setActiveModules: () => {},
});

// Custom hook to use the context
export const useModuleContext = () => useContext(ModuleContext);

export const ModuleProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [activeModules, setActiveModules] = useState<string[]>(() => {
    // Retrieve active modules from local storage on page reload
    const savedModules = localStorage.getItem("activeModules");
    return savedModules ? JSON.parse(savedModules) : [];
  });

  useEffect(() => {
    // Store active modules in local storage
    localStorage.setItem("activeModules", JSON.stringify(activeModules));
  }, [activeModules]);

  return (
    <ModuleContext.Provider value={{ activeModules, setActiveModules }}>
      {children}
    </ModuleContext.Provider>
  );
};