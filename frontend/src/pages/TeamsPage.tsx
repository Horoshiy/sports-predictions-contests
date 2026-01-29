import React, { useState } from 'react'
import { Space, Typography, Tabs, Input, Button, Modal, Descriptions, Popconfirm } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { joinTeamSchema, type JoinTeamSchemaType } from '../utils/team-validation'
import { useAuth } from '../contexts/AuthContext'
import { useJoinTeam, useLeaveTeam, useDeleteTeam } from '../hooks/use-teams'
import { showError, showSuccess } from '../utils/notification'
import TeamList from '../components/teams/TeamList'
import TeamForm from '../components/teams/TeamForm'
import TeamMembers from '../components/teams/TeamMembers'
import TeamInvite from '../components/teams/TeamInvite'
import type { Team } from '../types/team.types'

const { Title } = Typography

export const TeamsPage: React.FC = () => {
  const { user } = useAuth()
  const [activeTab, setActiveTab] = useState('my')
  const [formOpen, setFormOpen] = useState(false)
  const [editTeam, setEditTeam] = useState<Team | null>(null)
  const [viewTeam, setViewTeam] = useState<Team | null>(null)

  const joinTeamMutation = useJoinTeam()
  const leaveTeamMutation = useLeaveTeam()
  const deleteTeamMutation = useDeleteTeam()

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
      setActiveTab('my')
      showSuccess('Successfully joined team!')
    } catch (error: any) {
      showError(error?.message || 'Failed to join team')
    }
  }

  const handleLeaveTeam = async () => {
    if (!viewTeam) return
    try {
      await leaveTeamMutation.mutateAsync(viewTeam.id)
      setViewTeam(null)
      showSuccess('Successfully left team!')
    } catch (error: any) {
      showError(error?.message || 'Failed to leave team')
    }
  }

  const handleDeleteTeam = async () => {
    if (!viewTeam) return
    try {
      await deleteTeamMutation.mutateAsync(viewTeam.id)
      setViewTeam(null)
      showSuccess('Team deleted successfully!')
    } catch (error: any) {
      showError(error?.message || 'Failed to delete team')
    }
  }

  const isCaptain = viewTeam && user && viewTeam.captainId === user.id

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', padding: '24px' }} data-testid="teams-page">
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
                <form onSubmit={handleSubmit(handleJoinTeam)} data-testid="join-team-form">
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
          title={viewTeam.name}
          onCancel={() => setViewTeam(null)}
          width={800}
          data-testid="team-details-modal"
          footer={[
            !isCaptain ? (
              <Popconfirm
                key="leave"
                title="Leave team?"
                description="Are you sure you want to leave this team?"
                onConfirm={handleLeaveTeam}
                okText="Yes"
                cancelText="No"
              >
                <Button loading={leaveTeamMutation.isPending}>Leave Team</Button>
              </Popconfirm>
            ) : null,
            isCaptain ? (
              <Popconfirm
                key="delete"
                title="Delete team?"
                description="This action cannot be undone. All members will be removed."
                onConfirm={handleDeleteTeam}
                okText="Yes"
                cancelText="No"
              >
                <Button danger loading={deleteTeamMutation.isPending}>Delete Team</Button>
              </Popconfirm>
            ) : null,
            <Button key="close" type="primary" onClick={() => setViewTeam(null)}>
              Close
            </Button>,
          ].filter(Boolean)}
        >
          <Tabs
            items={[
              {
                key: 'info',
                label: 'Team Info',
                children: (
                  <Descriptions column={1} bordered>
                    <Descriptions.Item label="Description">{viewTeam.description || 'No description'}</Descriptions.Item>
                    <Descriptions.Item label="Members">{viewTeam.currentMembers} / {viewTeam.maxMembers}</Descriptions.Item>
                    <Descriptions.Item label="Status">{viewTeam.isActive ? 'Active' : 'Inactive'}</Descriptions.Item>
                  </Descriptions>
                ),
              },
              {
                key: 'members',
                label: 'Members',
                children: <TeamMembers team={viewTeam} />,
              },
              {
                key: 'invite',
                label: 'Invite Code',
                children: <TeamInvite teamId={viewTeam.id} inviteCode={viewTeam.inviteCode} isCaptain={!!isCaptain} />,
              },
            ]}
          />
        </Modal>
      )}
    </Space>
  )
}

export default TeamsPage
