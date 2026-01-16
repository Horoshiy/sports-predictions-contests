import React from 'react'
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Typography, Box, Chip } from '@mui/material'
import { EmojiEvents as TrophyIcon, People as PeopleIcon } from '@mui/icons-material'
import { useTeamLeaderboard } from '../../hooks/use-teams'

interface TeamLeaderboardProps {
  contestId: number
  userTeamId?: number
}

const getRankColor = (rank: number) => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'silver'
  if (rank === 3) return '#cd7f32'
  return undefined
}

export const TeamLeaderboard: React.FC<TeamLeaderboardProps> = ({ contestId, userTeamId }) => {
  const { data: entries, isLoading, isError } = useTeamLeaderboard(contestId, 20)

  if (isLoading) return <Typography>Loading team leaderboard...</Typography>
  if (isError) return <Typography color="error">Failed to load team leaderboard</Typography>
  if (!entries || entries.length === 0) return <Typography color="text.secondary">No teams in this contest yet</Typography>

  return (
    <TableContainer component={Paper}>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell width={60}>Rank</TableCell>
            <TableCell>Team</TableCell>
            <TableCell align="right">Members</TableCell>
            <TableCell align="right">Points</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {entries.map((entry, index) => {
            const rank = entry.rank || index + 1
            const isUserTeam = entry.teamId === userTeamId
            return (
              <TableRow key={entry.teamId} sx={{ bgcolor: isUserTeam ? 'action.selected' : undefined }}>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    {rank <= 3 && <TrophyIcon sx={{ color: getRankColor(rank), mr: 0.5, fontSize: 18 }} />}
                    <Typography fontWeight={rank <= 3 ? 'bold' : 'normal'}>{rank}</Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography fontWeight={isUserTeam ? 'bold' : 'normal'}>{entry.teamName}</Typography>
                    {isUserTeam && <Chip label="Your Team" size="small" color="primary" />}
                  </Box>
                </TableCell>
                <TableCell align="right">
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }}>
                    <PeopleIcon sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                    {entry.memberCount}
                  </Box>
                </TableCell>
                <TableCell align="right">
                  <Typography fontWeight="bold">{entry.totalPoints.toFixed(1)}</Typography>
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </TableContainer>
  )
}

export default TeamLeaderboard
