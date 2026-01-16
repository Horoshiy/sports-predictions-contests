import React, { useState } from 'react'
import { Box, Container, Typography, Tabs, Tab, Paper, TextField, Button, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { joinTeamSchema, type JoinTeamSchemaType } from '../utils/team-validation'
import { useAuth } from '../contexts/AuthContext'
import { useJoinTeam, useLeaveTeam } from '../hooks/use-teams'
import TeamList from '../components/teams/TeamList'
import TeamForm from '../components/teams/TeamForm'
import TeamMembers from '../components/teams/TeamMembers'
import TeamInvite from '../components/teams/TeamInvite'
import type { Team } from '../types/team.types'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

const TabPanel: React.FC<TabPanelProps> = ({ children, value, index }) => (
  <div hidden={value !== index}>{value === index && <Box sx={{ py: 3 }}>{children}</Box>}</div>
)

export const TeamsPage: React.FC = () => {
  const { user } = useAuth()
  const [tabValue, setTabValue] = useState(0)
  const [formOpen, setFormOpen] = useState(false)
  const [editTeam, setEditTeam] = useState<Team | null>(null)
  const [viewTeam, setViewTeam] = useState<Team | null>(null)

  const joinTeamMutation = useJoinTeam()
  const leaveTeamMutation = useLeaveTeam()

  const { control, handleSubmit, reset, formState: { errors } } = useForm<JoinTeamSchemaType>({
    resolver: zodResolver(joinTeamSchema),
    defaultValues: { inviteCode: '' },
  })

  const handleCreateTeam = () => {
    setEditTeam(null)
    setFormOpen(true)
  }

  const handleEditTeam = (team: Team) => {
    setEditTeam(team)
    setFormOpen(true)
  }

  const handleViewMembers = (team: Team) => {
    setViewTeam(team)
  }

  const handleJoinTeam = async (data: JoinTeamSchemaType) => {
    try {
      await joinTeamMutation.mutateAsync(data)
      reset()
      setTabValue(0) // Switch to My Teams
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleLeaveTeam = (team: Team) => {
    if (window.confirm(`Are you sure you want to leave "${team.name}"?`)) {
      leaveTeamMutation.mutate(team.id)
      setViewTeam(null)
    }
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>Teams</Typography>

      <Paper sx={{ mb: 3 }}>
        <Tabs value={tabValue} onChange={(_, v) => setTabValue(v)}>
          <Tab label="My Teams" />
          <Tab label="All Teams" />
          <Tab label="Join Team" />
        </Tabs>
      </Paper>

      <TabPanel value={tabValue} index={0}>
        <TeamList onCreateTeam={handleCreateTeam} onEditTeam={handleEditTeam} onViewMembers={handleViewMembers} myTeamsOnly />
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        <TeamList onCreateTeam={handleCreateTeam} onEditTeam={handleEditTeam} onViewMembers={handleViewMembers} />
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        <Paper sx={{ p: 3, maxWidth: 400 }}>
          <Typography variant="h6" gutterBottom>Join a Team</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Enter the invite code shared by your team captain.
          </Typography>
          <form onSubmit={handleSubmit(handleJoinTeam)}>
            <Controller
              name="inviteCode"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Invite Code"
                  fullWidth
                  placeholder="e.g., A1B2C3D4"
                  error={!!errors.inviteCode}
                  helperText={errors.inviteCode?.message}
                  sx={{ mb: 2 }}
                />
              )}
            />
            <Button type="submit" variant="contained" fullWidth disabled={joinTeamMutation.isPending}>
              {joinTeamMutation.isPending ? 'Joining...' : 'Join Team'}
            </Button>
          </form>
        </Paper>
      </TabPanel>

      <TeamForm open={formOpen} onClose={() => setFormOpen(false)} team={editTeam} />

      <Dialog open={!!viewTeam} onClose={() => setViewTeam(null)} maxWidth="sm" fullWidth>
        {viewTeam && (
          <>
            <DialogTitle>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="h6">{viewTeam.name}</Typography>
                <Typography variant="body2" color="text.secondary">
                  {viewTeam.currentMembers} / {viewTeam.maxMembers} members
                </Typography>
              </Box>
            </DialogTitle>
            <DialogContent>
              {viewTeam.description && (
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>{viewTeam.description}</Typography>
              )}
              <TeamInvite teamId={viewTeam.id} inviteCode={viewTeam.inviteCode} isCaptain={viewTeam.captainId === user?.id} />
              <Typography variant="subtitle2" sx={{ mt: 3, mb: 1 }}>Members</Typography>
              <TeamMembers team={viewTeam} />
            </DialogContent>
            <DialogActions>
              {viewTeam.captainId !== user?.id && (
                <Button color="error" onClick={() => handleLeaveTeam(viewTeam)} disabled={leaveTeamMutation.isPending}>
                  Leave Team
                </Button>
              )}
              <Button onClick={() => setViewTeam(null)}>Close</Button>
            </DialogActions>
          </>
        )}
      </Dialog>
    </Container>
  )
}

export default TeamsPage
