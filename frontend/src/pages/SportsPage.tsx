import React, { useState } from 'react'
import { Box, Typography, Paper, Tabs, Tab } from '@mui/material'
import SportList from '../components/sports/SportList'
import SportForm from '../components/sports/SportForm'
import LeagueList from '../components/sports/LeagueList'
import LeagueForm from '../components/sports/LeagueForm'
import TeamList from '../components/sports/TeamList'
import TeamForm from '../components/sports/TeamForm'
import MatchList from '../components/sports/MatchList'
import MatchForm from '../components/sports/MatchForm'
import {
  useCreateSport, useUpdateSport,
  useCreateLeague, useUpdateLeague,
  useCreateTeam, useUpdateTeam,
  useCreateMatch, useUpdateMatch,
} from '../hooks/use-sports'
import { generateSlug } from '../utils/sports-validation'
import type { Sport, League, Team, Match } from '../types/sports.types'
import type { SportFormData, LeagueFormData, TeamFormData, MatchFormData } from '../utils/sports-validation'

type EntityType = 'sport' | 'league' | 'team' | 'match'

export const SportsPage: React.FC = () => {
  const [tabValue, setTabValue] = useState(0)
  const [formOpen, setFormOpen] = useState(false)
  const [entityType, setEntityType] = useState<EntityType>('sport')
  const [selectedSport, setSelectedSport] = useState<Sport | null>(null)
  const [selectedLeague, setSelectedLeague] = useState<League | null>(null)
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null)
  const [selectedMatch, setSelectedMatch] = useState<Match | null>(null)

  const createSport = useCreateSport()
  const updateSport = useUpdateSport()
  const createLeague = useCreateLeague()
  const updateLeague = useUpdateLeague()
  const createTeam = useCreateTeam()
  const updateTeam = useUpdateTeam()
  const createMatch = useCreateMatch()
  const updateMatch = useUpdateMatch()

  const openForm = (type: EntityType, entity?: Sport | League | Team | Match) => {
    setEntityType(type)
    if (type === 'sport') setSelectedSport(entity as Sport || null)
    else if (type === 'league') setSelectedLeague(entity as League || null)
    else if (type === 'team') setSelectedTeam(entity as Team || null)
    else if (type === 'match') setSelectedMatch(entity as Match || null)
    setFormOpen(true)
  }

  const closeForm = () => {
    setFormOpen(false)
    setSelectedSport(null)
    setSelectedLeague(null)
    setSelectedTeam(null)
    setSelectedMatch(null)
  }

  const handleSportSubmit = async (data: SportFormData) => {
    try {
      const slug = data.slug || generateSlug(data.name)
      if (selectedSport) {
        await updateSport.mutateAsync({ id: selectedSport.id, ...data, slug, isActive: selectedSport.isActive })
      } else {
        await createSport.mutateAsync({ ...data, slug })
      }
      closeForm()
    } catch {
      // Error already shown via toast in hook
    }
  }

  const handleLeagueSubmit = async (data: LeagueFormData) => {
    try {
      const slug = data.slug || generateSlug(data.name)
      if (selectedLeague) {
        await updateLeague.mutateAsync({ id: selectedLeague.id, ...data, slug, isActive: selectedLeague.isActive })
      } else {
        await createLeague.mutateAsync({ ...data, slug })
      }
      closeForm()
    } catch {
      // Error already shown via toast in hook
    }
  }

  const handleTeamSubmit = async (data: TeamFormData) => {
    try {
      const slug = data.slug || generateSlug(data.name)
      if (selectedTeam) {
        await updateTeam.mutateAsync({ id: selectedTeam.id, ...data, slug, isActive: selectedTeam.isActive })
      } else {
        await createTeam.mutateAsync({ ...data, slug })
      }
      closeForm()
    } catch {
      // Error already shown via toast in hook
    }
  }

  const handleMatchSubmit = async (data: MatchFormData) => {
    try {
      const scheduledAt = data.scheduledAt.toISOString()
      if (selectedMatch) {
        await updateMatch.mutateAsync({
          id: selectedMatch.id,
          leagueId: data.leagueId,
          homeTeamId: data.homeTeamId,
          awayTeamId: data.awayTeamId,
          scheduledAt,
          status: data.status || selectedMatch.status,
          homeScore: data.homeScore,
          awayScore: data.awayScore,
          resultData: data.resultData,
        })
      } else {
        await createMatch.mutateAsync({ leagueId: data.leagueId, homeTeamId: data.homeTeamId, awayTeamId: data.awayTeamId, scheduledAt })
      }
      closeForm()
    } catch {
      // Error already shown via toast in hook
    }
  }

  const isLoading = createSport.isPending || updateSport.isPending || createLeague.isPending || updateLeague.isPending ||
    createTeam.isPending || updateTeam.isPending || createMatch.isPending || updateMatch.isPending

  return (
    <Box>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>Sports Management</Typography>
        <Typography variant="body1" color="text.secondary">Manage sports, leagues, teams, and matches</Typography>
      </Box>

      <Paper sx={{ mb: 2 }}>
        <Tabs value={tabValue} onChange={(_, v) => setTabValue(v)}>
          <Tab label="Sports" />
          <Tab label="Leagues" />
          <Tab label="Teams" />
          <Tab label="Matches" />
        </Tabs>
      </Paper>

      <Paper sx={{ p: 0, overflow: 'hidden' }}>
        {tabValue === 0 && <SportList onCreateSport={() => openForm('sport')} onEditSport={(s) => openForm('sport', s)} />}
        {tabValue === 1 && <LeagueList onCreateLeague={() => openForm('league')} onEditLeague={(l) => openForm('league', l)} />}
        {tabValue === 2 && <TeamList onCreateTeam={() => openForm('team')} onEditTeam={(t) => openForm('team', t)} />}
        {tabValue === 3 && <MatchList onCreateMatch={() => openForm('match')} onEditMatch={(m) => openForm('match', m)} />}
      </Paper>

      <SportForm open={formOpen && entityType === 'sport'} onClose={closeForm} onSubmit={handleSportSubmit} sport={selectedSport} loading={isLoading} />
      <LeagueForm open={formOpen && entityType === 'league'} onClose={closeForm} onSubmit={handleLeagueSubmit} league={selectedLeague} loading={isLoading} />
      <TeamForm open={formOpen && entityType === 'team'} onClose={closeForm} onSubmit={handleTeamSubmit} team={selectedTeam} loading={isLoading} />
      <MatchForm open={formOpen && entityType === 'match'} onClose={closeForm} onSubmit={handleMatchSubmit} match={selectedMatch} loading={isLoading} />
    </Box>
  )
}

export default SportsPage
