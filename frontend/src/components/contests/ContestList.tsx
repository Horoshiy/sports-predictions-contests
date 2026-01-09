import React, { useMemo, useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
  type MRT_Row,
} from 'material-react-table'
import {
  Box,
  Button,
  IconButton,
  Tooltip,
  Chip,
  Typography,
} from '@mui/material'
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  People as PeopleIcon,
} from '@mui/icons-material'
import { useContests, useDeleteContest } from '../../hooks/use-contests'
import type { Contest } from '../../types/contest.types'
import { formatDate, formatRelativeTime } from '../../utils/date-utils'

interface ContestListProps {
  onCreateContest: () => void
  onEditContest: (contest: Contest) => void
  onViewParticipants: (contest: Contest) => void
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'draft':
      return 'default'
    case 'active':
      return 'success'
    case 'completed':
      return 'info'
    case 'cancelled':
      return 'error'
    default:
      return 'default'
  }
}

export const ContestList: React.FC<ContestListProps> = ({
  onCreateContest,
  onEditContest,
  onViewParticipants,
}) => {
  const [searchParams, setSearchParams] = useSearchParams()
  
  const [pagination, setPagination] = useState({
    pageIndex: parseInt(searchParams.get('page') || '0'),
    pageSize: parseInt(searchParams.get('limit') || '10'),
  })

  // Sync pagination with URL
  useEffect(() => {
    const params = new URLSearchParams(searchParams)
    params.set('page', pagination.pageIndex.toString())
    params.set('limit', pagination.pageSize.toString())
    setSearchParams(params, { replace: true })
  }, [pagination, searchParams, setSearchParams])

  const { data, isLoading, isError, error } = useContests({
    pagination: {
      page: pagination.pageIndex + 1,
      limit: pagination.pageSize,
    },
  })

  const deleteContestMutation = useDeleteContest()

  const handleDeleteContest = (contest: Contest) => {
    if (window.confirm(`Are you sure you want to delete "${contest.title}"?`)) {
      deleteContestMutation.mutate(contest.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<Contest>[]>(
    () => [
      {
        accessorKey: 'id',
        header: 'ID',
        size: 80,
      },
      {
        accessorKey: 'title',
        header: 'Title',
        size: 200,
        Cell: ({ cell }) => (
          <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
            {cell.getValue<string>()}
          </Typography>
        ),
      },
      {
        accessorKey: 'sportType',
        header: 'Sport',
        size: 120,
      },
      {
        accessorKey: 'status',
        header: 'Status',
        size: 100,
        Cell: ({ cell }) => (
          <Chip
            label={cell.getValue<string>()}
            color={getStatusColor(cell.getValue<string>())}
            size="small"
          />
        ),
      },
      {
        accessorKey: 'currentParticipants',
        header: 'Participants',
        size: 120,
        Cell: ({ row }) => (
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <PeopleIcon sx={{ mr: 0.5, fontSize: 16, color: 'text.secondary' }} />
            <Typography variant="body2">
              {row.original.currentParticipants}
              {row.original.maxParticipants > 0 && ` / ${row.original.maxParticipants}`}
            </Typography>
          </Box>
        ),
      },
      {
        accessorKey: 'startDate',
        header: 'Start Date',
        size: 120,
        Cell: ({ cell }) => formatDate(cell.getValue<string>()),
      },
      {
        accessorKey: 'endDate',
        header: 'End Date',
        size: 120,
        Cell: ({ cell }) => formatDate(cell.getValue<string>()),
      },
      {
        accessorKey: 'createdAt',
        header: 'Created',
        size: 120,
        Cell: ({ cell }) => formatRelativeTime(cell.getValue<string>()),
      },
    ],
    []
  )

  const table = useMaterialReactTable({
    columns,
    data: data?.contests ?? [],
    enableRowSelection: false,
    enableColumnOrdering: true,
    enableGlobalFilter: true,
    enableSorting: true,
    enablePagination: true,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: {
      isLoading,
      pagination,
      showAlertBanner: isError,
      showProgressBars: isLoading,
    },
    muiToolbarAlertBannerProps: isError
      ? {
          color: 'error',
          children: `Error loading contests: ${error?.message}`,
        }
      : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '1rem' }}>
        <Tooltip title="View Participants">
          <IconButton
            color="info"
            onClick={() => onViewParticipants(row.original)}
          >
            <PeopleIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Edit Contest">
          <IconButton
            color="primary"
            onClick={() => onEditContest(row.original)}
          >
            <EditIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Contest">
          <IconButton
            color="error"
            onClick={() => handleDeleteContest(row.original)}
            disabled={deleteContestMutation.isPending}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Box>
    ),
    renderTopToolbarCustomActions: () => (
      <Button
        variant="contained"
        startIcon={<AddIcon />}
        onClick={onCreateContest}
      >
        Create Contest
      </Button>
    ),
    enableRowActions: true,
    positionActionsColumn: 'last',
    muiTableContainerProps: {
      sx: {
        minHeight: '500px',
      },
    },
  })

  if (isError) {
    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <Typography color="error" variant="h6">
          Failed to load contests
        </Typography>
        <Typography color="text.secondary">
          {error?.message || 'An unknown error occurred'}
        </Typography>
      </Box>
    )
  }

  return <MaterialReactTable table={table} />
}

export default ContestList
