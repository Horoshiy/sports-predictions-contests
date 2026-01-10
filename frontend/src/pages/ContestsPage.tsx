import React, { useState } from 'react'
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
} from '@mui/material'
import ContestList from '../components/contests/ContestList'
import ContestForm from '../components/contests/ContestForm'
import ParticipantList from '../components/contests/ParticipantList'
import { LeaderboardTable } from '../components/leaderboard/LeaderboardTable'
import { UserScore } from '../components/leaderboard/UserScore'
import {
  useCreateContest,
  useUpdateContest,
} from '../hooks/use-contests'
import { useToast } from '../contexts/ToastContext'
import { useAuth } from '../contexts/AuthContext'
import type { Contest, ContestFormData } from '../types/contest.types'
import { toISOString } from '../utils/date-utils'

export const ContestsPage: React.FC = () => {
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [isParticipantsOpen, setIsParticipantsOpen] = useState(false)
  const [selectedContest, setSelectedContest] = useState<Contest | null>(null)
  const [tabValue, setTabValue] = useState(0)

  const createContestMutation = useCreateContest()
  const updateContestMutation = useUpdateContest()
  const { showToast } = useToast()
  const { user } = useAuth()

  const handleCreateContest = () => {
    setSelectedContest(null)
    setIsFormOpen(true)
  }

  const handleEditContest = (contest: Contest) => {
    setSelectedContest(contest)
    setIsFormOpen(true)
  }

  const handleViewParticipants = (contest: Contest) => {
    setSelectedContest(contest)
    setIsParticipantsOpen(true)
  }

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
  }

  const handleEditContest = (contest: Contest) => {
    setSelectedContest(contest)
    setIsFormOpen(true)
  }

  const handleViewParticipants = (contest: Contest) => {
    setSelectedContest(contest)
    setIsParticipantsOpen(true)
  }

  const handleFormSubmit = async (data: ContestFormData) => {
    try {
      const contestData = {
        title: data.title,
        description: data.description || '',
        sportType: data.sportType,
        rules: data.rules || '',
        startDate: toISOString(data.startDate),
        endDate: toISOString(data.endDate),
        maxParticipants: data.maxParticipants,
      }

      if (selectedContest) {
        // Update existing contest
        await updateContestMutation.mutateAsync({
          id: selectedContest.id,
          ...contestData,
          status: selectedContest.status,
        })
      } else {
        // Create new contest
        await createContestMutation.mutateAsync(contestData)
      }

      setIsFormOpen(false)
      setSelectedContest(null)
    } catch (error) {
      console.error('Failed to save contest:', error)
      // Error handling is now done in the hooks with toast notifications
    }
  }

  const handleFormClose = () => {
    setIsFormOpen(false)
    setSelectedContest(null)
  }

  const handleParticipantsClose = () => {
    setIsParticipantsOpen(false)
    setSelectedContest(null)
  }

  return (
    <Box>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Contest Management
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Create and manage sports prediction contests
        </Typography>
      </Box>

      <Paper sx={{ mb: 2 }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="Contests" />
          <Tab label="Leaderboards" />
        </Tabs>
      </Paper>

      {tabValue === 0 && (
        <Paper sx={{ p: 0, overflow: 'hidden' }}>
          <ContestList
            onCreateContest={handleCreateContest}
            onEditContest={handleEditContest}
            onViewParticipants={handleViewParticipants}
          />
        </Paper>
      )}

      {tabValue === 1 && (
        <Box>
          <Typography variant="h6" gutterBottom>
            Contest Leaderboards
          </Typography>
          {selectedContest ? (
            <LeaderboardTable 
              contestId={selectedContest.id}
              currentUserId={user?.id || 0}
            />
          ) : (
            <Paper sx={{ p: 4, textAlign: 'center' }}>
              <Typography color="text.secondary">
                Select a contest to view its leaderboard
              </Typography>
            </Paper>
          )}
        </Box>
      )}

      <ContestForm
        open={isFormOpen}
        onClose={handleFormClose}
        onSubmit={handleFormSubmit}
        contest={selectedContest}
        loading={createContestMutation.isPending || updateContestMutation.isPending}
      />

      <ParticipantList
        open={isParticipantsOpen}
        onClose={handleParticipantsClose}
        contest={selectedContest}
      />
    </Box>
  )
}

export default ContestsPage
