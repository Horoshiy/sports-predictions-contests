import React, { useMemo, useState } from 'react'
import { MaterialReactTable, useMaterialReactTable, type MRT_ColumnDef } from 'material-react-table'
import { Box, Button, IconButton, Tooltip, Chip, FormControl, InputLabel, Select, MenuItem } from '@mui/material'
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material'
import { useLeagues, useDeleteLeague, useSports } from '../../hooks/use-sports'
import type { League } from '../../types/sports.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface LeagueListProps {
  onCreateLeague: () => void
  onEditLeague: (league: League) => void
}

export const LeagueList: React.FC<LeagueListProps> = ({ onCreateLeague, onEditLeague }) => {
  const [sportFilter, setSportFilter] = useState<number | ''>('')
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })

  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 } })
  const sportsMap = useMemo(() => new Map(sportsData?.sports?.map(s => [s.id, s.name]) || []), [sportsData])

  const { data, isLoading, isError, error } = useLeagues({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
    sportId: sportFilter || undefined,
  })

  const deleteMutation = useDeleteLeague()

  const handleDelete = (league: League) => {
    if (window.confirm(`Delete "${league.name}"?`)) {
      deleteMutation.mutate(league.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<League>[]>(() => [
    { accessorKey: 'id', header: 'ID', size: 60 },
    { accessorKey: 'name', header: 'Name', size: 180 },
    {
      accessorKey: 'sportId',
      header: 'Sport',
      size: 120,
      Cell: ({ cell }) => sportsMap.get(cell.getValue<number>()) || '-',
    },
    { accessorKey: 'country', header: 'Country', size: 100 },
    { accessorKey: 'season', header: 'Season', size: 100 },
    {
      accessorKey: 'isActive',
      header: 'Status',
      size: 100,
      Cell: ({ cell }) => (
        <Chip label={cell.getValue<boolean>() ? 'Active' : 'Inactive'} color={cell.getValue<boolean>() ? 'success' : 'default'} size="small" />
      ),
    },
    {
      accessorKey: 'createdAt',
      header: 'Created',
      size: 120,
      Cell: ({ cell }) => formatRelativeTime(cell.getValue<string>()),
    },
  ], [sportsMap])

  const table = useMaterialReactTable({
    columns,
    data: data?.leagues ?? [],
    enableRowSelection: false,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: { isLoading, pagination, showAlertBanner: isError },
    muiToolbarAlertBannerProps: isError ? { color: 'error', children: `Error: ${error?.message}` } : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '0.5rem' }}>
        <Tooltip title="Edit">
          <IconButton color="primary" onClick={() => onEditLeague(row.original)}><EditIcon /></IconButton>
        </Tooltip>
        <Tooltip title="Delete">
          <IconButton color="error" onClick={() => handleDelete(row.original)} disabled={deleteMutation.isPending}><DeleteIcon /></IconButton>
        </Tooltip>
      </Box>
    ),
    renderTopToolbarCustomActions: () => (
      <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
        <Button variant="contained" startIcon={<AddIcon />} onClick={onCreateLeague}>Add League</Button>
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

export default LeagueList
