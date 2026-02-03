import React, { useState } from 'react'
import { Typography, Card, Tabs, Space } from 'antd'
import ContestList from '../components/contests/ContestList'
import ContestForm from '../components/contests/ContestForm'
import ParticipantList from '../components/contests/ParticipantList'
import ContestEventsManager from '../components/contests/ContestEventsManager'
import { LeaderboardTable } from '../components/leaderboard/LeaderboardTable'
import TeamLeaderboard from '../components/teams/TeamLeaderboard'
import {
  useCreateContest,
  useUpdateContest,
} from '../hooks/use-contests'
import { useAuth } from '../contexts/AuthContext'
import { showError, showSuccess } from '../utils/notification'
import type { Contest, ContestFormData } from '../types/contest.types'
import { toISOString } from '../utils/date-utils'

const { Title, Text } = Typography

const ContestsPage: React.FC = () => {
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [isParticipantsOpen, setIsParticipantsOpen] = useState(false)
  const [selectedContest, setSelectedContest] = useState<Contest | null>(null)
  const [tabValue, setTabValue] = useState('1')

  const createContestMutation = useCreateContest()
  const updateContestMutation = useUpdateContest()
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

  const handleSelectContest = (contest: Contest) => {
    setSelectedContest(contest)
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

      console.log('Submitting contest data:', contestData)

      if (selectedContest) {
        await updateContestMutation.mutateAsync({
          id: selectedContest.id,
          ...contestData,
          status: selectedContest.status,
        })
        showSuccess('Contest updated successfully')
      } else {
        const result = await createContestMutation.mutateAsync(contestData)
        console.log('Create result:', result)
        showSuccess('Contest created successfully')
      }

      setIsFormOpen(false)
      setSelectedContest(null)
    } catch (error: any) {
      console.error('Form submit error:', error)
      showError(error?.message || 'Failed to save contest')
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
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div>
        <Title level={2}>Contest Management</Title>
        <Text type="secondary">Create and manage sports prediction contests</Text>
      </div>

      <Card>
        <Tabs
          activeKey={tabValue}
          onChange={setTabValue}
          items={[
            {
              key: '1',
              label: 'Contests',
              children: (
                <ContestList
                  onCreateContest={handleCreateContest}
                  onEditContest={handleEditContest}
                  onViewParticipants={handleViewParticipants}
                  onSelectContest={handleSelectContest}
                  selectedContestId={selectedContest?.id}
                />
              ),
            },
            {
              key: '2',
              label: 'Leaderboards',
              children: selectedContest ? (
                <LeaderboardTable 
                  contestId={selectedContest.id}
                  currentUserId={user?.id || 0}
                />
              ) : (
                <div style={{ padding: 48, textAlign: 'center' }}>
                  <Text type="secondary">Select a contest to view its leaderboard</Text>
                </div>
              ),
            },
            {
              key: '3',
              label: 'Team Leaderboard',
              children: selectedContest ? (
                <div data-testid="team-leaderboard">
                  <TeamLeaderboard contestId={selectedContest.id} />
                </div>
              ) : (
                <div style={{ padding: 48, textAlign: 'center' }}>
                  <Text type="secondary">Select a contest to view team leaderboard</Text>
                </div>
              ),
            },
            {
              key: '4',
              label: 'Events',
              children: selectedContest ? (
                <ContestEventsManager contest={selectedContest} />
              ) : (
                <div style={{ padding: 48, textAlign: 'center' }}>
                  <Text type="secondary">Select a contest to manage its events</Text>
                </div>
              ),
            },
          ]}
        />
      </Card>

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
    </Space>
  )
}

export { ContestsPage as default }
