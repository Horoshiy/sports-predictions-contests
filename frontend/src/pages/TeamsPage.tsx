import React, { useState } from 'react'
import { Space, Typography, Tabs, Input, Button, Modal } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { joinTeamSchema, type JoinTeamSchemaType } from '../utils/team-validation'
import { useAuth } from '../contexts/AuthContext'
import { useJoinTeam, useLeaveTeam } from '../hooks/use-teams'
import { showError } from '../utils/notification'
import TeamList from '../components/teams/TeamList'
import TeamForm from '../components/teams/TeamForm'
import TeamMembers from '../components/teams/TeamMembers'
import TeamInvite from '../components/teams/TeamInvite'
import type { Team } from '../types/team.types'

const { Title } = Typography

export const TeamsPage: React.FC = () => {
  const { user } = useAuth()
  const [activeTab, setActiveTab] = useState('all')
  const [formOpen, setFormOpen] = useState(false)
  const [editTeam, setEditTeam] = useState<Team | null>(null)
  const [viewTeam, setViewTeam] = useState<Team | null>(null)
  const [joinModalOpen, setJoinModalOpen] = useState(false)

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
      await joinTeamMutation.mutateAsync({ inviteCode: data.inviteCode })
      reset()
      setJoinModalOpen(false)
      setActiveTab('my')
    } catch (error: any) {
      showError(error?.message || 'Failed to join team')
    }
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', padding: '24px' }}>
      <Title level={2}>Teams</Title>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'my',
            label: 'My Teams',
            children: (
              <TeamList
                onCreateTeam={handleCreateTeam}
                onEditTeam={handleEditTeam}
                onViewMembers={handleViewMembers}
                myTeamsOnly={true}
              />
            ),
          },
          {
            key: 'all',
            label: 'All Teams',
            children: (
              <TeamList
                onCreateTeam={handleCreateTeam}
                onEditTeam={handleEditTeam}
                onViewMembers={handleViewMembers}
                myTeamsOnly={false}
              />
            ),
          },
          {
            key: 'join',
            label: 'Join Team',
            children: (
              <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                <Title level={4}>Join a Team with Invite Code</Title>
                <form onSubmit={handleSubmit(handleJoinTeam)}>
                  <Space direction="vertical" style={{ width: '100%' }}>
                    <Controller
                      name="inviteCode"
                      control={control}
                      render={({ field }) => (
                        <Input
                          {...field}
                          placeholder="Enter invite code"
                          status={errors.inviteCode ? 'error' : ''}
                          style={{ width: 300 }}
                        />
                      )}
                    />
                    {errors.inviteCode && (
                      <span style={{ color: '#ff4d4f' }}>{errors.inviteCode.message}</span>
                    )}
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={joinTeamMutation.isPending}
                    >
                      Join Team
                    </Button>
                  </Space>
                </form>
              </Space>
            ),
          },
        ]}
      />

      <TeamForm
        open={formOpen}
        onClose={() => {
          setFormOpen(false)
          setEditTeam(null)
        }}
        team={editTeam}
      />

      {viewTeam && (
        <Modal
          open={!!viewTeam}
          title={`${viewTeam.name} - Members`}
          onCancel={() => setViewTeam(null)}
          footer={null}
          width={800}
        >
          <TeamMembers team={viewTeam} />
          {/* TeamInvite component requires different props */}
        </Modal>
      )}
    </Space>
  )
}

export default TeamsPage
