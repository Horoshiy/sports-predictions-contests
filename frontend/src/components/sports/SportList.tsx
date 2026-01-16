import React, { useMemo, useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { MaterialReactTable, useMaterialReactTable, type MRT_ColumnDef } from 'material-react-table'
import { Box, Button, IconButton, Tooltip, Chip, Typography } from '@mui/material'
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material'
import { useSports, useDeleteSport } from '../../hooks/use-sports'
import type { Sport } from '../../types/sports.types'
import { formatRelativeTime } from '../../utils/date-utils'

interface SportListProps {
  onCreateSport: () => void
  onEditSport: (sport: Sport) => void
}

export const SportList: React.FC<SportListProps> = ({ onCreateSport, onEditSport }) => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [pagination, setPagination] = useState({
    pageIndex: parseInt(searchParams.get('page') || '0'),
    pageSize: parseInt(searchParams.get('limit') || '10'),
  })

  useEffect(() => {
    const params = new URLSearchParams()
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination.pageIndex, pagination.pageSize, setSearchParams])

  const { data, isLoading, isError, error } = useSports({
    pagination: { page: pagination.pageIndex + 1, limit: pagination.pageSize },
  })

  const deleteMutation = useDeleteSport()

  const handleDelete = (sport: Sport) => {
    if (window.confirm(`Delete "${sport.name}"?`)) {
      deleteMutation.mutate(sport.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<Sport>[]>(() => [
    { accessorKey: 'id', header: 'ID', size: 60 },
    { accessorKey: 'name', header: 'Name', size: 150 },
    { accessorKey: 'slug', header: 'Slug', size: 120 },
    {
      accessorKey: 'description',
      header: 'Description',
      size: 200,
      Cell: ({ cell }) => {
        const desc = cell.getValue<string>()
        return desc?.length > 50 ? `${desc.slice(0, 50)}...` : desc || '-'
      },
    },
    {
      accessorKey: 'isActive',
      header: 'Status',
      size: 100,
      Cell: ({ cell }) => (
        <Chip
          label={cell.getValue<boolean>() ? 'Active' : 'Inactive'}
          color={cell.getValue<boolean>() ? 'success' : 'default'}
          size="small"
        />
      ),
    },
    {
      accessorKey: 'createdAt',
      header: 'Created',
      size: 120,
      Cell: ({ cell }) => formatRelativeTime(cell.getValue<string>()),
    },
  ], [])

  const table = useMaterialReactTable({
    columns,
    data: data?.sports ?? [],
    enableRowSelection: false,
    enableColumnOrdering: true,
    enableGlobalFilter: true,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: { isLoading, pagination, showAlertBanner: isError },
    muiToolbarAlertBannerProps: isError ? { color: 'error', children: `Error: ${error?.message}` } : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '0.5rem' }}>
        <Tooltip title="Edit">
          <IconButton color="primary" onClick={() => onEditSport(row.original)}>
            <EditIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete">
          <IconButton color="error" onClick={() => handleDelete(row.original)} disabled={deleteMutation.isPending}>
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Box>
    ),
    renderTopToolbarCustomActions: () => (
      <Button variant="contained" startIcon={<AddIcon />} onClick={onCreateSport}>
        Add Sport
      </Button>
    ),
    enableRowActions: true,
    positionActionsColumn: 'last',
  })

  return <MaterialReactTable table={table} />
}

export default SportList
