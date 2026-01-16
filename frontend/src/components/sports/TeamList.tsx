import React, { useMemo, useState } from 'react'
import { MaterialReactTable, useMaterialReactTable, type MRT_ColumnDef } from 'material-react-table'
import { Box, Button, IconButton, Tooltip, Chip, FormControl, InputLabel, Select, MenuItem, Avatar } from '@mui/material'
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material'
import { useTeams, useDeleteTeam, useSports } from '../../hooks/use-sports'
import type { Team } from '../../types/sports.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface TeamListProps {
  onCreateTeam: () => void
  onEditTeam: (team: Team) => void
}

export const TeamList: React.FC<TeamListProps> = ({ onCreateTeam, onEditTeam }) => {
  const [sportFilter, setSportFilter] = useState<number | ''>('')
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 } })
  const sportsMap = useMemo(() => new Map(sportsData?.sports?.map(s => [s.id, s.name]) || []), [sportsData])

  const { data, isLoading, isError, error } = useTeams({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    sportId: sportFilter || undefined,
  })

  const deleteMutation = useDeleteTeam()

  const handleDelete = (team: Team) => {
    if (window.confirm(`Delete "${team.name}"?`)) {
      deleteMutation.mutate(team.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<Team>[]>(() => [
    { accessorKey: 'id', header: 'ID', size: 60 },
    {
      accessorKey: 'logoUrl',
      header: '',
      size: 50,
      Cell: ({ row }) => <Avatar src={row.original.logoUrl} sx={{ width: 32, height: 32 }}>{row.original.name[0]}</Avatar>,
    },
    { accessorKey: 'name', header: 'Name', size: 150 },
    { accessorKey: 'shortName', header: 'Short', size: 80 },
    {
      accessorKey: 'sportId',
      header: 'Sport',
      size: 100,
      Cell: ({ cell }) => sportsMap.get(cell.getValue<number>()) || '-',
    },
    { accessorKey: 'country', header: 'Country', size: 100 },
    {
      accessorKey: 'isActive',
      header: 'Status',
      size: 90,
      Cell: ({ cell }) => <Chip label={cell.getValue<boolean>() ? 'Active' : 'Inactive'} color={cell.getValue<boolean>() ? 'success' : 'default'} size="small" />,
    },
    {
      accessorKey: 'createdAt',
      header: 'Created',
      size: 110,
      Cell: ({ cell }) => formatRelativeTime(cell.getValue<string>()),
    },
  ], [sportsMap])

  const table = useMaterialReactTable({
    columns,
    data: data?.teams ?? [],
    enableRowSelection: false,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: { isLoading, pagination, showAlertBanner: isError },
    muiToolbarAlertBannerProps: isError ? { color: 'error', children: `Error: ${error?.message}` } : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '0.5rem' }}>
        <Tooltip title="Edit"><IconButton color="primary" onClick={() => onEditTeam(row.original)}><EditIcon /></IconButton></Tooltip>
        <Tooltip title="Delete"><IconButton color="error" onClick={() => handleDelete(row.original)} disabled={deleteMutation.isPending}><DeleteIcon /></IconButton></Tooltip>
      </Box>
    ),
    renderTopToolbarCustomActions: () => (
      <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
        <Button variant="contained" startIcon={<AddIcon />} onClick={onCreateTeam}>Add Team</Button>
        <FormControl size="small" sx={{ minWidth: 150 }}>
          <InputLabel>Filter by Sport</InputLabel>
          <Select value={sportFilter} label="Filter by Sport" onChange={(e) => setSportFilter(e.target.value as number | '')}>
            <MenuItem value="">All Sports</MenuItem>
            {sportsData?.sports?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}
          </Select>
        </FormControl>
      </Box>
    ),
    enableRowActions: true,
    positionActionsColumn: 'last',
  })

  return <MaterialReactTable table={table} />
}

export default TeamList
