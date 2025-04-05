import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Trophy, Clock, Users, Award } from 'lucide-react';
import Sidebar from './../components/Sidebar'; // Import the Sidebar component

import { 
  Card, 
  CardContent, 
  Typography, 
  Box, 
  Paper, 
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  SelectChangeEvent
} from '@mui/material';
import { fetchAllModules, fetchStudyGroupStats } from '../api/modules';
import { useAuth } from '../auth/useAuth';

interface Module {
  id: number;
  code: string;
  name: string;
}

interface StudyGroupStatsDTO {
  id: number;
  name: string;
  members: number;
  totalHours: number;
  weeklyHours: number[];
}

type StudyGroupStatsMap = Record<string, StudyGroupStatsDTO[]>;

interface WeekData {
  name: string;
  [groupName: string]: string | number;
}

const weekLabels: string[] = ['Week 1', 'Week 2', 'Week 3', 'Week 4'];

const LeadershipPage: React.FC = () => {
  const { token } = useAuth();
  const [selectedModule, setSelectedModule] = useState<string>('');
  const [modules, setModules] = useState<Module[]>([]);
  const [studyGroupStats, setStudyGroupStats] = useState<StudyGroupStatsMap>({});
  const [loading, setLoading] = useState<boolean>(true); // Added loading state

  useEffect(() => {
    const fetchData = async () => {
      const fetchedModules = await fetchAllModules(token);
      const fetchedStudyGroupStats = await fetchStudyGroupStats(token);

      setModules(fetchedModules);
      setStudyGroupStats(fetchedStudyGroupStats as StudyGroupStatsMap);

      if (fetchedModules.length > 0) {
        setSelectedModule(fetchedModules[0].code);
      }

      setLoading(false); // Set loading to false after data is fetched
    };

    fetchData();
  }, []);

  // Ensure currentGroups is populated before using it
  const currentGroups: StudyGroupStatsDTO[] = studyGroupStats[selectedModule] || [];

  // Format data for the chart
  const chartData: WeekData[] = weekLabels.map((week, index) => {
    const weekData: WeekData = { name: week };
    currentGroups.forEach(group => {
      weekData[group.name] = group.weeklyHours[index];
    });
    return weekData;
  });

  // Define medal colors
  const getMedalColor = (rank: number): string => {
    switch(rank) {
      case 0: return '#FFD700'; // gold
      case 1: return '#C0C0C0'; // silver
      case 2: return '#CD7F32'; // bronze
      default: return '#6B7280'; // gray
    }
  };

  // Get medal icon
  const getMedalIcon = (rank: number): React.ReactNode => {
    if (rank === 0) return <Trophy color={getMedalColor(rank)} size={24} />;
    if (rank === 1) return <Award color={getMedalColor(rank)} size={22} />;
    if (rank === 2) return <Award color={getMedalColor(rank)} size={20} />;
    return <Typography sx={{ color: '#6B7280', fontWeight: 'bold' }}>{rank + 1}</Typography>;
  };

  const handleModuleChange = (event: SelectChangeEvent) => {
    setSelectedModule(event.target.value);
  };

  if (loading) {
    return <div>Loading...</div>; // Loading message while fetching data
  }

  return (
    <Box sx={{ display: 'flex' }}>
      <Sidebar />
      
      <Box 
        sx={{ 
          bgcolor: '#F9FAFB', 
          minHeight: '100vh', 
          p: 3, 
          flexGrow: 1,
        }}
      >
        <Box sx={{ width: '100%', maxWidth: '100%', mx: 'auto', px: { xs: 2, md: 4 } }}>
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 3 }}>
            <Box>
              <Typography variant="h4" sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center' }}>
                <Trophy color="#FFD700" size={32} style={{ marginRight: '12px' }} />
                Study Group Leaderboard
              </Typography>
              <Typography variant="subtitle1" color="text.secondary">
                See which groups are putting in the most study hours
              </Typography>
            </Box>
          </Box>
          
          <Box sx={{ width: '100%' }}>
            {/* Module Selection Dropdown */}
            <Box sx={{ mb: 3, maxWidth: 400 }}>
              <FormControl fullWidth>
                <InputLabel id="module-select-label">Select Module</InputLabel>
                <Select
                  labelId="module-select-label"
                  id="module-select"
                  value={selectedModule}
                  label="Select Module"
                  onChange={handleModuleChange}
                >
                  {modules.map(module => (
                    <MenuItem key={module.id} value={module.code}>{module.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Box>
            
            <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', lg: '1fr 2fr' }, gap: 3 }}>
              {/* Podium - Top 3 */}
              <Card elevation={3}>
                <CardContent sx={{ pt: 3 }}>
                  <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, textAlign: 'center', color: '#4B5563' }}>
                    Top Performers
                  </Typography>
                  
                  <Box sx={{ position: 'relative', height: '256px', display: 'flex', alignItems: 'flex-end', justifyContent: 'center', pb: 2 }}>
                    {/* 2nd Place */}
                    {currentGroups.length > 1 && (
                      <Box sx={{ position: 'absolute', bottom: 0, left: 0, width: '112px', mx: 2 }}>
                        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                          <Box sx={{ bgcolor: '#E5E7EB', height: '128px', width: '100%', borderTopLeftRadius: '8px', borderTopRightRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <Award color="#9CA3AF" size={32} />
                          </Box>
                          <Box sx={{ bgcolor: '#4B5563', width: '100%', py: 1, textAlign: 'center', color: 'white', fontWeight: 600, borderBottomLeftRadius: '8px', borderBottomRightRadius: '8px' }}>
                            {currentGroups[1]?.totalHours ?? 0}h
                          </Box>
                          <Typography variant="body2" sx={{ mt: 1, fontWeight: 500, textAlign: 'center' }}>
                            {currentGroups[1]?.name ?? 'N/A'}
                          </Typography>
                        </Box>
                      </Box>
                    )}
                    
                    {/* 1st Place */}
                    <Box sx={{ position: 'absolute', bottom: 0, left: '50%', transform: 'translateX(-50%)', width: '128px', zIndex: 1 }}>
                      <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                        <Box sx={{ bgcolor: '#FEF3C7', height: '176px', width: '100%', borderTopLeftRadius: '8px', borderTopRightRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center', border: 2, borderColor: '#FCD34D' }}>
                          <Trophy color="#F59E0B" size={40} />
                        </Box>
                        <Box sx={{ bgcolor: '#F59E0B', width: '100%', py: 1.5, textAlign: 'center', color: 'white', fontWeight: 700, borderBottomLeftRadius: '8px', borderBottomRightRadius: '8px', fontSize: '1.125rem' }}>
                          {currentGroups[0]?.totalHours ?? 0}h
                        </Box>
                        <Typography variant="body2" sx={{ mt: 1, fontWeight: 500, textAlign: 'center' }}>
                          {currentGroups[0]?.name ?? 'N/A'}
                        </Typography>
                      </Box>
                    </Box>
                    
                    {/* 3rd Place */}
                    {currentGroups.length > 2 && (
                      <Box sx={{ position: 'absolute', bottom: 0, right: 0, width: '112px', mx: 2 }}>
                        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                          <Box sx={{ bgcolor: '#FEF3C7', height: '96px', width: '100%', borderTopLeftRadius: '8px', borderTopRightRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <Award color="#D97706" size={28} />
                          </Box>
                          <Box sx={{ bgcolor: '#D97706', width: '100%', py: 1, textAlign: 'center', color: 'white', fontWeight: 600, borderBottomLeftRadius: '8px', borderBottomRightRadius: '8px' }}>
                            {currentGroups[2]?.totalHours ?? 0}h
                          </Box>
                          <Typography variant="body2" sx={{ mt: 1, fontWeight: 500, textAlign: 'center' }}>
                            {currentGroups[2]?.name ?? 'N/A'}
                          </Typography>
                        </Box>
                      </Box>
                    )}
                  </Box>
                </CardContent>
              </Card>
              
              {/* Weekly Progress Chart */}
              <Card elevation={3}>
                <CardContent sx={{ pt: 3 }}>
                  <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, color: '#4B5563' }}>
                    Weekly Progress
                  </Typography>
                  <Box sx={{ height: 300, width: '100%' }}>
                    <ResponsiveContainer width="100%" height="100%">
                      <BarChart data={chartData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="name" />
                        <YAxis label={{ value: 'Hours', angle: -90, position: 'insideLeft' }} />
                        <Tooltip />
                        {currentGroups.slice(0, 5).map((group, index) => (
                          <Bar 
                            key={group.id} 
                            dataKey={group.name} 
                            stackId="a"
                            fill={[ "#3B82F6", "#10B981", "#F43F5E", "#8B5CF6", "#F59E0B" ][index % 5]} 
                          />
                        ))}
                      </BarChart>
                    </ResponsiveContainer>
                  </Box>
                </CardContent>
              </Card>
              
              {/* Full Rankings List */}
              <Box sx={{ gridColumn: { xs: 'span 1', lg: 'span 2' } }}>
                <Card elevation={3}>
                  <CardContent sx={{ pt: 3 }}>
                    <Typography variant="h6" sx={{ fontWeight: 600, mb: 2, color: '#4B5563' }}>
                      Complete Rankings
                    </Typography>
                    <TableContainer component={Paper} sx={{ border: 1, borderColor: 'divider', borderRadius: 1 }}>
                      <Table sx={{ minWidth: 650 }}>
                        <TableHead sx={{ bgcolor: '#F9FAFB' }}>
                          <TableRow>
                            <TableCell width="10%" align="center">Rank</TableCell>
                            <TableCell width="40%">Group</TableCell>
                            <TableCell width="20%" align="center">Members</TableCell>
                            <TableCell width="30%" align="center">Total Hours</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {currentGroups.map((group, index) => (
                            <TableRow key={group.id}>
                              <TableCell align="center">
                                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', width: 32, height: 32, mx: 'auto' }}>
                                  {getMedalIcon(index)}
                                </Box>
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                  {group.name}
                                </Typography>
                              </TableCell>
                              <TableCell align="center">
                                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                                  <Users size={18} color="#6B7280" />
                                  <Typography sx={{ ml: 1 }}>{group.members}</Typography>
                                </Box>
                              </TableCell>
                              <TableCell align="center">
                                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                                  <Clock size={18} color="#6B7280" />
                                  <Typography sx={{ ml: 1 }}>{group.totalHours}h</Typography>
                                </Box>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </CardContent>
                </Card>
              </Box>
            </Box>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default LeadershipPage;
