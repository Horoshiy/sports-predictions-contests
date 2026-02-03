import React, { useState } from 'react'
import { Space, Typography, Tabs } from 'antd'
import SportList from '../components/sports/SportList'
import SportForm from '../components/sports/SportForm'
import LeagueList from '../components/sports/LeagueList'
import LeagueForm from '../components/sports/LeagueForm'
import TeamList from '../components/sports/TeamList'
import TeamForm from '../components/sports/TeamForm'
import MatchList from '../components/sports/MatchList'
import MatchForm from '../components/sports/MatchForm'
import RiskyEventTypesManager from '../components/admin/RiskyEventTypesManager'
import {
  useCreateSport, useUpdateSport,
  useCreateLeague, useUpdateLeague,
  useCreateTeam, useUpdateTeam,
  useCreateMatch, useUpdateMatch,
} from '../hooks/use-sports'
import { generateSlug } from '../utils/sports-validation'
import { showSuccess, showError } from '../utils/notification'
import type { Sport, League, Team, Match, CreateSportRequest, UpdateSportRequest, CreateLeagueRequest, UpdateLeagueRequest, CreateTeamRequest, UpdateTeamRequest, CreateMatchRequest, UpdateMatchRequest } from '../types/sports.types'
import type { SportFormData, LeagueFormData, TeamFormData, MatchFormData } from '../utils/sports-validation'

const { Title } = Typography

type EntityType = 'sport' | 'league' | 'team' | 'match'

export const SportsPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('sports')
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
    setFormOpen(true)
    
    if (!entity) {
      setSelectedSport(null)
      setSelectedLeague(null)
      setSelectedTeam(null)
      setSelectedMatch(null)
      return
    }
    
    switch (type) {
      case 'sport':
        setSelectedSport(entity as Sport)
        break
      case 'league':
        setSelectedLeague(entity as League)
        break
      case 'team':
        setSelectedTeam(entity as Team)
        break
      case 'match':
        setSelectedMatch(entity as Match)
        break
    }
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
        const request: UpdateSportRequest = {
          id: selectedSport.id,
          name: data.name,
          slug,
          description: data.description,
          iconUrl: data.iconUrl,
          isActive: true,
        }
        await updateSport.mutateAsync(request)
        showSuccess('Sport updated successfully')
      } else {
        const request: CreateSportRequest = {
          name: data.name,
          slug,
          description: data.description,
          iconUrl: data.iconUrl,
        }
        await createSport.mutateAsync(request)
        showSuccess('Sport created successfully')
      }
      closeForm()
    } catch (error: any) {
      showError(error?.message || 'Failed to save sport')
    }
  }

  const handleLeagueSubmit = async (data: LeagueFormData) => {
    try {
      const slug = data.slug || generateSlug(data.name)
      if (selectedLeague) {
        const request: UpdateLeagueRequest = {
          id: selectedLeague.id,
          name: data.name,
          slug,
          sportId: data.sportId,
          country: data.country,
          season: data.season,
          isActive: true,
        }
        await updateLeague.mutateAsync(request)
        showSuccess('League updated successfully')
      } else {
        const request: CreateLeagueRequest = {
          name: data.name,
          slug,
          sportId: data.sportId,
          country: data.country,
          season: data.season,
        }
        await createLeague.mutateAsync(request)
        showSuccess('League created successfully')
      }
      closeForm()
    } catch (error: any) {
      showError(error?.message || 'Failed to save league')
    }
  }

  const handleTeamSubmit = async (data: TeamFormData) => {
    try {
      const slug = data.slug || generateSlug(data.name)
      if (selectedTeam) {
        const request: UpdateTeamRequest = {
          id: selectedTeam.id,
          name: data.name,
          slug,
          sportId: data.sportId,
          country: data.country,
          shortName: data.shortName,
          logoUrl: data.logoUrl,
          isActive: true,
        }
        await updateTeam.mutateAsync(request)
        showSuccess('Team updated successfully')
      } else {
        const request: CreateTeamRequest = {
          name: data.name,
          slug,
          sportId: data.sportId,
          country: data.country,
          shortName: data.shortName,
          logoUrl: data.logoUrl,
        }
        await createTeam.mutateAsync(request)
        showSuccess('Team created successfully')
      }
      closeForm()
    } catch (error: any) {
      showError(error?.message || 'Failed to save team')
    }
  }

  const handleMatchSubmit = async (data: MatchFormData) => {
    try {
      const scheduledAt = data.scheduledAt instanceof Date 
        ? data.scheduledAt.toISOString() 
        : data.scheduledAt

      if (selectedMatch) {
        const request: UpdateMatchRequest = {
          id: selectedMatch.id,
          leagueId: data.leagueId,
          homeTeamId: data.homeTeamId,
          awayTeamId: data.awayTeamId,
          scheduledAt,
          status: data.status,
          homeScore: data.homeScore,
          awayScore: data.awayScore,
          resultData: data.resultData,
        }
        await updateMatch.mutateAsync(request)
        showSuccess('Match updated successfully')
      } else {
        const request: CreateMatchRequest = {
          leagueId: data.leagueId,
          homeTeamId: data.homeTeamId,
          awayTeamId: data.awayTeamId,
          scheduledAt,
        }
        await createMatch.mutateAsync(request)
        showSuccess('Match created successfully')
      }
      closeForm()
    } catch (error: any) {
      showError(error?.message || 'Failed to save match')
    }
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', padding: '24px' }}>
      <Title level={2}>Sports Management</Title>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'sports',
            label: 'Sports',
            children: (
              <SportList
                onCreateSport={() => openForm('sport')}
                onEditSport={(sport) => openForm('sport', sport)}
              />
            ),
          },
          {
            key: 'leagues',
            label: 'Leagues',
            children: (
              <LeagueList
                onCreateLeague={() => openForm('league')}
                onEditLeague={(league) => openForm('league', league)}
              />
            ),
          },
          {
            key: 'teams',
            label: 'Teams',
            children: (
              <TeamList
                onCreateTeam={() => openForm('team')}
                onEditTeam={(team) => openForm('team', team)}
              />
            ),
          },
          {
            key: 'matches',
            label: 'Matches',
            children: (
              <MatchList
                onCreateMatch={() => openForm('match')}
                onEditMatch={(match) => openForm('match', match)}
              />
            ),
          },
          {
            key: 'risky-events',
            label: 'Risky Events',
            children: <RiskyEventTypesManager />,
          },
        ]}
      />

      {entityType === 'sport' && (
        <SportForm
          open={formOpen}
          onClose={closeForm}
          onSubmit={handleSportSubmit}
          sport={selectedSport}
          loading={createSport.isPending || updateSport.isPending}
        />
      )}

      {entityType === 'league' && (
        <LeagueForm
          open={formOpen}
          onClose={closeForm}
          onSubmit={handleLeagueSubmit}
          league={selectedLeague}
          loading={createLeague.isPending || updateLeague.isPending}
        />
      )}

      {entityType === 'team' && (
        <TeamForm
          open={formOpen}
          onClose={closeForm}
          onSubmit={handleTeamSubmit}
          team={selectedTeam}
          loading={createTeam.isPending || updateTeam.isPending}
        />
      )}

      {entityType === 'match' && (
        <MatchForm
          open={formOpen}
          onClose={closeForm}
          onSubmit={handleMatchSubmit}
          match={selectedMatch}
          loading={createMatch.isPending || updateMatch.isPending}
        />
      )}
    </Space>
  )
}

export default SportsPage
