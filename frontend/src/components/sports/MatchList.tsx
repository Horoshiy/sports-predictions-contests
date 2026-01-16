import React, { useMemo, useState } from 'react'
import { MaterialReactTable, useMaterialReactTable, type MRT_ColumnDef } from 'material-react-table'
import { Box, Button, IconButton, Tooltip, Chip, FormControl, InputLabel, Select, MenuItem, Typography } from '@mui/material'
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material'
import { useMatches, useDeleteMatch, useLeagues, useTeams } from '../../hooks/use-sports'
import type { Match, MatchStatus } from '../../types/sports.types'
import { formatDateTime } from '../../utils/date-utils'

interface MatchListProps {
  onCreateMatch: () => void
  onEditMatch: (match: Match) => void
}

const statusColors: Record<MatchStatus, 'default' | 'warning' | 'success' | 'error' | 'info'> = {
  scheduled: 'default',
  live: 'warning',
  finished: 'success',
  cancelled: 'error',
  postponed: 'info',
}

export const MatchList: React.FC<MatchListProps> = ({ onCreateMatch, onEditMatch }) => {
  const [leagueFilter, setLeagueFilter] = useState<number | ''>('')
  const [statusFilter, setStatusFilter] = useState<MatchStatus | ''>('')
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data: leaguesData } = useLeagues({ pagination: { page: 1, limit: 100 } })
  const { data: teamsData } = useTeams({ pagination: { page: 1, limit: 200 } })

  const leaguesMap = useMemo(() => new Map(leaguesData?.leagues?.map(l => [l.id, l.name]) || []), [leaguesData])
  const teamsMap = useMemo(() => new Map(teamsData?.teams?.map(t => [t.id, t.name]) || []), [teamsData])

  const { data, isLoading, isError, error } = useMatches({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    leagueId: leagueFilter || undefined,
    status: statusFilter || undefined,
  })

  const deleteMutation = useDeleteMatch()

  const handleDelete = (match: Match) => {
    if (window.confirm('Delete this match?')) {
      deleteMutation.mutate(match.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<Match>[]>(() => [
    { accessorKey: 'id', header: 'ID', size: 60 },
    {
      id: 'matchup',
      header: 'Match',
      size: 250,
      Cell: ({ row }) => (
        <Typography variant="body2">
          {teamsMap.get(row.original.homeTeamId) || 'TBD'} vs {teamsMap.get(row.original.awayTeamId) || 'TBD'}
        </Typography>
      ),
    },
    {
      accessorKey: 'leagueId',
      header: 'League',
      size: 150,
      Cell: ({ cell }) => leaguesMap.get(cell.getValue<number>()) || '-',
    },
    {
      accessorKey: 'scheduledAt',
      header: 'Scheduled',
      size: 150,
      Cell: ({ cell }) => formatDateTime(cell.getValue<string>()),
    },
    {
      accessorKey: 'status',
      header: 'Status',
      size: 100,
      Cell: ({ cell }) => {
        const status = cell.getValue<MatchStatus>()
        return <Chip label={status} color={statusColors[status]} size="small" />
      },
    },
    {
      id: 'score',
      header: 'Score',
      size: 80,
      Cell: ({ row }) => row.original.status === 'finished' || row.original.status === 'live'
        ? `${row.original.homeScore} - ${row.original.awayScore}`
        : '-',
    },
  ], [leaguesMap, teamsMap])

  const table = useMaterialReactTable({
    columns,
    data: data?.matches ?? [],
    enableRowSelection: false,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: { isLoading, pagination, showAlertBanner: isError },
    muiToolbarAlertBannerProps: isError ? { color: 'error', children: `Error: ${error?.message}` } : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '0.5rem' }}>
        <Tooltip title="Edit"><IconButton color="primary" onClick={() => onEditMatch(row.original)}><EditIcon /></IconButton></Tooltip>
        <Tooltip title="Delete"><IconButton color="error" onClick={() => handleDelete(row.original)} disabled={deleteMutation.isPending}><DeleteIcon /></IconButton></Tooltip>
      </Box>
    ),
    renderTopToolbarCustomActions: () => (
      <Box sx={{ display: 'flex', gap: 2, alignItems: 'center', flexWrap: 'wrap' }}>
        <Button variant="contained" startIcon={<AddIcon />} onClick={onCreateMatch}>Add Match</Button>
        <FormControl size="small" sx={{ minWidth: 150 }}>
          <InputLabel>League</InputLabel>
          <Select value={leagueFilter} label="League" onChange={(e) => setLeagueFilter(e.target.value as number | '')}>
            <MenuItem value="">All Leagues</MenuItem>
            {leaguesData?.leagues?.map(l => <MenuItem key={l.id} value={l.id}>{l.name}</MenuItem>)}
          </Select>
        </FormControl>
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Status</InputLabel>
          <Select value={statusFilter} label="Status" onChange={(e) => setStatusFilter(e.target.value as MatchStatus | '')}>
            <MenuItem value="">All</MenuItem>
            <MenuItem value="scheduled">Scheduled</MenuItem>
            <MenuItem value="live">Live</MenuItem>
            <MenuItem value="finished">Finished</MenuItem>
            <MenuItem value="cancelled">Cancelled</MenuItem>
            <MenuItem value="postponed">Postponed</MenuItem>
          </Select>
        </FormControl>
      </Box>
    ),
    enableRowActions: true,
    positionActionsColumn: 'last',
  })

  return <MaterialReactTable table={table} />
}

export default MatchList
