import React from 'react'
import { Space, Typography, Input, Button, Tooltip, Popconfirm } from 'antd'
import { CopyOutlined, ReloadOutlined } from '@ant-design/icons'
import { showSuccess, showWarning } from '../../utils/antd-helpers'
import { useRegenerateInviteCode } from '../../hooks/use-teams'

const { Text } = Typography

interface TeamInviteProps {
  teamId: number
  inviteCode: string
  isCaptain: boolean
}

export const TeamInvite: React.FC<TeamInviteProps> = ({ teamId, inviteCode, isCaptain }) => {
  const regenerateMutation = useRegenerateInviteCode()

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(inviteCode)
      showSuccess('Invite code copied!')
    } catch {
      showWarning('Failed to copy - please copy manually')
    }
  }

  const handleRegenerate = () => {
    regenerateMutation.mutate(teamId)
  }

  return (
    <div style={{ padding: 16, backgroundColor: '#f5f5f5', borderRadius: 4 }}>
      <Text type="secondary" strong>Invite Code</Text>
      <Space.Compact style={{ width: '100%', marginTop: 8 }}>
        <Input value={inviteCode} readOnly style={{ fontFamily: 'monospace' }} />
        <Tooltip title="Copy code">
          <Button icon={<CopyOutlined />} onClick={handleCopy} />
        </Tooltip>
        {isCaptain && (
          <Popconfirm
            title="Regenerate invite code?"
            description="The old code will stop working."
            onConfirm={handleRegenerate}
            okText="Yes"
            cancelText="No"
          >
            <Tooltip title="Regenerate code">
              <Button icon={<ReloadOutlined />} loading={regenerateMutation.isPending} danger />
            </Tooltip>
          </Popconfirm>
        )}
      </Space.Compact>
    </div>
  )
}

export default TeamInvite
