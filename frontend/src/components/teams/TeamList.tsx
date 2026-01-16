import React, { useMemo, useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { MaterialReactTable, useMaterialReactTable, type MRT_ColumnDef } from 'material-react-table'
import { Box, Button, IconButton, Tooltip, Chip, Typography } from '@mui/material'
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon, People as PeopleIcon } from '@mui/icons-material'
import { useTeams, useDeleteTeam } from '../../hooks/use-teams'
import { useAuth } from '../../contexts/AuthContext'
import type { Team, ListTeamsRequest } from '../../types/team.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface TeamListProps {
  onCreateTeam: () => void
  onEditTeam: (team: Team) => void
  onViewMembers: (team: Team) => void
  myTeamsOnly?: boolean
}

export const TeamList: React.FC<TeamListProps> = ({ onCreateTeam, onEditTeam, onViewMembers, myTeamsOnly = false }) => {
  const { user } = useAuth()
  const [searchParams, setSearchParams] = useSearchParams()
  const [pagination, setPagination] = useState({
    pageIndex: parseInt(searchParams.get('page') || '0', 10) || 0,
    pageSize: parseInt(searchParams.get('limit') || '10', 10) || 10,
  })

  useEffect(() => {
    const params = new URLSearchParams(searchParams)
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination.pageIndex, pagination.pageSize, setSearchParams])

  const request: ListTeamsRequest = {
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    myTeamsOnly,
  }

  const { data, isLoading, isError, error } = useTeams(request)
  const deleteTeamMutation = useDeleteTeam()

  const handleDeleteTeam = (team: Team) => {
    if (window.confirm(`Are you sure you want to delete "${team.name}"?`)) {
      deleteTeamMutation.mutate(team.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<Team>[]>(() => [
    { accessorKey: 'id', header: 'ID', size: 60 },
    {
      accessorKey: 'name',
      header: 'Team Name',
      size: 200,
      Cell: ({ row }) => (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Typography variant="body2" fontWeight="medium">{row.original.name}</Typography>
          {row.original.captainId === user?.id && <Chip label="Captain" size="small" color="primary" />}
        </Box>
      ),
    },
    {
      accessorKey: 'currentMembers',
      header: 'Members',
      size: 100,
      Cell: ({ row }) => (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <PeopleIcon sx={{ mr: 0.5, fontSize: 16, color: 'text.secondary' }} />
          <Typography variant="body2">{row.original.currentMembers} / {row.original.maxMembers}</Typography>
        </Box>
      ),
    },
    { accessorKey: 'inviteCode', header: 'Invite Code', size: 100, Cell: ({ cell }) => <code>{cell.getValue<string>()}</code> },
    { accessorKey: 'createdAt', header: 'Created', size: 120, Cell: ({ cell }) => formatRelativeTime(cell.getValue<string>()) },
  ], [user?.id])

  const table = useMaterialReactTable({
    columns,
    data: data?.teams ?? [],
    enableRowSelection: false,
    enableColumnOrdering: true,
    enableGlobalFilter: true,
    enablePagination: true,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: { isLoading, pagination, showAlertBanner: isError },
    muiToolbarAlertBannerProps: isError ? { color: 'error', children: `Error: ${error?.message}` } : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '0.5rem' }}>
        <Tooltip title="View Members">
          <IconButton color="info" onClick={() => onViewMembers(row.original)}><PeopleIcon /></IconButton>
        </Tooltip>
        {row.original.captainId === user?.id && (
          <>
            <Tooltip title="Edit Team">
              <IconButton color="primary" onClick={() => onEditTeam(row.original)}><EditIcon /></IconButton>
            </Tooltip>
            <Tooltip title="Delete Team">
              <IconButton color="error" onClick={() => handleDeleteTeam(row.original)} disabled={deleteTeamMutation.isPending}><DeleteIcon /></IconButton>
            </Tooltip>
          </>
        )}
      </Box>
    ),
    renderTopToolbarCustomActions: () => (
      <Button variant="contained" startIcon={<AddIcon />} onClick={onCreateTeam}>Create Team</Button>
    ),
    enableRowActions: true,
    positionActionsColumn: 'last',
  })

  return <MaterialReactTable table={table} />
}

export default TeamList
